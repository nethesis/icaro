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
	"github.com/gin-gonic/gin"
	"sun-api/configuration"
	"wax/methods"
)


func Init(testMode bool) *gin.Engine {
	var r *gin.Engine
	if (testMode == true) {
		gin.SetMode(gin.TestMode)
		r = gin.New()
	} else {
		r = gin.Default()
	}

	r.GET("/aaa", methods.Dispatch)

        return r
}

func main() {
	// read and init configuration
	configuration.Init()

	r := Init(false)

	r.Run(":8181")
}
