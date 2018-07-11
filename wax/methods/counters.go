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
		db := database.Instance()
		db.Save(&device)
	} else {
		if device.IpAddress != deviceIp {
			db := database.Instance()
			db.Model(&device).Update("IpAddress", deviceIp)
		}
	}

	currentHotspot := utils.GetHotspotById(unit.HotspotId)
	session := utils.GetSessionByKeyAndUnitId(sessionId, unit.Id)
	session.UnitId = unit.Id
	session.UnitMac = unit.MacAddress
	session.HotspotId = unit.HotspotId
	session.HotspotDesc = currentHotspot.Name + " - " + currentHotspot.Description
	session.DeviceId = device.Id
	session.DeviceMAC = device.MacAddress
	session.UserId = user.Id
	session.Username = user.Username
	session.BytesUp = 0
	session.BytesDown = 0
	session.StartTime = time.Now().UTC()
	session.UpdateTime = time.Now().UTC()
	session.SessionKey = sessionId

	db := database.Instance()
	db.Save(&session)

	// verified email only if user is authenticated
	if user.AccountType == "email" {
		user.EmailVerified = true
	}
	db.Save(&user)

	return 1
}

func stopSession(sessionId string, unitMacAddress string, bytesDown string, bytesUp string, duration string) int {
	unit := utils.GetUnitByMacAddress(unitMacAddress)
	if unit.Id <= 0 {
		return 0
	}
	session := utils.GetSessionByKeyAndUnitId(sessionId, unit.Id)
	if session.Id <= 0 {
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

	sessionHistory := models.SessionHistory{
		SessionId:   session.Id,
		UnitId:      session.UnitId,
		UnitMac:     session.UnitMac,
		HotspotId:   session.HotspotId,
		HotspotDesc: session.HotspotDesc,
		DeviceId:    session.DeviceId,
		DeviceMAC:   session.DeviceMAC,
		UserId:      session.UserId,
		Username:    session.Username,
		BytesUp:     session.BytesUp,
		BytesDown:   session.BytesDown,
		Duration:    session.Duration,
		AuthTime:    session.AuthTime,
		StartTime:   session.StartTime,
		UpdateTime:  session.UpdateTime,
		StopTime:    session.StopTime,
		SessionKey:  session.SessionKey,
	}

	db := database.Instance()

	// save to session_histories table
	db.Save(&sessionHistory)

	// delete from sessions table
	db.Delete(&session)

	return 1
}

func updateSession(sessionId string, unitMacAddress string, bytesDown string, bytesUp string, duration string) int {
	unit := utils.GetUnitByMacAddress(unitMacAddress)
	if unit.Id <= 0 {
		return 0
	}
	session := utils.GetSessionByKeyAndUnitId(sessionId, unit.Id)
	if session.Id <= 0 {
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

	db := database.Instance()
	db.Save(&session)
	return 1
}

func accountingOn(ap string) int {

	var sessions []models.Session

	unit := utils.GetUnitByMacAddress(ap)
	if unit.Id <= 0 {
		return 0
	}

	db := database.Instance()
	db.Where("unit_id = ? and stop_time = 0", unit.Id).Find(&sessions)

	for _, s := range sessions {
		sessionHistory := models.SessionHistory{
			SessionId:   s.Id,
			UnitId:      s.UnitId,
			UnitMac:     s.UnitMac,
			HotspotId:   s.HotspotId,
			HotspotDesc: s.HotspotDesc,
			DeviceId:    s.DeviceId,
			DeviceMAC:   s.DeviceMAC,
			UserId:      s.UserId,
			Username:    s.Username,
			BytesUp:     s.BytesUp,
			BytesDown:   s.BytesDown,
			Duration:    s.Duration,
			AuthTime:    s.AuthTime,
			StartTime:   s.StartTime,
			UpdateTime:  s.UpdateTime,
			StopTime:    time.Now().UTC(),
			SessionKey:  s.SessionKey,
		}
		// save to session_histories table
		db.Save(&sessionHistory)

		// delete from sessions table
		db.Delete(&s)
	}

	return 1
}

func Counters(c *gin.Context, parameters url.Values) {
	status := parameters.Get("status")

	var username string

	switch status {
	case "start":
		if strings.Compare(c.Query("user"), "temporary") != 0 {
			if strings.Compare(c.Query("user"), c.Query("mac")) == 0 {
				//autologin
				_, users := utils.GetUsersByMacAddressAndunitMacAddress(c.Query("mac"), c.Query("ap"))
				user := utils.FindAutoLoginUser(users)
				username = user.Username
			} else {
				username = c.Query("user")
			}
			Ack(c, startSession(username, c.Query("mac"), c.Query("ip"), c.Query("sessionid"), c.Query("nasid"), c.Query("ap")))
		} else {
			Ack(c, 1)
		}
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
