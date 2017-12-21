/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Icaro project.
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
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"sun-api/configuration"
	"sun-api/models"
	"sun-api/utils"
)

func respondWithError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"message": message})
	c.Abort()
}

func Authorization(role string, route models.Route) bool {

	// extract authorizations from configs
	auths := reflect.ValueOf(configuration.Config.Authorizations)

	// loop and check current route with authorized routes
	for i := 0; i < auths.NumField(); i++ {
		roleConfig := strings.ToLower(auths.Type().Field(i).Name)
		routesConfig := auths.Field(i).Interface().([]models.Route)

		for _, routeConfig := range routesConfig {
			if roleConfig == role {
				matched, _ := regexp.MatchString(routeConfig.Endpoint, route.Endpoint)

				if routeConfig.Verb == route.Verb && matched {
					return false
				}
			}
		}
	}

	return true
}

func AAWall(c *gin.Context) {
	var accessToken models.AccessToken
	token := c.GetHeader("Token")

	// check authentication
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
		if accessToken.Expires.Before(time.Now().UTC()) {
			respondWithError(http.StatusUnauthorized, "API token is expired", c)
			return
		}

		// check authorization
		route := models.Route{
			Verb:     c.Request.Method,
			Endpoint: c.Request.URL.Path,
		}

		authorized := Authorization(accessToken.Role, route)
		if !authorized {
			respondWithError(http.StatusForbidden, "Unauthorized action", c)
			return
		}

		// Refresh token and go ahead
		utils.RefreshToken(accessToken.Token)
		c.Set("token", accessToken)
		c.Next()
	}
}
