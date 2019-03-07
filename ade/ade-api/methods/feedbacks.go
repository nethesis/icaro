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

	if !adeToken.FeedbackSendTime.IsZero() {
		c.JSON(http.StatusForbidden, gin.H{"message": "Token expired"})
		return
	}

	hotspotPerfs := utils.GetHotspotPrefs(adeToken.HotspotId)

	FeedbackPage.HotspotName = hotspotPerfs["captive_2_title"]
	FeedbackPage.HotspotLogo = hotspotPerfs["captive_3_logo"]
	FeedbackPage.BgColor = hotspotPerfs["captive_7_background"]

	c.JSON(http.StatusOK, FeedbackPage)
}

func PostFeedbackResult(c *gin.Context) {

	var feedbackResult models.FeedbackResult

	token := c.Param("token")
	adeToken := utils.GetAdeTokenFromToken(token)

	if adeToken.Id <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No token found!"})
		return
	}

	if !adeToken.FeedbackSendTime.IsZero() {
		c.JSON(http.StatusForbidden, gin.H{"message": "Token expired"})
		return
	}

	if err := c.BindJSON(&feedbackResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	if feedbackResult.Message != "" {
		utils.SendFeedBackMessage(adeToken.HotspotId, feedbackResult.Message)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
