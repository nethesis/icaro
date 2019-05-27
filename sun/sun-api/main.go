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
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package main

import (
	"flag"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/methods"
	"github.com/nethesis/icaro/sun/sun-api/middleware"
)

func DefineAPI(router *gin.Engine) {
	// cors
	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = configuration.Config.Cors.Origins
	corsConf.AllowHeaders = configuration.Config.Cors.Headers
	corsConf.AllowMethods = configuration.Config.Cors.Methods
	router.Use(cors.New(corsConf))

	health := router.Group("/health")
	health.GET("/check", methods.HealthCheck)

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

		integrations := api.Group("/integrations")
		{
			integrations.GET("", methods.GetIntegrations)
		}

		hotspot_integrations := api.Group("/hotspot_integrations")
		{
			hotspot_integrations.GET("", methods.GetHotspotIntegrations)
			hotspot_integrations.PUT("/:hotspot_id", methods.UpdateHotspotIntegrations)
		}

		marketings := api.Group("/marketing")
		{
			marketings.GET("/:hotspot_id", methods.GetHotspotMarketing)
			marketings.PUT("/:hotspot_id", methods.UpdateHotspotMarketing)
			marketings.POST("/testmail/:hotspot_id/feedback", methods.SendTestFeedbackEmail)
			marketings.POST("/testmail/:hotspot_id/review", methods.SendTestReviewEmail)
		}

		preferences := api.Group("/preferences")
		{
			resellersPref := preferences.Group("/accounts")
			resellersPref.GET("", methods.GetAccountPrefs)
			resellersPref.PUT("", methods.UpdateAccountPrefs)

			hotspotsPref := preferences.Group("/hotspots")
			hotspotsPref.GET("/:hotspot_id", methods.GetHotspotPrefs)
			hotspotsPref.PUT("/:hotspot_id", methods.UpdateHotspotPrefs)
		}

		sessions := api.Group("/sessions")
		{
			sessions.GET("", methods.GetSessions)
			sessions.GET("/:session_id", methods.GetSession)
		}

		histories := api.Group("/histories")
		{
			histories.GET("", methods.GetSessionsHistory)
			histories.GET("/:history_id", methods.GetSessionHistory)
		}

		stats := api.Group("/stats")
		{
			hotspotsStats := stats.Group("/hotspots")
			hotspotsStats.GET("/total", methods.StatsHotspotTotal)

			unitsStats := stats.Group("/units")
			unitsStats.GET("/total", methods.StatsUnitTotal)

			accountsStats := stats.Group("/accounts")
			accountsStats.GET("/total", methods.StatsAccountTotal)

			usersStats := stats.Group("/users")
			usersStats.GET("/total", methods.StatsUserTotal)

			devicesStats := stats.Group("/devices")
			devicesStats.GET("/total", methods.StatsDeviceTotal)

			sessionsStats := stats.Group("/sessions")
			sessionsStats.GET("/total", methods.StatsSessionTotal)

			smsStats := stats.Group("/sms")
			smsStats.GET("/accounts", methods.StatsSMSTotalForAccount)
			smsStats.GET("/accounts/:account_id", methods.StatsSMSTotalForAccount)
			smsStats.POST("/accounts/:account_id", methods.UpdateSMSTotalForAccount)
			smsStats.GET("/hotspots", methods.StatsSMSTotalSentForHotspot)
			smsStats.GET("/hotspots/:hotspot_id", methods.StatsSMSTotalSentForHotspotByHotspot)

			reportStats := stats.Group("/reports")
			reportStats.GET("/current/graph", methods.GetCurrentSessions)
			reportStats.GET("/sessions/graph", methods.GetHistorySessions)
			reportStats.GET("/traffic/graph", methods.GetHistoryTraffic)
			reportStats.GET("/avg_user_traffic/graph", methods.GetHistoryAvgUserTraffic)
			reportStats.GET("/avg_user_duration/graph", methods.GetHistoryAvgUserDuration)
			reportStats.GET("/avg_conn_traffic/graph", methods.GetHistoryAvgConnTraffic)
			reportStats.GET("/avg_conn_duration/graph", methods.GetHistoryAvgConnDuration)
			reportStats.GET("/sms_year/graph", methods.GetHistorySMSYear)
			reportStats.GET("/sms_history/graph", methods.GetHistorySMSHistory)
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

		usersExpired := api.Group("/users_expired")
		{
			usersExpired.GET("", methods.GetUsersExpired)
			usersExpired.PUT("/:user_id", methods.UpdateUserExpired)
			usersExpired.DELETE("/:user_id", methods.DeleteUserExpired)
		}

		vouchers := api.Group("/vouchers")
		{
			vouchers.GET("/:hotspot_id", methods.GetVouchers)
			vouchers.POST("", methods.CreateVoucher)
			vouchers.PUT("/:voucher_id", methods.UpdateVoucher)
			vouchers.DELETE("/:voucher_id", methods.DeleteVoucher)
		}

		subscriptionsPlans := api.Group("/subscription_plans")
		{
			subscriptionsPlans.GET("/", methods.GetSubscriptionPlans)
		}
	}

	// handle missing endpoint
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "API not found"})
	})
}

func main() {
	// read and init configuration
	ConfigFilePtr := flag.String("c", "/opt/icaro/sun-api/conf.json", "Path to configuration file")
	flag.Parse()
	configuration.Init(ConfigFilePtr)

	// init database
	db := database.Init()
	defer db.Close()

	// init routers
	router := gin.Default()

	// define API
	DefineAPI(router)

	router.Run()
}
