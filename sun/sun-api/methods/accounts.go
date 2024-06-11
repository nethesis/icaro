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
	role := c.MustGet("token").(models.AccessToken).Role

	var subscriptionPlan models.SubscriptionPlan
	var accountSMS models.AccountSmsCount
	var hotspot models.Hotspot

	var json models.AccountJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	// Valid account types are: reseller, desk, customer
	if json.Type != "reseller" && json.Type != "desk" && json.Type != "customer" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unsupported account type"})
		return
	}

	// Only admin and reseller can create new users
	if role == "desk" || role == "customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Permission denied"})
		return
	}

	// reseller can create only customer and desk accounts
	if role == "reseller" && json.Type != "desk" && json.Type != "customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Permission denied"})
		return
	}

	h := md5.New()
	h.Write([]byte(json.Password))
	passwordHash := fmt.Sprintf("%x", h.Sum(nil))

	account := models.Account{
		CreatorId:      creatorId,
		Uuid:           json.Uuid,
		Type:           json.Type,
		Name:           json.Name,
		Username:       json.Username,
		Password:       passwordHash,
		Email:          json.Email,
		PrivacyName:    "",
		PrivacyVAT:     "",
		PrivacyAddress: "",
		PrivacyEmail:   "",
		PrivacyDPO:     "",
		PrivacyDPOMail: "",
		Created:        time.Now().UTC(),
	}

	// create account record
	db := database.Instance()
	db.Save(&account)

	if json.Type == "customer" || json.Type == "desk" {
		accountHotspot := models.AccountsHotspot{
			AccountId: account.Id,
			HotspotId: json.HotspotId,
		}
		db.Save(&accountHotspot)

		// force creator from hotspot account_id
		if creatorId == 1 {
			db.Where("id = ?", json.HotspotId).First(&hotspot)
			account.CreatorId = hotspot.AccountId
			db.Save(&account)
		}
	} else if json.Type == "reseller" {
		// retrieve subscription plain

		db.Where("id = ?", json.SubscriptionPlanId).First(&subscriptionPlan)

		// create new subscription
		subscription := models.Subscription{
			AccountID:          account.Id,
			SubscriptionPlanID: json.SubscriptionPlanId,
			ValidFrom:          time.Now().UTC(),
			ValidUntil:         time.Now().UTC().AddDate(0, 0, subscriptionPlan.Period),
			Created:            time.Now().UTC(),
		}
		db.Save(&subscription)

		// create SMS accounting
		accountSMS.AccountId = account.Id
		accountSMS.SmsMaxCount = subscriptionPlan.IncludedSMS
		db.Save(&accountSMS)
	}

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

	db := database.Instance()

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

	// privacy update
	if len(json.PrivacyName) > 0 {
		account.PrivacyName = json.PrivacyName
	}
	if len(json.PrivacyVAT) > 0 {
		account.PrivacyVAT = json.PrivacyVAT
	}
	if len(json.PrivacyAddress) > 0 {
		account.PrivacyAddress = json.PrivacyAddress
	}
	if len(json.PrivacyEmail) > 0 {
		account.PrivacyEmail = json.PrivacyEmail
	}
	account.PrivacyDPO = json.PrivacyDPO
	account.PrivacyDPOMail = json.PrivacyDPOMail

	// NOTE: Subscription plan change is not supported

	db.Save(&account)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetAccounts(c *gin.Context) {
	var accounts []models.AccountJSON
	var total int
	creatorId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")
	hotspotId := c.Query("hotspot")
	q := c.Query("q")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Instance()
	selectQuery := "accounts.*, hotspots.id as hotspot_id, hotspots.name as hotspot_name"
	joinQuery := "LEFT JOIN accounts_hotspots on accounts_hotspots.account_id = accounts.id LEFT JOIN hotspots on accounts_hotspots.hotspot_id = hotspots.id"

	chain := db.Select(selectQuery).Joins(joinQuery)

	if len(q) > 0 {
		chain = chain.Where("accounts.username LIKE ? OR accounts.name LIKE ? OR accounts.email LIKE ? OR accounts.type LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	if creatorId == 1 {
		// Queries for admin user
		if hotspotId != "" {
			chain = chain.Where("hotspots.id = ?", hotspotId).Find(&accounts).Count(&total)
			chain = chain.Where("hotspots.id = ?", hotspotId).Offset(offsets[0]).Limit(offsets[1]).Find(&accounts)
		} else {
			chain = chain.Find(&accounts).Count(&total)
			chain = chain.Offset(offsets[0]).Limit(offsets[1]).Find(&accounts)
		}
	} else {
		if hotspotId != "" {
			chain = chain.Where("creator_id = ? AND hotspots.id = ?", creatorId, hotspotId).Find(&accounts).Count(&total)
			chain = chain.Where("creator_id = ? AND hotspots.id = ?", creatorId, hotspotId).Offset(offsets[0]).Limit(offsets[1]).Find(&accounts)
		} else {
			chain = chain.Where("creator_id = ?", creatorId).Find(&accounts).Count(&total)
			chain = chain.Where("creator_id = ?", creatorId).Offset(offsets[0]).Limit(offsets[1]).Find(&accounts)
		}
	}

	if len(accounts) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No accounts found!"})
		return
	}

	db.Set("gorm:auto_preload", true)
	for i, account := range accounts {
		var subscription models.Subscription
		db.Preload("SubscriptionPlan").Where("account_id = ?", account.Id).First(&subscription)
		accounts[i].Subscription = subscription
		accounts[i].Subscription.Expired = subscription.IsExpired()

		var tokens []models.AccessToken
		db.Where("account_id = ? AND type = 'api'", account.Id).Find(&tokens)
		accounts[i].Tokens = tokens
	}

	c.JSON(http.StatusOK, gin.H{"data": accounts, "total": total})
}

