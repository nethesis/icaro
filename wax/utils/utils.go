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

package utils

import (
	"crypto/md5"
	"fmt"
	"io"

	"sun-api/database"
	"sun-api/models"
)

func ExtractUnit(uuid string) models.Unit {
	var unit models.Unit
	db := database.Database()
	db.Where("uuid = ?", uuid).First(&unit)
	db.Close()

	return unit
}

func CalcDigest(unit models.Unit) string {
	h := md5.New()
	io.WriteString(h, unit.Secret+unit.Uuid)
	digest := fmt.Sprintf("%x", h.Sum(nil))

	return digest
}

func Contains(intSlice []int, searchInt int) bool {
	for _, value := range intSlice {
		if value == searchInt {
			return true
		}
	}
	return false
}
