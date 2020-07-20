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
 * author: Matteo Valentini <matteo.valentini@nethesis.it>
 */

package models

type ReviewPage struct {
	HotspotName      string   `json:"hotspot_name"`
	BgColor          string   `json:"bg_color"`
	TitleColor       string   `json:"title_color"`
	TextColor        string   `json:"text_color"`
	ContainerBgColor string   `json:"container_bg_color"`
	Urls             []string `json:"urls"`
	Threshold        int      `json:"threshold"`
}

type ReviewResult struct {
	Stars   int    `json:"stars"`
	Message string `json:"message"`
}
