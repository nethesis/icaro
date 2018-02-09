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

package models

type WingsPrefs struct {
	HotspotId   int               `json:"hotspot_id"`
	HotspotName string            `json:"hotspot_name"`
	Preferences map[string]string `json:"preferences"`
	Socials     struct {
		FacebookClientId  string `json:"facebook_client_id"`
		GoogleClientId    string `json:"google_client_id"`
		LinkedInClientId  string `json:"linkedin_client_id"`
		InstagramClientId string `json:"instagram_client_id"`
	} `json:"socials"`
}
