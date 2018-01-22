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
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func Login(c *gin.Context) {
	var account models.Account

	var json models.Login
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	username := json.Username
	password := json.Password

	db := database.Database()
	db.Where("username = ?", username).First(&account)

	if account.Id == 0 {
		db.Close()
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
				AccountId: account.Id,
				Token:     token,
				Role:      account.Type,
				Expires:   expires,
			}

			db.Save(&accessToken)
			db.Close()

			c.JSON(http.StatusCreated, gin.H{
				"account_type": account.Type,
				"status":       "success",
				"token":        token,
				"expires":      expires.String(),
				"id":           account.Id,
			})

		} else {
			db.Close()
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Password is invalid"})
		}
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
