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
 * along with NethServer.  If not, see COPYING.
 *
 * author: Giacomo Sanchietti <giacomo.sanchietti@nethesis.it>
 */

package main

import (
	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)


func startupEnv() {
}

func destroyEnv() {
}



func TestNoStage(t *testing.T) {
        r := gofight.New()

        r.GET("/aaa").SetDebug(true).
                Run(Init(true), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "No stage provided", r.Body.String())
                        assert.Equal(t, http.StatusBadRequest, r.Code)
                })
}

func TestRegisterStage(t *testing.T) {
        r := gofight.New()

        r.GET("/aaa?stage=register").SetDebug(true).
                Run(Init(true), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
                        assert.Equal(t, "Not implemented: register", r.Body.String())
                        assert.Equal(t, http.StatusNotImplemented, r.Code)
                })
}


func TestLoginStage(t *testing.T) {
        r := gofight.New()

        r.GET("/aaa?stage=login").SetDebug(true).
                Run(Init(true), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
                        assert.Equal(t, "login", r.Body.String())
                        assert.Equal(t, http.StatusOK, r.Code)
                })
}

func TestCountersStage(t *testing.T) {
        r := gofight.New()

        r.GET("/aaa?stage=counters").SetDebug(true).
                Run(Init(true), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
                        assert.Equal(t, "counters", r.Body.String())
                        assert.Equal(t, http.StatusOK, r.Code)
                })
}

func TestInvalidStage(t *testing.T) {
        r := gofight.New()

        r.GET("/aaa?stage=BAD").SetDebug(true).
                Run(Init(true), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Invalid stage: 'BAD'", r.Body.String())
                        assert.Equal(t, http.StatusNotFound, r.Code)
                })
}

