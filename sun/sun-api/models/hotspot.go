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

type Hotspot struct {
	Id               int       `db:"id" json:"id"`
	Uuid             string    `db:"uuid" json:"uuid"`
	AccountId        int       `db:"account_id" json:"account_id"`
	Name             string    `db:"name" json:"name"`
	Description      string    `db:"description" json:"description"`
	BusinessName     string    `db:"business_name" json:"business_name"`
	BusinessVAT      string    `db:"business_vat" json:"business_vat"`
	BusinessAddress  string    `db:"business_address" json:"business_address"`
	BusinessEmail    string    `db:"business_email" json:"business_email"`
	BusinessDPO      string    `db:"business_dpo" json:"business_dpo"`
	BusinessDPOMail  string    `db:"business_dpo_mail" json:"business_dpo_mail"`
	IntegrationTerms string    `db:"integration_terms" json:"integration_terms"`
	Created          time.Time `db:"created" json:"created"`

	Account Account `gorm:"PRELOAD:false json:"account"`
}

type HotspotJSON struct {
	Id               int       `db:"id" json:"id"`
	Uuid             string    `db:"uuid" json:"uuid"`
	AccountId        int       `db:"account_id" json:"account_id"`
	Name             string    `db:"name" json:"name"`
	Description      string    `db:"description" json:"description"`
	BusinessName     string    `db:"business_name" json:"business_name"`
	BusinessVAT      string    `db:"business_vat" json:"business_vat"`
	BusinessAddress  string    `db:"business_address" json:"business_address"`
	BusinessEmail    string    `db:"business_email" json:"business_email"`
	BusinessDPO      string    `db:"business_dpo" json:"business_dpo"`
	BusinessDPOMail  string    `db:"business_dpo_mail" json:"business_dpo_mail"`
	IntegrationTerms string    `db:"integration_terms" json:"integration_terms"`
	Created          time.Time `db:"created" json:"created"`
	AccountName      string    `db:"account_name" json:"account_name"`
}

func (HotspotJSON) TableName() string {
	return "hotspots"
}
