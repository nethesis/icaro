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
	"manager-api/utils"
)

func CreateUnit(c *gin.Context) {
	macAddress := c.PostForm("mac_address")
	description := c.PostForm("description")
	uuid := c.PostForm("uuid")
	hotspotId := c.PostForm("hotspot_id")

	unit := models.Unit{
		MacAddress:  macAddress,
		Description: description,
		Uuid:        uuid,
	}

	if hotspotIdInt, err := strconv.Atoi(hotspotId); err == nil {
		unit.HotspotId = hotspotIdInt
	}

	db := database.Database()
	db.Save(&unit)

	db.Close()

	c.JSON(http.StatusCreated, gin.H{"id": unit.Id, "status": "success"})
}

func GetUnits(c *gin.Context) {
	var units []models.Unit

	page := c.Query("page")
	limit := c.Query("limit")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Offset(offsets[0]).Limit(offsets[1]).Find(&units)

	if len(units) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No units found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, units)
}

func GetUnit(c *gin.Context) {
	var unit models.Unit
	unitId := c.Param("unit_id")

	db := database.Database()
	db.Where("id = ?", unitId).First(&unit)

	if unit.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No unit found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, unit)
}

func DeleteUnit(c *gin.Context) {
	var unit models.Unit
	unitId := c.Param("unit_id")

	db := database.Database()
	db.Where("id = ?", unitId).First(&unit)

	if unit.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No unit found!"})
		return
	}

	db.Delete(&unit)

	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
