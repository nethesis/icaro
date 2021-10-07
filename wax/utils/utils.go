/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Icaro project.
 *
 * Icaro is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Icaro is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Icaro.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	gomail "gopkg.in/gomail.v2"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

func GetHotspotPreferences(hotspotId int) []models.HotspotPreference {
	var prefs []models.HotspotPreference
	db := database.Instance()
	db.Where("hotspot_id = ?", hotspotId).Find(&prefs)

	return prefs
}

func GetHotspotPreferencesByKey(hotspotId int, key string) models.HotspotPreference {
	var pref models.HotspotPreference
	db := database.Instance()
	db.Where("`key` = ? and hotspot_id = ?", key, hotspotId).First(&pref)

	return pref
}

func GetHotspotPreferencesByKeys(hotspotId int, keys []string) []models.HotspotPreference {
	var prefs []models.HotspotPreference
	db := database.Instance()
	db.Where("`key` in (?) and hotspot_id = ?", keys, hotspotId).Find(&prefs)

	return prefs
}

func GetHotspotIntegrations(hotspotId int) []models.IntegrationWings {
	var integrations []models.IntegrationWings
	var hotspotIntegrations []models.HotspotIntegration

	db := database.Instance()
	db.Where("hotspot_id = ?", hotspotId).Find(&hotspotIntegrations)

	for _, i := range hotspotIntegrations {
		var integration models.Integration

		db := database.Instance()
		db.Where("id = ?", i.IntegrationId).Find(&integration)

		integrations = append(integrations, models.IntegrationWings{
			PreAuthRedirectUrl:  integration.PreAuthRedirectUrl,
			PostAuthRedirectUrl: integration.PostAuthRedirectUrl,
		})
	}

	return integrations
}

func CreateUserSession(userId int, sessionKey string) {
	userSession := models.UserSession{
		UserId:     userId,
		SessionKey: sessionKey,
		Created:    time.Now().UTC(),
	}

	// save user session
	db := database.Instance()
	db.Save(&userSession)
}

func CreateUserMarketing(userId int, data interface{}, accountType string) {
	var userMarketing models.UserMarketing
	serialization, _ := json.Marshal(data)

	db := database.Instance()
	db.Where("user_id = ?", userId).First(&userMarketing)

	// create or update info
	if userMarketing.Id == 0 {
		userMarketing = models.UserMarketing{
			UserId:      userId,
			AccountType: accountType,
			Data:        string(serialization),
			Created:     time.Now().UTC(),
		}
	} else {
		userMarketing.AccountType = accountType
		userMarketing.Data = string(serialization)
	}

	// save user marketing
	db.Save(&userMarketing)

}

func CheckUserSession(userId int, sessionKey string) bool {
	var userSession models.UserSession

	// check if user session exists
	db := database.Instance()
	db.Where("user_id = ? AND session_key = ?", userId, sessionKey).First(&userSession)

	if userSession.Id == 0 {
		return false
	}

	return true
}

func CheckTempUserSession(userId int, deviceMac string, sessionKey string) bool {
	var check = true
	var userTempSession models.UserTempSession

	// check if user session exists
	db := database.Instance()
	db.Where("user_id = ? AND device_mac = ? AND session_key = ?", userId, deviceMac, sessionKey).First(&userTempSession)

	// if not exists create one
	if userTempSession.Id == 0 {
		userTempSessionNew := models.UserTempSession{
			UserId:     userId,
			DeviceMac:  deviceMac,
			SessionKey: sessionKey,
			Created:    time.Now().UTC(),
		}

		// save user temp session
		db := database.Instance()
		db.Save(&userTempSessionNew)
		check = true
	} else {
		check = false
	}

	return check
}

func CheckOtherUnitLogin(mac string, unitId int, hostspotId int) bool {
	var session models.Session
	db := database.Instance()
	db.Where("device_mac = ? and unit_id != ? and hotspot_id = ? and stop_time = 0", mac, unitId, hostspotId).First(&session)

	if session.Id > 0 {
		return true
	}

	return false
}

func DeleteUserSession(userId int, sessionKey string) bool {
	var userSession models.UserSession

	// check if user session exists
	db := database.Instance()
	db.Where("user_id = ? AND session_key = ?", userId, sessionKey).First(&userSession)

	if userSession.Id == 0 {
		return false
	}

	// delete user session
	db.Delete(&userSession)
	return true
}

