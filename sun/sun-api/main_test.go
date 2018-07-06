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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/appleboy/gofight"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

type LoginResponse struct {
	AccountType  string              `json:"account_type"`
	Status       string              `json:"status"`
	Token        string              `json:"token"`
	Expires      string              `json:"expires"`
	AccountID    int                 `json:"id"`
	Subscription models.Subscription `json:"subscription"`
}

type CreationResponse struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

func startupEnv() (*gofight.RequestConfig, *gin.Engine) {
	// set test mode to test suite
	fight := gofight.New()
	fight.SetDebug(false)

	// set test mode to gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// define API
	DefineAPI(router)

	// init database
	database.Init()

	return fight, router
}

func getAdminToken(f *gofight.RequestConfig, r *gin.Engine) string {
	return getToken(f, r, "admin", "admin")
}

func getResellerToken(f *gofight.RequestConfig, r *gin.Engine) string {
	return getToken(f, r, "firstuser", "password")
}

func getToken(f *gofight.RequestConfig, r *gin.Engine, user string, password string) string {
	var lr LoginResponse
	var token string

	f.POST("/api/login").
		SetJSON(gofight.D{
			"username": user,
			"password": password,
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &lr)
			if err == nil {
				token = lr.Token
			} else {
				fmt.Println(err)
			}
		})

	return token
}

func TestMain(m *testing.M) {
	// read and init configuration
	c := "config.json"
	configuration.Init(&c)

	// run tester
	os.Exit(m.Run())
}

/** Initialization **/

func TestAdminLogin(t *testing.T) {
	var lr LoginResponse
	f, r := startupEnv()

	f.POST("/api/login").
		SetJSON(gofight.D{
			"username": "firstuser",
			"password": "password",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &lr)
			assert.Equal(t, nil, err)
			assert.NotEqual(t, "", lr.Token)
			assert.NotEqual(t, "", lr.Expires)
			assert.Equal(t, "reseller", lr.AccountType)
			assert.Equal(t, "success", lr.Status)
			assert.Equal(t, "application/json; charset=utf-8", f.HeaderMap.Get("Content-Type"))
			assert.Equal(t, http.StatusCreated, f.Code)
		})
}

func TestLogin(t *testing.T) {
	var lr LoginResponse
	f, r := startupEnv()

	f.POST("/api/login").
		SetJSON(gofight.D{
			"username": "firstuser",
			"password": "password",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &lr)
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

/** Hotspots **/

func TestHotspotCreation(t *testing.T) {
	f, r := startupEnv()
	var cr CreationResponse

	token := getResellerToken(f, r)

	f.POST("/api/hotspots").
		SetHeader(gofight.H{
			"Token": token,
		}).
		SetJSON(gofight.D{
			"name":        "MyTestHotspot",
			"description": "Generated from go test",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &cr)
			assert.Equal(t, nil, err)
			assert.NotZero(t, cr.Id)
			assert.Equal(t, "success", cr.Status)
			assert.Equal(t, http.StatusCreated, f.Code)
		})

	// Cleanup
	db := database.Instance()
	db.Delete(models.Hotspot{Id: cr.Id})
}

/** Account **/

func TestResellerAccountCreation(t *testing.T) {
	f, r := startupEnv()
	var cr CreationResponse
	var subPlan models.SubscriptionPlan
	var sub models.Subscription
	var hs models.Hotspot
	var account models.Account
	var asms models.AccountSmsCount

	token := getAdminToken(f, r)

	db := database.Instance()
	db.Where("id = 4").First(&subPlan)
	db.Where("id = 1").First(&hs)

	f.POST("/api/accounts").
		SetHeader(gofight.H{
			"Token": token,
		}).
		SetJSON(gofight.D{
			"uuid":                 fmt.Sprintf("123reseller%d", time.Now().Nanosecond()),
			"type":                 "reseller",
			"name":                 "Test reseller",
			"username":             fmt.Sprintf("reseller%d", time.Now().Nanosecond()),
			"password":             "testpassword",
			"email":                "testreseller@nethserver.org",
			"subscription_plan_id": subPlan.Id,
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &cr)
			assert.Equal(t, nil, err)
			assert.Equal(t, http.StatusCreated, f.Code)

			db.Where("id = ?", cr.Id).First(&account)
			db.Where("account_id = ?", account.Id).Find(&asms)
			db.Where("account_id = ?", account.Id).Find(&sub)

			assert.WithinDuration(t, time.Now().UTC(), sub.Created, 10*time.Second)
			assert.WithinDuration(t, time.Now().UTC(), sub.ValidFrom, 10*time.Second)
			assert.WithinDuration(t, time.Now().UTC().AddDate(0, 0, 3650), sub.ValidUntil, (10 * time.Second))

			assert.Equal(t, asms.SmsMaxCount, subPlan.IncludedSMS)
			assert.Equal(t, 0, asms.SmsCount)
		})

	// Cleanup
	db.Delete(&account)
	db.Delete(&asms)
	db.Delete(&sub)
}

/** Units **/

