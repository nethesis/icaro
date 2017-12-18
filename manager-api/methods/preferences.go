/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Windmill project.
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
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"manager-api/database"
	"manager-api/models"
)

func CreateAccountPrefs(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")

	accountId := 1 // TODO get from token

	preference := models.AccountPreference{
		AccountId: accountId,
		Key:       key,
		Value:     value,
	}

	db := database.Database()
	db.Save(&preference)

	db.Close()

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func GetAccountPrefs(c *gin.Context) {
	var preferences []models.AccountPreference
	accountId := 1 // TODO get from token

	db := database.Database()
	db.Where("account_id = ?", accountId).Find(&preferences)

	if len(preferences) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No preferences found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, preferences)
}

func CreateHotspotPrefs(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")
	hotspotId := c.PostForm("hotspot_id")

	preference := models.HotspotPreference{
		Key:   key,
		Value: value,
	}

	if hotspotIdInt, err := strconv.Atoi(hotspotId); err == nil {
		preference.HotspotId = hotspotIdInt
	}

	db := database.Database()
	db.Save(&preference)

	db.Close()

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func GetHotspotPrefs(c *gin.Context) {
	var preferences []models.HotspotPreference
	hotspotId := 1 // TODO get from token

	db := database.Database()
	db.Where("hotspot_id = ?", hotspotId).Find(&preferences)

	if len(preferences) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No preferences found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, preferences)
}
