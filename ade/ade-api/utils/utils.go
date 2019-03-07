package utils

import (
	"fmt"
	"strconv"

	gomail "gopkg.in/gomail.v2"

	"github.com/nethesis/icaro/ade/ade-api/models"
	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	wax_utils "github.com/nethesis/icaro/wax/utils"
)

func GetHotspotPrefs(hotspotId int) map[string]string {
	hotspotPerfs := []string{"captive_2_title",
		"captive_3_logo",
		"captive_7_background",
		"marketing_9_threshold",
		"marketing_10_first_url",
		"marketing_11_second_url",
		"marketing_12_third_url",
	}

	prefsMap := make(map[string]string)

	prefs := wax_utils.GetHotspotPreferencesByKeys(hotspotId, hotspotPerfs)

	for i := 0; i < len(prefs); i++ {
		prefsMap[prefs[i].Key] = prefs[i].Value
	}

	return prefsMap
}

func GetAdeTokenFromToken(token string) models.AdeToken {
	var adeToken models.AdeToken

	db := database.Instance()
	db.Where("token = ?", token).First(&adeToken)

	return adeToken
}

func SendFeedBackMessage(hotspotId int, message string) bool {
	return SendEmail(hotspotId, "User Feedback", message)
}

func SendReviewMessage(hotspotId int, stars int, message string) bool {
	stars_s := strconv.Itoa(stars)
	return SendEmail(hotspotId, "Review: "+stars_s+"/5", message)
}

func SendEmail(hotspotId int, subject string, message string) bool {

	status := true
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.Endpoints.Email.From)

	mail_to := wax_utils.GetHotspotPreferencesByKey(hotspotId, "marketing_2_feedback_email").Value

	m.SetHeader("To", mail_to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	d := gomail.NewDialer(
		configuration.Config.Endpoints.Email.SMTPHost,
		configuration.Config.Endpoints.Email.SMTPPort,
		configuration.Config.Endpoints.Email.SMTPUser,
		configuration.Config.Endpoints.Email.SMTPPassword,
	)

	// send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		status = false
	}

	return status
}
