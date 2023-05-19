package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	JwtIssuer   = "github.com/krissukoco/go-gin-chat"
	JwtAudience = "github.com/krissukoco/go-gin-chat"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
)

func JwtFromUsername(username string, secret string, durationHour ...int) (string, error) {
	durHour := 24 * 7
	if len(durationHour) > 0 {
		durHour = durationHour[0]
	}
	now := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": JwtIssuer,
		"sub": username,
		"aud": JwtAudience,
		"exp": now + (60 * 60 * int64(durHour)),
		"nbf": now,
		"iat": now,
		"jti": fmt.Sprintf("go-gin-chat-token_%s", uuid.NewString()),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetUsernameFromJwt(token, secret string) (string, error) {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return claims["sub"].(string), nil
	}
	return "", ErrTokenInvalid
}
