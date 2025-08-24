package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
)

type JWTManager interface {
	NewToken(id uint, username string) (string, error)
	ParseToken(token string) (jwt.MapClaims, error)
	AddToBlacklist(ctx context.Context, jti string, expiry int64) error
	IsBlacklisted(ctx context.Context, jti string) bool
}

type jwtManager struct {
	secret []byte
	expiry time.Duration
	redis  redis.RedisService
}

func NewJWTManager(cfg *config.JWT, redis redis.RedisService) JWTManager {
	return &jwtManager{
		secret: []byte(cfg.Secret),
		expiry: cfg.Expiry,
		redis:  redis,
	}
}

func (j *jwtManager) NewToken(id uint, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expiry := time.Now().Add(j.expiry).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["username"] = username
	claims["exp"] = expiry
	claims["jti"] = uuid.New().String()

	return token.SignedString(j.secret)
}

func (j *jwtManager) ParseToken(token string) (jwt.MapClaims, error) {
	decoded, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := decoded.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token parse claims failed")
	}

	return claims, nil
}

func (j *jwtManager) AddToBlacklist(ctx context.Context, jti string, expiry int64) error {
	key := fmt.Sprintf("%s:%s", constant.RedisTokenBlacklist, jti)
	exp := expiry - time.Now().Unix()
	if exp <= 0 {
		return nil // already expired
	}
	return j.redis.Set(ctx, key, "", time.Duration(exp)*time.Second)
}

func (j *jwtManager) IsBlacklisted(ctx context.Context, jti string) bool {
	key := fmt.Sprintf("%s:%s", constant.RedisTokenBlacklist, jti)
	_, err := j.redis.Get(ctx, key)
	return err == nil
}
