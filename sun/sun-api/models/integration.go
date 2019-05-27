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
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package models

type Integration struct {
	Id                  int    `db:"id" json:"id"`
	Name                string `db:"name" json:"name"`
	Description         string `db:"description" json:"description"`
	Site                string `db:"site" json:"site"`
	Logo                string `db:"logo" json:"logo"`
	WebHookUrl          string `db:"webhook_url" json:"webhook_url"`
	WebHookToken        string `db:"webhook_token" json:"webhook_token"`
	PreAuthRedirectUrl  string `db:"pre_auth_redirect_url" json:"pre_auth_redirect_url"`
	PostAuthRedirectUrl string `db:"post_auth_redirect_url" json:"post_auth_redirect_url"`
}
