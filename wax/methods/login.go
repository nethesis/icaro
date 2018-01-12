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
 * author: Giacomo Sanchietti <giacomo.sanchietti@nethesis.it>
 */

package methods

import (
	"net/http"
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"

	"wax/utils"
)


func AuthAccept(c *gin.Context, prefs string) {
	fmt.Println("Auth: 1\n" + prefs)
	c.String(http.StatusOK, "Auth: 1\n" + prefs)
}

func AuthReject(c *gin.Context) {
	c.String(http.StatusForbidden, "Auth: 0")
}

func Login(c *gin.Context, unitMacAddress string, user string, password string, challenge string) {
	//TODO
	/*
		Social:
		- check the user exists on db
		- check valid until field isn't expired
		- validate with session id

		SMS/Email:
		- check the user exists on db
		- check user/password
		- validate with session id

		Example:
		service=login&user=aaa&chap_chal=c49f036173dbc004943a5a0aaabdbb96&chap_pass=47bce5c74f589f4867dbd57e9ca9f808&chap_id=0&ap=00-0D-B9-41-7C-F8&mac=88-C9-D0-FD-52-93&ip=10.1.0.2&sessionid=151318020000000003&nasid=hs-test
	*/
	unit := utils.GetUnitByMacAddress(unitMacAddress)
	if (unit.Id <= 0) {
		AuthReject(c)
	}
	prefs := utils.GetHotspotPreferencesByKeys(unit.HotspotId, []string{"Idle-Timeout", "Acct-Session-Time", "Session-Timeout", "CoovaChilli-Bandwidth-Max-Up", "CoovaChilli-Bandwidth-Max-Down"})
	var outPrefs bytes.Buffer
	for _, pref := range prefs {
		outPrefs.WriteString(fmt.Sprintf("%s:%s\n",pref.Key, pref.Value))
	}

	AuthAccept(c, outPrefs.String())


}
