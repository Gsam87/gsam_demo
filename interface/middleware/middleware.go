package middlerware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/qww83728/gsam_demo/util"
)

type Middleware interface {
	GenerateToken(userEmail string, updateTime time.Time) (string, error)
	JWTMiddleware() gin.HandlerFunc
}

type MiddlewareImpl struct {
}

func NewMiddleware() Middleware {
	return &MiddlewareImpl{}
}

var jwtSecret = []byte("your-secret-key") // JWT 密鑰

// 產生 JWT Token
func (m *MiddlewareImpl) GenerateToken(
	userEmail string,
	updateTime time.Time,
) (string, error) {
	claims := jwt.MapClaims{
		"email":   userEmail,
		"updated": updateTime,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWT 驗證 Middleware
func (m *MiddlewareImpl) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, util.MakeFailResponse(
				http.StatusUnauthorized,
				"Authorization header required",
				errors.New("authorization header required"),
			))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, util.MakeFailResponse(
				http.StatusUnauthorized,
				"Authorization header format must be Bearer {token}",
				errors.New("authorization header format must be Bearer {token}"),
			))
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, util.MakeFailResponse(
				http.StatusUnauthorized,
				"Invalid token",
				errors.New("invalid token"),
			))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, util.MakeFailResponse(
				http.StatusUnauthorized,
				"Invalid token claims",
				errors.New("invalid token claims"),
			))
			c.Abort()
			return
		}

		// 可以把 email 存到 context
		c.Set("email", claims["email"])
		c.Next()
	}
}
