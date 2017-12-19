/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Windmill project.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with NethServer.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package methods

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"sun-api/configuration"
	"sun-api/database"
	"sun-api/models"
)

func Login(c *gin.Context) {
	var account models.Account

	username := c.PostForm("username")
	password := c.PostForm("password")

	db := database.Database()
	db.Where("username = ?", username).First(&account)

	if account.Id == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "No username found!"})
		return
	} else {
		// check password
		h := md5.New()
		io.WriteString(h, password)
		digest := fmt.Sprintf("%x", h.Sum(nil))

		if account.Password == digest {
			// create authorization token
			h := sha256.New()
			h.Write([]byte(time.Now().UTC().String()))
			token := fmt.Sprintf("%x", h.Sum(nil))

			// set expiration date
			expires := time.Now().UTC().AddDate(0, 0, configuration.Config.TokenExpiresDay)

			accessToken := models.AccessToken{
				AccountId: account.Id,
				Token:     token,
				Role:      account.Type,
				Expires:   expires,
			}

			db.Save(&accessToken)
			db.Close()

			c.JSON(http.StatusCreated, gin.H{"account_type": account.Type, "status": "success", "token": token, "expires": expires.String()})

		} else {
			db.Close()
			c.JSON(http.StatusNotFound, gin.H{"message": "Password is invalid"})
			return
		}
	}

}
