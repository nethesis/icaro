/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Icaro project.
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

	"sun-api/configuration"
	"sun-api/methods"
	"sun-api/middleware"
)

func main() {
	// read and init configuration
	configuration.Init()

	// init routers
	router := gin.Default()

	api := router.Group("/api")

	api.POST("/login", methods.Login)
	api.POST("/logout", methods.Logout)

	api.Use(middleware.Authentication)
	{
		accounts := api.Group("/accounts")
		{
			accounts.GET("", methods.GetAccounts)
			accounts.GET("/:account_id", methods.GetAccount)
			accounts.POST("", methods.CreateAccount)
			accounts.PUT("/:account_id", methods.UpdateAccount)
			accounts.DELETE("/:account_id", methods.DeleteAccount)
		}

		devices := api.Group("/devices")
		{
			devices.GET("", methods.GetDevices)
			devices.GET("/:device_id", methods.GetDevice)
		}

		hotspots := api.Group("/hotspots")
		{
			hotspots.GET("", methods.GetHotspots)
			hotspots.GET("/:hotspot_id", methods.GetHotspot)
			hotspots.POST("", methods.CreateHotspot)
			hotspots.PUT("/:hotspot_id", methods.UpdateHotspot)
			hotspots.DELETE("/:hotspot_id", methods.DeleteHotspot)
		}

		preferences := api.Group("/preferences")
		{
			resellers_pref := preferences.Group("/accounts")
			resellers_pref.GET("", methods.GetAccountPrefs)
			resellers_pref.POST("", methods.CreateAccountPrefs)

			hotspots_pref := preferences.Group("/hotspots")
			hotspots_pref.GET("", methods.GetHotspotPrefs)
			hotspots_pref.POST("", methods.CreateHotspotPrefs)
		}

		sessions := api.Group("/sessions")
		{
			sessions.GET("", methods.GetSessions)
			sessions.GET("/:session_id", methods.GetSession)
		}

		units := api.Group("/units")
		{
			units.GET("", methods.GetUnits)
			units.GET("/:unit_id", methods.GetUnit)
			units.POST("", methods.CreateUnit)
			units.DELETE("/:unit_id", methods.DeleteUnit)
		}

		users := api.Group("/users")
		{
			users.GET("", methods.GetUsers)
			users.GET("/:user_id", methods.GetUser)
			users.POST("", methods.CreateUser)
			users.PUT("/:user_id", methods.UpdateUser)
			users.DELETE("/:user_id", methods.DeleteUser)
		}
	}

	router.Run()

}
