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
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package main

import (
	"flag"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/wax/methods"
	"github.com/nethesis/icaro/wax/middleware"
)

func DefineAPI(router *gin.Engine) {
	// cors
	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = configuration.Config.Cors.Origins
	corsConf.AllowHeaders = configuration.Config.Cors.Headers
	corsConf.AllowMethods = configuration.Config.Cors.Methods
	router.Use(cors.New(corsConf))

	wax := router.Group("/wax")

	wax.Use(middleware.WaxWall)
	{
		// handle AAA requests
		wax.GET("/aaa", methods.Dispatch)

		// handle voucher control
		wax.GET("/register/voucher/:code", methods.VoucherAuth)

		// handle social logins
		wax.GET("/register/social/facebook/:code", methods.FacebookAuth)
		wax.GET("/register/social/linkedin/:code", methods.LinkedInAuth)
		wax.GET("/register/social/instagram/:code", methods.InstagramAuth)

		// handle sms login
		wax.GET("/register/sms/:number", methods.SMSAuth)

		// handle email login
		wax.GET("/register/email/:email", methods.EmailAuth)

		// handle wings preferences
		wax.GET("/preferences", methods.GetWingsPrefs)
	}

	// handle missing endpoint
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "API not found"})
	})
}

func main() {
	// read and init configuration
	ConfigFilePtr := flag.String("c", "/opt/icaro/wax/conf.json", "Path to configuration file")
	flag.Parse()
	configuration.Init(ConfigFilePtr)

	// init routers
	router := gin.Default()

	// define API
	DefineAPI(router)

	router.Run()
}
