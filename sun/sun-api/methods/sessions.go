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
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func GetSessions(c *gin.Context) {
	var sessions []models.Session
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")
	hotspotId := c.Query("hotspot")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Offset(offsets[0]).Limit(offsets[1]).Find(&sessions)
	db.Close()

	if len(sessions) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No sessions found!"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func GetSession(c *gin.Context) {
	var session models.Session
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	sessionId := c.Param("session_id")

	db := database.Database()
	db.Where("id = ? AND hotspot_id in (?)", sessionId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&session)
	db.Close()

	if session.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No session found!"})
		return
	}

	c.JSON(http.StatusOK, session)
}
