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

func GetDevices(c *gin.Context) {
	var devices []models.Device
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1))).Offset(offsets[0]).Limit(offsets[1]).Find(&devices)
	db.Close()

	if len(devices) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No devices found!"})
		return
	}

	c.JSON(http.StatusOK, devices)
}

func GetDevice(c *gin.Context) {
	var device models.Device
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	deviceId := c.Param("device_id")

	db := database.Database()
	db.Where("id = ? AND hotspot_id in (?)", deviceId, utils.ExtractHotspotIds(accountId, (accountId == 1))).First(&device)
	db.Close()

	if device.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No device found!"})
		return
	}

	c.JSON(http.StatusOK, device)
}
