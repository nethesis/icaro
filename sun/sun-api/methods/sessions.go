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
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"sun-api/database"
	"sun-api/models"
	"sun-api/utils"
)

func GetSessions(c *gin.Context) {
	var sessions []models.Session
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId)).Offset(offsets[0]).Limit(offsets[1]).Find(&sessions)
	db.Close()

	if len(sessions) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No sessions found!"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func GetSession(c *gin.Context) {
	var session models.Session
	sessionId := c.Param("session_id")

	db := database.Database()
	db.Where("id = ?", sessionId).First(&session)
	db.Close()

	if session.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No session found!"})
		return
	}

	c.JSON(http.StatusOK, session)
}