func GetAccountSMSByAccountId(accountId int) models.AccountSmsCount {
	var accountSMS models.AccountSmsCount
	db := database.Instance()
	db.Where("account_id = ?", accountId).First(&accountSMS)

	return accountSMS
}

func GetAccountByAccountId(accountId int) models.Account {
	var account models.Account
	db := database.Instance()
	db.Where("id = ?", accountId).First(&account)

	return account
}

func GetSessionByKeyAndUnitId(key string, unitId int) models.Session {
	var session models.Session
	db := database.Instance()
	db.Where("session_key = ? and unit_id = ?", key, unitId).First(&session)

	return session
}

func GetDeviceByMacAddressAndUserId(mac string, userId int) models.Device {
	var device models.Device
	db := database.Instance()
	db.Where("mac_address = ? and user_id = ?", mac, userId).First(&device)

	return device
}

func GetDevicesByHotspotidAndMacAddress(hotspot_id int, mac string) []models.Device {
	var devices []models.Device
	db := database.Instance()
	db.Where("hotspot_id = ? and mac_address = ?", hotspot_id, mac).Find(&devices)

	return devices
}

func GetUnitByMacAddress(mac string) models.Unit {
	var unit models.Unit
	db := database.Instance()
	db.Where("mac_address = ?", mac).First(&unit)

	return unit
}

func GetUserByNameAndHotspotId(name string, hotspotId int) models.User {
	var unit models.User
	db := database.Instance()
	db.Where("username = ? and hotspot_id = ?", name, hotspotId).First(&unit)

	return unit
}

func GetUnitByUuid(uuid string) models.Unit {
	var unit models.Unit
	db := database.Instance()
	db.Where("uuid = ?", uuid).First(&unit)

	return unit
}

func GetHotspotById(id int) models.Hotspot {
	var hotspot models.Hotspot
	db := database.Instance()
	db.Where("id = ?", id).First(&hotspot)

	return hotspot
}

func GetHotspotByUuid(id string) models.Hotspot {
	var hotspot models.Hotspot
	db := database.Instance()
	db.Where("uuid = ?", id).First(&hotspot)

	return hotspot
}

func GetUserById(id int) models.User {
	var user models.User
	db := database.Instance()
	db.Where("id = ?", id).First(&user)

	return user
}

func GetUserByUsernameAndHotspot(username string, hotspotId int) models.User {
	var user models.User
	db := database.Instance()
	db.Where("username LIKE ? AND hotspot_id = ?", "%"+username, hotspotId).First(&user)

	return user
}

func HotspotHasValidSubscription(hotspotId int) bool {
	var hotspot models.Hotspot
	var subscription models.Subscription
	db := database.Instance()
	db.Set("gorm:auto_preload", true)
	db.Preload("Account").Where("id = ?", hotspotId).First(&hotspot)
	db.Preload("SubscriptionPlan").Where("account_id = ?", hotspot.Account.Id).First(&subscription)

	return subscription.ValidFrom.Before(time.Now().UTC()) && subscription.ValidUntil.After(time.Now().UTC())
}

func GetVoucherByCode(code string, hotspotId int) models.HotspotVoucher {
	var hotspotVoucher models.HotspotVoucher
	db := database.Instance()
	db.Where("binary code = ? AND hotspot_id = ?", strings.ToLower(code), hotspotId).First(&hotspotVoucher)

	return hotspotVoucher
}

func CalcUnitDigest(unit models.Unit) string {
	h := md5.New()
	io.WriteString(h, unit.Secret+unit.Uuid)
	digest := fmt.Sprintf("%x", h.Sum(nil))

	return digest
}

func CheckUnitSecret(uuid string, secret string) bool {
	// get unit
	unit := GetUnitByUuid(uuid)

	// check secret
	if unit.Secret == secret {
		return true
	}

	return false
}

func CalcUserDigest(user models.User, challenge string) string {
	h := md5.New()
	io.WriteString(h, "00"+user.Password+challenge)
	digest := fmt.Sprintf("%x", h.Sum(nil))

	return digest
}

