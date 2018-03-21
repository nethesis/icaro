/*
 * Copyright (C) 2018 Nethesis S.r.l.
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
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Icaro.  If not, see COPYING.
 *
 * author: Giacomo Sanchietti <giacomo.sanchietti@nethesis.it>
 */

package main

import (
	"net/http"
	"os"
	"testing"
	"encoding/json"

	"github.com/appleboy/gofight"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/models"
)


type LoginResponse struct {
	AccountType string `json:"account_type"`
	Status string `json:"status"`
	Token string `json:"token"`
	Expires string `json:"expires"`
	AccountID      int `json:"id"`
	Subscription models.Subscription `json:"subscription"`
}

func startupEnv() (*gofight.RequestConfig, *gin.Engine) {
	// set test mode to test suite
	fight := gofight.New()
	fight.SetDebug(true)

	// set test mode to gin
	gin.SetMode(gin.TestMode)
	router := gin.New()


	// define API
	DefineAPI(router)

	return fight, router
}

func destroyEnv() {
}


func TestMain(m *testing.M) {
	// read and init configuration
	c := "config.json"
	configuration.Init(&c)

	// run tester
	os.Exit(m.Run())
}

/** Login **/

func TestLogin(t *testing.T) {
	var lr LoginResponse
	f, r := startupEnv()

	f.POST("/api/login").
		SetJSON(gofight.D{
			"username" : "firstuser",
			"password" : "password",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &lr);
			assert.Equal(t, nil, err)
			assert.NotEqual(t, "", lr.Token)
			assert.NotEqual(t, "", lr.Expires)
			assert.Equal(t, 2, lr.AccountID)
			assert.Equal(t, "reseller", lr.AccountType)
			assert.Equal(t, "success", lr.Status)
			assert.Equal(t, 1000, lr.Subscription.SubscriptionPlan.IncludedSMS)
			assert.Equal(t, false, lr.Subscription.SubscriptionPlan.SocialAnalytics)
			assert.NotNil(t, lr.Subscription.ValidFrom)
			assert.Equal(t, "application/json; charset=utf-8", f.HeaderMap.Get("Content-Type"))
			assert.Equal(t, http.StatusCreated, f.Code)
		})
}