func GetAccount(c *gin.Context) {
	var account models.AccountJSON
	var subscription models.Subscription
	var tokens []models.AccessToken

	creatorId := c.MustGet("token").(models.AccessToken).AccountId

	accountId := c.Param("account_id")
	accountIdInt, err := strconv.Atoi(accountId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()

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

	if account.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No account found!"})
		return
	}

	db.Preload("SubscriptionPlan").Where("account_id = ?", account.Id).First(&subscription)
	account.Subscription = subscription
	account.Subscription.Expired = subscription.IsExpired()

	db.Where("account_id = ? AND type = 'api'", account.Id).Find(&tokens)
	account.Tokens = tokens

	c.JSON(http.StatusOK, account)
}

func DeleteAccount(c *gin.Context) {
	var account models.Account
	creatorId := c.MustGet("token").(models.AccessToken).AccountId
	accountId := c.Param("account_id")

	db := database.Instance()
	if creatorId == 1 {
		db.Where("id = ?", accountId).First(&account)
	} else {
		db.Where("id = ? AND creator_id = ?", accountId, creatorId).First(&account)
	}

	if account.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No account found!"})
		return
	}

	db.Delete(&account)

	// delete also all accounts created from deleted account
	db.Where("creator_id = ?", accountId).Delete(models.Account{})

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func StatsAccountTotal(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId
	var count int

	db := database.Instance()

	if accountId == 1 {
		db.Table("accounts").Count(&count)
	} else {
		db.Table("accounts").Where("creator_id = ?", accountId).Count(&count)
	}

	c.JSON(http.StatusOK, gin.H{"total": count})
}

func StatsSMSTotalForAccount(c *gin.Context) {
	var accountSMS models.AccountSmsCount
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	db := database.Instance()

	if accountId == 1 {
		destAccountId := c.Param("account_id")
		db.Where("account_id = ?", destAccountId).First(&accountSMS)
	} else {
		db.Where("account_id = ?", accountId).First(&accountSMS)
	}

	if accountSMS.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No sms account found!"})
		return
	}

	c.JSON(http.StatusOK, accountSMS)
}

func UpdateSMSThresholdForAccount(c *gin.Context) {
	var accountSMS models.AccountSmsCount

	var json models.AccountSmsThresholdJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()

	destAccountId := c.Param("account_id")
	db.Where("account_id = ?", destAccountId).First(&accountSMS)

	accountSMS.SmsThreshold = json.SmsThreshold

	db.Save(&accountSMS)
}

func UpdateSMSTotalForAccount(c *gin.Context) {
	var accountSMS models.AccountSmsCount
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	if accountId != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation not permitted!"})
		return
	}

	var json models.AccountSmsCountJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()

	destAccountId := c.Param("account_id")
	db.Where("account_id = ?", destAccountId).First(&accountSMS)

	accountSMS.SmsMaxCount = accountSMS.SmsMaxCount + json.SmsToAdd

	db.Save(&accountSMS)
}
