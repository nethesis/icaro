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

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func UpdateAccountPrefs(c *gin.Context) {
	var accountPref models.AccountPreference
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	var json models.AccountPreference
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()
	if accountId == 1 {
		db.Where("`key` = ?", json.Key).First(&accountPref)
	} else {
		db.Where("`key` = ? AND account_id = ?", json.Key, accountId).First(&accountPref)
	}

	if accountPref.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No preference found!"})
		return
	}

	accountPref.Value = json.Value

	db.Save(&accountPref)

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func GetAccountPrefs(c *gin.Context) {
	var preferences []models.AccountPreference
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	db := database.Instance()
	db.Where("account_id = ?", accountId).Find(&preferences)

	if len(preferences) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No preferences found!"})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

func UpdateHotspotPrefs(c *gin.Context) {
	var hsPref models.HotspotPreference
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	var json models.HotspotPreference
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	// convert hotspot id to int
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Hotspot id error", "error": err.Error()})
		return
	}

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt), hotspotIdInt) {
		db := database.Instance()
		db.Where("`key` = ? AND hotspot_id = ?", json.Key, hotspotIdInt).First(&hsPref)

		if hsPref.Id == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No preference found!"})
			return
		}

		// enforce authorization on captive portal configuration
		if strings.Index(json.Key, "captive") == 0 {
			if !utils.CanChangeCaptivePortalOptions(accountId) {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Can't update captive portal preference"})
				return
			}
		}

		hsPref.Value = json.Value

		db.Save(&hsPref)

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}

func GetHotspotPrefs(c *gin.Context) {
	var preferences []models.HotspotPreference
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	db := database.Instance()
	if accountId == 1 {
		db.Where("hotspot_id = ?", hotspotId).Order("`key` asc").Find(&preferences)
	} else {
		db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Order("`key` asc").Find(&preferences)
	}

	if len(preferences) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No preferences found!"})
		return
	}

	c.JSON(http.StatusOK, preferences)
}
