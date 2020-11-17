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

type HotspotVoucher struct {
	Id            int       `db:"id" json:"id"`
	HotspotId     int       `db:"hotspot_id" json:"hotspot_id"`
	Code          string    `db:"code" json:"code"`
	AutoLogin     bool      `db:"auto_login" json:"auto_login"`
	BandwidthUp   int       `db:"bandwidth_up" json:"bandwidth_up"`
	BandwidthDown int       `db:"bandwidth_down" json:"bandwidth_down"`
	Duration      int       `db:"duration" json:"duration"`
	MaxTraffic    int64     `db:"max_traffic" json:"max_traffic"`
	MaxTime       int       `db:"max_time" json:"max_time"`
	RemainUse     int       `db:"remain_use" json:"remain_use"`
	Expires       time.Time `db:"expires" json:"expires"`
	Type          string    `db:"type" json:"type"`
	UserName      string    `db:"user_name" json:"user_name"`
	UserMail      string    `db:"user_mail" json:"user_mail"`
	Printed       bool      `db:"printed" json:"printed"`
	OwnerId       int       `db:"owner_id" json:"owner_id"`
	Created       time.Time `db:"created" json:"created"`
}

type HotspotVoucherJSON struct {
	Id            int       `db:"id" json:"id"`
	HotspotId     int       `db:"hotspot_id" json:"hotspot_id"`
	Code          string    `db:"code" json:"code"`
	AutoLogin     bool      `db:"auto_login" json:"auto_login"`
	BandwidthUp   int       `db:"bandwidth_up" json:"bandwidth_up"`
	BandwidthDown int       `db:"bandwidth_down" json:"bandwidth_down"`
	Time          string    `json:"time"`
	Duration      int       `db:"duration" json:"duration"`
	Expiration    int       `json:"expiration"`
	MaxTraffic    int64     `db:"max_traffic" json:"max_traffic"`
	MaxTime       int       `db:"max_time" json:"max_time"`
	RemainUse     int       `db:"remain_use" json:"remain_use"`
	Expires       time.Time `db:"expires" json:"expires"`
	Type          string    `db:"type" json:"type"`
	UserName      string    `db:"user_name" json:"user_name"`
	UserMail      string    `db:"user_mail" json:"user_mail"`
	Printed       bool      `db:"printed" json:"printed"`
	OwnerId       int       `db:"owner_id" json:"owner_id"`
	Created       time.Time `db:"created" json:"created"`
	NumVouchers   int       `json:"num_vouchers"`
	VoucherIds    []int     `json:"voucher_ids"`
}
