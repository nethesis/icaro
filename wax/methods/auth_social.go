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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nethesis/icaro/wax/utils"

	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/methods"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

func FacebookAuth(c *gin.Context) {
	code := c.Param("code")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")
	voucherCode := c.Query("voucher_code")

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
		days := utils.GetHotspotPreferencesByKey(unit.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)

		down := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Down")
		downInt, _ := strconv.Atoi(down.Value)
		up := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Up")
		upInt, _ := strconv.Atoi(up.Value)

		autoLogin := utils.GetHotspotPreferencesByKey(unit.HotspotId, "auto_login")
		autoLoginBool, _ := strconv.ParseBool(autoLogin.Value)

		// retrieve voucher
		if len(voucherCode) > 0 {
			voucher := utils.GetVoucherByCode(voucherCode, unit.HotspotId)

			if voucher.Id > 0 {
				daysInt = int(voucher.Expires.Sub(time.Now().UTC()).Hours() / 24)
				downInt = voucher.BandwidthDown
				upInt = voucher.BandwidthUp
				autoLoginBool = voucher.AutoLogin
			}
		}

		newUser := models.User{
			HotspotId:   unit.HotspotId,
			Name:        fbUserDetail.Name,
			Username:    fbInspectToken.Data.UserId,
			Password:    "",
			Email:       fbUserDetail.Email,
			AccountType: "facebook",
			KbpsDown:    downInt,
			KbpsUp:      upInt,
			AutoLogin:   autoLoginBool,
			ValidFrom:   time.Now().UTC(),
			ValidUntil:  time.Now().UTC().AddDate(0, 0, daysInt+1),
		}
		newUser.Id = methods.CreateUser(newUser)

		// create user session check
		utils.CreateUserSession(newUser.Id, sessionId)

		// create marketing info with user infos
		utils.CreateUserMarketing(newUser.Id, fbUserDetail, "facebook")
	} else {
		// update user info
		days := utils.GetHotspotPreferencesByKey(user.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, daysInt)
		db := database.Instance()
		db.Save(&user)

		// create user session check
		utils.CreateUserSession(user.Id, sessionId)

		// create marketing info with user infos
		utils.CreateUserMarketing(user.Id, fbUserDetail, "facebook")

		// retrieve voucher
		if len(voucherCode) > 0 {
			voucher := utils.GetVoucherByCode(voucherCode, user.HotspotId)

			if voucher.Id > 0 {
				user.ValidUntil = time.Now().UTC().AddDate(0, 0, voucher.Duration)
				user.KbpsDown = voucher.BandwidthDown
				user.KbpsUp = voucher.BandwidthUp
				user.AutoLogin = voucher.AutoLogin
			}
		}
	}

	// response to client
	c.JSON(http.StatusOK, gin.H{"user_id": fbInspectToken.Data.UserId, "user_db_id": user.Id})

}

func LinkedInAuth(c *gin.Context) {
	code := c.Param("code")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")
	voucherCode := c.Query("voucher_code")

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
	urlAPI := "https://api.linkedin.com/v1/people/~:(id,first-name,last-name,email-address,headline,current-share,num-connections,location,positions)?format=json"

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
		days := utils.GetHotspotPreferencesByKey(unit.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)

		down := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Down")
		downInt, _ := strconv.Atoi(down.Value)
		up := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Up")
		upInt, _ := strconv.Atoi(up.Value)

		autoLogin := utils.GetHotspotPreferencesByKey(unit.HotspotId, "auto_login")
		autoLoginBool, _ := strconv.ParseBool(autoLogin.Value)

		// retrieve voucher
		if len(voucherCode) > 0 {
			voucher := utils.GetVoucherByCode(voucherCode, unit.HotspotId)

			if voucher.Id > 0 {
				daysInt = int(voucher.Expires.Sub(time.Now().UTC()).Hours() / 24)
				downInt = voucher.BandwidthDown
				upInt = voucher.BandwidthUp
				autoLoginBool = voucher.AutoLogin
			}
		}

		newUser := models.User{
			HotspotId:   unit.HotspotId,
			Name:        liUserDetail.FirstName + " " + liUserDetail.LastName,
			Username:    liUserDetail.Id,
			Password:    "",
			Email:       liUserDetail.Email,
			AccountType: "linkedin",
			KbpsDown:    downInt,
			KbpsUp:      upInt,
			AutoLogin:   autoLoginBool,
			ValidFrom:   time.Now().UTC(),
			ValidUntil:  time.Now().UTC().AddDate(0, 0, daysInt+1),
		}
		newUser.Id = methods.CreateUser(newUser)

		// create user session check
		utils.CreateUserSession(newUser.Id, sessionId)

		// create marketing info with user infos
		utils.CreateUserMarketing(newUser.Id, liUserDetail, "linkedin")
	} else {
		// update user info
		days := utils.GetHotspotPreferencesByKey(user.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, daysInt)
		db := database.Instance()
		db.Save(&user)

		// create user session check
		utils.CreateUserSession(user.Id, sessionId)

		// create marketing info with user infos
		utils.CreateUserMarketing(user.Id, liUserDetail, "linkedin")

		// retrieve voucher
		if len(voucherCode) > 0 {
			voucher := utils.GetVoucherByCode(voucherCode, user.HotspotId)

			if voucher.Id > 0 {
				user.ValidUntil = time.Now().UTC().AddDate(0, 0, voucher.Duration)
				user.KbpsDown = voucher.BandwidthDown
				user.KbpsUp = voucher.BandwidthUp
				user.AutoLogin = voucher.AutoLogin
			}
		}
	}

	// response to client
	c.JSON(http.StatusOK, gin.H{"user_id": liUserDetail.Id, "user_db_id": user.Id})
}

