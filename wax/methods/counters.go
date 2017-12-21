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
 * author: Giacomo Sanchietti <giacomo.sanchietti@nethesis.it>
 */

package methods

import (
	"net/http"
	"net/url"
        "github.com/gin-gonic/gin"
)

func Ack(c *gin.Context) {
	c.String(http.StatusOK, "Ack: 1")
}

func Counters(c *gin.Context, parameters url.Values) {
	status := parameters.Get("status")

	switch status {
        case "start":
		Ack(c)
        case "stop":
		Ack(c)
        case "update":
		Ack(c)
        case "":
                c.String(http.StatusBadRequest, "No status provided")
        default:
                c.String(http.StatusNotImplemented, "Invalid status: '%s'", status)
	}
}
