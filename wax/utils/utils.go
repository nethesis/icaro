/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Icaro project.
 *
 * Icaro is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Icaro is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Icaro.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package utils

import (
	"crypto/md5"
	"crypto/rand"
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
	db := database.Database()
	db.Where("hotspot_id = ?", hotspotId).Find(&prefs)
	db.Close()

	return prefs
}

func GetHotspotPreferencesByKey(hotspotId int, key string) models.HotspotPreference {
	var pref models.HotspotPreference
	db := database.Database()
	db.Where("`key` = ? and hotspot_id = ?", key, hotspotId).First(&pref)
	db.Close()

	return pref
}

func GetHotspotPreferencesByKeys(hotspotId int, keys []string) []models.HotspotPreference {
	var prefs []models.HotspotPreference
	db := database.Database()
	db.Where("`key` in (?) and hotspot_id = ?", keys, hotspotId).Find(&prefs)
	db.Close()

	return prefs
}

func CreateUserSession(userId int, sessionKey string) {
	userSession := models.UserSession{
		UserId:     userId,
		SessionKey: sessionKey,
		Created:    time.Now().UTC(),
	}

	// save user session
	db := database.Database()
	db.Save(&userSession)
	db.Close()
}

func CheckUserSession(userId int, sessionKey string) bool {
	var userSession models.UserSession

	// check if user session exists
	db := database.Database()
	db.Where("user_id = ? AND session_key = ?", userId, sessionKey).First(&userSession)
	db.Close()

	if userSession.Id == 0 {
		return false
	}

	return true
}

func CheckTempUserSession(userId int, deviceMac string, sessionKey string) bool {
	var check = true
	var userTempSession models.UserTempSession

	// check if user session exists
	db := database.Database()
	db.Where("user_id = ? AND device_mac = ? AND session_key = ?", userId, deviceMac, sessionKey).First(&userTempSession)
	db.Close()

	// if not exists create one
	if userTempSession.Id == 0 {
		userTempSessionNew := models.UserTempSession{
			UserId:     userId,
			DeviceMac:  deviceMac,
			SessionKey: sessionKey,
			Created:    time.Now().UTC(),
		}

		// save user temp session
		db := database.Database()
		db.Save(&userTempSessionNew)
		db.Close()
		check = true
	} else {
		check = false
	}

	return check
}

func DeleteUserSession(userId int, sessionKey string) bool {
	var userSession models.UserSession

	// check if user session exists
	db := database.Database()
	db.Where("user_id = ? AND session_key = ?", userId, sessionKey).First(&userSession)

	if userSession.Id == 0 {
		db.Close()
		return false
	}

	// delete user session
	db.Delete(&userSession)
	db.Close()
	return true
}

func GetSessionByKeyAndUnitId(key string, unitId int) models.Session {
	var session models.Session
	db := database.Database()
	db.Where("session_key = ? and unit_id = ?", key, unitId).First(&session)
	db.Close()

	return session
}

func GetDeviceByMacAddress(mac string) models.Device {
	var unit models.Device
	db := database.Database()
	db.Where("mac_address = ?", mac).First(&unit)
	db.Close()

	return unit
}

func GetUnitByMacAddress(mac string) models.Unit {
	var unit models.Unit
	db := database.Database()
	db.Where("mac_address = ?", mac).First(&unit)
	db.Close()

	return unit
}

func GetUserByNameAndHotspotId(name string, hotspotId int) models.User {
	var unit models.User
	db := database.Database()
	db.Where("username = ? and hotspot_id = ?", name, hotspotId).First(&unit)
	db.Close()

	return unit
}

func GetUnitByUuid(uuid string) models.Unit {
	var unit models.Unit
	db := database.Database()
	db.Where("uuid = ?", uuid).First(&unit)
	db.Close()

	return unit
}

func GetHotspotById(id int) models.Hotspot {
	var hotspot models.Hotspot
	db := database.Database()
	db.Where("id = ?", id).First(&hotspot)
	db.Close()

	return hotspot
}

func GetUserByUsername(username string) models.User {
	var user models.User
	db := database.Database()
	db.Where("username = ?", username).First(&user)
	db.Close()

	return user
}

func GetVoucherByCode(code string, hotspotId int) models.HotspotVoucher {
	var hotspotVoucher models.HotspotVoucher
	db := database.Database()
	db.Where("binary code = ? AND hotspot_id = ?", code, hotspotId).First(&hotspotVoucher)
	db.Close()

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

func SendSMSCode(number string, code string) int {
	// retrieve account info and token
	accountSid := configuration.Config.Endpoints.Sms.AccountSid
	authToken := configuration.Config.Endpoints.Sms.AuthToken
	urlAPI := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	// compose message data
	msgData := url.Values{}
	msgData.Set("To", number)
	msgData.Set("From", configuration.Config.Endpoints.Sms.Number)
	msgData.Set("Body", "Icaro - SMS Login code: "+code) // TODO: get message from hotspot preferences
	msgDataReader := *strings.NewReader(msgData.Encode())

	// create HTTP request client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlAPI, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make HTTP POST request
	resp, _ := client.Do(req)
	return resp.StatusCode
}

func SendEmailCode(email string, code string) bool {
	status := true
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.Endpoints.Email.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Icaro")
	m.SetBody("text/html", "Email Login code: "+code)

	d := gomail.NewDialer(
		configuration.Config.Endpoints.Email.SMTPHost,
		configuration.Config.Endpoints.Email.SMTPPort,
		configuration.Config.Endpoints.Email.SMTPUser,
		configuration.Config.Endpoints.Email.SMTPPassword,
	)

	// send the email
	if err := d.DialAndSend(m); err != nil {
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
