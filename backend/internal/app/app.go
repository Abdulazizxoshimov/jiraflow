package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"

	api "github.com/jira-backend/jiraflow-backend/api"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/email"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/minio"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/rabbitmq"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/redis"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/websocket"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/casbin"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/config"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/hasher"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/token"
	"github.com/jira-backend/jiraflow-backend/internal/storage"
	"github.com/jira-backend/jiraflow-backend/internal/usecase"
	"github.com/jira-backend/jiraflow-backend/internal/worker"
)

func Run(cfg *config.Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// Gin logs → stdout (terminal). Release mode in production.
	gin.DefaultWriter = os.Stdout
	gin.DefaultErrorWriter = os.Stderr
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.New(
		cfg.App.LogLevel, "jiraflow", "v1",
		logger.WithLoki(cfg.Loki.URL, cfg.Loki.User, cfg.Loki.Password, map[string]string{
			"app":         "jiraflow-backend",
			"environment": cfg.App.Env,
		}),
	)
	defer logger.Cleanup(log)

	ctx := context.Background()
	if cfg.Loki.URL != "" {
		log.Info(ctx, "loki configured", logger.String("url", cfg.Loki.URL), logger.String("user", cfg.Loki.User))
	} else {
		log.Warn(ctx, "loki not configured — set LOKI_URL to enable log shipping")
	}

	// ── Sentry ────────────────────────────────────────────────────────────────

	if cfg.Sentry.DSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Sentry.DSN,
			Environment:      cfg.App.Env,
			TracesSampleRate: cfg.Sentry.TracesSampleRate,
			AttachStacktrace: true,
		}); err != nil {
			log.Warn(context.Background(), "sentry init failed", logger.SafeString("err", err.Error()))
		} else {
			defer sentry.Flush(3 * time.Second)
		}
	}

	// ── Infrastructure ────────────────────────────────────────────────────────

	pg, err := postgres.New(context.Background(), cfg, log)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer pg.DB.Close()

	if err := postgres.RunMigrations(cfg.Postgres, "migrations"); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}
	log.Info(context.Background(), "migrations applied")

	cache, err := redis.New(cfg.Redis)
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	minioClient, err := minio.New(cfg.Minio)
	if err != nil {
		log.Warn(context.Background(), "minio unavailable — file uploads disabled", logger.SafeString("err", err.Error()))
		minioClient = minio.NewNop()
	} else if err := minioClient.EnsureBucket(context.Background()); err != nil {
		log.Warn(context.Background(), "minio bucket check failed — file uploads disabled", logger.SafeString("err", err.Error()))
		minioClient = minio.NewNop()
	}

	var directSender email.Sender
	if cfg.Email.Enabled {
		directSender, err = email.New(cfg.Email, log)
		if err != nil {
			return fmt.Errorf("email: %w", err)
		}
		if pingErr := email.Ping(cfg.Email); pingErr != nil {
			log.Warn(context.Background(), "smtp unreachable at startup — emails will fail until fixed",
				logger.SafeString("err", pingErr.Error()),
				logger.String("host", cfg.Email.Host),
			)
		} else {
			log.Info(context.Background(), "smtp ping ok", logger.String("host", cfg.Email.Host))
		}
	} else {
		log.Warn(context.Background(), "email disabled (EMAIL_ENABLED=false) — all email sends will be no-ops")
		directSender = email.NewNopSender()
	}

	mq, err := rabbitmq.New(cfg.RabbitMQ.URL)
	if err != nil {
		log.Warn(context.Background(), "rabbitmq unavailable, falling back to direct email send",
			logger.SafeString("err", err.Error()))
		mq = nil
	} else {
		defer mq.Close()
	}

	var emailSender email.Sender
	if mq != nil {
		emailSender = email.NewQueuedSender(mq)
	} else {
		emailSender = directSender
	}

	hub := websocket.NewHub(log)

	tokenMaker := token.NewJWTMaker(
		[]byte(cfg.JWT.Secret),
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
		cache.Client(),
		"session:",
		log,
	)

	enforcer, err := casbin.NewEnforcer("internal/pkg/casbin/model.conf", "internal/pkg/casbin/policy.csv")
	if err != nil {
		log.Warn(context.Background(), "casbin enforcer not loaded, RBAC disabled")
		enforcer = nil
	}

	// ── Storage ───────────────────────────────────────────────────────────────

	store := storage.New(pg)

	// ── Use cases ─────────────────────────────────────────────────────────────

	uc := usecase.New(usecase.Deps{
		Store:                 store,
		TokenMaker:            tokenMaker,
		Hasher:                hasher.New(12),
		Minio:                 minioClient,
		Cache:                 cache,
		Log:                   log,
		Hub:                   hub,
		EmailSender:           emailSender,
		FrontendBaseURL:       cfg.App.FrontendBaseURL,
		GoogleClientID:        cfg.OAuth.GoogleClientID,
		GoogleClientSecret:    cfg.OAuth.GoogleClientSecret,
		GoogleRedirectURL:     cfg.OAuth.GoogleRedirectURL,
		TelegramBotToken:      cfg.Telegram.BotToken,
		TelegramBotUsername:   cfg.Telegram.BotUsername,
		TelegramWebhookURL:    cfg.Telegram.WebhookURL,
		TelegramWebhookSecret: cfg.Telegram.WebhookSecret,
	})

	// ── Handler ───────────────────────────────────────────────────────────────

	h := &handlers.Handler{
		Auth:                  uc.Auth,
		User:                  uc.User,
		Project:               uc.Project,
		ProjectMember:         uc.ProjectMember,
		Invite:                uc.Invite,
		Workflow:              uc.Workflow,
		Sprint:                uc.Sprint,
		Issue:                 uc.Issue,
		IssueLink:             uc.IssueLink,
		Worklog:               uc.Worklog,
		Component:             uc.Component,
		Version:               uc.Version,
		Label:                 uc.Label,
		CustomField:           uc.CustomField,
		Board:                 uc.Board,
		Comment:               uc.Comment,
		Attachment:            uc.Attachment,
		File:                  uc.File,
		Space:                 uc.Space,
		Page:                  uc.Page,
		PageVersion:           uc.PageVersion,
		PageTag:               uc.PageTag,
		PageView:              uc.PageView,
		PageTemplate:          uc.PageTemplate,
		PageRestriction:       uc.PageRestriction,
		PageExport:            uc.PageExport,
		InlineComment:         uc.InlineComment,
		Favorite:              uc.Favorite,
		Notification:          uc.Notification,
		Search:                uc.Search,
		Audit:                 uc.Audit,
		Hub:                   hub,
		IssueAssignee:         uc.IssueAssignee,
		PageReaction:          uc.PageReaction,
		Webhook:               uc.Webhook,
		PageLock:              uc.PageLock,
		PageMacro:             uc.PageMacro,
		IssuePageLink:         uc.IssuePageLink,
		ActivityFeed:          uc.ActivityFeed,
		OAuth:                 uc.OAuth,
		APIKey:                uc.APIKey,
		PermissionScheme:      uc.PermissionScheme,
		DataImport:            uc.DataImport,
		IssueVote:             uc.IssueVote,
		BlogPost:              uc.BlogPost,
		SavedFilter:           uc.SavedFilter,
		SpaceExport:           uc.SpaceExport,
		SpaceCategory:         uc.SpaceCategory,
		ContentProperty:       uc.ContentProperty,
		Blueprint:             uc.Blueprint,
		IssueType:             uc.IssueType,
		NotificationScheme:    uc.NotificationScheme,
		ProjectTemplate:       uc.ProjectTemplate,
		FieldConfiguration:    uc.FieldConfiguration,
		SecurityScheme:        uc.SecurityScheme,
		Automation:            uc.Automation,
		Telegram:              uc.Telegram,
		GitHub:                uc.GitHub,
		TelegramWebhookSecret: cfg.Telegram.WebhookSecret,
		GitHubWebhookSecret:   cfg.GitHub.WebhookSecret,
		AllowOpenRegistration: cfg.App.AllowOpenRegistration,
		FrontendBaseURL:       cfg.App.FrontendBaseURL,
	}

	// ── Background workers ────────────────────────────────────────────────────

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	digest := worker.NewDailyDigestWorker(store.Notification, store.User, emailSender, log)
	go digest.Run(workerCtx, 8) // fires at 08:00 UTC daily

	if mq != nil {
		emailWorker := worker.NewEmailWorker(mq, directSender, log)
		go emailWorker.Run(workerCtx)
	}

	// ── Server ────────────────────────────────────────────────────────────────

	readyFn := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := pg.DB.Ping(ctx); err != nil {
			return fmt.Errorf("postgres: %w", err)
		}
		if err := cache.Ping(ctx); err != nil {
			return fmt.Errorf("redis: %w", err)
		}
		return nil
	}

	if cfg.Telegram.BotToken != "" && cfg.Telegram.WebhookURL != "" {
		go func() {
			if err := uc.Telegram.SetupWebhook(context.Background()); err != nil {
				log.Warn(context.Background(), "telegram: webhook setup failed", logger.SafeString("err", err.Error()))
			} else {
				log.Info(context.Background(), "telegram: webhook registered", logger.String("url", cfg.Telegram.WebhookURL))
			}
		}()
	} else if cfg.Telegram.BotToken != "" {
		// No public URL — use long-polling (local development mode)
		log.Info(context.Background(), "telegram: no webhook URL — starting polling mode (local dev)")
		go uc.Telegram.StartPolling(workerCtx)
	}

	router := api.NewRouter(h, tokenMaker, enforcer, log, readyFn, store.Audit)
	srv := api.NewServer(":"+cfg.App.Port, router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(context.Background(), "server error", logger.SafeString("err", err.Error()))
		}
	}()

	log.Info(context.Background(), fmt.Sprintf("server started on :%s", cfg.App.Port))

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hub.Stop()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	log.Info(context.Background(), "server stopped")
	return nil
}
