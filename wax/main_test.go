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

	"github.com/appleboy/gofight"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
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

	// calculate URI with correct MD
	uri := calculateUri(endpoint, query)

	return fight, router, uri
}

func destroyEnv() {
}

func calculateUri(endpoint string, query string) string {
	md := calculateMd(endpoint + "?" + query)
	return endpoint + "?" + query + "&md=" + md
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
	configuration.Init()

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
	f, r, uri := startupEnv("/wax/aaa", "stage=login&nasid=HSTest&ap=00-00-00-00-00-00")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			//assert.Equal(t, "login", f.Body.String())
			assert.Equal(t, http.StatusOK, f.Code)
		})
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

func TestCountersInvalidHotspot(t *testing.T) {
	f, r, uri := startupEnv("/wax/aaa", "stage=counters&nasid=BAD&ap=00-00-00-00-00-00&status=start")

	f.GET(uri).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, f.Code)
		})
}

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
