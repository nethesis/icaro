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
 * author: Giacomo Sanchietti <giacomo.sanchietti@nethesis.it>
 */

package methods

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/wax/utils"
)

func AuthAccept(c *gin.Context, prefs string) {
	c.String(http.StatusOK, "Auth: 1\n"+prefs)
}

func AuthReject(c *gin.Context, description string) {
	message := "Auth: 0"

	if len(description) > 0 {
		message = message + " \nReply-Message: " + description
	}

	c.String(http.StatusForbidden, message)
	c.Abort()
}

func calculatePreferences(c *gin.Context, unit models.Unit, user models.User, timezone string) (bytes.Buffer, string) {
	errorMessage := ""

	// extract preferences
	prefs := utils.GetHotspotPreferencesByKeys(
		unit.HotspotId,
		[]string{"Idle-Timeout", "Acct-Session-Time"},
	)
	var outPrefs bytes.Buffer
	for _, pref := range prefs {
		outPrefs.WriteString(fmt.Sprintf("%s:%s\n", pref.Key, pref.Value))
	}

	outPrefs.WriteString(fmt.Sprintf("%s:%d\n", "CoovaChilli-Bandwidth-Max-Up", user.KbpsUp))
	outPrefs.WriteString(fmt.Sprintf("%s:%d\n", "CoovaChilli-Bandwidth-Max-Down", user.KbpsDown))

	sessionTimeout := utils.CalculateRemainTime(user, timezone)
	if sessionTimeout <= 0 {
		errorMessage = "max navigation time reached"

	}
	outPrefs.WriteString(fmt.Sprintf("%s:%d\n", "Session-Timeout", sessionTimeout))

	if user.MaxNavigationTraffic != 0 {
		remainTraffic := utils.CalculateRemainTraffic(user)

		if remainTraffic <= 0 {
			errorMessage = "max navigation traffic reached"
		}
		outPrefs.WriteString(fmt.Sprintf("%s:%d\n", "CoovaChilli-Max-Total-Octets", remainTraffic))
	}

	return outPrefs, errorMessage
}

func autoLogin(c *gin.Context, unitMacAddress string, username string, userMac string, sessionId string, timezone string) {
	// extract unit and user
	unit := utils.GetUnitByMacAddress(unitMacAddress)
	isValid, users := utils.GetUsersByMacAddressAndunitMacAddress(userMac, unitMacAddress)

	// check if user exists
	if !isValid {
		AuthReject(c, "user account not found")
		return
	}

	user := utils.FindAutoLoginUser(users)
	if user.Id <= 0 {
		AuthReject(c, "auto login not permitted")
		return
	}

	// check if user account is not expired
	if user.ValidUntil.Before(time.Now().UTC()) {
		AuthReject(c, "user account is expired")
		return
	}

	// check if current device is already logged in other units
	// check also if prop "bypass_macaddress_check" is enabled
	byPassMacAddress := utils.GetHotspotPreferencesByKey(unit.HotspotId, "bypass_macaddress_check")
	byPassMacAddressBool, _ := strconv.ParseBool(byPassMacAddress.Value)

	if(!byPassMacAddressBool) {
		flagOtherLogin := utils.CheckOtherUnitLogin(userMac, unit.Id, unit.HotspotId)
		if flagOtherLogin {
			AuthReject(c, "device already logged in an other unit")
			return
		}
	}

	// extract preferences
	outPrefs, errorMessage := calculatePreferences(c, unit, user, timezone)

	// response to dedalo
	if len(errorMessage) > 0 {
		AuthReject(c, errorMessage)
		return
	}
	AuthAccept(c, outPrefs.String())
}

func Login(c *gin.Context, unitMacAddress string, username string, userMac string, chapPass string, chapChal string, sessionId string, timezone string) {
	// check if unit exists
	unit := utils.GetUnitByMacAddress(unitMacAddress)
	if unit.Id <= 0 {
		AuthReject(c, "unit not found")
		return
	}

	// check if user exists
	user := utils.GetUserByUsernameAndHotspot(username, unit.HotspotId)
	if user.Id <= 0 {
		AuthReject(c, "user not found")
		return
	}

	// check if user-sessions exists
	if user.AccountType != "email" && user.AccountType != "sms" && user.AccountType != "voucher" {
		valid := utils.CheckUserSession(user.Id, sessionId)
		if !valid {
			AuthReject(c, "user-session not found")
			return
		}
	}

	// check if user credentials are valid
	if chapPass != utils.CalcUserDigest(user, chapChal) {
		AuthReject(c, "password mismatch")
		return
	}

	// check if user account is not expired
	if user.ValidUntil.Before(time.Now().UTC()) {
		AuthReject(c, "user account is expired")
		return
	}

	// check if hotspot belongs to a reseller with a valid plan
	if !utils.HotspotHasValidSubscription(user.HotspotId) {
		AuthReject(c, "reseller account is expired")
		return
	}

	// check if current device is already logged in other units
	// check also if prop "bypass_macaddress_check" is enabled
	byPassMacAddress := utils.GetHotspotPreferencesByKey(unit.HotspotId, "bypass_macaddress_check")
	byPassMacAddressBool, _ := strconv.ParseBool(byPassMacAddress.Value)

	if(!byPassMacAddressBool) {
		flagOtherLogin := utils.CheckOtherUnitLogin(userMac, unit.Id, unit.HotspotId)
		if flagOtherLogin {
			AuthReject(c, "device already logged in an other unit")
			return
		}
	}

	// extract preferences
	outPrefs, errorMessage := calculatePreferences(c, unit, user, timezone)

	// response to dedalo
	if len(errorMessage) > 0 {
		AuthReject(c, errorMessage)
		return
	}
	AuthAccept(c, outPrefs.String())

}

func Logins(c *gin.Context) {

	service := c.Query("service")
	switch service {
	case "framed":
		unitMacAddress := c.Query("ap")
		user := c.Query("user")
		userMac := c.Query("mac")
		sessionId := c.Query("sessionid")
		timezone := c.Query("timezone")
		autoLogin(c, unitMacAddress, user, userMac, sessionId, timezone)

	case "login":
		unitMacAddress := c.Query("ap")
		user := c.Query("user")
		userMac := c.Query("mac")
		chapPass := c.Query("chap_pass")
		chapChal := c.Query("chap_chal")
		sessionId := c.Query("sessionid")
		timezone := c.Query("timezone")
		Login(c, unitMacAddress, user, userMac, chapPass, chapChal, sessionId, timezone)

	default:
		c.String(http.StatusNotFound, "Invalid login service: '%s'", service)
	}
}
