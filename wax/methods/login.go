package methods

import (
	"net/http"
        "github.com/gin-gonic/gin"
)

func AuthAccept(c *gin.Context) {
	c.String(http.StatusOK, "Auth: 1")
}

func AuthReject(c *gin.Context) {
	c.String(http.StatusForbidden, "Auth: 0")
}

func Login(c *gin.Context, hotspotId string, user string, password string, challenge string) {
	AuthAccept(c)
}

