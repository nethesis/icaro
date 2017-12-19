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
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"sun-api/database"
	"sun-api/models"
	"sun-api/utils"
)

func CreateAccount(c *gin.Context) {
	uuid := c.PostForm("uuid")
	typeField := c.PostForm("type")
	name := c.PostForm("name")
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	account := models.Account{
		Uuid:     uuid,
		Type:     typeField,
		Name:     name,
		Username: username,
		Password: password,
		Email:    email,
		Created:  time.Now(),
	}

	db := database.Database()
	db.Save(&account)

	db.Close()

	c.JSON(http.StatusCreated, gin.H{"id": account.Id, "status": "success"})
}

func UpdateAccount(c *gin.Context) {
	var account models.Account
	accountId := c.Param("account_id")

	name := c.PostForm("name")
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	db := database.Database()
	db.Where("id = ?", accountId).First(&account)

	if account.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No account found!"})
		return
	}

	account.Name = name
	account.Username = username
	account.Password = password
	account.Email = email

	db.Save(&account)

	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetAccounts(c *gin.Context) {
	var accounts []models.Account

	page := c.Query("page")
	limit := c.Query("limit")

	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Offset(offsets[0]).Limit(offsets[1]).Find(&accounts)

	if len(accounts) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No accounts found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, accounts)
}

func GetAccount(c *gin.Context) {
	var account models.Account
	accountId := c.Param("account_id")

	db := database.Database()
	db.Where("id = ?", accountId).First(&account)

	if account.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No account found!"})
		return
	}

	db.Close()

	c.JSON(http.StatusOK, account)
}

func DeleteAccount(c *gin.Context) {
	var account models.Account
	accountId := c.Param("account_id")

	db := database.Database()
	db.Where("id = ?", accountId).First(&account)

	if account.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No account found!"})
		return
	}

	db.Delete(&account)

	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
