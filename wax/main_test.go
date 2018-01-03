/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Icaro project.
 *
 * Icaro is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Icaro is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Icaro.  If not, see COPYING.
 *
 * author: Giacomo Sanchietti <giacomo.sanchietti@nethesis.it>
 */

package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"sun-api/configuration"
)

func startupEnv() (*gofight.RequestConfig, *gin.Engine) {
	// set test mode to test suite
	fight := gofight.New()
	fight.SetDebug(false)

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
	configuration.Init()

	// run tester
	os.Exit(m.Run())
}

/** Stage **/

func TestNoStage(t *testing.T) {
	f, r := startupEnv()

	f.GET("/wax/aaa?nasid=HSTest&ap=00-00-00-00-00-00").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "No stage provided", f.Body.String())
			assert.Equal(t, http.StatusBadRequest, f.Code)
		})
}

func TestRegisterStage(t *testing.T) {
	f, r := startupEnv()

	f.GET("/wax/aaa?stage=register&nasid=HSTest&ap=00-00-00-00-00-00").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Not implemented: register", f.Body.String())
			assert.Equal(t, http.StatusNotImplemented, f.Code)
		})
}

func TestLoginStage(t *testing.T) {
	f, r := startupEnv()

	f.GET("/wax/aaa?stage=login&nasid=HSTest&ap=00-00-00-00-00-00").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			//assert.Equal(t, "login", f.Body.String())
			assert.Equal(t, http.StatusOK, f.Code)
		})
}

func TestInvalidStage(t *testing.T) {
	f, r := startupEnv()

	f.GET("/wax/aaa?stage=BAD&nasid=HSTest&ap=00-00-00-00-00-00").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Invalid stage: 'BAD'", f.Body.String())
			assert.Equal(t, http.StatusNotFound, f.Code)
		})
}

/** Counters **/

func TestCountersInvalidHotspot(t *testing.T) {
	f, r := startupEnv()

	f.GET("/wax/aaa?stage=counters&nasid=BAD&ap=00-00-00-00-00-00&status=start").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, f.Code)
		})
}

func TestCountersInvalidUnit(t *testing.T) {
	f, r := startupEnv()

	f.GET("/wax/aaa?stage=counters&nasid=HSTest&ap=xx-00-00-00-00-00&status=start").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, f.Code)
		})
}

func TestCountersStage(t *testing.T) {
	f, r := startupEnv()

	f.GET("/wax/aaa?stage=counters&nasid=HSTest&ap=00-00-00-00-00-00&status=start").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, f.Code)
		})
}

func TestCountersInvalid(t *testing.T) {
	f, r := startupEnv()

	f.GET("/wax/aaa?stage=counters&nasid=HSTest&ap=00-00-00-00-00-00").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "No status provided", f.Body.String())
			assert.Equal(t, http.StatusBadRequest, f.Code)
		})

	f.GET("/wax/aaa?stage=counters&nasid=HSTest&ap=00-00-00-00-00-00&status=invalid").
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Invalid status: 'invalid'", f.Body.String())
			assert.Equal(t, http.StatusNotImplemented, f.Code)

		})
}
