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

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sun-api/configuration"
	"wax/utils"
)

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func respondWithError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"message": message})
	c.Abort()
}

func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func CheckAuth(digest string, uuid string, c *gin.Context) (bool, string) {
	// check if digest exists
	if digest == "" {
		return false, "digest required"
	}
	// check if uuid exists
	if uuid == "" {
		return false, "uuid required"
	}

	// check if uuid is valid
	unit := utils.ExtractUnit(uuid)
	if unit.Id == 0 {
		return false, "uuid is invalid"
	}

	// check if digest is valid
	calcDigest := utils.CalcDigest(unit)
	if digest != calcDigest {
		return false, "digest is invalid"
	}

	return true, ""
}

func WaxWall(c *gin.Context) {
	digest := c.Param("digest")
	uuid := c.Param("uuid")

	check, message := CheckAuth(digest, uuid, c)

	if !check {
		respondWithError(http.StatusBadRequest, message, c)
		return
	}

	// go ahead
	c.Next()

}

func CaptiveWings(c *gin.Context) {
	digest := c.Request.URL.Query().Get("digest")
	uuid := c.Request.URL.Query().Get("uuid")

	check, message := CheckAuth(digest, uuid, c)

	if !check {
		respondWithError(http.StatusBadRequest, message, c)
		return
	}

	// server static files of captive portal
	fileserver := http.FileServer(LocalFile(configuration.Config.CaptivePath, true))
	fileserver = http.StripPrefix("/wax/captive", fileserver)
	fileserver.ServeHTTP(c.Writer, c.Request)
	c.Abort()

}
