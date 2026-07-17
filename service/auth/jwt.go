package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/justKody/ecom-go/config"
	"github.com/justKody/ecom-go/types"
	"github.com/justKody/ecom-go/utils"
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

	tokenString, err := token.SignedString([]byte(config.Envs.JWTSecret))
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

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token from the user request
		tokenString, err := getTokenFromRequest(r)

		// validate the jwt
		claims, err := VerifyJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		// get the user id from the token
		userId := claims.UserId

		// get the user from the store
		user, err := store.GetUserById(userId)
		if err != nil {
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", user.ID)

		r = r.WithContext(ctx)
		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("unauthorized")
	}
	return tokenString, nil
}

func GetUserIdFromContext(ctx context.Context) (int, error) {
	userId, ok := ctx.Value("userId").(int)
	if !ok {
		return 0, errors.New("userId not found")
	}
	return userId, nil
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, errors.New("permission denied"))
}
