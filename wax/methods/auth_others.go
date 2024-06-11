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

package methods

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nethesis/icaro/wax/utils"

	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/methods"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

type smsMarketingData struct {
	Number string `json:"number"`
}
type emailMarketingData struct {
	Email string `json:"email"`
}

func SMSAuth(c *gin.Context) {
	number := c.Param("number")
	digest := c.Query("digest")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")
	reset := c.Query("reset")
	uamip := c.Query("uamip")
	uamport := c.Query("uamport")
	voucherCode := c.Query("voucher_code")

	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "number is required"})
		return
	}

	// check if user exists
	// get unit
	unit := utils.GetUnitByUuid(uuid)
	user := utils.GetUserByUsernameAndHotspot(number, unit.HotspotId)
	if user.Id == 0 {
		// create user
		days := utils.GetHotspotPreferencesByKey(unit.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)

		down := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Down")
		downInt, _ := strconv.Atoi(down.Value)
		up := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Up")
		upInt, _ := strconv.Atoi(up.Value)

		maxTraffic := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Max-Total-Octets")
		maxTrafficInt, _ := strconv.ParseInt(maxTraffic.Value, 10, 64)

		maxTime := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Max-Navigation-Time")
		maxTimeInt, _ := strconv.Atoi(maxTime.Value)

		autoLogin := utils.GetHotspotPreferencesByKey(unit.HotspotId, "auto_login")
		autoLoginBool, _ := strconv.ParseBool(autoLogin.Value)

		// retrieve voucher
		if len(voucherCode) > 0 {
			voucher := utils.GetVoucherByCode(voucherCode, unit.HotspotId)

			if voucher.Id > 0 {
				daysInt = int(voucher.Expires.Sub(time.Now().UTC()).Hours() / 24)
				downInt = voucher.BandwidthDown
				upInt = voucher.BandwidthUp
				autoLoginBool = voucher.AutoLogin
				maxTrafficInt = voucher.MaxTraffic
				maxTimeInt = voucher.MaxTime
			}
		}

		// generate code
		code := utils.GenerateCode(6)

		newUser := models.User{
			HotspotId:            unit.HotspotId,
			Name:                 number,
			Username:             number,
			Password:             code,
			Email:                "",
			AccountType:          "sms",
			Reason:               "",
			Country:              "",
			MarketingAuth:        true,
			SurveyAuth:           true,
			KbpsDown:             downInt,
			KbpsUp:               upInt,
			MaxNavigationTraffic: maxTrafficInt,
			MaxNavigationTime:    maxTimeInt,
			AutoLogin:            autoLoginBool,
			ValidFrom:            time.Now().UTC(),
			ValidUntil:           time.Now().UTC().AddDate(0, 0, daysInt+1),
		}
		newUser.Id = methods.CreateUser(newUser)

		// send sms with code
		userIdStr := strconv.Itoa(newUser.Id)
		status := utils.SendSMSCode(number, code, unit, "digest="+digest+"&uuid="+uuid+"&uamip="+uamip+"&uamport="+uamport+"&user="+userIdStr, uamip, uamport)

		// check response
		if status != 201 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not send"})
			return
		}

		// add sms statistics
		hotspotSmsCount := models.HotspotSmsCount{
			HotspotId: unit.HotspotId,
			UnitId:    unit.Id,
			Number:    number,
			Reset:     false,
			Sent:      time.Now().UTC(),
		}
		utils.SaveHotspotSMSCount(hotspotSmsCount)

		// create marketing info with user infos
		utils.CreateUserMarketing(newUser.Id, smsMarketingData{Number: number}, "sms")

		// create user auth
		utils.CreateUserAuth(sessionId, 0, uuid, newUser.Id, newUser.Username, newUser.Password, "created")

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": number, "user_db_id": newUser.Id})
	} else {
		// update user info
		days := utils.GetHotspotPreferencesByKey(user.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, daysInt)

		// create user session check
		utils.CreateUserSession(user.Id, sessionId)

		// create marketing info with user infos
		utils.CreateUserMarketing(user.Id, smsMarketingData{Number: number}, "sms")

		// retrieve voucher
		if len(voucherCode) > 0 {
			voucher := utils.GetVoucherByCode(voucherCode, user.HotspotId)

			if voucher.Id > 0 {
				duration := voucher.Duration

				if duration == 0 {
					duration = int(voucher.Expires.Sub(time.Now().UTC()).Hours()/24) + 1
				}
				user.ValidUntil = time.Now().UTC().AddDate(0, 0, duration)
				user.KbpsDown = voucher.BandwidthDown
				user.KbpsUp = voucher.BandwidthUp
				user.AutoLogin = voucher.AutoLogin
			}
		}

		// check if is reset
		if reset == "true" {
			// get unit
			unit := utils.GetUnitByUuid(uuid)

			// generate code
			code := utils.GenerateCode(6)

			// send sms with code
			userIdStr := strconv.Itoa(user.Id)
			status := utils.SendSMSCode(number, code, unit, "digest="+digest+"&uuid="+uuid+"&uamip="+uamip+"&uamport="+uamport+"&user="+userIdStr, uamip, uamport)

			// check response
			if status != 201 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not send"})
				return
			}

			// add sms statistics
			hotspotSmsCount := models.HotspotSmsCount{
				HotspotId: unit.HotspotId,
				UnitId:    unit.Id,
				Number:    number,
				Reset:     true,
				Sent:      time.Now().UTC(),
			}
			utils.SaveHotspotSMSCount(hotspotSmsCount)

			// update code
			user.Password = code
		}

		db := database.Instance()
		db.Save(&user)

		// create user auth
		utils.CreateUserAuth(sessionId, 0, uuid, user.Id, user.Username, user.Password, "updated")

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": number, "exists": true, "reset": reset, "user_db_id": user.Id})
	}
}

