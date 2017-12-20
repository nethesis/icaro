package methods

import (
	"net/http"
        "github.com/gin-gonic/gin"
)

func Login(c *gin.Context, hotspotId string, user string, password string, challenge string) {
	c.String(http.StatusOK, "login")
}

