package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

// Maker is the interface for token generation and validation.
type Maker interface {
	Generate(ctx context.Context, sub, sessionID, role string) (access, refresh string, err error)
	ValidateAccess(ctx context.Context, token string) (*Claims, error)
	Rotate(ctx context.Context, refreshToken string) (newAccess, newRefresh string, err error)
	Revoke(ctx context.Context, jti string) error
	StoreSession(ctx context.Context, sessionID, payload string, ttl time.Duration) error
	RevokeSession(ctx context.Context, sessionID string) error
}

type JWTMaker struct {
	signingKey  []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
	redis       *redis.Client
	redisPrefix string
	log         logger.Logger
}

func NewJWTMaker(
	signingKey []byte,
	accessTTL, refreshTTL time.Duration,
	rdb *redis.Client,
	prefix string,
	log logger.Logger,
) Maker {
	if prefix == "" {
		prefix = "auth"
	}
	return &JWTMaker{
		signingKey:  signingKey,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
		redis:       rdb,
		redisPrefix: prefix,
		log:         log,
	}
}

// Generate issues a new access+refresh token pair for the given subject/session/role.
// The refresh token's JTI is stored in Redis so it can be rotated or revoked.
func (m *JWTMaker) Generate(ctx context.Context, sub, sessionID, role string) (string, string, error) {
	now := time.Now().UTC()

	access, err := m.sign(jwt.MapClaims{
		"sub":  sub,
		"sid":  sessionID,
		"role": role,
		"type": "access",
		"iat":  now.Unix(),
		"exp":  now.Add(m.accessTTL).Unix(),
		"jti":  uuid.NewString(),
	})
	if err != nil {
		m.log.Error(ctx, "failed to sign access token", logger.Error(err))
		return "", "", fmt.Errorf("sign access token: %w", err)
	}

	jti := uuid.NewString()
	refresh, err := m.sign(jwt.MapClaims{
		"sub":  sub,
		"sid":  sessionID,
		"role": role,
		"type": "refresh",
		"iat":  now.Unix(),
		"exp":  now.Add(m.refreshTTL).Unix(),
		"jti":  jti,
	})
	if err != nil {
		m.log.Error(ctx, "failed to sign refresh token", logger.Error(err))
		return "", "", fmt.Errorf("sign refresh token: %w", err)
	}

	if err := m.redis.Set(ctx, m.refreshKey(jti), sessionID, m.refreshTTL).Err(); err != nil {
		m.log.Error(ctx, "failed to persist refresh jti", logger.Error(err))
		return "", "", fmt.Errorf("persist refresh jti: %w", err)
	}

	return access, refresh, nil
}

// ValidateAccess parses and validates an access token, then confirms the session
// is still active in Redis. Returns parsed Claims on success.
func (m *JWTMaker) ValidateAccess(ctx context.Context, tokenStr string) (*Claims, error) {
	raw, err := m.parse(tokenStr, "access")
	if err != nil {
		return nil, err
	}

	sid, _ := raw["sid"].(string)
	if sid == "" {
		return nil, errors.New("token: missing sid claim")
	}

	exists, err := m.redis.Exists(ctx, m.sessionKey(sid)).Result()
	if err != nil {
		m.log.Error(ctx, "redis: session check failed", logger.Error(err))
		return nil, fmt.Errorf("session check: %w", err)
	}
	if exists == 0 {
		return nil, errors.New("token: session revoked or not found")
	}

	return mapToClaims(raw), nil
}

// Rotate validates a refresh token, atomically swaps the JTI in Redis, and returns
// a new access+refresh pair. Signing happens before Redis mutation to avoid partial state.
func (m *JWTMaker) Rotate(ctx context.Context, oldRefresh string) (string, string, error) {
	raw, err := m.parse(oldRefresh, "refresh")
	if err != nil {
		return "", "", err
	}

	jti, _ := raw["jti"].(string)
	sid, _ := raw["sid"].(string)
	sub, _ := raw["sub"].(string)
	role, _ := raw["role"].(string)

	if jti == "" || sid == "" || sub == "" {
		return "", "", errors.New("token: refresh token missing required claims")
	}

	storedSid, err := m.redis.Get(ctx, m.refreshKey(jti)).Result()
	if errors.Is(err, redis.Nil) {
		return "", "", errors.New("token: refresh token already used or revoked")
	}
	if err != nil {
		m.log.Error(ctx, "redis: get refresh jti failed", logger.Error(err))
		return "", "", fmt.Errorf("redis get: %w", err)
	}
	if storedSid != sid {
		// sid mismatch is a sign of token theft — log it prominently.
		m.log.Warn(ctx, "token: jti/sid mismatch — possible token theft",
			logger.String("sub", sub), logger.String("sid", sid))
		return "", "", errors.New("token: session mismatch")
	}

	// Sign before touching Redis: if signing fails, Redis state stays consistent.
	now := time.Now().UTC()
	newJTI := uuid.NewString()

	access, err := m.sign(jwt.MapClaims{
		"sub":  sub,
		"sid":  sid,
		"role": role,
		"type": "access",
		"iat":  now.Unix(),
		"exp":  now.Add(m.accessTTL).Unix(),
		"jti":  uuid.NewString(),
	})
	if err != nil {
		return "", "", fmt.Errorf("sign access token: %w", err)
	}

	refresh, err := m.sign(jwt.MapClaims{
		"sub":  sub,
		"sid":  sid,
		"role": role,
		"type": "refresh",
		"iat":  now.Unix(),
		"exp":  now.Add(m.refreshTTL).Unix(),
		"jti":  newJTI,
	})
	if err != nil {
		return "", "", fmt.Errorf("sign refresh token: %w", err)
	}

	// Atomically set new JTI and delete old JTI.
	_, err = m.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, m.refreshKey(newJTI), sid, m.refreshTTL)
		pipe.Del(ctx, m.refreshKey(jti))
		return nil
	})
	if err != nil {
		m.log.Error(ctx, "redis: jti rotation failed", logger.Error(err))
		return "", "", fmt.Errorf("jti rotation: %w", err)
	}

	return access, refresh, nil
}

// Revoke deletes a refresh token's JTI from Redis (use on logout).
func (m *JWTMaker) Revoke(ctx context.Context, jti string) error {
	return m.redis.Del(ctx, m.refreshKey(jti)).Err()
}

// StoreSession persists a session marker in Redis with the given TTL.
func (m *JWTMaker) StoreSession(ctx context.Context, sessionID, payload string, ttl time.Duration) error {
	return m.redis.Set(ctx, m.sessionKey(sessionID), payload, ttl).Err()
}

// RevokeSession removes the session key from Redis, invalidating all access tokens
// that reference it.
func (m *JWTMaker) RevokeSession(ctx context.Context, sessionID string) error {
	return m.redis.Del(ctx, m.sessionKey(sessionID)).Err()
}

// sign creates and signs a JWT with HS256.
func (m *JWTMaker) sign(claims jwt.MapClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.signingKey)
}

// parse verifies signature, expiry, and token type.
func (m *JWTMaker) parse(tokenStr, wantType string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.signingKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parse: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token: invalid")
	}

	if typ, _ := claims["type"].(string); typ != wantType {
		return nil, fmt.Errorf("token: expected type %q, got %q", wantType, typ)
	}

	return claims, nil
}

func (m *JWTMaker) refreshKey(jti string) string { return m.redisPrefix + ":refresh:" + jti }
func (m *JWTMaker) sessionKey(sid string) string  { return m.redisPrefix + ":sess:" + sid }
