package pkg

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/sagarmaheshwary/microservices-authentication-service/config"
)

func CreateJwt(id uint, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["username"] = username
	claims["exp"] = config.Getjwt().Expiry

	return token.SignedString([]byte(config.Getjwt().Secret))
}
