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
 * author: Matteo  Valentini <matteo.valentini@nethesis.it>
 */

package models

type WebHook struct {
	Status   bool             `json:"status"`
	Token    string           `json:"token,omitempty"`
	Hotspot  HotspotWebHook   `json:"hotspot"`
	Accounts []AccountWebHook `json:"accounts,omitempty"`
	Units    []UnitWebHook    `json:"units,omitempty"`
}

type HotspotWebHook struct {
	Id              int    `db:"id" json:"id"`
	Uuid            string `db:"uuid" json:"uuid"`
	Name            string `db:"name" json:"name"`
	Description     string `db:"description" json:"description"`
	BusinessName    string `db:"business_name" json:"business_name"`
	BusinessVAT     string `db:"business_vat" json:"business_vat"`
	BusinessAddress string `db:"business_address" json:"business_address"`
	BusinessEmail   string `db:"business_email" json:"business_email"`
}

func (HotspotWebHook) TableName() string {
	return "hotspots"
}

type AccountWebHook struct {
	Id       int    `db:"id" json:"id"`
	Uuid     string `db:"uuid" json:"uuid"`
	Type     string `db:"type" json:"type"`
	Name     string `db:"name" json:"name"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password_hash"`
	Email    string `db:"email" json:"email"`
}

func (AccountWebHook) TableName() string {
	return "accounts"
}

type UnitWebHook struct {
	Id          int    `db:"id" json:"id"`
	MacAddress  string `db:"mac_address" json:"mac_address"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Uuid        string `db:"uuid" json:"uuid"`
}

func (UnitWebHook) TableName() string {
	return "units"
}