func GenerateCode(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func GenerateShortURL(longURL string) string {

	var shortUrl models.ShortUrl
	h := sha1.New()
	db := database.Instance()

	io.WriteString(h, longURL)
	//Calculate sha-1 hash, convert to hexadecimal representation and take only first 7 digits
	s := fmt.Sprintf("%.7s", fmt.Sprintf("%x", h.Sum(nil)))
	//Encode the first 7 digits in Base64 without padding and url safe
	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(s))

	db.Where("hash = ? ", encoded).First(&shortUrl)

	if shortUrl.Id == 0 {
		shortUrl.Hash = encoded
		shortUrl.CreatedAt = time.Now().UTC()
		shortUrl.LongUrl = longURL

		db.Save(&shortUrl)
	}

	return configuration.Config.Shortener.BaseUrl + encoded
}

func GetShortUrlByHash(hash string) models.ShortUrl {
	var shortUrl models.ShortUrl
	db := database.Instance()
	db.Where("hash = ? ", hash).First(&shortUrl)

	return shortUrl
}

func SendSMSCode(number string, code string, unit models.Unit, auth string) int {
	// get account sms count
	db := database.Instance()
	hotspot := GetHotspotById(unit.HotspotId)
	accountSMS := GetAccountSMSByAccountId(hotspot.AccountId)

	// count sms by hotspot
	var hotspotSmsCount []models.HotspotSmsCount
	db.Where("hotspot_id in (?)", hotspot.Id).Find(&hotspotSmsCount)

	hotspotCount := len(hotspotSmsCount)
	hotspotMaxCount := GetHotspotPreferencesByKey(hotspot.Id, "sms_login_max")
	hotspotMaxCountInt, err := strconv.Atoi(hotspotMaxCount.Value)
	if err != nil {
		hotspotMaxCountInt = 0
	}

	// check if exists an account for sms
	if accountSMS.Id == 0 {
		return http.StatusPaymentRequired
	}

	if accountSMS.SmsCount <= accountSMS.SmsMaxCount {

		if hotspotCount <= hotspotMaxCountInt || hotspotMaxCountInt == 0 {
			// check account and hotspot SMS thresholds
			numSMSLeftAccount := accountSMS.SmsMaxCount - accountSMS.SmsCount
			numSMSLeftHotspot := hotspotMaxCountInt - hotspotCount
			hotspotSMSThreshold, err := strconv.Atoi(GetHotspotPreferencesByKey(hotspot.Id, "sms_login_threshold").Value)

			if accountSMS.SmsThreshold > 0 && numSMSLeftAccount <= accountSMS.SmsThreshold {
				resellerAccount := GetAccountByAccountId(hotspot.AccountId)
				SendSmsAccountThresholdAlert(resellerAccount, numSMSLeftAccount)
			}

			if hotspotSMSThreshold > 0 && numSMSLeftHotspot <= hotspotSMSThreshold {
				resellerAccount := GetAccountByAccountId(hotspot.AccountId)
				SendSmsHotspotThresholdAlert(resellerAccount, hotspot, numSMSLeftHotspot)
			}

			// retrieve account info and token
			accountSid := configuration.Config.Endpoints.Sms.AccountSid
			authToken := configuration.Config.Endpoints.Sms.AuthToken
			urlAPI := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

			// compose message data
			msgData := url.Values{}
			msgData.Set("To", number)
			msgData.Set("MessagingServiceSid", configuration.Config.Endpoints.Sms.ServiceSid)
			msgData.Set("Body", `Login Link: `+GenerateShortURL(configuration.Config.Endpoints.Sms.Link+"?"+auth+"&code="+code+"&num="+url.QueryEscape(number))+`

Codice/Code: `+code+`

Logout Link: http://logout`)
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

			return resp.StatusCode
		}
	}

	if configuration.Config.Endpoints.Sms.SendQuotaAlert {
		resellerAccount := GetAccountByAccountId(hotspot.AccountId)
		SendSmsQuotaLimitAlert(resellerAccount)
	}

	return 500

}

func SaveHotspotSMSCount(hotspotSmsCount models.HotspotSmsCount) {
	// save hotspot sms count
	db := database.Instance()
	db.Save(&hotspotSmsCount)
}

