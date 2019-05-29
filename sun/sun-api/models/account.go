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

type Account struct {
	Id        int       `db:"id" json:"id"`
	CreatorId int       `db:"creator_id" json:"creator_id"`
	Uuid      string    `db:"uuid" json:"uuid"`
	Type      string    `db:"type" json:"type"`
	Name      string    `db:"name" json:"name"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"-"`
	Email     string    `db:"email" json:"email"`
	Created   time.Time `db:"created" json:"created"`
}

type AccountJSON struct {
	Id                   int           `json:"id"`
	CreatorId            int           `json:"creator_id"`
	Uuid                 string        `json:"uuid"`
	Type                 string        `json:"type"`
	Name                 string        `json:"name"`
	Username             string        `json:"username"`
	Password             string        `json:"password"`
	Email                string        `json:"email"`
	Created              time.Time     `json:"created"`
	HotspotId            int           `json:"hotspot_id"`
	HotspotName          string        `json:"hotspot_name"`
	SubscriptionPlanId   int           `json:"subscription_plan_id"`
	SubscriptionPlanName string        `json:"subscription_plan_name"`
	Subscription         Subscription  `json:"subscription"`
	Tokens               []AccessToken `json:"tokens"`
}

func (AccountJSON) TableName() string {
	return "accounts"
}
