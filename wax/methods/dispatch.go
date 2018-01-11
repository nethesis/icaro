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
 * author: Giacomo Sanchietti <giacomo.sanchietti@nethesis.it>
 */

package methods

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"sun-api/database"
	"sun-api/models"
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

func isValidSecret(c *gin.Context, hotspotUnitMac string, md string) bool {
	var unit models.Unit
	var stripper = regexp.MustCompile(`&md=[^&=]+$`)
	var strippedQuery = stripper.ReplaceAllString(c.Request.URL.RawQuery, "")
	var uri = fmt.Sprintf("%s%s?%s", "http://", c.Request.Host, strippedQuery)

	db := database.Database()
	db.Where("mac_address = ?", hotspotUnitMac).First(&unit)
	db.Close()

	h := md5.New()
	io.WriteString(h, uri)
	io.WriteString(h, unit.Secret)
	var check = strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))

	return (check == md)
}

func Dispatch(c *gin.Context) {
	stage := c.Query("stage")
	hotspotName := c.Query("nasid")
	hotspotUnitMac := c.Query("ap")
	md := c.Query("md")

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

	if !isValidSecret(c, hotspotUnitMac, md) {
		Reply(c, http.StatusForbidden, "Invalid unit secret")
		return
	}

	switch stage {
	case "login":
		hotspotId := c.Query("nasid")
		user := c.Query("user")
		password := c.Query("chap_pass")
		challenge := c.Query("chap_chal")
		Login(c, hotspotId, user, password, challenge)

	case "counters":
		parameters := c.Request.URL.Query()
		delete(parameters, "stage")

		// Pass all parameters except 'stage' key
		Counters(c, parameters)

	case "register":
		c.String(http.StatusNotImplemented, "Not implemented: %s", stage)
	default:
		c.String(http.StatusNotFound, "Invalid stage: '%s'", stage)
	}

}
