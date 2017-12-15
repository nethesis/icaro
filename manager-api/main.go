/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Windmill project.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with NethServer.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package main

import (
	"github.com/gin-gonic/gin"

	"manager-api/configuration"
	"manager-api/methods"
)

func main() {
	// read and init configuration
	configuration.Init()

	// init routers
	router := gin.Default()

	sessions := router.Group("/api/sessions")
	{
		sessions.GET("/", methods.GetSessions)
		sessions.GET("/:session_id", methods.GetSession)
	}

	units := router.Group("/api/units")
	{
		units.GET("/", methods.GetUnits)
		units.GET("/:unit_id", methods.GetUnit)
	}

	devices := router.Group("/api/devices")
	{
		devices.GET("/", methods.GetDevices)
		devices.GET("/:device_id", methods.GetDevice)
	}

	users := router.Group("/api/users")
	{
		users.GET("/", methods.GetUsers)
		users.GET("/:user_id", methods.GetUser)
		users.POST("/", methods.CreateUser)
		users.PUT("/:user_id", methods.UpdateUser)
		users.DELETE("/:user_id", methods.DeleteUser)
	}

	router.Run()

}
