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

func CreateUnit(c *gin.Context) {
	accountId := c.MustGet("token").(*models.AccessToken).AccountId

	macAddress := c.PostForm("mac_address")
	description := c.PostForm("description")
	uuid := c.PostForm("uuid")
	hotspotId := c.PostForm("hotspot_id")

	unit := models.Unit{
		MacAddress:  macAddress,
		Description: description,
		Uuid:        uuid,
		Created:     time.Now().UTC(),
	}

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		fmt.Println(err.Error())
	}
	unit.HotspotId = hotspotIdInt

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId), hotspotIdInt) {
		db := database.Database()
		db.Save(&unit)
		db.Close()

		c.JSON(http.StatusCreated, gin.H{"id": unit.Id, "status": "success"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}

func GetUnits(c *gin.Context) {
	var units []models.Unit
	accountId := c.MustGet("token").(*models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId)).Offset(offsets[0]).Limit(offsets[1]).Find(&units)

	if len(units) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No units found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, units)
}

func GetUnit(c *gin.Context) {
	var unit models.Unit
	accountId := c.MustGet("token").(*models.AccessToken).AccountId

	unitId := c.Param("unit_id")

	db := database.Database()
	db.Where("id = ? AND hotspot_id in (?)", unitId, utils.ExtractHotspotIds(accountId)).First(&unit)

	if unit.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No unit found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, unit)
}

func DeleteUnit(c *gin.Context) {
	var unit models.Unit
	accountId := c.MustGet("token").(*models.AccessToken).AccountId

	unitId := c.Param("unit_id")

	db := database.Database()
	db.Where("id = ? AND hotspot_id in (?)", unitId, utils.ExtractHotspotIds(accountId)).First(&unit)

	if unit.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No unit found!"})
		return
	}

	db.Delete(&unit)

	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
