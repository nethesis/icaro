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

type AccountWhatsappCount struct {
	Id                int `db:"id" json:"id"`
	AccountId         int `db:"account_id" json:"account_id"`
	WhatsappMaxCount  int `db:"whatsapp_max_count" json:"whatsapp_max_count"`
	WhatsappCount     int `db:"whatsapp_count" json:"whatsapp_count"`
	WhatsappThreshold int `db:"whatsapp_threshold" json:"whatsapp_threshold"`
}

type AccountWhatsappCountJSON struct {
	WhatsappToAdd int `json:"whatsapp_to_add"`
}

type AccountWhatsappThresholdJSON struct {
	WhatsappThreshold int `json:"whatsapp_threshold"`
}
