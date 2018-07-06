/*
 * Copyright (C) 2018 Nethesis S.r.l.
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

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

/*
  Alghoritm for VAT, this applies to EU companies selling services

  - customer is a UE company: no VAT
  - customer is a UE private: VAT from country of selling company

  - customer is a non-UE (Others countries) company or private: no VAT

  - customer is private or company from selling company country: VAT from country of selling company
*/
func GetVatPercentage(customerCountry string, customerVat string) int {
	var tax models.Tax

	// Customer is a company not is the same country of billing, no VAT applied
	if customerVat != "" && !(configuration.Config.Billing.Country == customerCountry) {
		return 0
	}

	db := database.Instance()
	db.Where("country = ?", customerCountry).First(&tax)

	// Customer is from non-UE countries, no VAT applied
	if tax.Country == "Other" {
		return 0
	}

	// Customer is a UE pyhsical person: VAT from company which deployed the application
	db.Where("country = ?", configuration.Config.Billing.Country).First(&tax)
	return tax.Percentage
}

func GetBilling(c *gin.Context) {
	var billing models.Billing
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	db := database.Instance()
	db.Where("account_id = ?", accountId).First(&billing)

	if billing.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no billing information found!"})
		return
	}

	billing.Tax = GetVatPercentage(billing.Country, billing.Vat)

	c.JSON(http.StatusOK, billing)
}

func CreateBilling(c *gin.Context) {
	var json models.BillingJSON
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	billing := models.Billing{
		AccountId:  accountId,
		Name:       json.Name,
		Address:    json.Address,
		Country:    json.Country,
		City:       json.City,
		PostalCode: json.PostalCode,
		Vat:        json.Vat,
	}

	db := database.Instance()
	if err := db.Create(&billing).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "billing not saved", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UpdateBilling(c *gin.Context) {
	var json models.BillingJSON
	var billing models.Billing
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()
	db.Where("account_id = ?", accountId).First(&billing)

	if billing.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no billing found!"})
		return
	}

	billing.AccountId = accountId
	billing.Name = json.Name
	billing.Address = json.Address
	billing.Country = json.Country
	billing.City = json.City
	billing.PostalCode = json.PostalCode
	billing.Vat = json.Vat

	if err := db.Save(&billing).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "billing not saved", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
