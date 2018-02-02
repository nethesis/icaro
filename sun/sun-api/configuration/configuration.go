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
	"os"

	"github.com/nethesis/icaro/sun/sun-api/models"
)

type Configuration struct {
	DbHost            string            `json:"db_host"`
	DbPort            string            `json:"db_port"`
	DbUser            string            `json:"db_user"`
	DbName            string            `json:"db_name"`
	DbPassword        string            `json:"db_password"`
	AuthSocial        models.AuthSocial `json:"auth_social"`
	PageLimit         string            `json:"page_limit"`
	TokenExpiresDays  int               `json:"token_expires_days"`
	RouteBlocked      models.AuthMaps   `json:"route_blocked"`
	Endpoints         models.Endpoints  `json:"endpoints"`
	Cors              struct {
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
}
