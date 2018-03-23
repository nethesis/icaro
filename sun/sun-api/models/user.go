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

package models

import "time"

type User struct {
	Id          int       `db:"id" json:"id"`
	HotspotId   int       `db:"hotspot_id" json:"hotspot_id"`
	Name        string    `db:"name" json:"name"`
	Username    string    `db:"username" json:"username"`
	Password    string    `db:"password" json:"-"`
	Email       string    `db:"email" json:"email"`
	AccountType string    `db:"account_type" json:"account_type"`
	KbpsDown    int       `db:"kbps_down" json:"kbps_down"`
	KbpsUp      int       `db:"kbps_up" json:"kbps_up"`
	AutoLogin   bool      `db:"auto_login" json:"auto_login"`
	ValidFrom   time.Time `db:"valid_from" json:"valid_from"`
	ValidUntil  time.Time `db:"valid_until" json:"valid_until"`
	Created     time.Time `db:"created" json:"created"`
}
