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
	uuid "github.com/satori/go.uuid"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func CreateHotspot(c *gin.Context) {
	var json models.Hotspot
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspot := models.Hotspot{
		Uuid:            uuid.Must(uuid.NewV4()).String(),
		AccountId:       accountId,
		Name:            json.Name,
		Description:     json.Description,
		BusinessName:    json.BusinessName,
		BusinessVAT:     json.BusinessVAT,
		BusinessAddress: json.BusinessAddress,
		BusinessEmail:   json.BusinessEmail,
		BusinessDPO:     json.BusinessDPO,
		Created:         time.Now().UTC(),
	}

	db := database.Instance()
	db.Save(&hotspot)

	if hotspot.Id == 0 {
		c.JSON(http.StatusConflict, gin.H{"id": hotspot.Id, "status": "hotspot already exists"})
	} else {
		// initialize hotspot preferences
		utils.SetDefaultHotspotPreferences(hotspot.Id)

		c.JSON(http.StatusCreated, gin.H{"id": hotspot.Id, "status": "success"})
	}
}

func UpdateHotspot(c *gin.Context) {
	var hotspot models.Hotspot
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	var json models.Hotspot
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()

	if accountId == 1 {
		db.Where("id = ?", hotspotId).First(&hotspot)
	} else {
		db.Where("id = ? AND account_id = ?", hotspotId, accountId).First(&hotspot)
	}

	if hotspot.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot found!"})
		return
	}

	if len(json.Description) > 0 {
		hotspot.Description = json.Description
	}
	if len(json.BusinessName) > 0 {
		hotspot.BusinessName = json.BusinessName
	}
	if len(json.BusinessVAT) > 0 {
		hotspot.BusinessVAT = json.BusinessVAT
	}
	if len(json.BusinessAddress) > 0 {
		hotspot.BusinessAddress = json.BusinessAddress
	}
	if len(json.BusinessEmail) > 0 {
		hotspot.BusinessEmail = json.BusinessEmail
	}
	if len(json.BusinessDPO) > 0 {
		hotspot.BusinessDPO = json.BusinessDPO
	}

	db.Save(&hotspot)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetHotspots(c *gin.Context) {
	var hotspots []models.Hotspot
	var total int
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	page := c.Query("page")
	limit := c.Query("limit")
	q := c.Query("q")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Instance()
	chain := db.Order("name asc, description")

	if len(q) > 0 {
		chain = chain.Where("business_dpo LIKE ? OR business_address LIKE ? OR business_email LIKE ? OR business_name LIKE ? OR business_vat LIKE ? OR description LIKE ? OR name LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	if accountId == 1 {
		chain = chain.Find(&hotspots).Count(&total)
		chain = chain.Offset(offsets[0]).Limit(offsets[1]).Find(&hotspots)
	} else {
		chain = chain.Where("account_id = ?", accountId).Find(&hotspots).Count(&total)
		chain = chain.Where("account_id = ?", accountId).Offset(offsets[0]).Limit(offsets[1]).Find(&hotspots)
	}

	if len(hotspots) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspots found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": hotspots, "total": total})
}

func GetHotspot(c *gin.Context) {
	var hotspot models.HotspotJSON
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	db := database.Instance()
	if accountId == 1 {
		db.Select("hotspots.*, accounts.name as account_name").Where("hotspots.id = ?", hotspotId).Joins("JOIN accounts on accounts.id = hotspots.account_id").First(&hotspot)
	} else {
		var account models.Account
		db.Where("id = ?", accountId).First(&account)

		if account.Type == "customer" || account.Type == "desk" {
			accountId = account.CreatorId
		}

		db.Select("hotspots.*, accounts.name as account_name").Where("hotspots.id = ? AND account_id = ?", hotspotId, accountId).Joins("JOIN accounts on accounts.id = hotspots.account_id").First(&hotspot)
	}

	if hotspot.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot found!"})
		return
	}

	c.JSON(http.StatusOK, hotspot)
}

func DeleteHotspot(c *gin.Context) {
	var hotspot models.Hotspot
	var accountsHotspot []models.AccountsHotspot

	accountId := c.MustGet("token").(models.AccessToken).AccountId
	hotspotId := c.Param("hotspot_id")

	db := database.Instance()
	if accountId == 1 {
		db.Where("id = ?", hotspotId).First(&hotspot)
	} else {
		db.Where("id = ? AND account_id = ?", hotspotId, accountId).First(&hotspot)
	}

	if hotspot.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot found!"})
		return
	}

	// delete all accounts related
	db.Where("hotspot_id = ?", hotspotId).Find(&accountsHotspot)
	for _, accountHotspot := range accountsHotspot {
		var account models.Account

		db.Where("id = ?", accountHotspot.AccountId).First(&account)
		db.Delete(&account)
	}

	// delete hotspot
	db.Delete(&hotspot)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func StatsHotspotTotal(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId
	var count int

	db := database.Instance()
	if accountId == 1 {
		db.Table("hotspots").Count(&count)
	} else {
		db.Table("hotspots").Where("account_id = ?", accountId).Count(&count)
	}

	c.JSON(http.StatusOK, gin.H{"total": count})
}

func StatsSMSTotalSentForHotspot(c *gin.Context) {
	var hotspotSmsCount []models.HotspotSmsCount
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	db := database.Instance()
	db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).Find(&hotspotSmsCount)

	if len(hotspotSmsCount) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No sms stats found!"})
		return
	}

	c.JSON(http.StatusOK, hotspotSmsCount)
}

func StatsSMSTotalSentForHotspotByHotspot(c *gin.Context) {
	var hotspotSmsCount []models.HotspotSmsCount
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	db := database.Instance()
	db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Find(&hotspotSmsCount)

	if len(hotspotSmsCount) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No sms stats found!"})
		return
	}

	c.JSON(http.StatusOK, hotspotSmsCount)
}
