package utils

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"strconv"
	"time"

	gomail "gopkg.in/gomail.v2"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	wax_utils "github.com/nethesis/icaro/wax/utils"
)

type reviewSend struct {
	HotspotName    string
	HotspotLogo    string
	HotspotDetails string
	BgColor        string
	Url            string
}

type feedbackSend struct {
	HotspotName    string
	HotspotLogo    string
	HotspotDetails string
	BgColor        string
	Url            string
}

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

func CreateAdeToken(user models.User) models.AdeToken {
	h := sha256.New()
	h.Write([]byte(time.Now().UTC().String() + user.Username + user.Password))
	token := fmt.Sprintf("%x", h.Sum(nil))

	adeToken := models.AdeToken{
		Token:     token,
		UserId:    user.Id,
		HotspotId: user.HotspotId,
	}

	db := database.Instance()
	db.Save(&adeToken)

	return adeToken
}

func SendFeedBackMessageToUser(adeToken models.AdeToken, user models.User, hotspotName string, hotspotLogo string, bgColor string, hotspot models.Hotspot) bool {
	var userMessage bytes.Buffer

	rp := feedbackSend{
		HotspotName:    hotspotName,
		HotspotLogo:    hotspotLogo[22:],
		BgColor:        bgColor,
		Url:            wax_utils.GenerateShortURL(configuration.Config.Survey.Url + "feedbacks/" + adeToken.Token),
		HotspotDetails: hotspot.BusinessName + " • " + hotspot.BusinessAddress + " • " + hotspot.BusinessEmail,
	}

	t := template.Must(template.ParseFiles("templates/feedback_user.tpl"))

	err := t.Execute(&userMessage, &rp)
	if err != nil {
		fmt.Println(err)
		return false
	}

	status := SendEmail("Feedback", userMessage.String(), user.Email)

	db := database.Instance()
	db.Model(&adeToken).Update("feedback_sent_time", time.Now())

	return status
}

func SendReviewMessageToUser(adeToken models.AdeToken, user models.User, hotspotName string, hotspotLogo string, bgColor string, hotspot models.Hotspot) bool {
	var userMessage bytes.Buffer

	rp := reviewSend{
		HotspotName:    hotspotName,
		HotspotLogo:    hotspotLogo[22:],
		BgColor:        bgColor,
		Url:            wax_utils.GenerateShortURL(configuration.Config.Survey.Url + "reviews/" + adeToken.Token),
		HotspotDetails: hotspot.BusinessName + " • " + hotspot.BusinessAddress + " • " + hotspot.BusinessEmail,
	}

	t, err := template.ParseFiles("templates/review_user.tpl")
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = t.Execute(&userMessage, &rp)
	if err != nil {
		fmt.Println(err)
		return false
	}

	status := SendEmail("Review", userMessage.String(), user.Email)

	db := database.Instance()
	db.Model(&adeToken).Update("review_sent_time", time.Now())

	return status
}

func SendFeedBackMessageToOwner(adeToken models.AdeToken, message string, bgColor string) bool {
	var ownerMessage bytes.Buffer

	rp := feedbackResponse{
		Message: message,
		BgColor: bgColor,
	}

	t, err := template.ParseFiles("templates/feedback_owner.tpl")
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = t.Execute(&ownerMessage, &rp)
	if err != nil {
		fmt.Println(err)
		return false
	}

	mailTo := wax_utils.GetHotspotPreferencesByKey(adeToken.HotspotId, "marketing_2_feedback_email").Value
	status := SendEmail("User Feedback", ownerMessage.String(), mailTo)

	db := database.Instance()
	adeToken.FeedbackLeftTime = time.Now()
	db.Save(&adeToken)

	return status
}

func SendReviewMessageToOwner(adeToken models.AdeToken, stars int, message string, bgColor string) bool {
	var ownerMessage bytes.Buffer
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

	err = t.Execute(&ownerMessage, &rp)
	if err != nil {
		fmt.Println(err)
		return false
	}

	mailTo := wax_utils.GetHotspotPreferencesByKey(adeToken.HotspotId, "marketing_2_feedback_email").Value
	status := SendEmail("Review: "+stars_s+"/5", ownerMessage.String(), mailTo)

	db := database.Instance()
	adeToken.ReviewLeftTime = time.Now()
	db.Save(&adeToken)

	return status
}

func SendEmail(subject string, message string, mailTo string) bool {

	status := true
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.Endpoints.Email.From)

	m.SetHeader("To", mailTo)
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
