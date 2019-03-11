package methods

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nethesis/icaro/ade/ade-api/utils"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

func GetFeedbackPage(c *gin.Context) {
	var FeedbackPage models.FeedbackPage

	token := c.Param("token")
	adeToken := utils.GetAdeTokenFromToken(token)

	if adeToken.Id <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No token found!"})
		return
	}

	if !adeToken.FeedbackLeftTime.IsZero() {
		c.JSON(http.StatusForbidden, gin.H{"message": "Token expired"})
		return
	}

	hotspotPrefs := utils.GetHotspotPrefs(adeToken.HotspotId)

	FeedbackPage.HotspotName = hotspotPrefs["captive_2_title"]
	FeedbackPage.HotspotLogo = hotspotPrefs["captive_3_logo"]
	FeedbackPage.BgColor = hotspotPrefs["captive_7_background"]

	c.JSON(http.StatusOK, FeedbackPage)
}

func PostFeedbackResult(c *gin.Context) {

	var feedbackResult models.FeedbackResult

	token := c.Param("token")
	adeToken := utils.GetAdeTokenFromToken(token)
	hotspotPrefs := utils.GetHotspotPrefs(adeToken.HotspotId)

	if adeToken.Id <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No token found!"})
		return
	}

	if !adeToken.FeedbackLeftTime.IsZero() {
		c.JSON(http.StatusForbidden, gin.H{"message": "Token expired"})
		return
	}

	if err := c.BindJSON(&feedbackResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	if feedbackResult.Message != "" {
		if !utils.SendFeedBackMessageToOwner(adeToken, feedbackResult.Message, hotspotPrefs["captive_7_background"]) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Feedback submission failed."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