func SendEmailCode(email string, code string, unit models.Unit, auth string) bool {
	hotspot := GetHotspotById(unit.HotspotId)

	status := true
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.Endpoints.Email.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Wi-Fi: "+hotspot.Description)
	m.SetBody("text/plain", `Benvenut*
Clicca sul link sottostante per accedere subito ad internet:
Login Link: `+GenerateShortURL(configuration.Config.Endpoints.Email.Link+"?"+auth+"&code="+code+"&email="+url.QueryEscape(email))+`


In alternativa puoi accedere usando la tua email assieme al seguente codice di accesso:
Codice: `+code+`


Se vuoi disconnetterti dal servizio clicca sul seguente link:
Logout link: http://logout

----------------------------------

Welcome
Click on the link below to immediately access the internet:
Login Link: `+GenerateShortURL(configuration.Config.Endpoints.Email.Link+"?"+auth+"&code="+code+"&email="+url.QueryEscape(email))+`


Otherwise you can log in using your email together with this access code:
Code: `+code+`


If you want to disconnect from the service click on this link:
Logout link: http://logout`)

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

func Contains(intSlice []int, searchInt int) bool {
	for _, value := range intSlice {
		if value == searchInt {
			return true
		}
	}
	return false
}

func ContainsS(stringSlice []string, searchString string) bool {
	for _, value := range stringSlice {
		if value == searchString {
			return true
		}
	}
	return false
}

func GetUsersByMacAddressAndunitMacAddress(mac string, unitMacAddress string) (bool, []models.User) {

	var users []models.User

	unit := GetUnitByMacAddress(unitMacAddress)
	if unit.Id <= 0 {
		return false, users
	}

	devices := GetDevicesByHotspotidAndMacAddress(unit.HotspotId, mac)
	if len(devices) < 0 {
		return false, users
	}

	for _, device := range devices {

		user := GetUserById(device.UserId)

		if user.Id > 0 {
			users = append(users, user)
		}
	}

	if len(users) <= 0 {
		return false, users
	}

	return true, users
}

func GetTodaySessionTrafficByUser(user models.User) int64 {
	// calculate today midnight
	now := time.Now().UTC()
	midnightToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var sessions []models.Session
	var sessionHistories []models.SessionHistory

	db := database.Instance()
	db.Where("update_time >= ? AND user_id = ?", midnightToday, user.Id).Find(&sessions)
	db.Where("update_time >= ? AND user_id = ?", midnightToday, user.Id).Find(&sessionHistories)

	var todayTraffic int64
	todayTraffic = 0

	for _, session := range sessions {
		todayTraffic += session.BytesDown + session.BytesUp
	}
	for _, history := range sessionHistories {
		todayTraffic += history.BytesDown + history.BytesUp
	}

	return todayTraffic
}

func GetTodaySessionTimeByUser(user models.User) int {
	// calculate today midnight
	now := time.Now().UTC()
	midnightToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var sessions []models.Session
	var sessionHistories []models.SessionHistory

	db := database.Instance()
	db.Where("update_time >= ? AND user_id = ?", midnightToday, user.Id).Find(&sessions)
	db.Where("update_time >= ? AND user_id = ?", midnightToday, user.Id).Find(&sessionHistories)

	var todayTime = 0

	for _, session := range sessions {
		todayTime += session.Duration
	}
	for _, history := range sessionHistories {
		todayTime += history.Duration
	}

	return todayTime
}

func CalculateRemainTraffic(user models.User) int64 {
	// get today total navigation traffic for user
	totalTraffic := GetTodaySessionTrafficByUser(user)
	remainTraffic := user.MaxNavigationTraffic - totalTraffic

	return remainTraffic
}

func CalculateRemainTime(user models.User, timezone string) int {
	// get today total navigation time for user
	totalTime := GetTodaySessionTimeByUser(user)
	remainTime := user.MaxNavigationTime - totalTime

	loc, _ := time.LoadLocation(timezone)
	now := time.Now().In(loc)
	midnightNow := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	secondsToMidnight := int(midnightNow.Sub(now).Seconds()) // Session-Timeout

	if user.MaxNavigationTime == 0 && user.MaxNavigationTraffic == 0 {
		return int(user.ValidUntil.Sub(time.Now().UTC()).Seconds())
	} else if user.MaxNavigationTime != 0 {
		if remainTime < secondsToMidnight {
			return remainTime
		}
	}

	return secondsToMidnight
}

func FindAutoLoginUser(users []models.User) models.User {

	var users_enabled []models.User
	var user models.User

	//find users with autologin enabled
	for _, _user := range users {
		if _user.AutoLogin {
			users_enabled = append(users_enabled, _user)
		}
	}

	// check if there is any autologin enabled users
	if len(users_enabled) <= 0 {
		return user
	}

	//find the most recent created user
	for _, user_enabled := range users_enabled {
		if user_enabled.Created.After(user.Created) {
			user = user_enabled
		}
	}

	return user
}

func SendSmsAccountThresholdAlert(reseller models.Account, remaining int) bool {
	subject := "Hotspot Alert: SMS threshold reached"
	body := fmt.Sprintf("You have left %d SMS in your account.\nPlease buy an additional SMS quota soon, or disable sms login/feedback from your hotspots.\n", remaining)
	return SendSmsAlert(reseller, subject, body)
}

func SendSmsHotspotThresholdAlert(reseller models.Account, hotspot models.Hotspot, remaining int) bool {
	subject := "Hotspot Alert: SMS threshold reached"
	body := fmt.Sprintf("You have left %d SMS in your hotspot %s.\nPlease buy an additional SMS quota soon, or disable sms login/feedback from your hotspots.\n", remaining, hotspot.Name)
	return SendSmsAlert(reseller, subject, body)
}

func SendSmsQuotaLimitAlert(reseller models.Account) bool {
	subject := "Hotspot Alert: SMS quota limit exceeded"
	body := "You do not have any more SMS to send in your account,\n" +
		"please buy an additional SMS quota or disable sms login/feedback from your hotspots.\n"
	return SendSmsAlert(reseller, subject, body)
}

func SendSmsAlert(reseller models.Account, subject string, body string) bool {
	status := true

	if reseller.Type == "reseller" {
		m := gomail.NewMessage()
		m.SetHeader("From", configuration.Config.Endpoints.Email.From)
		m.SetHeader("To", reseller.Email)
		m.SetHeader("Subject", subject)
		m.SetBody("text/plain", body)

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
	} else {
		status = false
	}

	return status
}

func CheckUserPassword(user models.User, password string) bool {
	// search user
	var userSearch models.User
	db := database.Instance()
	db.Where("id = ?", user.Id).First(&userSearch)

	// verify password
	if password == userSearch.Password {
		return true
	}

	return false
}

func CreateUserAuth(sessionId string, sessionTimeout int, unitUuid string, userId int, username string, password string, typeAuth string) {
	// create db instance
	db := database.Instance()

	// switch typeAuth
	switch typeAuth {
	case "created", "updated":
		// create record
		daemonAuth := models.DaemonAuth{
			SessionId:      sessionId,
			SessionTimeout: sessionTimeout,
			UnitUuid:       unitUuid,
			UserId:         userId,
			Username:       username,
			Password:       password,
			Type:           typeAuth,
			Updated:        time.Now().UTC(),
		}

		// save record
		db.Save(&daemonAuth)

	case "login", "logout":
		// search record
		var daemonAuth models.DaemonAuth
		db.Where("session_id = ? AND unit_uuid = ? AND username = ?", sessionId, unitUuid, username).First(&daemonAuth)

		if daemonAuth.Id == 0 {
			// create new record
			daemonAuth = models.DaemonAuth{
				SessionId:      sessionId,
				SessionTimeout: sessionTimeout,
				UnitUuid:       unitUuid,
				UserId:         userId,
				Username:       username,
				Password:       password,
				Type:           typeAuth,
				Updated:        time.Now().UTC(),
			}
		} else {
			// update record
			daemonAuth.Password = password
			daemonAuth.Type = typeAuth
			daemonAuth.Updated = time.Now().UTC()
		}

		// save record
		db.Save(&daemonAuth)

	case "temporary":
		// search record
		var daemonAuth models.DaemonAuth
		db.Where("session_id = ? AND unit_uuid = ? AND username = ?", sessionId, unitUuid, username).First(&daemonAuth)

		if daemonAuth.Id == 0 {
			// create new record
			daemonAuth = models.DaemonAuth{
				SessionId:      sessionId,
				SessionTimeout: sessionTimeout,
				UnitUuid:       unitUuid,
				UserId:         userId,
				Username:       username,
				Password:       password,
				Type:           typeAuth,
				Updated:        time.Now().UTC(),
			}
		} else {
			// update record
			daemonAuth.SessionTimeout = sessionTimeout
			daemonAuth.Type = typeAuth
			daemonAuth.Updated = time.Now().UTC()
		}

		// save record
		db.Save(&daemonAuth)
	}
}
