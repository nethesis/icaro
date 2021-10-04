/*
 * Copyright (C) 2021 Nethesis S.r.l.
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

type DaemonAuth struct {
	Id             int    `db:"id" json:"id"`
	SessionId      string `db:"session_id" json:"session_id"`
	SessionTimeout int `db:"session_timeout" json:"session_timeout"`
	UnitUuid       string `db:"unit_uuid" json:"unit_uuid"`
	UserId         int    `db:"user_id" json:"user_id"`
	Username       string `db:"username" json:"username"`
	Password       string `db:"password" json:"password"`
	Type           string `db:"type" json:"type"`
}
