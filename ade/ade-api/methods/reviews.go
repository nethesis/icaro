/*
 * Copyright (C) 2019 Nethesis S.r.l.
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
 * author: Matteo Valentini <matteo.valentini@nethesis.it>
 */

package methods

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nethesis/icaro/ade/ade-api/utils"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

func GetReviewPage(c *gin.Context) {
	var reviewPage models.ReviewPage
	var urls []string

	token := c.Param("token")
	adeToken := utils.GetAdeTokenFromToken(token)

	if adeToken.Id <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No token found!"})
		return
	}

	if !adeToken.ReviewLeftTime.IsZero() {
		c.JSON(http.StatusForbidden, gin.H{"message": "Token expired"})
		return
	}

	hotspotPrefs := utils.GetHotspotPrefs(adeToken.HotspotId)

	reviewPage.HotspotName = hotspotPrefs["captive_2_title"]
	reviewPage.HotspotLogo = hotspotPrefs["captive_3_logo"]
	reviewPage.BgColor = hotspotPrefs["captive_7_background"]

	reviewPage.Threshold, _ = strconv.Atoi(hotspotPrefs["marketing_9_threshold"])

	if hotspotPrefs["marketing_10_first_url"] != "" {
		urls = append(urls, hotspotPrefs["marketing_10_first_url"])
	}
	if hotspotPrefs["marketing_11_second_url"] != "" {
		urls = append(urls, hotspotPrefs["marketing_11_second_url"])
	}
	if hotspotPrefs["marketing_12_third_url"] != "" {
		urls = append(urls, hotspotPrefs["marketing_12_third_url"])
	}

	reviewPage.Urls = urls

	c.JSON(http.StatusOK, reviewPage)
}

func PostReviewResult(c *gin.Context) {
	var reviewResult models.ReviewResult

	token := c.Param("token")
	adeToken := utils.GetAdeTokenFromToken(token)
	hotspotPrefs := utils.GetHotspotPrefs(adeToken.HotspotId)

	if adeToken.Id <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No token found!"})
		return
	}

	if !adeToken.ReviewLeftTime.IsZero() {
		c.JSON(http.StatusForbidden, gin.H{"message": "Token expired"})
		return
	}

	if err := c.BindJSON(&reviewResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	if reviewResult.Stars > 0 && reviewResult.Stars <= 5 {
		if !utils.SendReviewMessageToOwner(adeToken, reviewResult.Stars, reviewResult.Message, hotspotPrefs["captive_7_background"]) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Feedback submission failed."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
