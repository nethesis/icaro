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

func SMSAuth(c *gin.Context) {
	number := c.Param("number")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")

	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "number is required"})
		return
	}

	// check if user exists
	user := utils.GetUserByUsername(number)
	if user.Id == 0 {
		// generate code
		code := utils.GenerateCode(6)

		// send sms with code
		status := utils.SendSMSCode(number, code)
		// check response
		if status != 201 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not send"})
			return
		}

		// create user
		unit := utils.GetUnitByUuid(uuid)
		days := utils.GetHotspotPreferencesByKey(unit.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)
		newUser := models.User{
			HotspotId:   unit.HotspotId,
			Name:        number, // TODO: how we can get the name?
			Username:    number,
			Password:    code,
			Email:       "",
			AccountType: "sms",
			KbpsDown:    0,
			KbpsUp:      0,
			ValidFrom:   time.Now().UTC(),
			ValidUntil:  time.Now().UTC().AddDate(0, 0, daysInt),
		}
		newUser.Id = methods.CreateUser(newUser)

		// create user session check
		utils.CreateUserSession(newUser.Id, sessionId)

		// TODO: create marketing info with user infos and birthday

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": number})
	} else {
		// update user info
		days := utils.GetHotspotPreferencesByKey(user.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, daysInt)
		db := database.Database()
		db.Save(&user)
		db.Close()

		// create user session check
		utils.CreateUserSession(user.Id, sessionId)

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": number, "password": user.Password})
	}
}

func EmailAuth(c *gin.Context) {
	email := c.Param("email")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is required"})
		return
	}

	// check if user exists
	user := utils.GetUserByUsername(email)
	if user.Id == 0 {
		// generate code
		code := utils.GenerateCode(6)

		// send email with code
		status := utils.SendEmailCode(email, code)

		// check response
		if !status {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not send"})
			return
		}

		// create user
		unit := utils.GetUnitByUuid(uuid)
		days := utils.GetHotspotPreferencesByKey(unit.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)
		newUser := models.User{
			HotspotId:   unit.HotspotId,
			Name:        email, // TODO: how we can get the name?
			Username:    email,
			Password:    code,
			Email:       "",
			AccountType: "email",
			KbpsDown:    0,
			KbpsUp:      0,
			ValidFrom:   time.Now().UTC(),
			ValidUntil:  time.Now().UTC().AddDate(0, 0, daysInt),
		}
		newUser.Id = methods.CreateUser(newUser)

		// create user session check
		utils.CreateUserSession(newUser.Id, sessionId)

		// TODO: create marketing info with user infos and birthday

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": email})
	} else {
		// update user info
		days := utils.GetHotspotPreferencesByKey(user.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, daysInt)
		db := database.Database()
		db.Save(&user)
		db.Close()

		// create user session check
		utils.CreateUserSession(user.Id, sessionId)

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": email, "password": user.Password})
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
		if voucher.Expires.Before(time.Now().UTC()) {
			c.JSON(http.StatusOK, gin.H{"message": "Voucher is expired"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Voucher is valid"})
		}
	}

}
