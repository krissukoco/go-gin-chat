package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-gin-chat/schema"
	"github.com/krissukoco/go-gin-chat/security"
)

type AuthMiddleware struct {
	JwtSecret string
}

func (a *AuthMiddleware) AuthorizationHeader(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, &schema.ErrorResponse{
			Code:    schema.ErrAuthenticationRequired,
			Message: "Authentication required",
		})
		c.Abort()
		return
	}
	split := strings.Split(authHeader, " ")
	if len(split) != 2 {
		c.JSON(401, &schema.ErrorResponse{
			Code:    schema.ErrAuthenticationRequired,
			Message: "Invalid authentication header",
		})
		c.Abort()
		return
	}
	if split[0] != "Bearer" {
		c.JSON(401, &schema.ErrorResponse{
			Code:    schema.ErrAuthenticationRequired,
			Message: "Invalid authentication method",
		})
		c.Abort()
		return
	}
	token := split[1]
	// Get username from token
	username, err := security.GetUsernameFromJwt(token, a.JwtSecret)
	if err != nil {
		c.JSON(401, &schema.ErrorResponse{
			Code:    schema.ErrTokenInvalid,
			Message: "Invalid token",
		})
		c.Abort()
		return
	}

	c.Set("username", username)

	c.Next()
}
