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

type Endpoints struct {
	Sms struct {
		AccountSid     string `json:"account_sid"`
		AuthToken      string `json:"auth_token"`
		ServiceSid     string `json:"service_sid"`
		Link           string `json:"link"`
		SendQuotaAlert bool   `json:"send_quota_alert"`
	} `json:"sms"`
	Whatsapp struct {
		Number         string `json:"number"`
		AccountSid     string `json:"account_sid"`
		AuthToken      string `json:"auth_token"`
		ServiceSid     string `json:"service_sid"`
		Link           string `json:"link"`
		SendQuotaAlert bool   `json:"send_quota_alert"`
	} `json:"whatsapp"`
	Email struct {
		From         string `json:"from"`
		FromName     string `json:from_name`
		SMTPHost     string `json:"smtp_host"`
		SMTPPort     int    `json:"smtp_port"`
		SMTPUser     string `json:"smtp_user"`
		SMTPPassword string `json:"smtp_password"`
		Link         string `json:"link"`
	} `json:"email"`
}
