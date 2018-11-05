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
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func CreateUser(user models.User) int {
	user.Created = time.Now().UTC()

	// save new user
	db := database.Instance()
	db.Save(&user)

	return user.Id
}

func UpdateUser(c *gin.Context) {
	var user models.User
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	userId := c.Param("user_id")

	var json models.User
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()
	db.Where("id = ?", userId).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId, (accountId == 1), user.HotspotId), user.HotspotId) {
		if len(json.Name) > 0 {
			user.Name = json.Name
		}
		if len(json.Email) > 0 {
			user.Email = json.Email
		}

		if json.KbpsDown >= 0 {
			user.KbpsDown = json.KbpsDown
		}
		if json.KbpsUp >= 0 {
			user.KbpsUp = json.KbpsUp
		}

		if json.MaxNavigationTraffic >= 0 {
			user.MaxNavigationTraffic = json.MaxNavigationTraffic
		}
		if json.MaxNavigationTime >= 0 {
			user.MaxNavigationTime = json.MaxNavigationTime
		}

		if !json.ValidFrom.IsZero() {
			user.ValidFrom = json.ValidFrom
		}
		if !json.ValidUntil.IsZero() {
			user.ValidUntil = json.ValidUntil
		}

		user.AutoLogin = json.AutoLogin

		db.Save(&user)

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This user is not yours"})
	}
}

func UpdateUserExpired(c *gin.Context) {
	var user models.UserHistory
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	userId := c.Param("user_id")

	var json models.User
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()
	db.Where("id = ?", userId).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId, (accountId == 1), user.HotspotId), user.HotspotId) {
		if len(json.Name) > 0 {
			user.Name = json.Name
		}
		if len(json.Email) > 0 {
			user.Email = json.Email
		}

		if json.KbpsDown >= 0 {
			user.KbpsDown = json.KbpsDown
		}
		if json.KbpsUp >= 0 {
			user.KbpsUp = json.KbpsUp
		}

		if json.MaxNavigationTraffic >= 0 {
			user.MaxNavigationTraffic = json.MaxNavigationTraffic
		}
		if json.MaxNavigationTime >= 0 {
			user.MaxNavigationTime = json.MaxNavigationTime
		}

		if !json.ValidFrom.IsZero() {
			user.ValidFrom = json.ValidFrom
		}
		if !json.ValidUntil.IsZero() {
			user.ValidUntil = json.ValidUntil
		}

		user.AutoLogin = json.AutoLogin

		// check if user is not expired
		if user.ValidUntil.After(time.Now().UTC()) {
			newUser := models.User{
				HotspotId:            user.HotspotId,
				Name:                 user.Name,
				Username:             user.Username,
				Password:             user.Password,
				Email:                user.Email,
				EmailVerified:        user.EmailVerified,
				AccountType:          user.AccountType,
				MarketingAuth:        user.MarketingAuth,
				KbpsDown:             user.KbpsDown,
				KbpsUp:               user.KbpsUp,
				MaxNavigationTraffic: user.MaxNavigationTraffic,
				MaxNavigationTime:    user.MaxNavigationTime,
				AutoLogin:            user.AutoLogin,
				ValidFrom:            user.ValidFrom,
				ValidUntil:           user.ValidUntil,
				Created:              user.Created,
			}

			// save new user
			db.Save(&newUser)

			// delete the old one
			db.Delete(&user)
		} else {
			db.Save(&user)
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This user is not yours"})
	}
}

func GetUsers(c *gin.Context) {
	var users []models.User
	var total int
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")
	hotspotId := c.Query("hotspot")
	accountType := c.Query("type")
	marketing := c.Query("marketing")
	q := c.Query("q")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	offsets := utils.OffsetCalc(page, limit)

	db := database.Instance()
	chain := db.Where("hotspot_id in (?) AND valid_until > NOW()", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt))

	if len(accountType) > 0 {
		chain = chain.Where("account_type = ?", accountType)
	}

	if len(marketing) > 0 {
		chain = chain.Where("marketing_auth = 1")
	}

	if len(q) > 0 {
		chain = chain.Where("username LIKE ? OR name LIKE ? OR email LIKE ? OR account_type LIKE ? OR kbps_down LIKE ? OR kbps_up LIKE ? OR max_navigation_time LIKE ? OR max_navigation_traffic LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	chain.Find(&users).Count(&total)
	chain.Offset(offsets[0]).Limit(offsets[1]).Find(&users)

	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No users found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users, "total": total})
}

func GetUsersExpired(c *gin.Context) {
	var users []models.User
	var userHistories []models.UserHistory
	var totalUsers int
	var totalHistories int
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")
	hotspotId := c.Query("hotspot")
	accountType := c.Query("type")
	q := c.Query("q")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	offsets := utils.OffsetCalc(page, limit)

	db := database.Instance()

	// users expired
	chain := db.Where("hotspot_id in (?) AND valid_until <= NOW()", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt))
	if len(accountType) > 0 {
		chain = chain.Where("account_type = ?", accountType)
	}
	chain.Find(&users).Count(&totalUsers)
	chain.Offset(offsets[0]).Limit(offsets[1]).Find(&users)

	// users history
	chain = db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt))
	if len(accountType) > 0 {
		chain = chain.Where("account_type = ?", accountType)
	}

	if len(q) > 0 {
		chain = chain.Where("username LIKE ? OR name LIKE ? OR email LIKE ? OR account_type LIKE ? OR kbps_down LIKE ? OR kbps_up LIKE ? OR max_navigation_time LIKE ? OR max_navigation_traffic LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	chain.Find(&userHistories).Count(&totalHistories)
	chain.Offset(offsets[0]).Limit(offsets[1]).Find(&userHistories)

	if len(users) == 0 && len(userHistories) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user expired found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": totalHistories + totalUsers, "data_user_histories": userHistories, "data_users": users})
}

func GetUser(c *gin.Context) {
	var user models.User
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	userId := c.Param("user_id")

	db := database.Instance()
	db.Where("id = ? AND hotspot_id in (?)", userId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	userId := c.Param("user_id")

	db := database.Instance()
	db.Where("id = ? AND hotspot_id in (?)", userId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	// delete marketing info
	var userMarketing models.UserMarketing
	db.Where("user_id = ?", user.Id).First(&userMarketing)
	db.Delete(&userMarketing)

	// delete user
	db.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func DeleteUserExpired(c *gin.Context) {
	var user models.UserHistory
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	userId := c.Param("user_id")

	db := database.Instance()
	db.Where("id = ? AND hotspot_id in (?)", userId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user found!"})
		return
	}

	// delete marketing info
	var userMarketing models.UserMarketing
	db.Where("user_id = ?", user.UserId).First(&userMarketing)
	db.Delete(&userMarketing)

	// delete user
	db.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func StatsUserTotal(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId
	var count int

	db := database.Instance()
	db.Table("users").Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).Count(&count)

	c.JSON(http.StatusOK, gin.H{"total": count})
}
