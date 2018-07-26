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

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func GetDevices(c *gin.Context) {
	var devices []models.Device
	var total int
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")
	hotspotId := c.Query("hotspot")
	userId := c.Query("user")
	q := c.Query("q")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	offsets := utils.OffsetCalc(page, limit)

	db := database.Instance()
	chain := db.Preload("User").Preload("Hotspot").Where("devices.hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt))

	if len(userId) > 0 {
		chain = chain.Where("devices.user_id = ?", userId)
	}

	if len(q) > 0 {
		chain = chain.Select("users.*, devices.*").Joins("JOIN hotspots on hotspots.id = devices.hotspot_id JOIN users on users.id = devices.user_id").
			Where("devices.mac_address LIKE ? OR devices.ip_address LIKE ? OR users.name LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	chain.Find(&devices).Count(&total)
	chain.Offset(offsets[0]).Limit(offsets[1]).Find(&devices)

	if len(devices) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No devices found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": devices, "total": total})
}

func GetDevice(c *gin.Context) {
	var device models.Device
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	deviceId := c.Param("device_id")

	db := database.Instance()
	db.Where("id = ? AND hotspot_id in (?)", deviceId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&device)

	if device.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No device found!"})
		return
	}

	c.JSON(http.StatusOK, device)
}

func StatsDeviceTotal(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId
	var count int

	db := database.Instance()
	db.Table("devices").Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).Count(&count)

	c.JSON(http.StatusOK, gin.H{"total": count})
}
