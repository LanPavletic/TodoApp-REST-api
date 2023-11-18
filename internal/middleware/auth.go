package middleware

import (
	"os"
	"strings"

	"github.com/LanPavletic/go-rest-server/internal/responds"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT_SECRET is the secret used to sign the JWT token
// we initalize this variable twice. Not ideal.
var JWT_SECRET = os.Getenv("JWT_SECRET")

// AuthRequired is a middleware that checks if the request contains a valid token
// The token is in authorization part of the header with a "Bearer" prefix
// The token is signed with a secret that will be encoded in the future
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// first part of tokenString is "Bearer"
		// check for it and remove it
		tokenString = strings.Trim(strings.TrimPrefix(tokenString, "Bearer"), " ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { return []byte(JWT_SECRET), nil })
		if token != nil && token.Valid {
			c.Next()
			return
		}

		// Here we can add additional checks for token validity for different cases
		switch err {
		default:
			responds.Unauthorized(c)
		}
	}

}
