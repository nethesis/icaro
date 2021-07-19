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

type Social struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

type AuthSocial struct {
	Facebook  Social `json:"facebook"`
	LinkedIn  Social `json:"linkedin"`
	Instagram Social `json:"instagram"`
}

type FacebookRespToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type LinkedInRespToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type InstagramRespToken struct {
	AccessToken string `json:"access_token"`
	// User        struct {
	// 	Id             string `json:"id"`
	// 	Username       string `json:"username"`
	// 	FullName       string `json:"full_name"`
	// 	ProfilePicture string `json:"profile_picture"`
	// } `json:"user"`
	UserId int `json:"user_id"`
}

type FacebookInspectToken struct {
	Data struct {
		AppId       string   `json:"app_id"`
		Type        string   `json:"type"`
		Application string   `json:"application"`
		ExpiresAt   int      `json:"expires_at"`
		IsValid     bool     `json:"is_valid"`
		IssuedAt    int      `json:"issued_at"`
		Scopes      []string `json:"scopes"`
		UserId      string   `json:"user_id"`
	} `json:"data"`
}

type FacebookUserDetail struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	Likes    struct {
		Data []struct {
			Name        string `json:"name"`
			Id          string `json:"id"`
			CreatedTime string `json:"created_time"`
		} `json:"data"`
	} `json:"likes"`
}

type LinkedInUserDetail struct {
	Id        string `json:"id"`
	FirstName string `json:"localizedFirstName"`
	LastName  string `json:"localizedLastName"`
}

type LinkedinEmailDetail struct {
	Elements []struct {
		Handle        string `json:"handle"`
		HandleDetails struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"handle~"`
	} `json:"elements"`
}

type InstagramUserDetail struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	// Data struct {
	// 	Id             string `json:"id"`
	// 	Username       string `json:"username"`
	// 	FullName       string `json:"full_name"`
	// 	ProfilePicture string `json:"profile_picture"`
	// 	Bio            string `json:"bio"`
	// 	Website        string `json:"website"`
	// 	IsBusiness     bool   `json:"is_business"`
	// 	Counts         struct {
	// 		Media      int `json:"media"`
	// 		Follows    int `json:"follows"`
	// 		FollowedBy int `json:"followed_by"`
	// 	} `json:"counts"`
	// } `json:"data"`
}
