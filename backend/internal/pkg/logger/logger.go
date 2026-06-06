package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LogFilePath = "logs/app.log"

// OpenLogFile opens (or creates) the app log file for writing.
func OpenLogFile() (*os.File, error) {
	if err := os.MkdirAll("logs", 0o755); err != nil {
		return nil, err
	}
	return os.OpenFile(LogFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
}

// ─── Levels ───────────────────────────────────────────────────────────────────

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelPanic = "panic"
	LevelFatal = "fatal"
)

// LogLevelFromString ...
func LogLevelFromString(level string) int {
	switch level {
	case LevelDebug:
		return -1
	case LevelInfo:
		return 0
	case LevelWarn:
		return 1
	case LevelError:
		return 2
	case LevelPanic:
		return 4
	case LevelFatal:
		return 5
	default:
		return 0
	}
}

// ─── Field helpers ────────────────────────────────────────────────────────────

type Field = zapcore.Field

var (
	Int    = zap.Int
	Int64  = zap.Int64
	String = zap.String
	Error  = zap.Error
	Bool   = zap.Bool
	Any    = zap.Any
)

// SafeEmail masks email before logging.
// john@example.com → j***@example.com
func SafeEmail(key, email string) Field {
	return zap.String(key, MaskEmail(email))
}

// SafePhone masks phone before logging.
// +998901234567 → ***********67
func SafePhone(key, phone string) Field {
	return zap.String(key, MaskPhone(phone))
}

// SafeString applies full PII masking to any free-form string.
func SafeString(key, value string) Field {
	return zap.String(key, MaskPII(value))
}

// ─── Logger interface ─────────────────────────────────────────────────────────

// FIX 1: Barcha metodlar endi context.Context qabul qiladi.
// Bu trace_id, span_id, user_id ni har bir logda majburiy qiladi.
// Developer context uzatmasa — kompilyator xato beradi.
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)
}

type loggerImpl struct {
	zap *zap.Logger
}

// New creates a structured JSON logger.
//
//	log := logger.New("info", "grpc", "payment-service")
func New(level, namespace, serviceName string) Logger {
	if level == "" {
		level = LevelInfo
	}

	l := loggerImpl{
		zap: newZapLogger(level),
	}
	l.zap = l.zap.Named(namespace)

	if serviceName != "" {
		l.zap = l.zap.With(zap.String("service_name", serviceName))
	}

	zap.RedirectStdLog(l.zap)
	return &l
}

// FIX 1 + FIX 2: Har bir metod ctx dan trace_id, span_id, user_id ni
// avtomatik o'qiydi. Interceptor endi faqat contextga yozadi,
// logger o'zi o'qiydi — yagona manba (single source of truth).
func (l *loggerImpl) Debug(ctx context.Context, msg string, fields ...Field) {
	l.zap.Debug(msg, l.fromCtx(ctx, fields)...)
}

func (l *loggerImpl) Info(ctx context.Context, msg string, fields ...Field) {
	l.zap.Info(msg, l.fromCtx(ctx, fields)...)
}

func (l *loggerImpl) Warn(ctx context.Context, msg string, fields ...Field) {
	l.zap.Warn(msg, l.fromCtx(ctx, fields)...)
}

func (l *loggerImpl) Error(ctx context.Context, msg string, fields ...Field) {
	l.zap.Error(msg, l.fromCtx(ctx, fields)...)
}

func (l *loggerImpl) Fatal(ctx context.Context, msg string, fields ...Field) {
	l.zap.Fatal(msg, l.fromCtx(ctx, fields)...)
}

// fromCtx — contextdan majburiy fieldlarni olib, extra fieldlar bilan birlashtiradi.
// Bu FIX 2 ning yuragi: trace ma'lumoti faqat contextdan keladi.
func (l *loggerImpl) fromCtx(ctx context.Context, extra []Field) []Field {
	fields := make([]Field, 0, 3+len(extra))

	if v, _ := ctx.Value(traceIDKey).(string); v != "" {
		fields = append(fields, zap.String("trace_id", v))
	}
	if v, _ := ctx.Value(spanIDKey).(string); v != "" {
		fields = append(fields, zap.String("span_id", v))
	}
	if v, _ := ctx.Value(userIDKey).(string); v != "" {
		fields = append(fields, zap.String("user_id", v))
	}

	return append(fields, extra...)
}

// ─── Logger utility functions ─────────────────────────────────────────────────

// GetNamed ...
func GetNamed(l Logger, name string) Logger {
	switch v := l.(type) {
	case *loggerImpl:
		return &loggerImpl{zap: v.zap.Named(name)}
	default:
		l.Info(context.Background(), "logger.GetNamed: invalid logger type")
		return l
	}
}

// WithFields returns a child logger with permanent extra fields.
func WithFields(l Logger, fields ...Field) Logger {
	switch v := l.(type) {
	case *loggerImpl:
		return &loggerImpl{zap: v.zap.With(fields...)}
	default:
		l.Info(context.Background(), "logger.WithFields: invalid logger type")
		return l
	}
}

// Cleanup flushes buffered log entries. Call on shutdown.
func Cleanup(l Logger) error {
	switch v := l.(type) {
	case *loggerImpl:
		return v.zap.Sync()
	default:
		l.Info(context.Background(), "logger.Cleanup: invalid logger type")
		return nil
	}
}

// GetZapLogger extracts raw zap logger (for third-party integrations).
func GetZapLogger(l Logger) *zap.Logger {
	if l == nil {
		return newZapLogger(LevelInfo)
	}
	switch v := l.(type) {
	case *loggerImpl:
		return v.zap
	default:
		return newZapLogger(LevelInfo)
	}
}

// ─── Zap internals ────────────────────────────────────────────────────────────

func newZapLogger(level string) *zap.Logger {
	globalLevel := parseLevel(level)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.MessageKey = "message"
	encoderCfg.LevelKey = "level"
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(time.RFC3339Nano))
	}
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

	consoleCfg := encoderCfg
	consoleCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)

	atomicLevel := zap.NewAtomicLevelAt(globalLevel)

	cores := []zapcore.Core{
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), atomicLevel),
	}
	if f, err := OpenLogFile(); err == nil {
		cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(f), atomicLevel))
	}

	return zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelPanic:
		return zapcore.PanicLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
