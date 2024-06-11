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
 * author: Edoard Spadoni <edoardo.spadoni@nethesis.it>
 */

package methods

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/wax/utils"
)

func GetPrivacies(c *gin.Context) {
	hotspotUuid := c.Param("hotspot_uuid")
	hotspot := utils.GetHotspotByUuid(hotspotUuid)
	account := utils.GetAccountByAccountId(hotspot.AccountId)

	if hotspot.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"id": hotspot.Id, "status": "hotspot not found"})
		return
	}

	var hotspotIntegrations []models.HotspotIntegration
	var integrations []models.Integration

	// default disclaimers
	defaultPrivacyDisclaimer := models.Disclaimer{Id: 0, Type: "privacy", Title: "Default", Body: configuration.Config.Disclaimers.MarketingUse}
	defaultTosDisclaimer := models.Disclaimer{Id: 0, Type: "tos", Title: "Default", Body: configuration.Config.Disclaimers.TermsOfUse}

	var privacyDisclaimers []models.Disclaimer
	var tosDisclaimers []models.Disclaimer

	privacyDisclaimers = append(privacyDisclaimers, defaultPrivacyDisclaimer)
	tosDisclaimers = append(tosDisclaimers, defaultTosDisclaimer)

	var selectedDisclaimers []models.Disclaimer

	db := database.Instance()
	db.Select("disclaimers.*").Where("disclaimers_hotspots.hotspot_id = ?", hotspot.Id).Joins("JOIN disclaimers_hotspots on disclaimers_hotspots.disclaimer_id = disclaimers.id").Find(&selectedDisclaimers)

	var selectedPrivacyDisclaimerId = 0
	var selectedTosDisclaimerId = 0

	for _, disclaimer := range selectedDisclaimers {
		if disclaimer.Type == "privacy" {
			selectedPrivacyDisclaimerId = disclaimer.Id
		} else if disclaimer.Type == "tos" {
			selectedTosDisclaimerId = disclaimer.Id
		}
	}

	var disclaimers []models.Disclaimer

	// all custom disclaimers
	db.Select("disclaimers.*").Joins("JOIN disclaimers_accounts on disclaimers_accounts.disclaimer_id = disclaimers.id").Where("disclaimers_accounts.account_id = ?", hotspot.AccountId).Find(&disclaimers)

	for _, disclaimer := range disclaimers {
		if disclaimer.Type == "privacy" {
			privacyDisclaimers = append(privacyDisclaimers, disclaimer)
		} else if disclaimer.Type == "tos" {
			tosDisclaimers = append(tosDisclaimers, disclaimer)
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

	for idx := range privacyDisclaimers {
		var marketingMessage bytes.Buffer
		m := template.Must(template.New("marketings").Parse(privacyDisclaimers[idx].Body))

		errM := m.Execute(&marketingMessage, &account)
		if errM != nil {
			fmt.Println(errM)
		}
		privacyDisclaimers[idx].Body = marketingMessage.String()
	}

	for idx := range tosDisclaimers {
		var termsMessage bytes.Buffer
		t := template.Must(template.New("terms").Parse(tosDisclaimers[idx].Body))

		errT := t.Execute(&termsMessage, &hotspot)
		if errT != nil {
			fmt.Println(errT)
		}
		tosDisclaimers[idx].Body = termsMessage.String()
	}

	c.JSON(http.StatusOK, gin.H{"privacy_disclaimers": privacyDisclaimers, "tos_disclaimers": tosDisclaimers, "privacy_selected_id": selectedPrivacyDisclaimerId, "tos_selected_id": selectedTosDisclaimerId})
}
