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

package methods

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func CreateAccount(c *gin.Context) {
	creatorId := c.MustGet("token").(models.AccessToken).AccountId

	var json models.AccountJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	h := md5.New()
	h.Write([]byte(json.Password))
	passwordHash := fmt.Sprintf("%x", h.Sum(nil))

	account := models.Account{
		CreatorId: creatorId,
		Uuid:      json.Uuid,
		Type:      json.Type,
		Name:      json.Name,
		Username:  json.Username,
		Password:  passwordHash,
		Email:     json.Email,
		Created:   time.Now().UTC(),
	}

	db := database.Database()
	db.Save(&account)

	if json.Type == "customer" || json.Type == "desk" {
		accountHotspot := models.AccountsHotspot{
			AccountId: account.Id,
			HotspotId: json.HotspotId,
		}
		db.Save(&accountHotspot)
	}

	// TODO: init account preferences to database

	db.Close()

	if account.Id == 0 {
		c.JSON(http.StatusConflict, gin.H{"id": account.Id, "status": "account already exists"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"id": account.Id, "status": "success"})
	}

}

func UpdateAccount(c *gin.Context) {
	var account models.Account
	creatorId := c.MustGet("token").(models.AccessToken).AccountId

	accountId := c.Param("account_id")
	accountIdInt, err := strconv.Atoi(accountId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
	}

	var json models.AccountJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Database()

	// check if the user is me or not
	if accountIdInt == creatorId {
		db.Where("id = ?", accountId).First(&account)
	} else {
		if creatorId == 1 {
			db.Where("id = ?", accountId).First(&account)
		} else {
			db.Where("id = ? AND creator_id = ?", accountId, creatorId).First(&account)
		}
	}

	if account.Id == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "No account found!"})
		return
	}

	if len(json.Password) > 0 {
		h := md5.New()
		h.Write([]byte(json.Password))
		passwordHash := fmt.Sprintf("%x", h.Sum(nil))
		account.Password = passwordHash
	}
	if len(json.Name) > 0 {
		account.Name = json.Name
	}
	if len(json.Username) > 0 {
		account.Username = json.Username
	}
	if len(json.Email) > 0 {
		account.Email = json.Email
	}
	if len(json.Type) > 0 {
		account.Type = json.Type
	}

	db.Save(&account)
	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetAccounts(c *gin.Context) {
	var accounts []models.AccountJSON
	creatorId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")
	hotspotId := c.Query("hotspot")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	if creatorId == 1 {
		if hotspotId != "" {
			db.Select("accounts.*, hotspots.id as hotspot_id, hotspots.name as hotspot_name").Joins("LEFT JOIN accounts_hotspots on accounts_hotspots.account_id = accounts.id LEFT JOIN hotspots on accounts_hotspots.hotspot_id = hotspots.id").Where("hotspots.id = ?", hotspotId).Offset(offsets[0]).Limit(offsets[1]).Find(&accounts)
		} else {
			db.Select("accounts.*, hotspots.id as hotspot_id, hotspots.name as hotspot_name").Joins("LEFT JOIN accounts_hotspots on accounts_hotspots.account_id = accounts.id LEFT JOIN hotspots on accounts_hotspots.hotspot_id = hotspots.id").Offset(offsets[0]).Limit(offsets[1]).Find(&accounts)
		}
	} else {
		if hotspotId != "" {
			db.Select("accounts.*, hotspots.id as hotspot_id, hotspots.name as hotspot_name").Joins("LEFT JOIN accounts_hotspots on accounts_hotspots.account_id = accounts.id LEFT JOIN hotspots on accounts_hotspots.hotspot_id = hotspots.id").Where("creator_id = ? AND hotspots.id = ?", creatorId, hotspotId).Offset(offsets[0]).Limit(offsets[1]).Find(&accounts)
		} else {
			db.Select("accounts.*, hotspots.id as hotspot_id, hotspots.name as hotspot_name").Joins("LEFT JOIN accounts_hotspots on accounts_hotspots.account_id = accounts.id LEFT JOIN hotspots on accounts_hotspots.hotspot_id = hotspots.id").Where("creator_id = ?", creatorId).Offset(offsets[0]).Limit(offsets[1]).Find(&accounts)
		}
	}
	db.Close()

	if len(accounts) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No accounts found!"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func GetAccount(c *gin.Context) {
	var account models.AccountJSON
	creatorId := c.MustGet("token").(models.AccessToken).AccountId

	accountId := c.Param("account_id")
	accountIdInt, err := strconv.Atoi(accountId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Database()

	// check if the user is me or not
	if accountIdInt == creatorId {
		db.Select("accounts.*, hotspots.id as hotspot_id, hotspots.name as hotspot_name").Joins("LEFT JOIN accounts_hotspots on accounts_hotspots.account_id = accounts.id LEFT JOIN hotspots on accounts_hotspots.hotspot_id = hotspots.id").Where("accounts.id = ?", accountId).First(&account)
	} else {
		if creatorId == 1 {
			db.Select("accounts.*, hotspots.id as hotspot_id, hotspots.name as hotspot_name").Joins("LEFT JOIN accounts_hotspots on accounts_hotspots.account_id = accounts.id LEFT JOIN hotspots on accounts_hotspots.hotspot_id = hotspots.id").Where("accounts.id = ?", accountId).First(&account)
		} else {
			db.Select("accounts.*, hotspots.id as hotspot_id, hotspots.name as hotspot_name").Joins("LEFT JOIN accounts_hotspots on accounts_hotspots.account_id = accounts.id LEFT JOIN hotspots on accounts_hotspots.hotspot_id = hotspots.id").Where("accounts.id = ? AND creator_id = ?", accountId, creatorId).First(&account)
		}
	}

	db.Close()

	if account.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No account found!"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func DeleteAccount(c *gin.Context) {
	var account models.Account
	creatorId := c.MustGet("token").(models.AccessToken).AccountId
	accountId := c.Param("account_id")

	db := database.Database()
	if creatorId == 1 {
		db.Where("id = ?", accountId).First(&account)
	} else {
		db.Where("id = ? AND creator_id = ?", accountId, creatorId).First(&account)
	}

	if account.Id == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "No account found!"})
		return
	}

	db.Delete(&account)

	// delete also all accounts created from deleted account
	db.Where("creator_id = ?", accountId).Delete(models.Account{})

	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func StatsAccountTotal(c *gin.Context) {
	creatorId := c.MustGet("token").(models.AccessToken).AccountId
	var count int

	db := database.Database()
	db.Table("accounts").Where("creator_id = ?", creatorId).Count(&count)
	db.Close()

	c.JSON(http.StatusOK, gin.H{"total": count})
}

func StatsSMSTotal(c *gin.Context) {
	var accountSMS models.AccountSmsCount
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	db := database.Database()
	db.Where("account_id = ?", accountId).First(&accountSMS)
	db.Close()

	if accountSMS.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No sms account found!"})
		return
	}

	c.JSON(http.StatusOK, accountSMS)
}
