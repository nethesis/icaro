/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Icaro project.
 *
 * Icaro is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Icaro is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Icaro.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package configuration

import (
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"os"

	"github.com/nethesis/icaro/sun/sun-api/models"
)

type Configuration struct {
	DbHost           string            `json:"db_host"`
	DbPort           string            `json:"db_port"`
	DbUser           string            `json:"db_user"`
	DbName           string            `json:"db_name"`
	DbPassword       string            `json:"db_password"`
	AuthSocial       models.AuthSocial `json:"auth_social"`
	TokenExpiresDays int               `json:"token_expires_days"`
	RouteBlocked     models.AuthMaps   `json:"route_blocked"`
	Endpoints        models.Endpoints  `json:"endpoints"`
	Cors             struct {
		Headers []string `json:"headers"`
		Origins []string `json:"origins"`
		Methods []string `json:"methods"`
	} `json:"cors"`
}

var Config = Configuration{}

func Init(ConfigFilePtr *string) {
	// read configuration
	if _, err := os.Stat(*ConfigFilePtr); err == nil {
		file, _ := os.Open(*ConfigFilePtr)
		decoder := json.NewDecoder(file)
		// check errors or parse JSON
		err := decoder.Decode(&Config)
		if err != nil {
			fmt.Println("Configuration parsing error:", err)
		}
	}

	if os.Getenv("DB_USER") != "" {
		Config.DbUser = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		Config.DbPassword = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("DB_HOST") != "" {
		Config.DbHost = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		Config.DbPort = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_NAME") != "" {
		Config.DbName = os.Getenv("DB_NAME")
	}

	if os.Getenv("CORS_ORIGINS") != "" {
		Config.Cors.Origins = strings.Split(os.Getenv("CORS_ORIGINS"), " ")
	}

	if os.Getenv("FACEBOOK_CLIENT_ID") != "" {
		Config.AuthSocial.Facebook.ClientId = os.Getenv("FACEBOOK_CLIENT_ID")
	}
	if os.Getenv("FACEBOOK_CLIENT_SECRET") != "" {
		Config.AuthSocial.Facebook.ClientSecret = os.Getenv("FACEBOOK_CLIENT_SECRET")
	}
	if os.Getenv("FACEBOOK_REDIRECT_URL") != "" {
		Config.AuthSocial.Facebook.RedirectURI = os.Getenv("FACEBOOK_REDIRECT_URL")
	}
	if os.Getenv("GOOGLE_CLIENT_ID") != "" {
		Config.AuthSocial.Google.ClientId = os.Getenv("GOOGLE_CLIENT_ID")
	}
	if os.Getenv("GOOGLE_CLIENT_SECRET") != "" {
		Config.AuthSocial.Google.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	}
	if os.Getenv("GOOGLE_CLIENT_REDIRECT_URL") != "" {
		Config.AuthSocial.Google.RedirectURI = os.Getenv("GOOGLE_REDIRECT_URL")
	}
	if os.Getenv("LINKEDIN_CLIENT_ID") != "" {
		Config.AuthSocial.LinkedIn.ClientId = os.Getenv("LINKEDIN_CLIENT_ID")
	}
	if os.Getenv("LINKEDIN_CLIENT_SECRET") != "" {
		Config.AuthSocial.LinkedIn.ClientSecret = os.Getenv("LINKEDIN_CLIENT_SECRET")
	}
	if os.Getenv("LINKEDIN_REDIRECT_URL") != "" {
		Config.AuthSocial.LinkedIn.RedirectURI = os.Getenv("LINKEDIN_REDIRECT_URL")
	}
	if os.Getenv("INSTAGRAM_CLIENT_ID") != "" {
		Config.AuthSocial.Instagram.ClientId = os.Getenv("INSTAGRAM_CLIENT_ID")
	}
	if os.Getenv("INSTAGRAM_CLIENT_SECRET") != "" {
		Config.AuthSocial.Instagram.ClientSecret = os.Getenv("INSTAGRAM_CLIENT_SECRET")
	}
	if os.Getenv("INSTAGRAM_REDIRECT_URL") != "" {
		Config.AuthSocial.Instagram.RedirectURI = os.Getenv("INSTAGRAM_REDIRECT_URL")
	}

	if os.Getenv("SMS_ACCOUNT_SID") != "" {
		Config.Endpoints.Sms.AccountSid = os.Getenv("SMS_ACCOUNT_SID")
	}
	if os.Getenv("SMS_AUTH_TOKEN") != "" {
			Config.Endpoints.Sms.AuthToken = os.Getenv("SMS_AUTH_TOKEN")
	}
	if os.Getenv("SMS_NUMBER") != "" {
		Config.Endpoints.Sms.Number = os.Getenv("SMS_NUMBER")
	}

	if os.Getenv("EMAIL_FORM") != "" {
		Config.Endpoints.Email.From = os.Getenv("EMAIL_FROM")
	}
	if os.Getenv("EMAIL_SMTP_HOST") != "" {
		Config.Endpoints.Email.SMTPHost = os.Getenv("EMAIL_SMTP_HOST")
	}
	if os.Getenv("EMAIL_SMTP_PORT") != "" {
		Config.Endpoints.Email.SMTPPort, _ = strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
	}
	if os.Getenv("EMAIL_SMTP_USER") != "" {
		Config.Endpoints.Email.SMTPUser = os.Getenv("EMAIL_SMTP_USER")
	}
	if os.Getenv("EMAIL_SMTP_PASSWORD") != "" {
		Config.Endpoints.Email.SMTPPassword = os.Getenv("EMAIL_SMTP_PASSWORD")
	}
}
