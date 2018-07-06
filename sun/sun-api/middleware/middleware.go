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

package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/logpacker/PayPal-Go-SDK"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func respondWithError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"message": message})
	c.Abort()
}

func Authorization(role string, route models.Route) bool {

	// extract authorizations from configs
	auths := reflect.ValueOf(configuration.Config.RouteBlocked)

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

	// check if token exists
	if token == "" {
		respondWithError(http.StatusUnauthorized, "API token required", c)
		return
	}

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

	// refresh token and go ahead
	utils.RefreshToken(accessToken.Token)
	c.Set("token", accessToken)
	c.Next()

}

func PaymentCheck(paymentID string, planCode string, uuid string) bool {
	var apiBase string
	if configuration.Config.PayPal.Sandbox {
		apiBase = paypalsdk.APIBaseSandBox
	} else {
		apiBase = paypalsdk.APIBaseLive
	}
	c, errSDK := paypalsdk.NewClient(configuration.Config.PayPal.ClientID, configuration.Config.PayPal.ClientSecret, apiBase)
	if errSDK != nil {
		fmt.Println(errSDK.Error())
	}
	_, err := c.GetAccessToken()

	payment, err := c.GetPayment(paymentID)
	if err != nil {
		fmt.Println(err.Error())
	}

	SavePaymentDetails(paymentID, uuid)

	if payment.State == "approved" {
		if payment.Transactions[0].ItemList.Items[0].Name == planCode && payment.Transactions[0].ItemList.Items[0].SKU == uuid {
			return true
		}
		return false
	}
	return false
}

func SavePaymentDetails(paymentID string, accountUUID string) {
	var account models.Account

	db := database.Instance()
	db.Set("gorm:auto_preload", false).Where("uuid = ?", accountUUID).First(&account)
	payment := models.Payment{Payment: paymentID, AccountId: account.Id, Created: time.Now().UTC()}
	db.Create(&payment)
}
