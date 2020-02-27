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
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/wax/utils"
)

func Temporary(c *gin.Context, parameters url.Values) {
	username := parameters.Get("user")
	deviceMacAddress := parameters.Get("mac")
	sessionId := parameters.Get("sessionid")
	unitMacAddress := parameters.Get("ap")
	status := parameters.Get("status")

	var user models.User

	// check if unit exists
	unit := utils.GetUnitByMacAddress(unitMacAddress)
	if unit.Id <= 0 {
		c.String(http.StatusForbidden, "unit not found")
		return
	}

	if len(username) > 0 {
		// check if user exists
		user = utils.GetUserByUsernameAndHotspot(username, unit.HotspotId)
		if user.Id <= 0 {
			c.String(http.StatusForbidden, "user not found")
			return
		}
	}

	// check if the user already has a temporary session
	tempCheck := utils.CheckTempUserSession(user.Id, deviceMacAddress, sessionId)
	if !tempCheck {
		c.String(http.StatusForbidden, "available temporary sessions exceed")
		return
	}

	seconds := utils.GetHotspotPreferencesByKey(unit.HotspotId, "temp_session_duration")

	if status == "new-json" {

		c.JSON(http.StatusOK, gin.H{"sessiontimeout": seconds.Value})

	} else {
		c.String(http.StatusOK, seconds.Value)
	}

}
