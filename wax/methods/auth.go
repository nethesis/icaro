/*
 * Copyright (C) 2021 Nethesis S.r.l.
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

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/wax/utils"
)

func GetDaemonAuth(c *gin.Context) {
	var daemonAuth models.DaemonAuth

	// get params
	sessionId := c.Query("sessionid")
	unitUuid := c.Query("uuid")

	// search in database
	db := database.Instance()
	db.Where("session_id = ? AND unit_uuid = ?", sessionId, unitUuid).First(&daemonAuth)

	if daemonAuth.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No auth found!"})
		return
	}

	// return auth
	c.JSON(http.StatusOK, daemonAuth)
}

func GetDaemonLogin(c *gin.Context) {
	// get params
	sessionId := c.Query("sessionid")
	unitUuid := c.Query("uuid")
	username := c.Query("username")

	// create user auth
	utils.CreateUserAuth(sessionId, 0, unitUuid, 0, username, "", "login")

  // return result
  c.JSON(http.StatusCreated, gin.H{"clientState": "1"})
}

func GetDaemonTemporary(c *gin.Context) {
	// get params
	sessionId := c.Query("sessionid")
	unitUuid := c.Query("uuid")
	username := c.Query("username")

	// get unit
	unit := utils.GetUnitByUuid(unitUuid);

	// get session timeout pref
	seconds := utils.GetHotspotPreferencesByKey(unit.HotspotId, "temp_session_duration")
	secondsInt, _ := strconv.Atoi(seconds.Value)

	// create user auth
	utils.CreateUserAuth(sessionId, secondsInt, unitUuid, 0, username, "", "temporary")

  // return result
  c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func GetDaemonLogout(c *gin.Context) {
	// get params
	sessionId := c.Query("sessionid")
	unitUuid := c.Query("uuid")
	username := c.Query("username")

	// create user auth
	utils.CreateUserAuth(sessionId, 0, unitUuid, 0, username, "", "logout")

  // return result
  c.JSON(http.StatusCreated, gin.H{"status": "success"})
}
