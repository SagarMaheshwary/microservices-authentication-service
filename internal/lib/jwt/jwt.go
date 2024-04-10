package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sagarmaheshwary/microservices-authentication-service/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/constants"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
)

func New(id uint, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	jwtConfig := config.Getjwt()

	expiry := time.Now().Add(time.Duration(jwtConfig.Expiry) * time.Second).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["username"] = username
	claims["exp"] = expiry
	claims["jti"] = uuid.New().String()

	return token.SignedString([]byte(jwtConfig.Secret))
}

func Parse(token string) (jwt.MapClaims, error) {
	decoded, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Getjwt().Secret), nil
	})

	if err != nil {
		log.Error("Invalid jwt token %v", err)

		return nil, err
	}

	claims, ok := decoded.Claims.(jwt.MapClaims)

	if !ok {
		log.Error("Token parse claims failed %v", claims)

		return nil, errors.New("token parse claims failed")
	}

	return claims, nil
}

func AddToBlacklist(jti string, expiry int64) error {
	key := fmt.Sprintf("%s:%s", constants.RDS_TOKEN_BLACKLIST, jti)
	expiry = expiry - time.Now().Unix()

	err := redis.Set(key, "", time.Duration(expiry)*time.Second)

	return err
}

func IsBlacklisted(jti string) bool {
	key := fmt.Sprintf("%s:%s", constants.RDS_TOKEN_BLACKLIST, jti)

	_, err := redis.Get(key)

	return err == nil
}
