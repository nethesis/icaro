/*
 * Copyright (C) 2019 Nethesis S.r.l.
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
 * author: Matteo Valentinit <matteo.valentini@nethesis.it>
 */

package models

import "time"

type AdeToken struct {
	Id               int       `db:"id" json:"id"`
	Token            string    `db:"token" json:"token"`
	Type             string    `db:"type" json:"type"`
	Contact          string    `db:"contact" json:"contact"`
	FeedbackSentTime time.Time `db:"feedback_sent_time" json:"feedback_sent_time"`
	FeedbackSendTime time.Time `db:"feedback_send_time" json:"feedback_send_time"`
	ReviewSentTime   time.Time `db:"review_sent_time" json:"review_sent_time"`
	ReviewSendTime   time.Time `db:"review_send_time" json:"review_send_time"`
	UserId           int       `db:"user_id" json:"user_id"`
	HotspotId        int       `db:"hotspot_id" json:"hotspot_id"`
	Created          time.Time `db:"created" json:"created"`
}
