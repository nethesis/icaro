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
 * author: Giacomo Sanchietti <giacomo.sanchietti@nethesis.it>
 */

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/appleboy/gofight"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/wax/utils"
)

func startupEnv(endpoint string, query string) (*gofight.RequestConfig, *gin.Engine, string) {
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

	// calculate URI with correct MD
	uri := calculateUri(endpoint, query)

	return fight, router, uri
}

func destroyEnv() {
}

func calculateUri(endpoint string, query string) string {
	md := calculateMd(endpoint + "?" + query)
	testUuid := "1234-uuid-aaaa"
	unit := utils.GetUnitByUuid(testUuid)
	digest := utils.CalcUnitDigest(unit)
	return endpoint + "?" + query + "&md=" + md + "&digest=" + digest + "&uuid=" + testUuid + "&sessionid=3"
}

func calculateMd(query string) string {
	uri := "http://" + query
	h := md5.New()
	io.WriteString(h, uri)
	io.WriteString(h, "secret") // secret for AP HSTest
	check := strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))

	return check
}

func TestMain(m *testing.M) {
	// read and init configuration
	c := "config.json"
	configuration.Init(&c)

	// run tester
	os.Exit(m.Run())
}

/** Stage **/

func TestNoStage(t *testing.T) {
	f, r, uri := startupEnv("/wax/aaa", "nasid=HSTest&ap=00-00-00-00-00-00")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "No stage provided", f.Body.String())
			assert.Equal(t, http.StatusBadRequest, f.Code)
		})
}

func TestRegisterStage(t *testing.T) {
	f, r, uri := startupEnv("/wax/aaa", "stage=register&nasid=HSTest&ap=00-00-00-00-00-00")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Not implemented: register", f.Body.String())
			assert.Equal(t, http.StatusNotImplemented, f.Code)
		})
}

func TestLoginStage(t *testing.T) {
	f, r, uri := startupEnv("/wax/aaa", "stage=login&service=login&nasid=HSTest&ap=00-00-00-00-00-00&user=firstuser&chap_pass=9221f7e65679b0f49435707286920228&chap_chal=challange&sessionid=1234")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, f.Code)
		})
}

func TestFailLoginForExpiredReseller(t *testing.T) {
	var account models.Account
	var subscription models.Subscription
	f, r, uri := startupEnv("/wax/aaa", "stage=login&service=login&nasid=HSTest&ap=00-00-00-00-00-00&user=firstuser&chap_pass=9221f7e65679b0f49435707286920228&chap_chal=challange&sessionid=1234")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, f.Code)
		})

	db := database.Instance()
	db.Where("username = ?", "firstuser").First(&account)
	db.Where("account_id = ?", account.Id).First(&subscription)
	oldDate := subscription.ValidUntil
	subscription.ValidUntil = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	db.Save(&subscription)

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusForbidden, f.Code)
		})

	subscription.ValidUntil = oldDate
	db.Save(&subscription)

}

func TestInvalidStage(t *testing.T) {
	f, r, uri := startupEnv("/wax/aaa", "stage=BAD&nasid=HSTest&ap=00-00-00-00-00-00")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Invalid stage: 'BAD'", f.Body.String())
			assert.Equal(t, http.StatusNotFound, f.Code)
		})
}

/** Counters **/
func TestCountersInvalidUnit(t *testing.T) {
	f, r, uri := startupEnv("/wax/aaa", "stage=counters&nasid=HSTest&ap=xx-00-00-00-00-00&status=start")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, f.Code)
		})
}

func TestCountersStart(t *testing.T) {
	f, r, uri := startupEnv("/wax/aaa", "stage=counters&nasid=HSTest&ap=00-00-00-00-00-00&status=start&user=firstuser&mac=11-11-11-11-11&ip=10.1.0.3&sessionid=151318020800000001")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Ack: 1", f.Body.String())
			assert.Equal(t, http.StatusOK, f.Code)
		})
}

func TestCountersUpdate(t *testing.T) {
	f, r, uri := startupEnv("/wax/aaa", "stage=counters&nasid=HSTest&ap=00-00-00-00-00-00&status=update&user=firstuser&mac=11-11-11-11-11&ip=10.1.0.3&sessionid=151318020800000001&duration=222&bytes_down=80000&pkts_down=3000&bytes_up=189235&pkts_up=2071")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Ack: 1", f.Body.String())
			assert.Equal(t, http.StatusOK, f.Code)
		})
}

