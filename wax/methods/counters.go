package methods

import (
	"net/http"
	"net/url"
        "github.com/gin-gonic/gin"
)

func Ack(c *gin.Context) {
	c.String(http.StatusOK, "Ack: 1")
}

func Counters(c *gin.Context, parameters url.Values) {
	status := parameters.Get("status")

	switch status {
        case "start":
		Ack(c)
        case "stop":
		Ack(c)
        case "update":
		Ack(c)
        case "":
                c.String(http.StatusBadRequest, "No status provided")
        default:
                c.String(http.StatusNotImplemented, "Invalid status: '%s'", status)
	}
}
