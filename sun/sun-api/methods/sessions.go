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

func GetSessions(c *gin.Context) {
	var sessions []models.Session
	var total int
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")
	hotspotId := c.Query("hotspot")
	userId := c.Query("user")
	unitId := c.Query("unit")
	from := c.Query("from")
	to := c.Query("to")
	q := c.Query("q")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	offsets := utils.OffsetCalc(page, limit)

	db := database.Instance()
	chain := db.Preload("Unit").Preload("User").Where("sessions.hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt))

	if len(userId) > 0 {
		chain = chain.Where("sessions.user_id = ?", userId)
	}

	if len(unitId) > 0 {
		chain = chain.Where("sessions.unit_id = ?", unitId)
	}

	if len(from) > 0 {
		chain = chain.Where("sessions.start_time >= ?", from)
	}

	if len(to) > 0 {
		chain = chain.Where("sessions.start_time <= ?", to)
	}

	if len(q) > 0 {
		chain = chain.Select("sessions.*, units.*, users.*").Joins("JOIN units on units.id = sessions.unit_id").Joins("JOIN users on users.id = sessions.user_id").
			Where("units.name LIKE ? OR sessions.username LIKE ? OR sessions.ip_address LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	chain.Find(&sessions).Count(&total)
	chain.Offset(offsets[0]).Limit(offsets[1]).Find(&sessions)

	if len(sessions) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No sessions found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sessions, "total": total})
}

func GetSessionsHistory(c *gin.Context) {
	var sessions []models.SessionHistory
	var total int
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")
	hotspotId := c.Query("hotspot")
	userId := c.Query("user")
	unitId := c.Query("unit")
	from := c.Query("from")
	to := c.Query("to")
	q := c.Query("q")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	offsets := utils.OffsetCalc(page, limit)

	db := database.Instance()
	chain := db.Preload("Unit").Preload("User").Where("session_histories.hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt))

	if len(userId) > 0 {
		chain = chain.Where("session_histories.user_id = ?", userId)
	}

	if len(unitId) > 0 {
		chain = chain.Where("session_histories.unit_id = ?", unitId)
	}

	if len(from) > 0 {
		chain = chain.Where("session_histories.start_time >= ?", from)
	}

	if len(to) > 0 {
		chain = chain.Where("session_histories.start_time <= ?", to)
	}

	if len(q) > 0 {
		chain = chain.Select("session_histories.*, units.*, users.*").Joins("JOIN units on units.id = session_histories.unit_id").Joins("JOIN users on users.id = session_histories.user_id").
			Where("units.name LIKE ? OR session_histories.username LIKE ? OR session_histories.ip_address LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	chain.Find(&sessions).Count(&total)
	chain.Offset(offsets[0]).Limit(offsets[1]).Find(&sessions)

	if len(sessions) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No session histories found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sessions, "total": total})
}

func GetSession(c *gin.Context) {
	var session models.Session
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	sessionId := c.Param("session_id")

	db := database.Instance()
	db.Where("id = ? AND hotspot_id in (?)", sessionId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&session)

	if session.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No session found!"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func GetSessionHistory(c *gin.Context) {
	var session models.SessionHistory
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	historyId := c.Param("history_id")

	db := database.Instance()
	db.Where("id = ? AND hotspot_id in (?)", historyId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&session)

	if session.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No session history found!"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func StatsSessionTotal(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId
	var count int

	db := database.Instance()
	db.Table("sessions").Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).Count(&count)

	c.JSON(http.StatusOK, gin.H{"total": count})
}