func TestUnitRegistration(t *testing.T) {
	f, r := startupEnv()
	var cr CreationResponse

	token := getResellerToken(f, r)

	f.POST("/api/units").
		SetHeader(gofight.H{
			"Token": token,
		}).
		SetJSON(gofight.D{
			"mac_address": "43-C2-41-98-85-08",
			"description": "My test unit",
			"name":        "MyTestUnit",
			"uuid":        fmt.Sprintf("%d", time.Now().Nanosecond()),
			"secret":      "mysecret",
			"hotspot_id":  "1",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &cr)
			assert.Equal(t, nil, err)
			assert.Equal(t, http.StatusCreated, f.Code)
		})

	// Cleanup
	db := database.Instance()
	db.Delete(models.Unit{}, "id > 1")
}

func TestUnitRegistrationLimit(t *testing.T) {
	max := 10
	var cr CreationResponse
	f, r := startupEnv()
	token := getResellerToken(f, r)

	db := database.Instance()
	// create max-1 units because the first one already exists
	for i := 0; i < max-1; i++ {
		f.POST("/api/units").
			SetHeader(gofight.H{
				"Token": token,
			}).
			SetJSON(gofight.D{
				"mac_address": "42-C2-41-98-85-08",
				"description": fmt.Sprintf("My test unit %d", i),
				"name":        "MyTestUnit",
				"uuid":        fmt.Sprintf("uuid-test-%d", i),
				"secret":      fmt.Sprintf("mysecret%d", i),
				"hotspot_id":  "1",
			}).
			Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
				err := json.Unmarshal([]byte(f.Body.String()), &cr)
				assert.Equal(t, nil, err)
				assert.Equal(t, http.StatusCreated, f.Code)
			})
	}

	// this one must fail
	f.POST("/api/units").
		SetHeader(gofight.H{
			"Token": token,
		}).
		SetJSON(gofight.D{
			"mac_address": "44-C2-41-98-85-08",
			"description": fmt.Sprintf("My test unit %d", max),
			"uuid":        fmt.Sprintf("uuid-test-%d", max),
			"secret":      fmt.Sprintf("mysecret%d", max),
			"hotspot_id":  "1",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &cr)
			assert.Equal(t, nil, err)
			assert.Equal(t, http.StatusForbidden, f.Code)
		})

	// cleanup
	db.Delete(models.Unit{}, "id > 1")
}

/** Preferences **/

func TestFailingCaptiveConfiguration(t *testing.T) {
	f, r := startupEnv()
	var cr CreationResponse
	var sub models.Subscription
	var hs models.Hotspot
	var account models.Account
	var asms models.AccountSmsCount

	var subPlan models.SubscriptionPlan

	token := getAdminToken(f, r)
	db := database.Instance()
	db.Where("id = 1").First(&subPlan)

	// Create reseller with free plan
	f.POST("/api/accounts").
		SetHeader(gofight.H{
			"Token": token,
		}).
		SetJSON(gofight.D{
			"uuid":                 fmt.Sprintf("123reseller%d", time.Now().Nanosecond()),
			"type":                 "reseller",
			"name":                 "Test reseller",
			"username":             "restestcaptive",
			"password":             "restestcaptive",
			"email":                "testreseller@nethserver.org",
			"subscription_plan_id": subPlan.Id,
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &cr)
			assert.Equal(t, nil, err)
			assert.Equal(t, http.StatusCreated, f.Code)

			db.Where("id = ?", cr.Id).First(&account)
			db.Where("account_id = ?", account.Id).Find(&asms)
			db.Where("account_id = ?", account.Id).Find(&sub)

		})

	// login with created reseller
	token = getToken(f, r, "restestcaptive", "restestcaptive")

	f.POST("/api/hotspots").
		SetHeader(gofight.H{
			"Token": token,
		}).
		SetJSON(gofight.D{
			"name":        "MyTestHotspot",
			"description": "Generated from go test",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &cr)
			assert.Equal(t, nil, err)
			assert.Equal(t, http.StatusCreated, f.Code)

			db.Where("id = ?", cr.Id).First(&hs)

		})

	f.PUT(fmt.Sprintf("/api/preferences/hotspots/%d", hs.Id)).
		SetHeader(gofight.H{
			"Token": token,
		}).
		SetJSON(gofight.D{
			"hotspot_id": hs.Id,
			"key":        "email_login",
			"value":      "true",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &cr)
			assert.Equal(t, nil, err)
			assert.Equal(t, http.StatusOK, f.Code)
		})

	f.PUT(fmt.Sprintf("/api/preferences/hotspots/%d", hs.Id)).
		SetHeader(gofight.H{
			"Token": token,
		}).
		SetJSON(gofight.D{
			"key":   "captive_2_title",
			"value": "t123",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			err := json.Unmarshal([]byte(f.Body.String()), &cr)
			assert.Equal(t, nil, err)
			assert.Equal(t, http.StatusUnauthorized, f.Code)
		})

	// Cleanup
	db.Delete(&account)
	db.Delete(&asms)
	db.Delete(&sub)
	db.Delete(&hs)
}
