package pkg

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sagarmaheshwary/microservices-authentication-service/config"
)

func Createjwt(id uint, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	jwtConfig := config.Getjwt()

	expiry := time.Now().Add(time.Duration(jwtConfig.Expiry) * time.Second).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["username"] = username
	claims["exp"] = expiry

	return token.SignedString([]byte(jwtConfig.Secret))
}

func Parsejwt(token string) (jwt.MapClaims, error) {
	decoded, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Getjwt().Secret), nil
	})

	if err != nil {
		log.Println("Invalid token", err)

		return nil, err
	}

	claims, ok := decoded.Claims.(jwt.MapClaims)

	if !ok {
		fmt.Println("Unable to parse claims", claims)

		return nil, errors.New("token parse claims failed")
	}

	return claims, nil
}
