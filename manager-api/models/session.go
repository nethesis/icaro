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

import "time"

type Session struct {
	Id         int       `db:"id" json:"id"`
	UnitId     string    `db:"unit_id" json:"unit"`
	HotspotId  int       `db:"hotspot_id" json:"hotspot"`
	DeviceId   string    `db:"device_id" json:"device"`
	UserId     string    `db:"user_id" json:"user"`
	BytesUp    int       `db:"bytes_up" json:"bytes_up"`
	BytesDown  int       `db:"bytes_down" json:"bytes_down"`
	Duration   int       `db:"duration" json:"duration"`
	AuthTime   time.Time `db:"auth_time" json:"auth_time"`
	StartTime  time.Time `db:"start_time" json:"start_time"`
	UpdateTime time.Time `db:"update_time" json:"update_time"`
	StopTime   time.Time `db:"stop_time" json:"stop_time"`
	SessionKey string    `db:"session_key" json:"session_key"`
}
