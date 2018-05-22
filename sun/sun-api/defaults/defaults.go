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

package defaults

var HotspotPreferences = map[string]string{
	"user_expiration_days":            "30",
	"temp_session_duration":           "300",
	"facebook_login":                  "true",
	"facebook_login_page":             "",
	"linkedin_login":                  "true",
	"instagram_login":                 "true",
	"sms_login":                       "true",
	"email_login":                     "true",
	"voucher_login":                   "false",
	"auto_login":                      "true",
	"CoovaChilli-Bandwidth-Max-Down":  "0",
	"CoovaChilli-Bandwidth-Max-Up":    "0",
	"CoovaChilli-Max-Total-Octets":    "0",
	"CoovaChilli-Max-Navigation-Time": "0",
}

var SMSMaxCount = 500
