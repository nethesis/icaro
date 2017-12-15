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

package models

type User struct {
	Id          int    `db:"id" json:"id"`
	HotspotId   int    `db:"hotspot_id" json:"hotspot"`
	FirstName   string `db:"firstname" json:"firstname"`
	LastName    string `db:"lastname" json:"lastname"`
	UserName    string `db:"username" json:"username"`
	Email       string `db:"email" json:"email"`
	AccountType string `db:"account_type" json:"type"`
	KbpsDown    int    `db:"kbps_down" json:"kbps_down"`
	KbpsUp      int    `db:"kbps_up" json:"kbps_up"`
	ValidFrom   string `db:"valid_from" json:"valid_from"`
	ValidUntil  string `db:"valid_until" json:"valid_until"`
	Created     string `db:"created" json:"created"`
}
