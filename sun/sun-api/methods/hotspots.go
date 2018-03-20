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
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func CreateHotspot(c *gin.Context) {
	var json models.Hotspot
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspot := models.Hotspot{
		AccountId:   accountId,
		Name:        json.Name,
		Description: json.Description,
		Created:     time.Now().UTC(),
	}

	db := database.Database()
	db.Save(&hotspot)
	db.Close()

	if hotspot.Id == 0 {
		c.JSON(http.StatusConflict, gin.H{"id": hotspot.Id, "status": "hotspot already exists"})
	} else {
		// initialize hotspot preferences
		utils.SetDefaultHotspotPreferences(hotspot.Id)

		c.JSON(http.StatusCreated, gin.H{"id": hotspot.Id, "status": "success"})
	}
}

func UpdateHotspot(c *gin.Context) {
	var hotspot models.Hotspot
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	var json models.Hotspot
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Database()

	if accountId == 1 {
		db.Where("id = ?", hotspotId).First(&hotspot)
	} else {
		db.Where("id = ? AND account_id = ?", hotspotId, accountId).First(&hotspot)
	}

	if hotspot.Id == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot found!"})
		return
	}

	if len(json.Description) > 0 {
		hotspot.Description = json.Description
	}

	db.Save(&hotspot)
	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetHotspots(c *gin.Context) {
	var hotspots []models.Hotspot
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	if accountId == 1 {
		db.Offset(offsets[0]).Limit(offsets[1]).Find(&hotspots)
	} else {
		db.Where("account_id = ?", accountId).Offset(offsets[0]).Limit(offsets[1]).Find(&hotspots)
	}
	db.Close()

	if len(hotspots) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspots found!"})
		return
	}

	c.JSON(http.StatusOK, hotspots)
}

func GetHotspot(c *gin.Context) {
	var hotspot models.HotspotJSON
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	db := database.Database()
	if accountId == 1 {
		db.Select("hotspots.*, accounts.name as account_name").Where("hotspots.id = ?", hotspotId).Joins("JOIN accounts on accounts.id = hotspots.account_id").First(&hotspot)
	} else {
		db.Select("hotspots.*, accounts.name as account_name").Where("hotspots.id = ? AND account_id = ?", hotspotId, accountId).Joins("JOIN accounts on accounts.id = hotspots.account_id").First(&hotspot)
	}
	db.Close()

	if hotspot.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot found!"})
		return
	}

	c.JSON(http.StatusOK, hotspot)
}

func DeleteHotspot(c *gin.Context) {
	var hotspot models.Hotspot
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	db := database.Database()
	if accountId == 1 {
		db.Where("id = ?", hotspotId).First(&hotspot)
	} else {
		db.Where("id = ? AND account_id = ?", hotspotId, accountId).First(&hotspot)
	}

	if hotspot.Id == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot found!"})
		return
	}

	db.Delete(&hotspot)
	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func StatsHotspotTotal(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId
	var count int

	db := database.Database()
	if accountId == 1 {
		db.Table("hotspots").Count(&count)
	} else {
		db.Table("hotspots").Where("account_id = ?", accountId).Count(&count)
	}
	db.Close()

	c.JSON(http.StatusOK, gin.H{"total": count})
}
