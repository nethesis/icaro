/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Windmill project.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with NethServer.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"manager-api/models"
	"manager-api/utils"
)

func respondWithError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"message": message})
	c.Abort()
}

func Authentication(c *gin.Context) {
	var accessToken models.AccessToken
	token := c.GetHeader("Token")

	if token == "" {
		respondWithError(http.StatusUnauthorized, "API token required", c)
		return
	} else {
		// check token validation
		accessToken = utils.ExtractToken(token)

		if accessToken.Id == 0 {
			respondWithError(http.StatusUnauthorized, "API token is invalid", c)
			return
		}
		if accessToken.Expires.After(time.Now()) {
			respondWithError(http.StatusUnauthorized, "API token is expired", c)
			return
		}

		// Refresh token and go ahead
		utils.RefreshToken(accessToken.Token)
		c.Set("token", accessToken)
		c.Next()
	}
}
