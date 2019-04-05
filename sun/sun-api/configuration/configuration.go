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

package configuration

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nethesis/icaro/sun/sun-api/models"
)

type Configuration struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Name     string `json:"name"`
		Password string `json:"password"`
	} `json:"database"`
	AuthSocial       models.AuthSocial `json:"auth_social"`
	TokenExpiresDays int               `json:"token_expires_days"`
	RouteBlocked     models.AuthMaps   `json:"route_blocked"`
	Endpoints        models.Endpoints  `json:"endpoints"`
	Shortener        models.Shortener  `json:"shortener"`
	Cors             struct {
		Headers []string `json:"headers"`
		Origins []string `json:"origins"`
		Methods []string `json:"methods"`
	} `json:"cors"`
	Disclaimers struct {
		TermsOfUse   string `json:"terms_of_use"`
		MarketingUse string `json:"marketing_use"`
	} `json:"disclaimers"`
	CaptivePortal struct {
		Redirect       string `json:"redirect"`
		Title          string `json:"title"`
		Subtitle       string `json:"subtitle"`
		Description    string `json:"description"`
		Background     string `json:"background"`
		Logo           string `json:"logo"`   // Logo file name
		Banner         string `json:"banner"` // Banner file name
		LogoContents   string `json:"-"`      // base64 content of Logo
		BannerContents string `json:"-"`      //base64 content of Banner
	} `json:"captive_portal"`
	Survey struct {
		Url                 string `json:"url"`
		FeedbackBodyText    string `json:"feedback_body_text"`
		ReviewBodyText      string `json:"review_body_text"`
		FeedbackSubjectText string `json:"feedback_subject_text"`
		ReviewSubjectText   string `json:"review_subject_text"`
	} `json:"survey"`
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
		Config.Database.User = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		Config.Database.Password = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("DB_HOST") != "" {
		Config.Database.Host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		Config.Database.Port = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_NAME") != "" {
		Config.Database.Name = os.Getenv("DB_NAME")
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
	if os.Getenv("SMS_SERVICE_SID") != "" {
		Config.Endpoints.Sms.ServiceSid = os.Getenv("SMS_SERVICE_SID")
	}
	if os.Getenv("SMS_LOGIN_LINK") != "" {
		Config.Endpoints.Sms.Link = os.Getenv("SMS_LOGIN_LINK")
	}
	if os.Getenv("SMS_SEND_QUOTA_ALERT") != "" {
		Config.Endpoints.Sms.SendQuotaAlert, _ = strconv.ParseBool(os.Getenv("SMS_SEND_QUOTA_ALERT"))
	}

	if os.Getenv("EMAIL_FROM") != "" {
		Config.Endpoints.Email.From = os.Getenv("EMAIL_FROM")
	}
	if os.Getenv("EMAIL_FROM_NAME") != "" {
		Config.Endpoints.Email.From = os.Getenv("EMAIL_FROM_NAME")
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
	if os.Getenv("EMAIL_LOGIN_LINK") != "" {
		Config.Endpoints.Email.Link = os.Getenv("EMAIL_LOGIN_LINK")
	}

	if os.Getenv("SHORTENER_BASE_URL") != "" {
		Config.Shortener.BaseUrl = os.Getenv("SHORTENER_BASE_URL")
	}

	if os.Getenv("CAPTIVE_REDIRECT") != "" {
		Config.CaptivePortal.Redirect = os.Getenv("CAPTIVE_REDIRECT")
	}
	if os.Getenv("CAPTIVE_TITLE") != "" {
		Config.CaptivePortal.Title = os.Getenv("CAPTIVE_TITLE")
	}
	if os.Getenv("CAPTIVE_SUBTITLE") != "" {
		Config.CaptivePortal.Subtitle = os.Getenv("CAPTIVE_SUBTITLE")
	}
	if os.Getenv("CAPTIVE_BACKGROUND") != "" {
		Config.CaptivePortal.Background = os.Getenv("CAPTIVE_BACKGROUND")
	}
	if os.Getenv("CAPTIVE_DESCRIPTION") != "" {
		Config.CaptivePortal.Description = os.Getenv("CAPTIVE_DESCRIPTION")
	}
	if os.Getenv("CAPTIVE_LOGO") != "" {
		Config.CaptivePortal.Logo = os.Getenv("CAPTIVE_LOGO")
	}
	if os.Getenv("CAPTIVE_BANNER") != "" {
		Config.CaptivePortal.Banner = os.Getenv("CAPTIVE_BANNER")
	}

	if os.Getenv("SURVEY_URL") != "" {
		Config.Survey.Url = os.Getenv("SURVEY_URL")
	}

	Config.CaptivePortal.LogoContents = ""
	if _, err := os.Stat(Config.CaptivePortal.Logo); err == nil {
		if data, errRead := ioutil.ReadFile(Config.CaptivePortal.Logo); errRead == nil {
			mimeType := http.DetectContentType(data)
			Config.CaptivePortal.LogoContents = "data:" + mimeType + ";base64," + b64.StdEncoding.EncodeToString([]byte(data))
		}
	}

	Config.CaptivePortal.BannerContents = ""
	if _, err := os.Stat(Config.CaptivePortal.Banner); err == nil {
		if data, errRead := ioutil.ReadFile(Config.CaptivePortal.Banner); errRead == nil {
			mimeType := http.DetectContentType(data)
			Config.CaptivePortal.BannerContents = "data:" + mimeType + ";base64," + b64.StdEncoding.EncodeToString([]byte(data))
		}
	}

}
