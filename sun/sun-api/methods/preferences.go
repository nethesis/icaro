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

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func CreateAccountPrefs(c *gin.Context) {
	var json models.AccountPreference
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	accountId := c.MustGet("token").(models.AccessToken).AccountId

	preference := models.AccountPreference{
		AccountId: accountId,
		Key:       json.Key,
		Value:     json.Value,
	}

	db := database.Database()
	db.Save(&preference)
	db.Close()

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func GetAccountPrefs(c *gin.Context) {
	var preferences []models.AccountPreference
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	db := database.Database()
	db.Where("account_id = ?", accountId).Find(&preferences)
	db.Close()

	if len(preferences) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No preferences found!"})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

func CreateHotspotPrefs(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	var json models.HotspotPreference
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	preference := models.HotspotPreference{
		Key:   json.Key,
		Value: json.Value,
	}

	preference.HotspotId = json.HotspotId

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId), json.HotspotId) {
		db := database.Database()
		db.Save(&preference)
		db.Close()

		c.JSON(http.StatusCreated, gin.H{"status": "success"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}

func GetHotspotPrefs(c *gin.Context) {
	var preferences []models.HotspotPreference
	hotspotId := c.Param("hotspot_id")

	db := database.Database()
	db.Where("hotspot_id = ?", hotspotId).Find(&preferences)
	db.Close()

	if len(preferences) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No preferences found!"})
		return
	}

	c.JSON(http.StatusOK, preferences)
}
