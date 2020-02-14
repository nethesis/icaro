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

	var hotspotIntegration models.HotspotIntegration
	var integration models.Integration

	if hotspot.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"id": hotspot.Id, "status": "hotspot not found"})
		return
	}

	terms := configuration.Config.Disclaimers.TermsOfUse
	marketings := configuration.Config.Disclaimers.MarketingUse

	var termsMessage bytes.Buffer
	var marketingMessage bytes.Buffer

	t := template.Must(template.New("terms").Parse(terms))
	m := template.Must(template.New("marketings").Parse(marketings))

	errT := t.Execute(&termsMessage, &hotspot)
	if errT != nil {
		fmt.Println(errT)
	}
	errM := m.Execute(&marketingMessage, &hotspot)
	if errM != nil {
		fmt.Println(errM)
	}

	terms = termsMessage.String()
	marketings = marketingMessage.String()

	// get integration privacy text
	db := database.Instance()
	db.Where("hotspot_id in (?)", hotspot.Id).First(&hotspotIntegration)

	if hotspotIntegration.Id != 0 {
		db.Where("id = ?", hotspotIntegration.IntegrationId).First(&integration)

		if integration.Id != 0 {
			terms += "\n" + integration.Privacy
		}
	}

	c.JSON(http.StatusOK, gin.H{"terms": terms, "marketings": marketings})
}
