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
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
	waxUtils "github.com/nethesis/icaro/wax/utils"
)

type voucherMarketingData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func generateVoucherCode() string {
	length := 8
	charset := strings.ToLower("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	voucherCode := ""
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		voucherCode += string(charset[rand.Intn(len(charset))])
	}
	return voucherCode
}

func CreateVouchers(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	var json models.HotspotVoucherJSON
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
		Code:          strings.ToLower(json.Code),
		AutoLogin:     json.AutoLogin,
		BandwidthUp:   json.BandwidthUp,
		BandwidthDown: json.BandwidthDown,
		MaxTraffic:    json.MaxTraffic,
		MaxTime:       json.MaxTime,
		RemainUse:     json.RemainUse,
		Type:          json.Type,
		UserName:      json.UserName,
		UserMail:      json.UserMail,
		OwnerId:       ownerId,
		Created:       time.Now().UTC(),
	}

	if json.Time == "duration" {
		hotspotVoucher.Duration = json.Duration
	}

	if json.Time == "expiration" {
		hotspotVoucher.Expires = time.Now().UTC().AddDate(0, 0, json.Expiration)
	}

	hotspotVoucher.HotspotId = json.HotspotId

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId, (accountId == 1), json.HotspotId), json.HotspotId) {
		db := database.Instance()

		for i := 0; i < json.NumVouchers; i++ {
			voucherInstance := hotspotVoucher

			if json.Code == "" || json.NumVouchers > 1 {
				voucherInstance.Code = generateVoucherCode()
			}
			db.Save(&voucherInstance)

			if voucherInstance.Id == 0 {
				c.JSON(http.StatusConflict, gin.H{"id": voucherInstance.Id, "status": "voucher already exists"})
				return
			}

			var newUserId = 0
			if voucherInstance.Type == "auth" {
				userName := voucherInstance.UserName

				if i > 0 {
					userName += "-" + strconv.Itoa(i)
				}

				newUser := models.User{
					HotspotId:            json.HotspotId,
					Name:                 userName + " (" + strings.ToUpper(voucherInstance.Code) + ")",
					Username:             strings.ToUpper(voucherInstance.Code),
					Password:             voucherInstance.Code,
					Email:                voucherInstance.UserMail,
					AccountType:          "voucher",
					Reason:               "",
					Country:              "",
					MarketingAuth:        true,
					SurveyAuth:           true,
					KbpsDown:             voucherInstance.BandwidthDown,
					KbpsUp:               voucherInstance.BandwidthUp,
					MaxNavigationTraffic: voucherInstance.MaxTraffic,
					MaxNavigationTime:    voucherInstance.MaxTime,
					AutoLogin:            voucherInstance.AutoLogin,
				}

				// create user
				newUserId = CreateUser(newUser)

				// create marketing info with user infos
				waxUtils.CreateUserMarketing(newUserId, voucherMarketingData{Name: voucherInstance.UserName, Email: voucherInstance.UserMail}, "voucher")
			}
		}
		c.JSON(http.StatusCreated, gin.H{"status": "success"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}

func UpdateVouchers(c *gin.Context) {
	var hotspotVouchers []models.HotspotVoucher
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	var json models.HotspotVoucherJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	voucherIds := json.VoucherIds
	hotspotIds := utils.ExtractHotspotIds(accountId, (accountId == 1), 0)
	db := database.Instance()
	db.Where("id IN (?) AND hotspot_id IN (?)", voucherIds, hotspotIds).Find(&hotspotVouchers)

	if len(hotspotVouchers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hotspot voucher found!"})
		return
	}

	db.Model(models.HotspotVoucher{}).Where("id IN (?) AND hotspot_id IN (?)", voucherIds, hotspotIds).Updates(models.HotspotVoucher{Printed: json.Printed})

	c.JSON(http.StatusOK, gin.H{"status": "success"})

}

func GetVouchers(c *gin.Context) {
	var hotspotVouchers []models.HotspotVoucher
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	account := utils.GetAccountById(accountId)

	page := c.Query("page")
	limit := c.Query("limit")

	code := c.Query("code")
	duration := c.Query("duration")
	autoLogin := c.Query("auto_login")
	used := c.Query("used")
	reusable := c.Query("reusable")
	printed := c.Query("printed")
	voucherType := c.Query("type")
	expiredStart := c.Query("expiredStart")
	expiredEnd := c.Query("expiredEnd")
	createdStart := c.Query("createdStart")
	createdEnd := c.Query("createdEnd")

	hotspotId := c.Param("hotspot_id")

	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	offsets := utils.OffsetCalc(page, limit)
	timeLayout := "2006-01-02T15:04:05.000Z"

	db := database.Instance()
	chain := db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt))

	// if desk account, get only my voucher not voucher of other users
	if account.Type == "desk" {
		chain = chain.Where("owner_id = ?", accountId)
	}

	if len(code) > 0 {
		chain = chain.Where("code LIKE ?", "%"+strings.ToLower(code)+"%")
	}

	if len(duration) > 0 {
		chain = chain.Where("duration = ?", duration)
	}

	if len(autoLogin) > 0 {
		chain = chain.Where("auto_login = ?", autoLogin)
	}

	if len(used) > 0 {
		if used == "0" {
			chain = chain.Where("expires = 0")
		} else {
			chain = chain.Where("expires != 0")
		}
	}

	if len(reusable) > 0 {
		if reusable == "0" {
			chain = chain.Where("remain_use > 0")
		} else {
			chain = chain.Where("remain_use = -1")
		}
	}

	if len(printed) > 0 {
		chain = chain.Where("printed = ?", printed)
	}

	if len(voucherType) > 0 {
		chain = chain.Where("type = ?", voucherType)
	}

	if len(expiredStart) > 0 {
		t, _ := time.Parse(timeLayout, expiredStart)
		chain = chain.Where("expires >= ?", t)
	}

	if len(expiredEnd) > 0 {
		t, _ := time.Parse(timeLayout, expiredEnd)
		chain = chain.Where("expires <= ?", t)
	}

	if len(createdStart) > 0 {
		t, _ := time.Parse(timeLayout, createdStart)
		chain = chain.Where("created >= ?", t)
	}

	if len(createdEnd) > 0 {
		t, _ := time.Parse(timeLayout, createdEnd)
		chain = chain.Where("created <= ?", t)
	}

	chain.Order("created desc").Offset(offsets[0]).Limit(offsets[1]).Find(&hotspotVouchers)

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
