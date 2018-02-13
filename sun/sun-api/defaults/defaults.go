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
	"sms_login_message":       "Your Icaro SMS login code:",
	"captive_1_redir":         "https://nethesis.github.io/icaro",
	"captive_2_title":         "This is a title",
	"captive_3_logo":          "",
	"captive_4_subtitle":      "This is a subtitle",
	"captive_5_banner":        "",
	"captive_6_description":   "This is a description",
	"captive_7_background":    "#2a87be",
	"facebook_login":          "true",
	"linkedin_login":          "true",
	"instagram_login":         "true",
	"sms_login":               "true",
	"email_login":             "true",
	"voucher_login":           "true",
}

var SMSMaxCount = 500
