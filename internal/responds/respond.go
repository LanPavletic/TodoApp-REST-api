package responds

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Helper functions for responding with different status codes.
// In the future we should add more information to the response body
// (e.g. unique error codes, error descriptions, etc.)
func NotFound(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
}

func BadRequest(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
}

func InternalServerError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}

func Unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}
