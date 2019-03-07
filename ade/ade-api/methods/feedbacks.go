package methods

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFeedbackPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func PostFeedbackResult(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
