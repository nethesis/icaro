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

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"sun-api/database"
	"sun-api/models"
	"sun-api/utils"
)

func CreateUser(user models.User) {
	password := "" // TODO generate randomly

	user.Password = password
	user.Created = time.Now().UTC()

	// save new user
	db := database.Database()
	db.Save(&user)
	db.Close()
}

func UpdateUser(c *gin.Context) {
	var user models.User
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	userId := c.Param("user_id")

	var json models.User
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Database()
	db.Where("id = ?", userId).First(&user)

	if user.Id == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId), user.HotspotId) {
		user.Name = json.Name
		user.Email = json.Email

		user.KbpsDown = json.KbpsDown
		user.KbpsUp = json.KbpsUp

		user.ValidFrom = json.ValidFrom
		user.ValidUntil = json.ValidUntil

		db.Save(&user)
		db.Close()

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		db.Close()
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This user is not yours"})
	}
}

func GetUsers(c *gin.Context) {
	var users []models.User
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId)).Offset(offsets[0]).Limit(offsets[1]).Find(&users)
	db.Close()

	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No users found!"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	var user models.User
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	userId := c.Param("user_id")

	db := database.Database()
	db.Where("id = ? AND hotspot_id in (?)", userId, utils.ExtractHotspotIds(accountId)).First(&user)
	db.Close()

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	userId := c.Param("user_id")

	db := database.Database()
	db.Where("id = ? AND hotspot_id in (?)", userId, utils.ExtractHotspotIds(accountId)).First(&user)

	if user.Id == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	db.Delete(&user)
	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
