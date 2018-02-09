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

package defaults

var HotspotPreferences = map[string]string{
	"voucher_expiration_days": "1",
	"user_expiration_days":    "30",
	"temp_session_duration":   "300",
	"captive_title":           "This is a title",
	"captive_logo":            "",
	"captive_subtitle":        "This is a subtitle",
	"captive_banner":          "",
	"captive_description":     "This is a description",
	"captive_redir":           "https://nethesis.github.io/icaro",
	"facebook_login":          "true",
	"google_login":            "true",
	"linkedin_login":          "true",
	"sms_login":               "true",
	"email_login":             "true",
	"voucher_login":           "true",
	"sms_login_message":       "Your Icaro SMS login code:",
}

var SMSMaxCount = 500
