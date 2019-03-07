package methods

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nethesis/icaro/ade/ade-api/models"
	"github.com/nethesis/icaro/ade/ade-api/utils"
)

func GetFeedbackPage(c *gin.Context) {
	var FeedbackPage models.FeedbackPage

	token := c.Param("token")
	adeToken := utils.GetAdeTokenFromToken(token)

	if adeToken.Id <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No token found!"})
		return
	}

	hotspotPerfs := utils.GetHotspotPrefs(adeToken.HotspotId)

	FeedbackPage.HotspotName = hotspotPerfs["captive_2_title"]
	FeedbackPage.HotspotLogo = hotspotPerfs["captive_3_logo"]
	FeedbackPage.BgColor = hotspotPerfs["captive_7_background"]

	c.JSON(http.StatusOK, FeedbackPage)
}

func PostFeedbackResult(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
