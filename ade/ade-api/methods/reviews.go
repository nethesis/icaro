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

package methods

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReviewPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func PostReviewResult(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
