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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"wax/utils"

	"github.com/gin-gonic/gin"

	"sun-api/configuration"
	"sun-api/methods"
	"sun-api/models"
)

func GoogleAuth(c *gin.Context) {
	code := c.Param("code")
	uuid := c.Query("uuid")

	// retrieve access token
	url := "https://www.googleapis.com/oauth2/v3/tokeninfo?" +
		"&access_token=" + code

	fmt.Println("url", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var glRespToken models.GoogleRespToken
	err = json.Unmarshal(body, &glRespToken)
	if err != nil {
		fmt.Println(err.Error())
	}

	// extract user info
	url = "https://www.googleapis.com/plus/v1/people/me" +
		"?access_token=" + code

	fmt.Println("extract", url)

	resp, err = http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	var glUserDetail models.GoogleUserDetail
	err = json.Unmarshal(body, &glUserDetail)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("detail", glUserDetail)

	// create user
	unit := utils.ExtractUnit(uuid)
	user := models.User{
		HotspotId:   unit.HotspotId,
		Name:        glUserDetail.DisplayName,
		Username:    glUserDetail.Id,
		Email:       glRespToken.Email,
		AccountType: "google",
		KbpsDown:    0,
		KbpsUp:      0,
		ValidFrom:   time.Now().UTC(),
		ValidUntil:  time.Now().UTC().AddDate(0, 0, 30), // TODO: get days from hotspot account preferences
	}
	methods.CreateUser(user)

	// TODO: create marketing info with user likes and birthday

	// response to client
	c.JSON(http.StatusOK, gin.H{"user_id": glUserDetail.Id})
}

func FacebookAuth(c *gin.Context) {
	code := c.Param("code")
	uuid := c.Query("uuid")

	clientId := configuration.Config.AuthSocial.Facebook.ClientId
	clientSecret := configuration.Config.AuthSocial.Facebook.ClientSecret
	redirectURI := configuration.Config.AuthSocial.Facebook.RedirectURI

	// retrieve access token
	url := "https://graph.facebook.com/v2.11/oauth/access_token?" +
		"client_id=" + clientId +
		"&redirect_uri=" + redirectURI +
		"&client_secret=" + clientSecret +
		"&code=" + code

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var fbToken models.FacebookRespToken
	err = json.Unmarshal(body, &fbToken)
	if err != nil {
		fmt.Println(err.Error())
	}

	// retrive user id
	url = "https://graph.facebook.com/debug_token?" +
		"access_token=" + clientId + "|" + clientSecret +
		"&input_token=" + fbToken.AccessToken

	resp, err = http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	var fbInspectToken models.FacebookInspectToken
	err = json.Unmarshal(body, &fbInspectToken)
	if err != nil {
		fmt.Println(err.Error())
	}

	// extract user info
	url = "https://graph.facebook.com/" + fbInspectToken.Data.UserId + "?fields=" +
		"id,name,gender,email,birthday,likes" +
		"&access_token=" + fbToken.AccessToken

	resp, err = http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	var fbUserDetail models.FacebookUserDetail
	err = json.Unmarshal(body, &fbUserDetail)
	if err != nil {
		fmt.Println(err.Error())
	}

	// create user
	unit := utils.ExtractUnit(uuid)
	user := models.User{
		HotspotId:   unit.HotspotId,
		Name:        fbUserDetail.Name,
		Username:    fbInspectToken.Data.UserId,
		Email:       fbUserDetail.Email,
		AccountType: "facebook",
		KbpsDown:    0,
		KbpsUp:      0,
		ValidFrom:   time.Now().UTC(),
		ValidUntil:  time.Now().UTC().AddDate(0, 0, 30), // TODO: days info from hotspot account preferences
	}
	methods.CreateUser(user)

	// TODO: create marketing info with user likes and birthday

	// response to client
	c.JSON(http.StatusOK, gin.H{"user_id": fbInspectToken.Data.UserId})

}
