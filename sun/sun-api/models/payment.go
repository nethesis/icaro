/*
 * Copyright (C) 2018 Nethesis S.r.l.
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
 */

package models

import "time"

type Payment struct {
	ID        int       `db:"id" json:"id"`
	AccountId int       `db:"account_id" json:"account_id"`
	Payment   string    `db:"payment" json:"payment"`
	Created   time.Time `db:"created" json:"created"`
}

type PaypalPayment struct {
	Total           float64 `json:"total"`
	Subtotal        float64 `json:"subtotal"`
	Tax             float64 `json:"tax"`
	Currency        string  `json:"currency"`
	Item            string  `json:item`
	ItemDescription string  `json:"item_description"`
}
