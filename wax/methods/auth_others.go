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
	user := utils.GetUserByUsername(number)
	if user.Id == 0 {
		// get unit
		unit := utils.GetUnitByUuid(uuid)

		// create user
		days := utils.GetHotspotPreferencesByKey(unit.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)

		down := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Down")
		downInt, _ := strconv.Atoi(down.Value)
		up := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Up")
		upInt, _ := strconv.Atoi(up.Value)

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
			}
		}

		// generate code
		code := utils.GenerateCode(6)

		newUser := models.User{
			HotspotId:     unit.HotspotId,
			Name:          number,
			Username:      number,
			Password:      code,
			Email:         "",
			AccountType:   "sms",
			MarketingAuth: true,
			KbpsDown:      downInt,
			KbpsUp:        upInt,
			AutoLogin:     autoLoginBool,
			ValidFrom:     time.Now().UTC(),
			ValidUntil:    time.Now().UTC().AddDate(0, 0, daysInt+1),
		}
		newUser.Id = methods.CreateUser(newUser)

		// send sms with code
		userIdStr := strconv.Itoa(newUser.Id)
		status := utils.SendSMSCode(number, code, unit, "digest="+digest+"&uuid="+uuid+"&sessionid="+sessionId+"&uamip="+uamip+"&uamport="+uamport+"&user="+userIdStr)

		// check response
		if status != 201 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not send"})
			return
		}

		// add sms statistics
		hotspotSmsCount := models.HotspotSmsCount{
			HotspotId: unit.HotspotId,
			UniId:     unit.Id,
			Number:    number,
			Reset:     false,
			Sent:      time.Now().UTC(),
		}
		utils.SaveHotspotSMSCount(hotspotSmsCount)

		// create marketing info with user infos
		utils.CreateUserMarketing(newUser.Id, smsMarketingData{Number: number}, "sms")

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
				user.ValidUntil = time.Now().UTC().AddDate(0, 0, voucher.Duration)
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
			status := utils.SendSMSCode(number, code, unit, "digest="+digest+"&uuid="+uuid+"&sessionid="+sessionId+"&uamip="+uamip+"&uamport="+uamport+"&user="+userIdStr)

			// check response
			if status != 201 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not send"})
				return
			}

			// add sms statistics
			hotspotSmsCount := models.HotspotSmsCount{
				HotspotId: unit.HotspotId,
				UniId:     unit.Id,
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

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is required"})
		return
	}

	// check if user exists
	user := utils.GetUserByUsername(email)
	if user.Id == 0 {
		// get unit
		unit := utils.GetUnitByUuid(uuid)

		// create user
		days := utils.GetHotspotPreferencesByKey(unit.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)

		down := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Down")
		downInt, _ := strconv.Atoi(down.Value)
		up := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Up")
		upInt, _ := strconv.Atoi(up.Value)

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
			}
		}

		// generate code
		code := utils.GenerateCode(6)

		newUser := models.User{
			HotspotId:     unit.HotspotId,
			Name:          email,
			Username:      email,
			Password:      code,
			Email:         email,
			AccountType:   "email",
			MarketingAuth: true,
			KbpsDown:      downInt,
			KbpsUp:        upInt,
			AutoLogin:     autoLoginBool,
			ValidFrom:     time.Now().UTC(),
			ValidUntil:    time.Now().UTC().AddDate(0, 0, daysInt+1),
		}
		newUser.Id = methods.CreateUser(newUser)

		// send email with code
		userIdStr := strconv.Itoa(newUser.Id)
		status := utils.SendEmailCode(email, code, unit, "digest="+digest+"&uuid="+uuid+"&sessionid="+sessionId+"&uamip="+uamip+"&uamport="+uamport+"&user="+userIdStr)

		// check response
		if !status {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not sent"})
			return
		}

		// create marketing info with user infos
		utils.CreateUserMarketing(newUser.Id, emailMarketingData{Email: email}, "email")

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
				user.ValidUntil = time.Now().UTC().AddDate(0, 0, voucher.Duration)
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
			status := utils.SendEmailCode(email, code, unit, "digest="+digest+"&uuid="+uuid+"&sessionid="+sessionId+"&uamip="+uamip+"&uamport="+uamport+"&user="+userIdStr)

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
	user := utils.GetUserByUsername(mac)
	if user.Id == 0 {
		// get unit
		unit := utils.GetUnitByUuid(uuid)

		// create user
		downInt, _ := strconv.Atoi(kbps_down)
		upInt, _ := strconv.Atoi(kbps_up)

		newUser := models.User{
			HotspotId:   unit.HotspotId,
			Name:        name,
			Username:    mac,
			Password:    "",
			Email:       "",
			AccountType: "mac",
			KbpsDown:    downInt,
			KbpsUp:      upInt,
			AutoLogin:   true,
			ValidFrom:   time.Now().UTC(),
			ValidUntil:  time.Now().UTC().AddDate(0, 0, 3650), // ten years
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

	// extract unit
	unit := utils.GetUnitByUuid(uuid)

	// extract voucher
	voucher := utils.GetVoucherByCode(code, unit.HotspotId)

	// check voucher validity
	if voucher.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Voucher is invalid"})
	} else {
		// check if is expired
		if !voucher.Expires.IsZero() && voucher.Expires.Before(time.Now().UTC()) {
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

					// check if is not limitless
					if voucher.RemainUse != -1 {
						voucher.RemainUse--
					}

					db := database.Instance()
					db.Save(&voucher)
				}

				c.JSON(http.StatusOK, gin.H{"message": "Voucher is valid", "code": voucher.Code})
			}

		}
	}

}
