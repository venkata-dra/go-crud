package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/venkata-dra/crud-go/initializers"
	"github.com/venkata-dra/crud-go/models"
)

func RequiredAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Retrieve the secret key from environment variables
	secret := []byte(os.Getenv("SECRET"))

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is signed using HMAC algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	exp, ok := claims["time"].(float64)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if float64(time.Now().Unix()) > exp {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	// Retrieve the user from the database using the token subject ("sub")
	var user models.User
	initializers.DB.First(&user, claims["sub"])

	// Check if the user exists
	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Attach the user information to the request context
	c.Set("user", user)

	// Continue with the next middleware or handler
	c.Next()
}
