/*
 * Copyright (C) 2020 Nethesis S.r.l.
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
 */

package methods

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

func GetAccountDisclaimers(c *gin.Context) {
	var disclaimers []models.Disclaimer
	accountId := c.MustGet("token").(models.AccessToken).AccountId
	db := database.Instance()
	queryAccountId := accountId

	if accountId == 1 {
		queryAccountId, _ = strconv.Atoi(c.Param("account_id"))
	}

	selectQuery := "disclaimers.*"
	joinQuery := "JOIN disclaimers_accounts on disclaimers_accounts.disclaimer_id = disclaimers.id"
	db.Select(selectQuery).Joins(joinQuery).Where("disclaimers_accounts.account_id = ?", queryAccountId).Find(&disclaimers)

	if len(disclaimers) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No disclaimer found"})
		return
	}
	c.JSON(http.StatusOK, disclaimers)
}

func DeleteDisclaimer(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	if accountId != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation not permitted!"})
		return
	}

	var disclaimer models.Disclaimer
	disclaimerId := c.Param("disclaimer_id")
	db := database.Instance()
	db.Where("id = ?", disclaimerId).First(&disclaimer)

	if disclaimer.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No disclaimer found!"})
		return
	}

	db.Delete(&disclaimer)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func CreateAccountDisclaimer(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	if accountId != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation not permitted!"})
		return
	}

	var json models.Disclaimer
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	// create disclaimer

	disclaimer := models.Disclaimer{
		Type:  json.Type,
		Title: json.Title,
		Body:  json.Body,
	}

	db := database.Instance()
	db.Save(&disclaimer)

	if disclaimer.Id == 0 {
		c.JSON(http.StatusConflict, gin.H{"disclaimer_id": disclaimer.Id, "status": "disclaimer already exists"})
	}

	// associate disclaimer to account

	accountId, _ = strconv.Atoi(c.Param("account_id"))

	disclaimersAccounts := models.DisclaimersAccounts{
		DisclaimerId: disclaimer.Id,
		AccountId:    accountId,
	}

	db.Save(&disclaimersAccounts)

	if disclaimersAccounts.Id == 0 {
		c.JSON(http.StatusConflict, gin.H{"disclaimers_accounts_id": disclaimersAccounts.Id, "status": "disclaimers_accounts already exists"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"disclaimer_id": disclaimer.Id, "disclaimers_accounts_id": disclaimersAccounts.Id, "status": "success"})
	}
}
