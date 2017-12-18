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
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"manager-api/database"
	"manager-api/models"
	"manager-api/utils"
)

func CreateHotspot(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")

	accountId := 1 // TODO get from token

	hotspot := models.Hotspot{
		AccountId:   accountId,
		Name:        name,
		Description: description,
		Created:     time.Now().String(),
	}

	db := database.Database()
	db.Save(&hotspot)

	db.Close()

	c.JSON(http.StatusCreated, gin.H{"id": hotspot.Id, "status": "success"})
}

func UpdateHotspot(c *gin.Context) {
	var hotspot models.Hotspot
	hotspotId := c.Param("hotspot_id")

	name := c.PostForm("name")
	description := c.PostForm("description")

	db := database.Database()
	db.Where("id = ?", hotspotId).First(&hotspot)

	if hotspot.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot found!"})
		return
	}

	hotspot.Name = name
	hotspot.Description = description

	db.Save(&hotspot)

	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetHotspots(c *gin.Context) {
	var hotspots []models.Hotspot

	page := c.Query("page")
	limit := c.Query("limit")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Offset(offsets[0]).Limit(offsets[1]).Find(&hotspots)

	if len(hotspots) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspots found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, hotspots)
}

func GetHotspot(c *gin.Context) {
	var hotspot models.Hotspot
	hotspotId := c.Param("hotspot_id")

	db := database.Database()
	db.Where("id = ?", hotspotId).First(&hotspot)

	if hotspot.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, hotspot)
}

func DeleteHotspot(c *gin.Context) {
	var hotspot models.Hotspot
	hotspotId := c.Param("hotspot_id")

	db := database.Database()
	db.Where("id = ?", hotspotId).First(&hotspot)

	if hotspot.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot found!"})
		return
	}

	db.Delete(&hotspot)

	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
