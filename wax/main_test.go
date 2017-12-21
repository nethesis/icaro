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
	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"sun-api/configuration"
	"os"
)

var e = Init(true);

func startupEnv() *gofight.RequestConfig {
        r := gofight.New()
	r.SetDebug(false)
	return r
}

func destroyEnv() {
}

func TestMain(m *testing.M) {
	configuration.Init()


	os.Exit(m.Run())
}

/** Stage **/

func TestNoStage(t *testing.T) {
	r := startupEnv()

        r.GET("/?nasid=HSTest&ap=00-00-00-00-00-00").
                Run(Init(true), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "No stage provided", r.Body.String())
                        assert.Equal(t, http.StatusBadRequest, r.Code)
                })
}

func TestRegisterStage(t *testing.T) {
	r := startupEnv()

        r.GET("/?stage=register&nasid=HSTest&ap=00-00-00-00-00-00").
                Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
                        assert.Equal(t, "Not implemented: register", r.Body.String())
                        assert.Equal(t, http.StatusNotImplemented, r.Code)
                })
}


func TestLoginStage(t *testing.T) {
	r := startupEnv()

        r.GET("/?stage=login&nasid=HSTest&ap=00-00-00-00-00-00").
                Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
         //               assert.Equal(t, "login", r.Body.String())
                        assert.Equal(t, http.StatusOK, r.Code)
                })
}

func TestInvalidStage(t *testing.T) {
	r := startupEnv()

        r.GET("/?stage=BAD&nasid=HSTest&ap=00-00-00-00-00-00").
                Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Invalid stage: 'BAD'", r.Body.String())
                        assert.Equal(t, http.StatusNotFound, r.Code)
                })
}

/** Counters **/

func TestCountersInvalidHotspot(t *testing.T) {
	r := startupEnv()

        r.GET("/?stage=counters&nasid=BAD&ap=00-00-00-00-00-00&status=start").
                Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
                        assert.Equal(t, http.StatusNotFound, r.Code)
                })
}

func TestCountersInvalidUnit(t *testing.T) {
        r := startupEnv()

        r.GET("/?stage=counters&nasid=HSTest&ap=xx-00-00-00-00-00&status=start").
                Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
                        assert.Equal(t, http.StatusNotFound, r.Code)
                })
}


func TestCountersStage(t *testing.T) {
	r := startupEnv()

        r.GET("/?stage=counters&nasid=HSTest&ap=00-00-00-00-00-00&status=start").
                Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
                        assert.Equal(t, http.StatusOK, r.Code)
                })
}

func TestCountersInvalid(t *testing.T) {
	r := startupEnv()

        r.GET("/?stage=counters&nasid=HSTest&ap=00-00-00-00-00-00").
                Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "No status provided", r.Body.String())
                        assert.Equal(t, http.StatusBadRequest, r.Code)
                })

	r.GET("/?stage=counters&nasid=HSTest&ap=00-00-00-00-00-00&status=invalid").
                Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Invalid status: 'invalid'", r.Body.String())
                        assert.Equal(t, http.StatusNotImplemented, r.Code)

		})
}

