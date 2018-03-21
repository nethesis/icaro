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

	"github.com/appleboy/gofight"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
)


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
	f, r := startupEnv()

	f.POST("/api/login").
		SetJSON(gofight.D{
			"username" : "admin",
			"password" : "admin",
		}).
		Run(r, func(f gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, f.Code)
		})
}

