package methods

import (
	"net/http"
        "github.com/gin-gonic/gin"
)

func Dispatch(c *gin.Context) {
        stage := c.Query("stage")

	switch stage {
	case "login":
		hotspotId := c.Query("nasid")
		user := c.Query("user")
		password :=  c.Query("chap_pass")
		challenge := c.Query("chap_chal")
		Login(c, hotspotId, user, password, challenge)


	case "counters":
		parameters := c.Request.URL.Query()
		delete(parameters, "stage")

		// Pass all parameters except 'stage' key
		Counters(c, parameters)

	case "register":
		c.String(http.StatusNotImplemented, "Not implemented: %s", stage)

	case "":
		c.String(http.StatusBadRequest, "No stage provided")

	default:
		c.String(http.StatusNotFound, "Invalid stage: '%s'", stage)
	}

}
