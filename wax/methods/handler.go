package methods

import (
	"net/http"
        "github.com/gin-gonic/gin"
)


func Dispatch(c *gin.Context) {
        stage := c.Query("stage")

	c.String(http.StatusOK, stage)
}
