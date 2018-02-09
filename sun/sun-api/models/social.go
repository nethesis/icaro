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
	Facebook  Social `json:"facebook"`
	Google    Social `json:"google"`
	LinkedIn  Social `json:"linkedin"`
	Instagram Social `json:"instagram"`
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

type InstagramRespToken struct {
	AccessToken string `json:"access_token"`
	User        struct {
		Id             string `json:"id"`
		Username       string `json:"username"`
		FullName       string `json:"full_name"`
		ProfilePicture string `json:"profile_picture"`
	} `json:"user"`
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
	Email     string `json:"emailAddress"`
	FirstName string `json:"firstName"`
	Headline  string `json:"headline"`
	LastName  string `json:"lastName"`
	Location  struct {
		Country struct {
			Code string `json:"code"`
		} `json:"country"`
		Name string `json:"name"`
	} `json:"location"`
	NumConnections int `json:"numConnections"`
	Positions      struct {
		Total  int `json:"_total"`
		Values []struct {
			Company struct {
				ID       int    `json:"id"`
				Industry string `json:"industry"`
				Name     string `json:"name"`
				Size     string `json:"size"`
				Type     string `json:"type"`
			} `json:"company"`
			ID        int  `json:"id"`
			IsCurrent bool `json:"isCurrent"`
			Location  struct {
			} `json:"location"`
			Title string `json:"title"`
		} `json:"values"`
	} `json:"positions"`
}

type InstagramUserDetail struct {
	Data struct {
		Id             string `json:"id"`
		Username       string `json:"username"`
		FullName       string `json:"full_name"`
		ProfilePicture string `json:"profile_picture"`
		Bio            string `json:"bio"`
		Website        string `json:"website"`
		IsBusiness     bool   `json:"is_business"`
		Counts         struct {
			Media      int `json:"media"`
			Follows    int `json:"follows"`
			FollowedBy int `json:"followed_by"`
		} `json:"counts"`
	} `json:"data"`
}
