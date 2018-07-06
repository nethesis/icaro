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
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func Login(c *gin.Context) {
	var account models.Account
	var subscription models.Subscription

	var json models.Login
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	username := json.Username
	password := json.Password

	db := database.Instance()
	db.Where("username = ?", username).First(&account)

	if account.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No username found!"})
		return
	} else {
		// check password
		h := md5.New()
		h.Write([]byte(password))
		digest := fmt.Sprintf("%x", h.Sum(nil))

		if account.Password == digest {
			// create authorization token
			h := sha256.New()
			h.Write([]byte(time.Now().UTC().String() + username + password))
			token := fmt.Sprintf("%x", h.Sum(nil))

			// set expiration date
			expires := time.Now().UTC().AddDate(0, 0, configuration.Config.TokenExpiresDays)

			accessToken := models.AccessToken{
				AccountId:   account.Id,
				Token:       token,
				Role:        account.Type,
				Type:        "login",
				Expires:     expires,
				ACLs:        "full",
				Description: "",
			}

			db.Save(&accessToken)

			db.Set("gorm:auto_preload", true)
			if account.Type == "reseller" {
				db.Preload("SubscriptionPlan").Where("account_id = ?", account.Id).First(&subscription)
			} else {
				db.Preload("SubscriptionPlan").Where("account_id = ?", account.CreatorId).First(&subscription)
			}
			subscription.Expired = subscription.ValidUntil.Before(time.Now().UTC())

			c.JSON(http.StatusCreated, gin.H{
				"account_type": account.Type,
				"status":       "success",
				"token":        token,
				"expires":      expires.String(),
				"id":           account.Id,
				"subscription": subscription,
			})

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Password is invalid"})
		}
	}

}

func LoginAuth0(c *gin.Context) {
	var subscription models.Subscription

	// define api endpoint and audience
	AUTH0_DOMAIN := "https://" + configuration.Config.Auth0.Domain + "/"
	JWKS_URI := "https://" + configuration.Config.Auth0.Domain + "/.well-known/jwks.json"
	AUDIENCE := []string{configuration.Config.Auth0.Audience}

	// create client configuration instance to check jwt
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: JWKS_URI}, nil)
	configAuth0 := auth0.NewConfiguration(client, AUDIENCE, AUTH0_DOMAIN, jose.RS256)
	validator := auth0.NewValidator(configAuth0, nil)

	// check jwt validation
	token, err := validator.ValidateRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing or invalid token"})
		return
	} else {
		// extract claims from token
		claims := map[string]interface{}{}
		err := validator.Claims(c.Request, token, &claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Claims extraction failed"})
			return
		}

		// check if account exists
		account := utils.GetAccountOrCreate(claims)

		// create authorization token
		h := sha256.New()
		h.Write([]byte(time.Now().UTC().String() + account.Uuid))
		token := fmt.Sprintf("%x", h.Sum(nil))

		// set expiration date
		expires := time.Now().UTC().AddDate(0, 0, configuration.Config.TokenExpiresDays)

		accessToken := models.AccessToken{
			AccountId: account.Id,
			Token:     token,
			Role:      account.Type,
			Expires:   expires,
		}

		db := database.Instance()
		db.Save(&accessToken)

		db.Set("gorm:auto_preload", true)
		if account.Type == "reseller" {
			db.Preload("SubscriptionPlan").Where("account_id = ?", account.Id).First(&subscription)
		} else {
			db.Preload("SubscriptionPlan").Where("account_id = ?", account.CreatorId).First(&subscription)
		}
		subscription.Expired = subscription.ValidUntil.Before(time.Now().UTC())

		c.JSON(http.StatusCreated, gin.H{
			"account_type": account.Type,
			"status":       "success",
			"token":        token,
			"expires":      expires.String(),
			"id":           account.Id,
			"subscription": subscription,
		})
	}
}

func Logout(c *gin.Context) {
	token := c.GetHeader("Token")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
		return
	} else {
		// delete token
		utils.DeleteToken(token)

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