func EmailAuth(c *gin.Context) {
	email := c.Param("email")
	digest := c.Query("digest")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")
	reset := c.Query("reset")
	uamip := c.Query("uamip")
	uamport := c.Query("uamport")
	voucherCode := c.Query("voucher_code")
	mac := c.Query("mac")

	// get unit
	unit := utils.GetUnitByUuid(uuid)

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is required"})
		return
	}

	// check if email is in a valid domain
	checkEmailDomain := utils.GetHotspotPreferencesByKey(unit.HotspotId, "check_email_domain")
	checkEmailDomainBool, _ := strconv.ParseBool(checkEmailDomain.Value)

	if checkEmailDomainBool {
		// get allowed domains
		allowedDomainsStr := utils.GetHotspotPreferencesByKey(unit.HotspotId, "check_email_domain_list")
		allowedDomains := strings.Split(allowedDomainsStr.Value, ",")

		// extract domain from email
		parts := strings.Split(email, "@")

		// check if email is valid
		if len(parts) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "email is invalid"})
			return
		}

		// get domain
		domain := parts[1]

		if !utils.ContainsS(allowedDomains, domain) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "email domain is invalid"})
			return
		}
	}

	// check if there is skip_auth option
	skipVerification := utils.GetHotspotPreferencesByKey(unit.HotspotId, "email_login_skip_auth")

	// set default username to search
	usernameToCheck := email

	// compose email with or without mac_address
	if skipVerification.Value == "true" {
		usernameToCheck += ":" + mac
	}

	// check if user exists
	user := utils.GetUserByUsernameAndHotspot(usernameToCheck, unit.HotspotId)
	if user.Id == 0 {
		// create user
		days := utils.GetHotspotPreferencesByKey(unit.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)

		down := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Down")
		downInt, _ := strconv.Atoi(down.Value)
		up := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Up")
		upInt, _ := strconv.Atoi(up.Value)

		maxTraffic := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Max-Total-Octets")
		maxTrafficInt, _ := strconv.ParseInt(maxTraffic.Value, 10, 64)

		maxTime := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Max-Navigation-Time")
		maxTimeInt, _ := strconv.Atoi(maxTime.Value)

		autoLogin := utils.GetHotspotPreferencesByKey(unit.HotspotId, "auto_login")
		autoLoginBool, _ := strconv.ParseBool(autoLogin.Value)

		// retrieve voucher
		if len(voucherCode) > 0 {
			voucher := utils.GetVoucherByCode(voucherCode, unit.HotspotId)

			if voucher.Id > 0 {
				daysInt = int(voucher.Expires.Sub(time.Now().UTC()).Hours() / 24)
				downInt = voucher.BandwidthDown
				upInt = voucher.BandwidthUp
				autoLoginBool = voucher.AutoLogin
				maxTrafficInt = voucher.MaxTraffic
				maxTimeInt = voucher.MaxTime
			}
		}

		// generate code
		code := utils.GenerateCode(6)

		// define username
		newUsername := email

		// add mac address to username if skip_auth is enabled
		if skipVerification.Value == "true" {
			newUsername += ":" + mac
		}

		newUser := models.User{
			HotspotId:            unit.HotspotId,
			Name:                 newUsername,
			Username:             newUsername,
			Password:             code,
			Email:                email,
			AccountType:          "email",
			Reason:               "",
			Country:              "",
			MarketingAuth:        true,
			SurveyAuth:           true,
			KbpsDown:             downInt,
			KbpsUp:               upInt,
			MaxNavigationTraffic: maxTrafficInt,
			MaxNavigationTime:    maxTimeInt,
			AutoLogin:            autoLoginBool,
			ValidFrom:            time.Now().UTC(),
			ValidUntil:           time.Now().UTC().AddDate(0, 0, daysInt+1),
		}
		newUser.Id = methods.CreateUser(newUser)

		// send email with code, only if skip auth is false
		if skipVerification.Value == "false" {
			userIdStr := strconv.Itoa(newUser.Id)
			status := utils.SendEmailCode(email, code, unit, "digest="+digest+"&uuid="+uuid+"&uamip="+uamip+"&uamport="+uamport+"&user="+userIdStr, uamip, uamport)

			// check response
			if !status {
				c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not sent"})
				return
			}
		}

		// create marketing info with user infos
		utils.CreateUserMarketing(newUser.Id, emailMarketingData{Email: email}, "email")

		// create user auth
		utils.CreateUserAuth(sessionId, 0, uuid, newUser.Id, newUser.Username, newUser.Password, "created")

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": email, "user_db_id": newUser.Id})
	} else {
		// update user info
		days := utils.GetHotspotPreferencesByKey(user.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, daysInt)

		// create user session check
		utils.CreateUserSession(user.Id, sessionId)

		// create marketing info with user infos
		utils.CreateUserMarketing(user.Id, emailMarketingData{Email: email}, "email")

		// retrieve voucher
		if len(voucherCode) > 0 {
			voucher := utils.GetVoucherByCode(voucherCode, user.HotspotId)

			if voucher.Id > 0 {
				duration := voucher.Duration

				if duration == 0 {
					duration = int(voucher.Expires.Sub(time.Now().UTC()).Hours()/24) + 1
				}
				user.ValidUntil = time.Now().UTC().AddDate(0, 0, duration)
				user.KbpsDown = voucher.BandwidthDown
				user.KbpsUp = voucher.BandwidthUp
				user.AutoLogin = voucher.AutoLogin
			}
		}

		// check if is reset
		if reset == "true" {
			// get unit
			unit := utils.GetUnitByUuid(uuid)

			// generate code
			code := utils.GenerateCode(6)

			// send email with code
			userIdStr := strconv.Itoa(user.Id)
			status := utils.SendEmailCode(email, code, unit, "digest="+digest+"&uuid="+uuid+"&uamip="+uamip+"&uamport="+uamport+"&user="+userIdStr, uamip, uamport)

			// check response
			if !status {
				c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not send"})
				return
			}

			// update code
			user.Password = code
		}

		db := database.Instance()
		db.Save(&user)

		// create user auth
		utils.CreateUserAuth(sessionId, 0, uuid, user.Id, user.Username, user.Password, "updated")

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": email, "exists": true, "reset": reset, "user_db_id": user.Id})
	}
}

