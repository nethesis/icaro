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

package methods

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/wax/utils"
)

func DeleteMarketing(c *gin.Context) {
	var userMarketing models.UserMarketing

	userId := c.Param("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		userIdInt = 0
	}

	db := database.Instance()
	db.Where("user_id = ?", userId).First(&userMarketing)

	if userMarketing.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	db.Delete(&userMarketing)

	// remove marketing auth
	user := utils.GetUserById(userIdInt)
	user.MarketingAuth = false
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
