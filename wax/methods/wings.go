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
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/wax/utils"
)

func GetWingsPrefs(c *gin.Context) {
	var hotspotIntegrations []models.HotspotIntegration
	var integrations []models.Integration

	uuid := c.Query("uuid")

	// extract unit info
	unit := utils.GetUnitByUuid(uuid)

	// extract hotspot info
	hotspot := utils.GetHotspotById(unit.HotspotId)

	// extract account info
	account := utils.GetAccountByAccountId(hotspot.AccountId)

	// get hotspot preferences
	prefs := utils.GetHotspotPreferences(hotspot.Id)

	var wingsPrefs models.WingsPrefs
	wingsPrefs.HotspotId = hotspot.Id
	wingsPrefs.HotspotName = hotspot.Name

	prefsMap := make(map[string]string)
	for i := 0; i < len(prefs); i++ {
		prefsMap[prefs[i].Key] = prefs[i].Value
	}
	wingsPrefs.Preferences = prefsMap

	// get social ids
	wingsPrefs.Socials.FacebookClientId = configuration.Config.AuthSocial.Facebook.ClientId
	wingsPrefs.Socials.LinkedInClientId = configuration.Config.AuthSocial.LinkedIn.ClientId
	wingsPrefs.Socials.InstagramClientId = configuration.Config.AuthSocial.Instagram.ClientId

	// default disclaimers
	terms := configuration.Config.Disclaimers.TermsOfUse
	marketings := configuration.Config.Disclaimers.MarketingUse

	var disclaimers []models.Disclaimer
	db := database.Instance()

	// custom disclaimers
	db.Select("disclaimers.*").Where("disclaimers_hotspots.hotspot_id = ?", hotspot.Id).Joins("JOIN disclaimers_hotspots on disclaimers_hotspots.disclaimer_id = disclaimers.id").Find(&disclaimers)

	for _, disclaimer := range disclaimers {
		if disclaimer.Type == "privacy" {
			marketings = disclaimer.Body
		} else if disclaimer.Type == "tos" {
			terms = disclaimer.Body
		}
	}

	// get integration privacy text
	db.Where("hotspot_id in (?)", hotspot.Id).Find(&hotspotIntegrations)

	for _, hotspotIntegration := range hotspotIntegrations {
		db.Where("id = ?", hotspotIntegration.IntegrationId).Find(&integrations)

		for _, integration := range integrations {
			hotspot.IntegrationTerms += integration.Privacy + " "
		}
	}

	var termsMessage bytes.Buffer
	var marketingMessage bytes.Buffer

	t := template.Must(template.New("terms").Parse(terms))
	m := template.Must(template.New("marketings").Parse(marketings))

	errT := t.Execute(&termsMessage, &account)
	if errT != nil {
		fmt.Println(errT)
	}
	errM := m.Execute(&marketingMessage, &hotspot)
	if errM != nil {
		fmt.Println(errM)
	}

	terms = termsMessage.String()
	marketings = marketingMessage.String()

	wingsPrefs.Disclaimers.TermsOfUse = terms
	wingsPrefs.Disclaimers.MarketingUse = marketings

	// integrations
	wingsPrefs.Integrations = utils.GetHotspotIntegrations(hotspot.Id)

	c.JSON(http.StatusOK, wingsPrefs)
}

func AdditionalInfo(c *gin.Context) {
	userId := c.Param("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		userIdInt = 0
	}

	var json models.UserAdditional
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()

	// update fields
	user := utils.GetUserById(userIdInt)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	user.Reason = json.Reason
	user.Country = json.Country
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
