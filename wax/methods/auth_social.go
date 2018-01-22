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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"github.com/nethesis/icaro/wax/utils"

	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/methods"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

func GoogleAuth(c *gin.Context) {
	code := c.Param("code")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")

	// retrieve access token
	url := "https://www.googleapis.com/oauth2/v3/tokeninfo?" +
		"&access_token=" + code

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

	if glUserDetail.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "access token is invalid"})
		return
	}

	// check if user exists
	user := utils.GetUserByUsername(glUserDetail.Id)
	if user.Id == 0 {
		// create user
		unit := utils.GetUnitByUuid(uuid)
		newUser := models.User{
			HotspotId:   unit.HotspotId,
			Name:        glUserDetail.DisplayName,
			Username:    glUserDetail.Id,
			Password:    "",
			Email:       glRespToken.Email,
			AccountType: "google",
			KbpsDown:    0,
			KbpsUp:      0,
			ValidFrom:   time.Now().UTC(),
			ValidUntil:  time.Now().UTC().AddDate(0, 0, 30), // TODO: get days from hotspot account preferences
		}
		newUser.Id = methods.CreateUser(newUser)

		// create user session check
		utils.CreateUserSession(newUser.Id, sessionId)

		// TODO: create marketing info with user infos and birthday
	} else {
		// update user info
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, 30) // TODO: days info from hotspot account preferences
		db := database.Database()
		db.Save(&user)
		db.Close()

		// create user session check
		utils.CreateUserSession(user.Id, sessionId)
	}

	// response to client
	c.JSON(http.StatusOK, gin.H{"user_id": glUserDetail.Id})
}

func FacebookAuth(c *gin.Context) {
	code := c.Param("code")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")

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

	if fbInspectToken.Data.UserId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "access token is invalid"})
		return
	}

	// check if user exists
	user := utils.GetUserByUsername(fbInspectToken.Data.UserId)
	if user.Id == 0 {
		// create user
		unit := utils.GetUnitByUuid(uuid)
		newUser := models.User{
			HotspotId:   unit.HotspotId,
			Name:        fbUserDetail.Name,
			Username:    fbInspectToken.Data.UserId,
			Password:    "",
			Email:       fbUserDetail.Email,
			AccountType: "facebook",
			KbpsDown:    0,
			KbpsUp:      0,
			ValidFrom:   time.Now().UTC(),
			ValidUntil:  time.Now().UTC().AddDate(0, 0, 30), // TODO: days info from hotspot account preferences
		}
		newUser.Id = methods.CreateUser(newUser)

		// create user session check
		utils.CreateUserSession(newUser.Id, sessionId)

		// TODO: create marketing info with user likes and birthday
	} else {
		// update user info
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, 30) // TODO: days info from hotspot account preferences
		db := database.Database()
		db.Save(&user)
		db.Close()

		// create user session check
		utils.CreateUserSession(user.Id, sessionId)
	}

	// response to client
	c.JSON(http.StatusOK, gin.H{"user_id": fbInspectToken.Data.UserId})

}

func LinkedInAuth(c *gin.Context) {
	code := c.Param("code")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")

	clientId := configuration.Config.AuthSocial.LinkedIn.ClientId
	clientSecret := configuration.Config.AuthSocial.LinkedIn.ClientSecret
	redirectURI := configuration.Config.AuthSocial.LinkedIn.RedirectURI

	// retrieve access token
	urlEndpoint := "https://www.linkedin.com/oauth/v2/accessToken"

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}
	bodyByte := bytes.NewBufferString(form.Encode())

	resp, err := http.Post(urlEndpoint, "application/x-www-form-urlencoded", bodyByte)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var liRespToken models.LinkedInRespToken
	err = json.Unmarshal(body, &liRespToken)
	if err != nil {
		fmt.Println(err.Error())
	}

	// extract user info
	urlAPI := "https://api.linkedin.com/v1/people/~:(id,first-name,last-name,email-address)?format=json"

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlAPI, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+liRespToken.AccessToken)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	var liUserDetail models.LinkedInUserDetail
	err = json.Unmarshal(body, &liUserDetail)
	if err != nil {
		fmt.Println(err.Error())
	}

	if liUserDetail.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "access token is invalid"})
		return
	}

	// check if user exists
	user := utils.GetUserByUsername(liUserDetail.Id)
	if user.Id == 0 {
		// create user
		unit := utils.GetUnitByUuid(uuid)
		newUser := models.User{
			HotspotId:   unit.HotspotId,
			Name:        liUserDetail.FirstName + " " + liUserDetail.LastName,
			Username:    liUserDetail.Id,
			Password:    "",
			Email:       liUserDetail.Email,
			AccountType: "linkedin",
			KbpsDown:    0,
			KbpsUp:      0,
			ValidFrom:   time.Now().UTC(),
			ValidUntil:  time.Now().UTC().AddDate(0, 0, 30), // TODO: get days from hotspot account preferences
		}
		newUser.Id = methods.CreateUser(newUser)

		// create user session check
		utils.CreateUserSession(newUser.Id, sessionId)

		// TODO: create marketing info with user infos and birthday
	} else {
		// update user info
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, 30) // TODO: days info from hotspot account preferences
		db := database.Database()
		db.Save(&user)
		db.Close()

		// create user session check
		utils.CreateUserSession(user.Id, sessionId)
	}

	// response to client
	c.JSON(http.StatusOK, gin.H{"user_id": liUserDetail.Id})
}
