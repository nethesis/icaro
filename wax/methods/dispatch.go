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

	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

func Reply(c *gin.Context, httpCode int, message string) {
	c.String(httpCode, "Reply-Message: %s", message)
}

func isHotspotUnit(hotspotUnitMac string) bool {
	var unit models.Unit

	db := database.Instance()
	db.Where("mac_address = ?", hotspotUnitMac).First(&unit)

	return (unit.Id != 0)
}

func isHotspot(hotspotName string) bool {
	var hotspot models.Hotspot

	db := database.Instance()
	db.Where("name = ?", hotspotName).First(&hotspot)

	return (hotspot.Id != 0)
}

func Dispatch(c *gin.Context) {
	stage := c.Query("stage")
	hotspotUnitMac := c.Query("ap")

	if stage == "" {
		c.String(http.StatusBadRequest, "No stage provided")
		return
	}

	if !isHotspotUnit(hotspotUnitMac) {
		Reply(c, http.StatusNotFound, "Hotspot unit not found")
		return
	}

	switch stage {
	case "login":
		Logins(c)

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