func InstagramAuth(c *gin.Context) {
	code := c.Param("code")
	uuid := c.Query("uuid")
	sessionId := c.Query("sessionid")
	voucherCode := c.Query("voucher_code")

	clientId := configuration.Config.AuthSocial.Instagram.ClientId
	clientSecret := configuration.Config.AuthSocial.Instagram.ClientSecret
	redirectURI := configuration.Config.AuthSocial.Instagram.RedirectURI

	// retrieve access token
	urlEndpoint := "https://api.instagram.com/oauth/access_token"

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

	var inRespToken models.InstagramRespToken
	err = json.Unmarshal(body, &inRespToken)
	if err != nil {
		fmt.Println(err.Error())
	}

	// extract user info
	urlAPI := "https://api.instagram.com/v1/users/self/?access_token=" + inRespToken.AccessToken

	resp, err = http.Get(urlAPI)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	var inUserDetail models.InstagramUserDetail
	err = json.Unmarshal(body, &inUserDetail)
	if err != nil {
		fmt.Println(err.Error())
	}

	if inUserDetail.Data.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "access token is invalid"})
		return
	}

	// check if user exists
	user := utils.GetUserByUsername(inUserDetail.Data.Id)
	if user.Id == 0 {
		// create user
		unit := utils.GetUnitByUuid(uuid)
		days := utils.GetHotspotPreferencesByKey(unit.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)

		down := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Down")
		downInt, _ := strconv.Atoi(down.Value)
		up := utils.GetHotspotPreferencesByKey(unit.HotspotId, "CoovaChilli-Bandwidth-Max-Up")
		upInt, _ := strconv.Atoi(up.Value)

		autoLogin := utils.GetHotspotPreferencesByKey(unit.HotspotId, "auto_login")
		autoLoginBool, _ := strconv.ParseBool(autoLogin.Value)

		// retrieve voucher
		if len(voucherCode) > 0 {
			voucher := utils.GetVoucherByCode(voucherCode, unit.HotspotId)

			if voucher.Id > 0 {
				daysInt = int(voucher.Expires.Sub(time.Now().UTC()).Hours() / 24)
				downInt = voucher.BandwidthDown
				upInt = voucher.BandwidthUp
				autoLoginBool = voucher.AutoLogin
			}
		}

		newUser := models.User{
			HotspotId:   unit.HotspotId,
			Name:        inUserDetail.Data.FullName,
			Username:    inUserDetail.Data.Id,
			Password:    "",
			Email:       "",
			AccountType: "instagram",
			KbpsDown:    downInt,
			KbpsUp:      upInt,
			AutoLogin:   autoLoginBool,
			ValidFrom:   time.Now().UTC(),
			ValidUntil:  time.Now().UTC().AddDate(0, 0, daysInt+1),
		}
		newUser.Id = methods.CreateUser(newUser)

		// create user session check
		utils.CreateUserSession(newUser.Id, sessionId)

		// create marketing info with user infos
		utils.CreateUserMarketing(newUser.Id, inUserDetail.Data, "instagram")
	} else {
		// update user info
		days := utils.GetHotspotPreferencesByKey(user.HotspotId, "user_expiration_days")
		daysInt, _ := strconv.Atoi(days.Value)
		user.ValidUntil = time.Now().UTC().AddDate(0, 0, daysInt)
		db := database.Instance()
		db.Save(&user)

		// create user session check
		utils.CreateUserSession(user.Id, sessionId)

		// create marketing info with user infos
		utils.CreateUserMarketing(user.Id, inUserDetail.Data, "instagram")

		// retrieve voucher
		if len(voucherCode) > 0 {
			voucher := utils.GetVoucherByCode(voucherCode, user.HotspotId)

			if voucher.Id > 0 {
				user.ValidUntil = time.Now().UTC().AddDate(0, 0, voucher.Duration)
				user.KbpsDown = voucher.BandwidthDown
				user.KbpsUp = voucher.BandwidthUp
				user.AutoLogin = voucher.AutoLogin
			}
		}
	}

	// response to client
	c.JSON(http.StatusOK, gin.H{"user_id": inUserDetail.Data.Id, "user_db_id": user.Id})
}
