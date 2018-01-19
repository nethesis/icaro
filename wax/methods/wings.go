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
	"wax/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"sun-api/models"
)

func GetWingsPrefs(c *gin.Context) {
	uuid := c.Query("uuid")

	// extract unit info
	unit := utils.GetUnitByUuid(uuid)

	// extract hotspot info
	hotspot := utils.GetHotspotById(unit.HotspotId)

	// get hotspot preferences
	prefs := utils.GetHotspotPreferences(hotspot.Id)

	var wingsPrefs models.WingsPrefs
	wingsPrefs.HotspotId = hotspot.Id
	wingsPrefs.HotspotName = hotspot.Name

	prefsMap := make(map[string]string)
	for i := 0; i < len(prefs); i += 2 {
		prefsMap[prefs[i].Key] = prefs[i].Value
	}
	wingsPrefs.Prefs = prefsMap

	c.JSON(http.StatusOK, wingsPrefs)
}
