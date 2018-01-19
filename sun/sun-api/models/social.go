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

package models

type Social struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

type AuthSocial struct {
	Facebook Social `json:"facebook"`
	Google   Social `json:"google"`
	LinkedIn Social `json:"linkedin"`
}

type FacebookRespToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type GoogleRespToken struct {
	Aud       string `json:"aud"`
	UserId    string `json:"user_id"`
	Scope     string `json:"scope"`
	ExpiresIn string `json:"expires_in"`
	Email     string `json:"email"`
}

type LinkedInRespToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
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

type GoogleUserDetail struct {
	Birthday string `json:"birthday"`
	Emails   []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"emails"`
	ObjectType  string `json:"objectType"`
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Language    string `json:"language"`
	AgeRange    struct {
		Min int `json:"min"`
	} `json:"ageRange"`
}

type LinkedInUserDetail struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"emailAddress"`
}
