/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Icaro project.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with NethServer.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package methods

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"sun-api/database"
	"sun-api/models"
	"sun-api/utils"
)

func CreateUser(c *gin.Context) {
	accountId := c.MustGet("token").(*models.AccessToken).AccountId

	name := c.PostForm("name")
	username := c.PostForm("username")
	email := c.PostForm("email")
	accountType := c.PostForm("type")
	kbpsDown := c.PostForm("kbps_down")
	kbpsUp := c.PostForm("kbps_up")
	validFrom := c.PostForm("valid_from")
	validUntil := c.PostForm("valid_until")
	hotspotId := c.PostForm("hotspot_id")

	password := "password,1234" // TODO generate randomly

	user := models.User{
		Name:        name,
		Username:    username,
		Password:    password,
		Email:       email,
		AccountType: accountType,
		Created:     time.Now().UTC(),
	}

	fromT, _ := time.Parse("2017-12-31 00:00:00", validFrom)
	user.ValidFrom = fromT

	untilT, _ := time.Parse("2017-12-31 00:00:00", validUntil)
	user.ValidUntil = untilT

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		fmt.Println(err.Error())
	}
	user.HotspotId = hotspotIdInt

	kbpsDownInt, err := strconv.Atoi(kbpsDown)
	if err != nil {
		fmt.Println(err.Error())
	}
	user.KbpsDown = kbpsDownInt

	kbpsUpInt, err := strconv.Atoi(kbpsUp)
	if err != nil {
		fmt.Println(err.Error())
	}
	user.KbpsUp = kbpsUpInt

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId), hotspotIdInt) {
		db := database.Database()
		db.Save(&user)
		db.Close()

		c.JSON(http.StatusCreated, gin.H{"id": user.Id, "status": "success"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}

func UpdateUser(c *gin.Context) {
	var user models.User
	accountId := c.MustGet("token").(*models.AccessToken).AccountId

	userId := c.Param("user_id")

	name := c.PostForm("name")
	email := c.PostForm("email")
	kbpsDown := c.PostForm("kbps_down")
	kbpsUp := c.PostForm("kbps_up")
	validFrom := c.PostForm("valid_from")
	validUntil := c.PostForm("valid_until")

	db := database.Database()
	db.Where("id = ?", userId).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId), user.HotspotId) {
		user.Name = name
		user.Email = email

		kbpsDownInt, err := strconv.Atoi(kbpsDown)
		if err != nil {
			fmt.Println(err.Error())
		}
		user.KbpsDown = kbpsDownInt

		kbpsUpInt, err := strconv.Atoi(kbpsUp)
		if err != nil {
			fmt.Println(err.Error())
		}
		user.KbpsUp = kbpsUpInt

		fromT, _ := time.Parse("2017-12-31 00:00:00", validFrom)
		user.ValidFrom = fromT

		untilT, _ := time.Parse("2017-12-31 00:00:00", validUntil)
		user.ValidUntil = untilT

		db.Save(&user)
		db.Close()

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This user is not yours"})
	}
}

func GetUsers(c *gin.Context) {
	var users []models.User
	accountId := c.MustGet("token").(*models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId)).Offset(offsets[0]).Limit(offsets[1]).Find(&users)

	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No users found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	var user models.User
	accountId := c.MustGet("token").(*models.AccessToken).AccountId

	userId := c.Param("user_id")

	db := database.Database()
	db.Where("id = ? AND hotspot_id in (?)", userId, utils.ExtractHotspotIds(accountId)).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	accountId := c.MustGet("token").(*models.AccessToken).AccountId

	userId := c.Param("user_id")

	db := database.Database()
	db.Where("id = ? AND hotspot_id in (?)", userId, utils.ExtractHotspotIds(accountId)).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	db.Delete(&user)

	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
