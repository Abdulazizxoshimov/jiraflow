package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig      `json:"app"`
	Postgres PostgresConfig `json:"postgres"`
	Redis    RedisConfig    `json:"redis"`
	Minio    MinioConfig    `json:"minio"`
	RabbitMQ RabbitMQConfig `json:"rabbitmq"`
	JWT      JWTConfig      `json:"jwt"`
	Email    EmailConfig    `json:"email"`
	OAuth    OAuthConfig    `json:"oauth"`
	Telegram TelegramConfig `json:"telegram"`
	GitHub   GitHubConfig   `json:"github"`
}

type TelegramConfig struct {
	BotToken      string `json:"bot_token"`
	WebhookURL    string `json:"webhook_url"`
	WebhookSecret string `json:"webhook_secret"`
}

type GitHubConfig struct {
	WebhookSecret string `json:"webhook_secret"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
}

type OAuthConfig struct {
	GoogleClientID     string `json:"google_client_id"`
	GoogleClientSecret string `json:"google_client_secret"`
	GoogleRedirectURL  string `json:"google_redirect_url"`
}

type AppConfig struct {
	Port     string `json:"port"`
	Env      string `json:"env"`
	LogLevel string `json:"log_level"`
}

type PostgresConfig struct {
	Host              string        `json:"host"`
	Port              string        `json:"port"`
	Database          string        `json:"database"`
	Username          string        `json:"username"`
	Password          string        `json:"password"`
	MaxConns          int32         `json:"max_conns"`
	MinConns          int32         `json:"min_conns"`
	MaxConnIdleTime   time.Duration `json:"max_conn_idle_time"`
	MaxConnLifetime   time.Duration `json:"max_conn_lifetime"`
	HealthCheckPeriod time.Duration `json:"health_check_period"`
}

// DSN returns a pgxpool-compatible connection string.
func (p PostgresConfig) DSN() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return "host=" + p.Host +
		" port=" + p.Port +
		" dbname=" + p.Database +
		" user=" + p.Username +
		" password=" + p.Password +
		" sslmode=disable"
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Addr returns "host:port" for redis.NewClient.
func (r RedisConfig) Addr() string {
	return r.Host + ":" + r.Port
}

type MinioConfig struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	UseSSL    bool   `json:"use_ssl"`
}

type RabbitMQConfig struct {
	URL string `json:"url"`
}

type JWTConfig struct {
	Secret     string        `json:"secret"`
	AccessTTL  time.Duration `json:"access_ttl"`
	RefreshTTL time.Duration `json:"refresh_ttl"`
}

type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	From     string `json:"from"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Load reads environment variables into Config.
// In local development, set values via .env file loaded externally or shell exports.
func Load() *Config {
	return &Config{
		App: AppConfig{
			Port:     getEnv("APP_PORT", "8080"),
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Postgres: PostgresConfig{
			Host:              getEnv("DB_HOST", "localhost"),
			Port:              getEnv("DB_PORT", "5432"),
			Database:          getEnv("DB_NAME", "jiraflow"),
			Username:          getEnv("DB_USER", "postgres"),
			Password:          getEnv("DB_PASSWORD", "4444"),
			MaxConns:          int32(getEnvInt("DB_MAX_CONNS", 50)),
			MinConns:          int32(getEnvInt("DB_MIN_CONNS", 5)),
			MaxConnIdleTime:   getEnvDuration("DB_MAX_CONN_IDLE_TIME", 5*time.Minute),
			MaxConnLifetime:   getEnvDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute),
			HealthCheckPeriod: getEnvDuration("DB_HEALTH_CHECK_PERIOD", 1*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Minio: MinioConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "jiraflow"),
			UseSSL:    getEnvBool("MINIO_USE_SSL", false),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			AccessTTL:  getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL: getEnvDuration("JWT_REFRESH_TTL", 720*time.Hour),
		},
		Email: EmailConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnvInt("SMTP_PORT", 587),
			From:     getEnv("SMTP_FROM", ""),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
		},
		Telegram: TelegramConfig{
			BotToken:      getEnv("TELEGRAM_BOT_TOKEN", ""),
			WebhookURL:    getEnv("TELEGRAM_WEBHOOK_URL", ""),
			WebhookSecret: getEnv("TELEGRAM_WEBHOOK_SECRET", ""),
		},
		GitHub: GitHubConfig{
			WebhookSecret: getEnv("GITHUB_WEBHOOK_SECRET", ""),
			ClientID:      getEnv("GITHUB_CLIENT_ID", ""),
			ClientSecret:  getEnv("GITHUB_CLIENT_SECRET", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func getEnvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
