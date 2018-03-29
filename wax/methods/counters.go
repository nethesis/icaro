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
 * author: Giacomo Sanchietti <giacomo.sanchietti@nethesis.it>
 */

package methods

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/wax/utils"
)

func Ack(c *gin.Context, act int) {
	c.String(http.StatusOK, "Ack: %d", act)
}

func startSession(userName string, deviceMacAddress string, deviceIp string, sessionId string, hotspotName string, unitMacAddress string) int {
	unit := utils.GetUnitByMacAddress(unitMacAddress)
	if unit.Id <= 0 {
		return 0
	}
	user := utils.GetUserByNameAndHotspotId(userName, unit.HotspotId)
	if user.Id <= 0 {
		return 0
	}
	device := utils.GetDeviceByMacAddress(deviceMacAddress)
	if device.Id <= 0 {
		device.HotspotId = unit.HotspotId
		device.UserId = user.Id
		device.MacAddress = deviceMacAddress
		device.IpAddress = deviceIp
		device.Description = ""
		device.Created = time.Now().UTC()

		// save new device
		db := database.Database()
		db.Save(&device)
		db.Close()
	} else {
		if device.IpAddress != deviceIp {
			db := database.Database()
			db.Model(&device).Update("IpAddress", deviceIp)
			db.Close()
		}
	}
	session := utils.GetSessionByKeyAndUnitId(sessionId, unit.Id)
	session.UnitId = unit.Id
	session.HotspotId = unit.HotspotId
	session.DeviceId = device.Id
	session.DeviceMAC = device.MacAddress
	session.UserId = user.Id
	session.Username = user.Username
	session.BytesUp = 0
	session.BytesDown = 0
	session.StartTime = time.Now().UTC()
	session.UpdateTime = time.Now().UTC()
	session.SessionKey = sessionId

	db := database.Database()
	db.Save(&session)
	db.Close()
	return 1
}

func stopSession(sessionId string, unitMacAddress string, bytesDown string, bytesUp string, duration string) int {
	unit := utils.GetUnitByMacAddress(unitMacAddress)
	if unit.Id <= 0 {
		return 0
	}
	session := utils.GetSessionByKeyAndUnitId(sessionId, unit.Id)
	if session.Id < 0 {
		return 0
	}

	session.UpdateTime = time.Now().UTC()
	session.StopTime = time.Now().UTC()
	if d, err := strconv.Atoi(duration); err == nil {
		session.Duration = d
	}
	if bd, err := strconv.Atoi(bytesDown); err == nil {
		session.BytesDown = bd
	}
	if bu, err := strconv.Atoi(bytesUp); err == nil {
		session.BytesUp = bu
	}

	db := database.Database()
	db.Save(&session)
	db.Close()
	return 1
}

func updateSession(sessionId string, unitMacAddress string, bytesDown string, bytesUp string, duration string) int {
	unit := utils.GetUnitByMacAddress(unitMacAddress)
	if unit.Id <= 0 {
		return 0
	}
	session := utils.GetSessionByKeyAndUnitId(sessionId, unit.Id)
	if session.Id < 0 {
		return 0
	}

	session.UpdateTime = time.Now().UTC()
	if d, err := strconv.Atoi(duration); err == nil {
		session.Duration = d
	}
	if bd, err := strconv.Atoi(bytesDown); err == nil {
		session.BytesDown = bd
	}
	if bu, err := strconv.Atoi(bytesUp); err == nil {
		session.BytesUp = bu
	}

	db := database.Database()
	db.Save(&session)
	db.Close()
	return 1
}

func accountingOn(ap string) int {

	var sessions []models.Session

	unit := utils.GetUnitByMacAddress(ap)
	if unit.Id <= 0 {
		return 0
	}

	db := database.Database()
	db.Where("unit_id = ? and stop_time = 0", unit.Id).Find(&sessions)

	for i, _ := range sessions {
		sessions[i].StopTime = time.Now().UTC()
		db.Save(&sessions[i])
	}

	db.Close()

	return 1
}

func Counters(c *gin.Context, parameters url.Values) {
	status := parameters.Get("status")

	var username string

	switch status {
	case "start":
		if strings.Compare(c.Query("user"), c.Query("mac")) == 0 {
			//autologin
			_, user := utils.GetUserByMacAddressAndunitMacAddress(c.Query("mac"), c.Query("ap"))
			username = user.Username
		} else {
			username = c.Query("user")
		}
		Ack(c, startSession(username, c.Query("mac"), c.Query("ip"), c.Query("sessionid"), c.Query("nasid"), c.Query("ap")))
	case "stop":
		Ack(c, stopSession(c.Query("sessionid"), c.Query("ap"), c.Query("bytes_down"), c.Query("bytes_up"), c.Query("duration")))
	case "update":
		Ack(c, updateSession(c.Query("sessionid"), c.Query("ap"), c.Query("bytes_down"), c.Query("bytes_up"), c.Query("duration")))
	case "up":
		Ack(c, accountingOn(c.Query("ap")))
	case "":
		c.String(http.StatusBadRequest, "No status provided")
	default:
		c.String(http.StatusNotImplemented, "Invalid status: '%s'", status)
	}
}
