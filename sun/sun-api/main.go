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
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"sun-api/configuration"
	"sun-api/crons"
	"sun-api/methods"
	"sun-api/middleware"
)

func main() {
	// read and init configuration
	configuration.Init()

	// init crons
	crons.Init()

	// init routers
	router := gin.Default()

	// cors
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AddAllowHeaders("Token")
	router.Use(cors.New(corsConf))

	api := router.Group("/api")

	api.POST("/login", methods.Login)
	api.POST("/logout", methods.Logout)

	api.Use(middleware.AAWall)
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
			resellersPref := preferences.Group("/accounts")
			resellersPref.GET("", methods.GetAccountPrefs)
			resellersPref.POST("", methods.CreateAccountPrefs)

			hotspotsPref := preferences.Group("/hotspots")
			hotspotsPref.GET("/:hotspot_id", methods.GetHotspotPrefs)
			hotspotsPref.POST("", methods.CreateHotspotPrefs)
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
			users.PUT("/:user_id", methods.UpdateUser)
			users.DELETE("/:user_id", methods.DeleteUser)
		}

		vouchers := api.Group("/vouchers")
		{
			vouchers.GET("", methods.GetVouchers)
			vouchers.POST("", methods.CreateVoucher)
			vouchers.DELETE("/:voucher_id", methods.DeleteVoucher)
		}
	}

	// handle missing endpoint
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "API not found"})
	})

	router.Run()

}