func MACAuth(c *gin.Context) {
	mac := c.Param("mac")
	name := c.Query("name")
	kbps_down := c.Query("kbps_down")
	kbps_up := c.Query("kbps_up")
	uuid := c.Query("uuid")

	if mac == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "mac is required"})
		return
	}

	// check if user exists
	// get unit
	unit := utils.GetUnitByUuid(uuid)
	user := utils.GetUserByUsernameAndHotspot(mac, unit.HotspotId)
	if user.Id == 0 {
		// create user
		downInt, _ := strconv.Atoi(kbps_down)
		upInt, _ := strconv.Atoi(kbps_up)

		newUser := models.User{
			HotspotId:            unit.HotspotId,
			Name:                 name,
			Username:             mac,
			Password:             "",
			Email:                "",
			AccountType:          "mac",
			KbpsDown:             downInt,
			KbpsUp:               upInt,
			MaxNavigationTraffic: 0,
			MaxNavigationTime:    0,
			AutoLogin:            true,
			ValidFrom:            time.Now().UTC(),
			ValidUntil:           time.Now().UTC().AddDate(0, 0, 3650), // ten years
		}
		newUser.Id = methods.CreateUser(newUser)

		// create device
		newDevice := models.Device{
			HotspotId:  unit.HotspotId,
			UserId:     newUser.Id,
			MacAddress: mac,
			Created:    time.Now().UTC(),
		}

		db := database.Instance()
		db.Save(&newDevice)

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": mac, "device_id": newDevice.Id})
	} else {
		// response to client
		c.JSON(http.StatusConflict, gin.H{"message": "mac already used for this hotspot!"})
	}
}

