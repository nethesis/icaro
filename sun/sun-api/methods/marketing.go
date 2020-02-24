/*
 * Copyright (C) 2019 Nethesis S.r.l.
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

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	ade_utils "github.com/nethesis/icaro/ade/ade-api/utils"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func UpdateHotspotMarketing(c *gin.Context) {
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
			c.JSON(http.StatusNotFound, gin.H{"message": "No marketing preference found!"})
			return
		}

		hsPref.Value = json.Value

		db.Save(&hsPref)

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}

func GetHotspotMarketing(c *gin.Context) {
	var preferences []models.HotspotPreference
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	db := database.Instance()
	if accountId == 1 {
		db.Where("hotspot_id = ? AND `key` LIKE 'marketing_%'", hotspotId).Order("`key` asc").Find(&preferences)
	} else {
		db.Where("hotspot_id in (?) AND `key` LIKE 'marketing_%'", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Order("`key` asc").Find(&preferences)
	}

	if len(preferences) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No marketing preferences found!"})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

func SendTestFeedbackEmail(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	var json models.SurveyTestEmail
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
		prefs := ade_utils.GetHotspotPrefs(hotspotIdInt)
		hotspot := utils.GetHotspotById(hotspotId)

		status := ade_utils.SendFeedBackMessageToUser(models.AdeToken{}, json.To, prefs["captive_2_title"], prefs["captive_7_background"], prefs["captive_82_container_bg_color"], prefs["captive_83_title_color"], prefs["captive_84_text_color"], prefs["captive_85_text_style"], hotspot, json.Body)

		if status {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {

			c.JSON(http.StatusBadRequest, gin.H{"message": "Email send failed."})
		}

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}

func SendTestReviewEmail(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	var json models.SurveyTestEmail
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
		prefs := ade_utils.GetHotspotPrefs(hotspotIdInt)
		hotspot := utils.GetHotspotById(hotspotId)

		status := ade_utils.SendReviewMessageToUser(models.AdeToken{}, json.To, prefs["captive_2_title"], prefs["captive_7_background"], prefs["captive_82_container_bg_color"], prefs["captive_83_title_color"], prefs["captive_84_text_color"], prefs["captive_85_text_style"], hotspot, json.Body)

		if status {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {

			c.JSON(http.StatusBadRequest, gin.H{"message": "Email send failed."})
		}

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}