func TestCountersStop(t *testing.T) {
	f, r, uri := startupEnv("/wax/aaa", "stage=counters&nasid=HSTest&ap=00-00-00-00-00-00&status=stop&user=firstuser&mac=11-11-11-11-11&ip=10.1.0.3&sessionid=151318020800000001&duration=16036&bytes_down=727516&pkts_down=1906&bytes_up=189235&pkts_up=2071")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Ack: 1", f.Body.String())
			assert.Equal(t, http.StatusOK, f.Code)
		})
}

func TestCountersInvalid(t *testing.T) {
	f, r, uri := startupEnv("/wax/aaa", "stage=counters&nasid=HSTest&ap=00-00-00-00-00-00")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "No status provided", f.Body.String())
			assert.Equal(t, http.StatusBadRequest, f.Code)
		})

	uri = calculateUri("/wax/aaa", "stage=counters&nasid=HSTest&ap=00-00-00-00-00-00&status=invalid")
	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Invalid status: 'invalid'", f.Body.String())
			assert.Equal(t, http.StatusNotImplemented, f.Code)

		})
}

func TestFindAutoLoginUser(t *testing.T) {
	var user models.User
	time_now := time.Now().UTC()

	var test_users []models.User
	var test_users_empty = []models.User{}
	var test_users_autologin_false = []models.User{
		models.User{
			Id:                   11,
			HotspotId:            10,
			Name:                 "+555555555555555",
			Username:             "+555555555555555",
			Password:             "123456",
			Email:                "",
			EmailVerified:        false,
			AccountType:          "sms",
			MarketingAuth:        false,
			KbpsDown:             0,
			KbpsUp:               0,
			MaxNavigationTraffic: 0,
			MaxNavigationTime:    0,
			AutoLogin:            false,
			ValidFrom:            time_now,
			ValidUntil:           time_now.Add((time.Hour * 24) * 30),
			Created:              time_now,
		},
	}
	var test_users_autologin_true = []models.User{
		models.User{
			Id:                   12,
			HotspotId:            10,
			Name:                 "+555555555555555",
			Username:             "+555555555555555",
			Password:             "123456",
			Email:                "",
			EmailVerified:        false,
			AccountType:          "sms",
			MarketingAuth:        false,
			KbpsDown:             0,
			KbpsUp:               0,
			MaxNavigationTraffic: 0,
			MaxNavigationTime:    0,
			AutoLogin:            true,
			ValidFrom:            time_now,
			ValidUntil:           time_now.Add((time.Hour * 24) * 30),
			Created:              time_now,
		},
	}

	var test_users_autologin_true_new = []models.User{
		models.User{
			Id:                   13,
			HotspotId:            10,
			Name:                 "+555555555555555",
			Username:             "+555555555555555",
			Password:             "123456",
			Email:                "",
			EmailVerified:        false,
			AccountType:          "sms",
			MarketingAuth:        false,
			KbpsDown:             0,
			KbpsUp:               0,
			MaxNavigationTraffic: 0,
			MaxNavigationTime:    0,
			AutoLogin:            true,
			ValidFrom:            time_now.Add(time.Hour * 24),
			ValidUntil:           time_now.Add((time.Hour * 24) * 30),
			Created:              time_now.Add(time.Hour * 24),
		},
	}

	user = utils.FindAutoLoginUser(test_users_empty)
	if !assert.Equal(t, user.Id, 0) {
		t.Fatal("expected user.Id = 0, got:", user.Id)
	}

	user = utils.FindAutoLoginUser(test_users_autologin_false)
	if !assert.Equal(t, user.Id, 0) {
		t.Fatal("expected user.Id = 0, got:", user.Id)
	}

	user = utils.FindAutoLoginUser(test_users_autologin_true)
	if !assert.Equal(t, user.Id, 12) {
		t.Fatal("expected user.Id = 12, got:", user.Id)
	}

	test_users = append(test_users, test_users_autologin_false[0],
		test_users_autologin_true[0],
		test_users_autologin_true_new[0])

	user = utils.FindAutoLoginUser(test_users)
	if !assert.Equal(t, user.Id, 13) {
		t.Fatal("expected user.Id = 13, got:", user.Id)
	}

}
