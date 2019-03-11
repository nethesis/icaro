package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"time"

	gomail "gopkg.in/gomail.v2"

	"github.com/nethesis/icaro/ade/ade-api/models"
	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	wax_utils "github.com/nethesis/icaro/wax/utils"
)

type reviewResponse struct {
	Message string
	Stars   []int
	BgColor string
}

type feedbackResponse struct {
	Message string
	BgColor string
}

func GetHotspotPrefs(hotspotId int) map[string]string {
	hotspotPrefs := []string{"captive_2_title",
		"captive_3_logo",
		"captive_7_background",
		"marketing_9_threshold",
		"marketing_10_first_url",
		"marketing_11_second_url",
		"marketing_12_third_url",
	}

	prefsMap := make(map[string]string)

	prefs := wax_utils.GetHotspotPreferencesByKeys(hotspotId, hotspotPrefs)

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

func SendFeedBackMessage(adeToken models.AdeToken, message string, bgColor string) bool {
	var userMessage bytes.Buffer

	rp := feedbackResponse{
		Message: message,
		BgColor: bgColor,
	}

	t, err := template.ParseFiles("templates/feedback_owner.tpl")
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = t.Execute(&userMessage, &rp)
	if err != nil {
		fmt.Println(err)
		return false
	}

	status := SendEmail(adeToken.HotspotId, "User Feedback", userMessage.String())

	db := database.Instance()
	adeToken.FeedbackLeftTime = time.Now()
	db.Save(&adeToken)

	return status
}

func SendReviewMessage(adeToken models.AdeToken, stars int, message string, bgColor string) bool {
	var userMessage bytes.Buffer
	stars_s := strconv.Itoa(stars)

	starsFill := make([]int, stars)
	for i := range starsFill {
		starsFill[i] = 0
	}
	rp := reviewResponse{
		Message: message,
		Stars:   starsFill,
		BgColor: bgColor,
	}

	t, err := template.ParseFiles("templates/review_owner.tpl")
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = t.Execute(&userMessage, &rp)
	if err != nil {
		fmt.Println(err)
		return false
	}

	status := SendEmail(adeToken.HotspotId, "Review: "+stars_s+"/5", userMessage.String())

	db := database.Instance()
	adeToken.ReviewLeftTime = time.Now()
	db.Save(&adeToken)

	return status
}

func SendEmail(hotspotId int, subject string, message string) bool {

	status := true
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.Endpoints.Email.From)

	mail_to := wax_utils.GetHotspotPreferencesByKey(hotspotId, "marketing_2_feedback_email").Value

	m.SetHeader("To", mail_to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

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
