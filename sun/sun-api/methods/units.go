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

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func CreateUnit(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	var json models.UnitJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	hotspot := utils.GetHotspotById(json.HotspotId)

	if hotspot.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot found!"})
		return
	}

	// check if hotspot belongs to a reseller with a plan which is over quota
	if utils.HotspotIsOverQuota(hotspot.Id) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Units over quota"})
		return
	}

	unit := models.Unit{
		HotspotId:   hotspot.Id,
		MacAddress:  json.MacAddress,
		Name:        json.Name,
		Description: json.Description,
		Uuid:        json.Uuid,
		Secret:      json.Secret,
		Created:     time.Now().UTC(),
	}

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId, (accountId == 1), hotspot.Id), hotspot.Id) {
		db := database.Instance()
		db.Save(&unit)

		if unit.Id == 0 {
			c.JSON(http.StatusConflict, gin.H{"id": unit.Id, "status": "unit already exists"})
		} else {
			c.JSON(http.StatusCreated, gin.H{"id": unit.Id, "status": "success"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}

func GetUnits(c *gin.Context) {
	var units []models.Unit
	var total int
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")
	hotspotId := c.Query("hotspot")
	q := c.Query("q")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	offsets := utils.OffsetCalc(page, limit)

	db := database.Instance()
	chain := db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt))

	if len(q) > 0 {
		chain = chain.Where("name LIKE ? OR description LIKE ? OR mac_address LIKE ? OR uuid LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	chain = chain.Find(&units).Count(&total)
	chain = chain.Offset(offsets[0]).Limit(offsets[1]).Find(&units)

	if len(units) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No units found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": units, "total": total})
}

func GetUnit(c *gin.Context) {
	var unit models.Unit
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	unitId := c.Param("unit_id")

	db := database.Instance()
	db.Where("id = ? AND hotspot_id in (?)", unitId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&unit)

	if unit.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No unit found!"})
		return
	}

	c.JSON(http.StatusOK, unit)
}

func DeleteUnit(c *gin.Context) {
	var unit models.Unit
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	unitId := c.Param("unit_id")

	db := database.Instance()
	db.Where("id = ? AND hotspot_id in (?)", unitId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&unit)

	if unit.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No unit found!"})
		return
	}

	db.Delete(&unit)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func StatsUnitTotal(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId
	var count int

	db := database.Instance()
	db.Table("units").Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).Count(&count)

	c.JSON(http.StatusOK, gin.H{"total": count})
}
