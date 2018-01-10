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
	"time"
	"wax/utils"

	"github.com/gin-gonic/gin"

	"sun-api/database"
	"sun-api/methods"
	"sun-api/models"
)

func SMSAuth(c *gin.Context) {
	number := c.Param("number")
	uuid := c.Query("uuid")

	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "number is required"})
		return
	}

	// check if user exists
	user := utils.GetUserByUsername(number)
	if user.Id == 0 {
		// generate code
		code := utils.GenerateCode(6)

		// create user
		unit := utils.GetUnitByUuid(uuid)
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
			ValidUntil:  time.Now().UTC().AddDate(0, 0, 30), // TODO: get days from hotspot account preferences
		}
		methods.CreateUser(newUser)

		// send sms with code
		utils.SendSMSCode(number, code)

		// TODO: create marketing info with user infos and birthday

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": number})
	} else {
		// update user info
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, 30) // TODO: days info from hotspot account preferences
		db := database.Database()
		db.Save(&user)
		db.Close()

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": number, "password": user.Password})
	}
}

func EmailAuth(c *gin.Context) {
	email := c.Param("email")
	uuid := c.Query("uuid")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is required"})
		return
	}

	// check if user exists
	user := utils.GetUserByUsername(email)
	if user.Id == 0 {
		// generate code
		code := utils.GenerateCode(6)

		// create user
		unit := utils.GetUnitByUuid(uuid)
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
			ValidUntil:  time.Now().UTC().AddDate(0, 0, 30), // TODO: get days from hotspot account preferences
		}
		methods.CreateUser(newUser)

		// send email with code
		utils.SendEmailCode(email, code)

		// TODO: create marketing info with user infos and birthday

		// response to client
		c.JSON(http.StatusOK, gin.H{"user_id": email})
	} else {
		// update user info
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, 30) // TODO: days info from hotspot account preferences
		db := database.Database()
		db.Save(&user)
		db.Close()

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
