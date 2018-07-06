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
 */

package methods

import (
	"net/http"
	"time"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/middleware"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func GetSubscriptionPlans(c *gin.Context) {
	var subscriptionPlans []models.SubscriptionPlan

	db := database.Instance()
	db.Where("code != 'premium'").Find(&subscriptionPlans)

	if len(subscriptionPlans) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No subscription plans found!"})
		return
	}

	c.JSON(http.StatusOK, subscriptionPlans)
}

func RenewalPlan(c *gin.Context) {
	var account models.Account
	var subscription models.Subscription
	var accountSMS models.AccountSmsCount

	accountId := c.MustGet("token").(models.AccessToken).AccountId

	var json models.SubscriptionRenewalJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()

	// get account
	db.Where("id = ?", accountId).First(&account)
	if account.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no account found!"})
		return
	}

	// get subscription
	db.Preload("SubscriptionPlan").Where("account_id = ?", account.Id).First(&subscription)

	// check payment
	if middleware.PaymentCheck(json.PaymentId, subscription.SubscriptionPlan.Code, account.Uuid) {
		// update subscription
		subscription.ValidFrom = time.Now().UTC()
		subscription.ValidUntil = subscription.ValidUntil.AddDate(0, 0, subscription.SubscriptionPlan.Period)

		// update sms count
		db.Where("account_id = ?", accountId).First(&accountSMS)
		accountSMS.SmsMaxCount = accountSMS.SmsMaxCount + subscription.SubscriptionPlan.IncludedSMS
		db.Save(&accountSMS)

		// update subscription info
		if err := db.Save(&subscription).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "subscription plan not renewed", "error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "no payment related to this plan for this server found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UpgradePlanPrice(c *gin.Context) {
	var account models.Account
	var subscription models.Subscription
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	plan := c.Query("plan")

	newSubuscriptionPlan := utils.GetSubscriptionPlanByCode(plan)

	db := database.Instance()

	// get account
	db.Where("id = ?", accountId).First(&account)
	if account.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no account found!"})
		return
	}

	// get subscription
	db.Preload("SubscriptionPlan").Where("account_id = ?", account.Id).First(&subscription)

	// calculate discount upgrade
	daysDiff := subscription.ValidUntil.Sub(time.Now().UTC())
	discount := (daysDiff.Hours() / 24) * subscription.SubscriptionPlan.Price / float64(subscription.SubscriptionPlan.Period)
	finalPrice := newSubuscriptionPlan.Price - discount

	c.JSON(http.StatusOK, gin.H{"discount": discount, "full_price": newSubuscriptionPlan.Price, "price": utils.Round(finalPrice, 0.5, 2), "name": newSubuscriptionPlan.Code})
}

func UpgradePlan(c *gin.Context) {
	var account models.Account
	var subscription models.Subscription
	var accountSMS models.AccountSmsCount

	accountId := c.MustGet("token").(models.AccessToken).AccountId

	var json models.SubscriptionUpgradeJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()

	// get account
	db.Where("id = ?", accountId).First(&account)
	if account.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no account found!"})
		return
	}

	// get subscription
	db.Preload("SubscriptionPlan").Where("account_id = ?", account.Id).First(&subscription)

	if json.SubscriptionPlanId != 0 && json.SubscriptionPlanId != subscription.SubscriptionPlanId {
		// get subscription using id
		newSubscriptionPlan := utils.GetSubscriptionPlanById(json.SubscriptionPlanId)

		// check payment
		if middleware.PaymentCheck(json.PaymentId, newSubscriptionPlan.Code, account.Uuid) {
			// update subscription
			subscription.SubscriptionPlan = newSubscriptionPlan
			subscription.ValidFrom = time.Now().UTC()
			subscription.ValidUntil = time.Now().UTC().AddDate(0, 0, newSubscriptionPlan.Period)

			// update sms count
			db.Where("account_id = ?", accountId).First(&accountSMS)
			accountSMS.SmsMaxCount = accountSMS.SmsMaxCount + newSubscriptionPlan.IncludedSMS
			db.Save(&accountSMS)

			// update subscription info
			if err := db.Save(&subscription).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "subscription plan not updated", "error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "no payment related to this plan for this server found"})
			return
		}
	} else {
		c.JSON(http.StatusConflict, gin.H{"message": "this plan is already associated with this server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
