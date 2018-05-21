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

type HotspotSmsCount struct {
	Id        int       `db:"id" json:"id"`
	HotspotId int       `db:"hotspot_id" json:"hotspot_id"`
	UniId     int       `db:"unit_id" json:"unit_id"`
	Number    string    `db:"number" json:"number"`
	Reset     bool      `db:"reset" json:"reset"`
	Sent      time.Time `db:"sent" json:"sent"`
}
