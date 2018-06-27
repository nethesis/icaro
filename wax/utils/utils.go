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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func GetSessionByKeyAndUnitId(key string, unitId int) models.Session {
	var session models.Session
	db := database.Instance()
	db.Where("session_key = ? and unit_id = ?", key, unitId).First(&session)

	return session
}

func GetDeviceByMacAddress(mac string) models.Device {
	var unit models.Device
	db := database.Instance()
	db.Where("mac_address = ?", mac).First(&unit)

	return unit
}

func GetDeviceByHotspotidAndMacAddress(hotspot_id int, mac string) models.Device {
	var device models.Device
	db := database.Instance()
	db.Where("hotspot_id = ? and mac_address = ?", hotspot_id, mac).First(&device)

	return device
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

func GetUserById(id int) models.User {
	var user models.User
	db := database.Instance()
	db.Where("id = ?", id).First(&user)

	return user
}

func GetUserByUsername(username string) models.User {
	var user models.User
	db := database.Instance()
	db.Where("username = ?", username).First(&user)

	return user
}

func GetUserByUsernameAndHotspot(username string, hotspotId int) models.User {
	var user models.User
	db := database.Instance()
	db.Where("username = ? AND hotspot_id = ?", username, hotspotId).First(&user)

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
	db.Where("binary code = ? AND hotspot_id = ?", code, hotspotId).First(&hotspotVoucher)

	return hotspotVoucher
}

func CalcUnitDigest(unit models.Unit) string {
	h := md5.New()
	io.WriteString(h, unit.Secret+unit.Uuid)
	digest := fmt.Sprintf("%x", h.Sum(nil))

	return digest
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

func SendSMSCode(number string, code string, unit models.Unit, auth string) int {
	// get account sms count
	db := database.Instance()
	hotspot := GetHotspotById(unit.HotspotId)
	accountSMS := GetAccountSMSByAccountId(hotspot.AccountId)

	// check if exists an account for sms
	if accountSMS.Id == 0 {
		return http.StatusPaymentRequired
	}

	if accountSMS.SmsCount <= accountSMS.SmsMaxCount {
		// retrieve account info and token
		accountSid := configuration.Config.Endpoints.Sms.AccountSid
		authToken := configuration.Config.Endpoints.Sms.AuthToken
		urlAPI := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

		// compose message data
		msgData := url.Values{}
		msgData.Set("To", number)
		msgData.Set("MessagingServiceSid", configuration.Config.Endpoints.Sms.ServiceSid)
		msgData.Set("Body", "SMS login code: "+code+
			"\n\nLogin Link: "+configuration.Config.Endpoints.Sms.Link+
			"?"+auth+"&code="+code+"&num="+url.QueryEscape(number)+
			"\n\nLogout Link: http://logout")
		msgDataReader := *strings.NewReader(msgData.Encode())

		// create HTTP request client
		client := &http.Client{}
		req, _ := http.NewRequest("POST", urlAPI, &msgDataReader)
		req.SetBasicAuth(accountSid, authToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// make HTTP POST request
		resp, _ := client.Do(req)

		// update sms accounting table
		if resp.StatusCode == 201 {
			accountSMS.SmsCount = accountSMS.SmsCount + 1
			db.Save(&accountSMS)
		}

		return resp.StatusCode

	} else {
		return 500
	}

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
	m.SetBody("text/plain", "Email login code: "+code+
		"\n\nLogin Link: "+configuration.Config.Endpoints.Email.Link+
		"?"+auth+"&code="+code+"&email="+url.QueryEscape(email)+
		"\n\nLogout Link: http://logout")

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

func GetUserByMacAddressAndunitMacAddress(mac string, unitMacAddress string) (bool, models.User) {

	var user models.User

	unit := GetUnitByMacAddress(unitMacAddress)
	if unit.Id <= 0 {
		return false, user
	}

	device := GetDeviceByHotspotidAndMacAddress(unit.HotspotId, mac)
	if device.Id < 0 {
		return false, user
	}

	user = GetUserById(device.UserId)
	if user.Id < 0 {
		return false, user
	}

	return true, user
}

func GetTodaySessionTrafficByUser(user models.User) int {
	// calculate today midnight
	now := time.Now().UTC()
	midnightToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var sessions []models.Session
	var sessionHistories []models.SessionHistory

	db := database.Instance()
	db.Where("update_time >= ? AND user_id = ?", midnightToday, user.Id).Find(&sessions)
	db.Where("update_time >= ? AND user_id = ?", midnightToday, user.Id).Find(&sessionHistories)

	var todayTraffic = 0

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

func CalculateRemainTraffic(user models.User) int {
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
	midnightNow := time.Date(now.Year(), now.Month(), now.Add(time.Hour*24).Day(), 0, 0, 0, 0, now.Location())
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
