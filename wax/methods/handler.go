package methods

import (
	"net/http"
        "github.com/gin-gonic/gin"
)


func Dispatch(c *gin.Context) {
        stage := c.Query("stage")

	switch stage {
	case "login":
		c.String(http.StatusOK, stage)

	case "counters":
		c.String(http.StatusOK, stage)

	case "register":
		c.String(http.StatusNotImplemented, "Not implemented: %s", stage)

	case "":
		c.String(http.StatusBadRequest, "No stage provided")

	default:
		c.String(http.StatusNotFound, "Invalid stage: '%s'", stage)
	}

}
