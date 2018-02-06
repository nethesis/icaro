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

package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

var hotspot_preferences = map[string]string{
	"voucher_expiration_days": "1",
	"user_expiration_days":    "30",
	"temp_session_duration":   "300",
	"captive_title":           "This is a title",
	"captive_logo":            "",
	"captive_subtitle":        "This is a subtitle",
	"captive_banner":          "",
	"captive_description":     "This is a description",
	"facebook_login":          "true",
	"google_login":            "true",
	"linkedin_login":          "true",
	"sms_login":               "true",
	"email_login":             "true",
	"voucher_login":           "true",
	"sms_login_message":       "SMS Login code:",
}

func SetDefaultHotspotPreferences(hotspotId int) {
	db := database.Database()

	// iterate all default hotspot preferences
	for k, v := range hotspot_preferences {
		hsPreferences := models.HotspotPreference{
			HotspotId: hotspotId,
			Key:       k,
			Value:     v,
		}
		db.Save(&hsPreferences)
	}

	db.Close()
}

func OffsetCalc(page string, limit string) [2]int {
	var resLimit = 0
	var resOffset = 0

	if len(page) == 0 {
		page = "1"
	}
	if len(limit) == 0 {
		limit = "-1"
	}

	limitInt, errLimit := strconv.Atoi(limit)
	if errLimit != nil {
		fmt.Println(errLimit.Error())
	}

	pageInt, errPage := strconv.Atoi(page)
	if errPage != nil {
		fmt.Println(errPage.Error())
	}

	resLimit = limitInt
	resOffset = (pageInt - 1) * resLimit

	result := [2]int{resOffset, resLimit}
	return result
}

func ExtractToken(token string) models.AccessToken {
	var accessToken models.AccessToken
	db := database.Database()
	db.Where("token = ?", token).First(&accessToken)
	db.Close()

	return accessToken
}

func DeleteToken(token string) {
	var accessToken models.AccessToken
	db := database.Database()
	db.Where("token = ?", token).First(&accessToken)

	db.Delete(&accessToken)
	db.Close()
}

func RefreshToken(token string) {
	var accessToken models.AccessToken
	db := database.Database()
	db.Where("token = ?", token).First(&accessToken)

	// add 1 day to expiration date
	accessToken.Expires = time.Now().UTC().AddDate(0, 0, configuration.Config.TokenExpiresDays)
	db.Save(&accessToken)

	db.Close()
}

func ExtractHotspotIds(accountId int, admin bool, hotspotId int) []int {
	var hotspots []models.Hotspot
	account := GetAccountById(accountId)

	if account.Type == "customer" || account.Type == "desk" {
		accountId = account.CreatorId
	}

	db := database.Database()
	if admin {
		if hotspotId != 0 {
			db.Select("id").Where("id = ?", hotspotId).Find(&hotspots)
		} else {
			db.Select("id").Find(&hotspots)
		}
	} else {
		if hotspotId != 0 {
			db.Select("id").Where("account_id = ? AND id = ?", accountId, hotspotId).Find(&hotspots)
		} else {
			db.Select("id").Where("account_id = ?", accountId).Find(&hotspots)
		}
	}
	db.Close()

	result := []int{}

	for _, hotspot := range hotspots {
		result = append(result, hotspot.Id)
	}

	return result
}

func GetAccountById(id int) models.Account {
	var account models.Account
	db := database.Database()
	db.Where("id = ?", id).First(&account)
	db.Close()

	return account
}

func GetHotspotByName(name string) models.Hotspot {
	var hotspot models.Hotspot
	db := database.Database()
	db.Where("name = ?", name).First(&hotspot)
	db.Close()

	return hotspot
}

func Contains(intSlice []int, searchInt int) bool {
	for _, value := range intSlice {
		if value == searchInt {
			return true
		}
	}
	return false
}