func VoucherAuth(c *gin.Context) {
	code := c.Param("code")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")

	// extract unit
	unit := utils.GetUnitByUuid(uuid)
	user := utils.GetUserByUsernameAndHotspot(code, unit.HotspotId)

	// extract voucher
	voucher := utils.GetVoucherByCode(code, unit.HotspotId)

	// check voucher validity
	if voucher.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Voucher is invalid"})
	} else {
		// check if it's expired, considering end of day
		voucherExpired := false

		if !voucher.Expires.IsZero() {
			// consider end of day
			y, m, d := voucher.Expires.Date()
			expirationDate := time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)

			if expirationDate.Before(time.Now().UTC()) {
				voucherExpired = true
			}
		}

		if voucherExpired {
			// delete voucher
			db := database.Instance()
			db.Delete(&voucher)

			c.JSON(http.StatusUnauthorized, gin.H{"message": "Voucher is expired"})
		} else {
			// check max use
			if voucher.RemainUse != -1 && voucher.RemainUse == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Voucher limit reached"})
			} else {
				// update epiration date
				if voucher.Expires.IsZero() {
					voucher.Expires = time.Now().UTC().AddDate(0, 0, voucher.Duration)
				}

				// check if is not limitless
				if voucher.RemainUse != -1 {
					voucher.RemainUse--
				}

				db := database.Instance()
				db.Save(&voucher)

				// update user valid info
				if user.ValidFrom.IsZero() && user.ValidUntil.IsZero() {
					user.ValidFrom = time.Now().UTC()
					duration := voucher.Duration

					if duration == 0 {
						duration = int(voucher.Expires.Sub(time.Now().UTC()).Hours()/24) + 1
					}

					user.ValidUntil = time.Now().UTC().AddDate(0, 0, duration)
					db.Save(&user)
				}

				// create user auth
				utils.CreateUserAuth(sessionId, 0, uuid, user.Id, user.Username, user.Password, "created")

				c.JSON(http.StatusOK, gin.H{"message": "Voucher is valid", "code": voucher.Code, "type": voucher.Type, "user_db_id": user.Id})
			}

		}
	}

}
