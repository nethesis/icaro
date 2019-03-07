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
	"github.com/nethesis/icaro/ade/ade-api/models"
	"github.com/nethesis/icaro/ade/ade-api/utils"
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

	hotspotPerfs := utils.GetHotspotPrefs(adeToken.HotspotId)

	reviewPage.HotspotName = hotspotPerfs["captive_2_title"]
	reviewPage.HotspotLogo = hotspotPerfs["captive_3_logo"]
	reviewPage.BgColor = hotspotPerfs["captive_7_background"]

	reviewPage.Threshold, _ = strconv.Atoi(hotspotPerfs["marketing_9_threshold"])

	if hotspotPerfs["marketing_10_first_url"] != "" {
		urls = append(urls, hotspotPerfs["marketing_10_first_url"])
	}
	if hotspotPerfs["marketing_11_second_url"] != "" {
		urls = append(urls, hotspotPerfs["marketing_11_second_url"])
	}
	if hotspotPerfs["marketing_12_third_url"] != "" {
		urls = append(urls, hotspotPerfs["marketing_12_third_url"])
	}

	reviewPage.Urls = urls

	c.JSON(http.StatusOK, reviewPage)
}

func PostReviewResult(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
