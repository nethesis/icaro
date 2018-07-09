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
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func CreateVoucher(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	var json models.HotspotVoucher
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	// add ownership of account if type desk
	account := utils.GetAccountById(accountId)
	var ownerId int
	if account.Type == "desk" {
		ownerId = account.Id
	}

	hotspotVoucher := models.HotspotVoucher{
		Code:          json.Code,
		AutoLogin:     json.AutoLogin,
		BandwidthUp:   json.BandwidthUp,
		BandwidthDown: json.BandwidthDown,
		MaxTraffic:    json.MaxTraffic,
		MaxTime:       json.MaxTime,
		Duration:      json.Duration,
		RemainUse:     json.RemainUse,
		OwnerId:       ownerId,
	}

	hotspotVoucher.HotspotId = json.HotspotId

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId, (accountId == 1), json.HotspotId), json.HotspotId) {
		db := database.Instance()
		db.Save(&hotspotVoucher)

		if hotspotVoucher.Id == 0 {
			c.JSON(http.StatusConflict, gin.H{"id": hotspotVoucher.Id, "status": "voucher already exists"})
		} else {
			c.JSON(http.StatusCreated, gin.H{"id": hotspotVoucher.Id, "status": "success"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}

func UpdateVoucher(c *gin.Context) {
	var hotspotVoucher models.HotspotVoucher
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	var json models.HotspotVoucher
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	voucherId := c.Param("voucher_id")

	db := database.Instance()
	db.Where("id = ? AND hotspot_id in (?)", voucherId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&hotspotVoucher)

	if hotspotVoucher.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot voucher found!"})
		return
	}

	if json.Printed {
		hotspotVoucher.Printed = true
	}
	db.Save(&hotspotVoucher)

	c.JSON(http.StatusOK, gin.H{"status": "success"})

}

func GetVouchers(c *gin.Context) {
	var hotspotVouchers []models.HotspotVoucher
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	account := utils.GetAccountById(accountId)

	page := c.Query("page")
	limit := c.Query("limit")

	duration := c.Query("duration")
	autoLogin := c.Query("auto_login")
	used := c.Query("used")
	reusable := c.Query("reusable")
	printed := c.Query("printed")
	showAll := c.Query("show_all")

	hotspotId := c.Param("hotspot_id")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	offsets := utils.OffsetCalc(page, limit)

	db := database.Instance()
	chain := db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt))

	// if desk account, get only my voucher not voucher of other users
	if account.Type == "desk" {
		chain = chain.Where("owner_id = ?", accountId)
	}

	if len(showAll) == 0 {
		if len(duration) > 0 {
			chain = chain.Where("duration = ?", duration)
		}

		if len(autoLogin) > 0 {
			chain = chain.Where("auto_login = 1")
		} else {
			chain = chain.Where("auto_login = 0")
		}

		if len(used) > 0 {
			chain = chain.Where("expires != 0")
		} else {
			chain = chain.Where("expires = 0")
		}

		if len(reusable) > 0 {
			chain = chain.Where("remain_use = -1")
		} else {
			chain = chain.Where("remain_use > 0")
		}

		if len(printed) > 0 {
			chain = chain.Where("printed = 1")
		} else {
			chain = chain.Where("printed = 0")
		}
	}

	chain.Offset(offsets[0]).Limit(offsets[1]).Find(&hotspotVouchers)

	if len(hotspotVouchers) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot vouchers found!"})
		return
	}

	c.JSON(http.StatusOK, hotspotVouchers)
}

func DeleteVoucher(c *gin.Context) {
	var hotspotVoucher models.HotspotVoucher
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	voucherId := c.Param("voucher_id")

	db := database.Instance()
	db.Where("id = ? AND hotspot_id in (?)", voucherId, utils.ExtractHotspotIds(accountId, (accountId == 1), 0)).First(&hotspotVoucher)

	if hotspotVoucher.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot voucher found!"})
		return
	}

	db.Delete(&hotspotVoucher)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
