/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Dartagnan project.
 *
 * Dartagnan is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Dartagnan is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dartagnan.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package models

import (
	"time"
)

type SubscriptionPlan struct {
	ID          int     `db:"id" json:"id"`
	Code        string  `db:"code" json:"code"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Price       float64 `db:"price" json:"price"`
	Period      int     `db:"period" json:"period"`
	IncludedSMS         int  `db:"included_sms" json:"included_sms"`
	MaxUnits            int  `db:"max_units" json:"max_units"`
	AdvancedReport      bool `db:"advanced_report" json:"advanced_report"`
	WingsCustomization  bool `db:"wings_customization" json:"wings_customization"`
	SocialAnalytics     bool `db:"social_analytics" json:"social_analytics"`
}

type Subscription struct {
	ID                   int       `db:"id" json:"id"`
	UserID               string    `db:"account_id" json:"account_id"`
	ValidFrom            time.Time `db:"valid_from" json:"valid_from"`
	ValidUntil           time.Time `db:"valid_until" json:"valid_until"`
	Created              time.Time `db:"created" json:"created"`
        AccountID            int       `db:"account_id" json:"account_id"`

	SubscriptionPlanID   int              `db:"subscription_plan_id" json:"-"`
	SubscriptionPlan     SubscriptionPlan `json:"subscription_plan"`
}

