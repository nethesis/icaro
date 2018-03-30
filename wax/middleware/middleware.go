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

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/wax/utils"
)

func respondWithError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"message": message})
	c.Abort()
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
	unit := utils.GetUnitByUuid(uuid)
	if unit.Id == 0 {
		return false, "uuid is invalid"
	}

	// check if digest is valid
	calcDigest := utils.CalcUnitDigest(unit)
	if digest != calcDigest {
		return false, "digest is invalid"
	}

	return true, ""
}

func WaxWall(c *gin.Context) {
	digest := c.Query("digest")
	uuid := c.Query("uuid")

	check, message := CheckAuth(digest, uuid, c)

	if !check {
		respondWithError(http.StatusBadRequest, message, c)
		return
	}

	// go ahead
	c.Next()
}
