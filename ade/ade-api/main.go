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
 * author: Matteo Valentini <matteo.valentini@nethesis.it>
 */

package main

import (
	"flag"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/ade/ade-api/methods"
	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
)

func DefineAPI(router *gin.Engine) {
	// cors
	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = configuration.Config.Cors.Origins
	corsConf.AllowHeaders = configuration.Config.Cors.Headers
	corsConf.AllowMethods = configuration.Config.Cors.Methods
	router.Use(cors.New(corsConf))

	api := router.Group("/ade")

	health := api.Group("/health")
	health.GET("/check", methods.HealthCheck)

	feedbacks := api.Group("/feedbacks")
	feedbacks.GET("/:token", methods.GetFeedbackPage)
	feedbacks.POST("/:token", methods.PostFeedbackResult)

	reviews := api.Group("/reviews")
	reviews.GET("/:token", methods.GetReviewPage)
	reviews.POST("/:token", methods.PostReviewResult)

	api.GET("/short/:hash", methods.GetLongUrl)

	// handle missing endpoint
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "API not found"})
	})
}

func main() {
	// read and init configuration
	ConfigFilePtr := flag.String("c", "/opt/icaro/ade-api/conf.json", "Path to configuration file")
	flag.Parse()
	configuration.Init(ConfigFilePtr)

	// init Database
	db := database.Init()
	defer db.Close()

	// init routers
	router := gin.Default()

	// define API
	DefineAPI(router)

	router.Run()
}
