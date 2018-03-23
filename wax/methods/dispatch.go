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
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/wax/utils"
)

func Reply(c *gin.Context, httpCode int, message string) {
	c.String(httpCode, "Reply-Message: %s", message)
}

func isHotspotUnit(hotspotUnitMac string) bool {
	var unit models.Unit

	db := database.Database()
	db.Where("mac_address = ?", hotspotUnitMac).First(&unit)
	db.Close()

	return (unit.Id != 0)
}

func isHotspot(hotspotName string) bool {
	var hotspot models.Hotspot

	db := database.Database()
	db.Where("name = ?", hotspotName).First(&hotspot)
	db.Close()

	return (hotspot.Id != 0)
}

func isAutoLoginEnabled(hotspotId int) bool {
	var hotspot_pref models.HotspotPreference
	isEnabled := false

	db := database.Database()
	db.Where("hotspot_id = ? and `key` = ?", hotspotId, "auto_login").First(&hotspot_pref)
	db.Close()

	if hotspot_pref.Id != 0 {
		isEnabled, _ = strconv.ParseBool(hotspot_pref.Value)
	}

	return isEnabled
}

func Dispatch(c *gin.Context) {
	stage := c.Query("stage")
	hotspotName := c.Query("nasid")
	hotspotUnitMac := c.Query("ap")

	if stage == "" {
		c.String(http.StatusBadRequest, "No stage provided")
		return
	}

	if !isHotspot(hotspotName) {
		Reply(c, http.StatusNotFound, "Hotspot not found")
		return
	}

	if !isHotspotUnit(hotspotUnitMac) {
		Reply(c, http.StatusNotFound, "Hotspot unit not found")
		return
	}

	switch stage {
	case "login":
		service := c.Query("service")
		switch service {
		case "framed":
			unit := utils.GetUnitByMacAddress(c.Query("ap"))
			if isAutoLoginEnabled(unit.HotspotId) {
				unitMacAddress := c.Query("ap")
				user := c.Query("user")
				userMac := c.Query("mac")
				sessionId := c.Query("sessionid")
				autoLogin(c, unitMacAddress, user, userMac, sessionId)
			} else {
				c.String(http.StatusForbidden, "Autologin disabled")
			}

		case "login":
			unitMacAddress := c.Query("ap")
			user := c.Query("user")
			chapPass := c.Query("chap_pass")
			chapChal := c.Query("chap_chal")
			sessionId := c.Query("sessionid")
			Login(c, unitMacAddress, user, chapPass, chapChal, sessionId)

		default:
			c.String(http.StatusNotFound, "Invalid login service: '%s'", service)
		}

	case "counters":
		parameters := c.Request.URL.Query()
		delete(parameters, "stage")

		// Pass all parameters except 'stage' key
		Counters(c, parameters)

	case "register":
		c.String(http.StatusNotImplemented, "Not implemented: %s", stage)

	case "temporary":
		parameters := c.Request.URL.Query()

		Temporary(c, parameters)
	default:
		c.String(http.StatusNotFound, "Invalid stage: '%s'", stage)
	}

}
