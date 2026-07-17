package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/justKody/ecom-go/config"
)

type Claims struct {
	UserId int `json:"userId"`
	jwt.RegisteredClaims
}

func CreateJWT(userId int) (string, error) {
	expirationDuration := time.Duration(config.Envs.JWTExpiration) * time.Minute

	log.Println("expirationDuration:", expirationDuration)
	// create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
		},
	})

	log.Println("token:", token)
	tokenString, err := token.SignedString([]byte(config.Envs.JWTSecret))
	log.Println("tokenString:", tokenString)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyJWT verifies a JWT token by returning the claims
func VerifyJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Envs.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*Claims), nil
}
