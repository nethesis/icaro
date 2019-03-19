package utils

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	gomail "gopkg.in/gomail.v2"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	wax_utils "github.com/nethesis/icaro/wax/utils"
)

type surveySend struct {
	HotspotName           string
	HotspotLogo           string
	HotspotDetails        string
	BgColor               string
	HotspotSurveyBodyText string
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

func compileUserEmailTemplate(Type string, adeToken models.AdeToken, hotspotName string, hotspotLogo string, bgColor string, hotspot models.Hotspot, BodyText string) (string, bool) {
	var userMessage bytes.Buffer
	var templateFile string
	var Url string

	if adeToken.Id != 0 {
		Url = wax_utils.GenerateShortURL(configuration.Config.Survey.Url + Type + "s/" + adeToken.Token)
	} else {
		Url = "https://example.org/"
	}

	templateFile = "templates/" + Type + "_user.tpl"

	t, _ := template.ParseFiles(templateFile)

	rp := surveySend{
		HotspotName:           hotspotName,
		HotspotLogo:           hotspotLogo[22:],
		BgColor:               bgColor,
		HotspotDetails:        hotspot.BusinessName + " • " + hotspot.BusinessAddress + " • " + hotspot.BusinessEmail,
		HotspotSurveyBodyText: strings.Replace(BodyText, "$$URL$$", Url, -1),
	}

	err := t.Execute(&userMessage, &rp)
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	return userMessage.String(), true
}

func SendFeedBackMessageToUser(adeToken models.AdeToken, userEmail string, hotspotName string, hotspotLogo string, bgColor string, hotspot models.Hotspot, BodyText string) bool {

	userMessage, status := compileUserEmailTemplate("feedback", adeToken, hotspotName, hotspotLogo, bgColor, hotspot, BodyText)

	if status {
		status = SendEmail("Feedback", userMessage, userEmail)

		if adeToken.Id != 0 {
			db := database.Instance()
			db.Model(&adeToken).Update("feedback_sent_time", time.Now())
		}
	}

	return status
}

func SendReviewMessageToUser(adeToken models.AdeToken, userEmail string, hotspotName string, hotspotLogo string, bgColor string, hotspot models.Hotspot, BodyText string) bool {

	userMessage, status := compileUserEmailTemplate("review", adeToken, hotspotName, hotspotLogo, bgColor, hotspot, BodyText)

	if status {
		status = SendEmail("Review", userMessage, userEmail)

		if adeToken.Id != 0 {
			db := database.Instance()
			db.Model(&adeToken).Update("review_sent_time", time.Now())
		}
	}

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

func SendSMS(adeToken models.AdeToken, message string, survey string, smsTo string, hotspotId int) bool {
	// get account sms count
	db := database.Instance()
	hotspot := wax_utils.GetHotspotById(hotspotId)
	accountSMS := wax_utils.GetAccountSMSByAccountId(hotspot.AccountId)

	// check if exists an account for sms
	if accountSMS.Id == 0 {
		return false
	}

	if accountSMS.SmsCount <= accountSMS.SmsMaxCount {
		// retrieve account info and token
		accountSid := configuration.Config.Endpoints.Sms.AccountSid
		authToken := configuration.Config.Endpoints.Sms.AuthToken
		urlAPI := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

		// compose message data
		msgData := url.Values{}
		msgData.Set("To", smsTo)
		msgData.Set("MessagingServiceSid", configuration.Config.Endpoints.Sms.ServiceSid)
		msgData.Set("Body",
			hotspot.BusinessName+"\n"+
				message+"\n"+"Link: "+
				wax_utils.GenerateShortURL(configuration.Config.Survey.Url+survey+"/"+adeToken.Token))
		msgDataReader := *strings.NewReader(msgData.Encode())

		// create HTTP request client
		client := &http.Client{
			Timeout: time.Second * 30,
		}
		req, _ := http.NewRequest("POST", urlAPI, &msgDataReader)
		req.SetBasicAuth(accountSid, authToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// make HTTP POST request
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err.Error())
		}

		defer resp.Body.Close()

		// update sms accounting table
		if resp.StatusCode == 201 {
			accountSMS.SmsCount = accountSMS.SmsCount + 1
			db.Save(&accountSMS)
		}

		return true

	} else {

		if configuration.Config.Endpoints.Sms.SendQuotaAlert {
			resellerAccount := wax_utils.GetAccountByAccountId(hotspot.AccountId)
			wax_utils.SendSmsQuotaLimitAlert(resellerAccount)
		}

		return false
	}

}
