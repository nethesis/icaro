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

package methods

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

type GraphResponse struct {
	Labels []string `json:"labels"`
	Set0   struct {
		Labels []string `json:"labels"`
		Data   []int    `db:"data" json:"data"`
	} `json:"set0"`
	Set1 struct {
		Labels []string `json:"labels"`
		Data   []int    `json:"data"`
	} `json:"set1"`
	Avg []int `json:"avg"`
}

func GetCurrentSessions(c *gin.Context) {
	var response GraphResponse
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Query("hotspot")
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	for i := 0; i < 24; i++ {
		var hour = "0"
		if i > 9 {
			hour = strconv.Itoa(i)
		} else {
			hour = hour + strconv.Itoa(i)
		}
		response.Labels = append(response.Labels, hour)
	}

	db := database.Instance()
	chain := db.Table("users").Select("DATE_FORMAT(created, '%h') AS label, count(distinct id) AS data").Where("created >= CURRENT_DATE() AND created <= NOW()")
	rows, err := chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("label").Rows()

	var total = 0
	var num = 0
	for rows.Next() {
		var label = ""
		var data = 0
		rows.Scan(&label, &data)

		response.Set0.Labels = append(response.Set0.Labels, label)
		response.Set0.Data = append(response.Set0.Data, data)

		total += data
		num++
	}
	var res = 0
	if num > 0 {
		res = total / num
	}
	response.Avg = append(response.Avg, res)

	chain = db.Table("sessions").Select("DATE_FORMAT(start_time, '%h') AS label, count(distinct id) AS data").Where("start_time >= CURRENT_DATE() AND start_time <= NOW()")
	rows, err = chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("label").Rows()

	total = 0
	num = 0
	for rows.Next() {
		var label = ""
		var data = 0
		rows.Scan(&label, &data)

		response.Set1.Labels = append(response.Set1.Labels, label)
		response.Set1.Data = append(response.Set1.Data, data)

		total += data
		num++
	}
	res = 0
	if num > 0 {
		res = total / num
	}
	response.Avg = append(response.Avg, res)

	c.JSON(http.StatusOK, response)
}

func GetHistorySessions(c *gin.Context) {
	var response GraphResponse
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Query("hotspot")
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	rangeDate := c.Query("range")
	rangeDateInt, err := strconv.Atoi(rangeDate)
	if err != nil {
		rangeDateInt = 0
	}

	for i := rangeDateInt; i >= 0; i-- {
		now := time.Now().UTC()
		response.Labels = append(response.Labels, now.AddDate(0, 0, -i).Format("02 Jan"))
	}

	db := database.Instance()
	chain := db.Table("session_histories").Select("DATE_FORMAT(start_time, '%d %b') AS label, count(*) AS data").Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND stop_time <= NOW()", rangeDateInt)
	rows, err := chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(start_time)").Rows()

	for rows.Next() {
		var label = ""
		var data = 0
		rows.Scan(&label, &data)

		response.Set0.Labels = append(response.Set0.Labels, label)
		response.Set0.Data = append(response.Set0.Data, data)
	}

	chain = db.Table("users").Select("DATE_FORMAT(created, '%d %b') AS label, count(*) AS data").Where("created BETWEEN DATE_SUB(NOW(), INTERVAL ? DAY) AND NOW()", rangeDateInt)
	rows, err = chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(created)").Rows()

	for rows.Next() {
		var label = ""
		var data = 0
		rows.Scan(&label, &data)

		response.Set1.Labels = append(response.Set1.Labels, label)
		response.Set1.Data = append(response.Set1.Data, data)
	}

	c.JSON(http.StatusOK, response)
}

func GetHistoryTraffic(c *gin.Context) {
	var response GraphResponse
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Query("hotspot")
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	rangeDate := c.Query("range")
	rangeDateInt, err := strconv.Atoi(rangeDate)
	if err != nil {
		rangeDateInt = 0
	}

	for i := rangeDateInt; i >= 0; i-- {
		now := time.Now().UTC()
		response.Labels = append(response.Labels, now.AddDate(0, 0, -i).Format("02 Jan"))
	}

	db := database.Instance()
	chain := db.Table("session_histories").Select("DATE_FORMAT(start_time, '%d %b') AS label, sum(bytes_up) AS data").Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND stop_time <= NOW()", rangeDateInt)
	rows, err := chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(start_time)").Rows()

	for rows.Next() {
		var label = ""
		var data = 0
		rows.Scan(&label, &data)

		response.Set0.Labels = append(response.Set0.Labels, label)
		response.Set0.Data = append(response.Set0.Data, data)
	}

	chain = db.Table("session_histories").Select("DATE_FORMAT(start_time, '%d %b') AS label, sum(bytes_down) AS data").Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND stop_time <= NOW()", rangeDateInt)
	rows, err = chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(start_time)").Rows()

	for rows.Next() {
		var label = ""
		var data = 0
		rows.Scan(&label, &data)

		response.Set1.Labels = append(response.Set1.Labels, label)
		response.Set1.Data = append(response.Set1.Data, data)
	}

	c.JSON(http.StatusOK, response)
}

func GetHistoryAvgUserTraffic(c *gin.Context) {
	var response GraphResponse
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Query("hotspot")
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	rangeDate := c.Query("range")
	rangeDateInt, err := strconv.Atoi(rangeDate)
	if err != nil {
		rangeDateInt = 0
	}

	for i := rangeDateInt; i >= 0; i-- {
		now := time.Now().UTC()
		response.Labels = append(response.Labels, now.AddDate(0, 0, -i).Format("02 Jan"))
	}

	db := database.Instance()
	chain := db.Table("session_histories").Select("DATE_FORMAT(start_time, '%d %b') AS label, sum(bytes_up)/count(distinct user_id) AS data").Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND stop_time <= NOW()", rangeDateInt)
	rows, err := chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(start_time)").Rows()

	var total = 0
	var num = 0
	for rows.Next() {
		var label = ""
		var data = 0.0
		rows.Scan(&label, &data)

		response.Set0.Labels = append(response.Set0.Labels, label)
		response.Set0.Data = append(response.Set0.Data, int(math.Round(data)))

		total += int(math.Round(data))
		num++
	}
	var res = 0
	if num > 0 {
		res = total / num
	}
	response.Avg = append(response.Avg, res)

	chain = db.Table("session_histories").Select("DATE_FORMAT(start_time, '%d %b') AS label, sum(bytes_down)/count(distinct user_id) AS data").Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND stop_time <= NOW()", rangeDateInt)
	rows, err = chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(start_time)").Rows()

	total = 0
	num = 0
	for rows.Next() {
		var label = ""
		var data = 0.0
		rows.Scan(&label, &data)

		response.Set1.Labels = append(response.Set1.Labels, label)
		response.Set1.Data = append(response.Set1.Data, int(math.Round(data)))

		total += int(math.Round(data))
		num++
	}
	res = 0
	if num > 0 {
		res = total / num
	}
	response.Avg = append(response.Avg, res)

	c.JSON(http.StatusOK, response)
}

func GetHistoryAvgUserDuration(c *gin.Context) {
	var response GraphResponse
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Query("hotspot")
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	rangeDate := c.Query("range")
	rangeDateInt, err := strconv.Atoi(rangeDate)
	if err != nil {
		rangeDateInt = 0
	}

	for i := rangeDateInt; i >= 0; i-- {
		now := time.Now().UTC()
		response.Labels = append(response.Labels, now.AddDate(0, 0, -i).Format("02 Jan"))
	}

	db := database.Instance()
	chain := db.Table("session_histories").Select("DATE_FORMAT(start_time, '%d %b') AS label, sum(duration)/count(distinct user_id) AS data").Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND stop_time <= NOW()", rangeDateInt)
	rows, err := chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(start_time)").Rows()

	var total = 0
	var num = 0
	for rows.Next() {
		var label = ""
		var data = 0.0
		rows.Scan(&label, &data)

		response.Set0.Labels = append(response.Set0.Labels, label)
		response.Set0.Data = append(response.Set0.Data, int(math.Round(data)))

		total += int(math.Round(data))
		num++
	}
	var res = 0
	if num > 0 {
		res = total / num
	}
	response.Avg = append(response.Avg, res)

	c.JSON(http.StatusOK, response)
}

func GetHistoryAvgConnTraffic(c *gin.Context) {
	var response GraphResponse
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Query("hotspot")
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	rangeDate := c.Query("range")
	rangeDateInt, err := strconv.Atoi(rangeDate)
	if err != nil {
		rangeDateInt = 0
	}

	for i := rangeDateInt; i >= 0; i-- {
		now := time.Now().UTC()
		response.Labels = append(response.Labels, now.AddDate(0, 0, -i).Format("02 Jan"))
	}

	db := database.Instance()
	chain := db.Table("session_histories").Select("DATE_FORMAT(start_time, '%d %b') AS label, sum(bytes_up)/count(id) AS data").Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND stop_time <= NOW()", rangeDateInt)
	rows, err := chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(start_time)").Rows()

	var total = 0
	var num = 0
	for rows.Next() {
		var label = ""
		var data = 0.0
		rows.Scan(&label, &data)

		response.Set0.Labels = append(response.Set0.Labels, label)
		response.Set0.Data = append(response.Set0.Data, int(math.Round(data)))

		total += int(math.Round(data))
		num++
	}
	var res = 0
	if num > 0 {
		res = total / num
	}
	response.Avg = append(response.Avg, res)

	chain = db.Table("session_histories").Select("DATE_FORMAT(start_time, '%d %b') AS label, sum(bytes_down)/count(id) AS data").Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND stop_time <= NOW()", rangeDateInt)
	rows, err = chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(start_time)").Rows()

	total = 0
	num = 0
	for rows.Next() {
		var label = ""
		var data = 0.0
		rows.Scan(&label, &data)

		response.Set1.Labels = append(response.Set1.Labels, label)
		response.Set1.Data = append(response.Set1.Data, int(math.Round(data)))

		total += int(math.Round(data))
		num++
	}
	res = 0
	if num > 0 {
		res = total / num
	}
	response.Avg = append(response.Avg, res)

	c.JSON(http.StatusOK, response)
}

func GetHistoryAvgConnDuration(c *gin.Context) {
	var response GraphResponse
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Query("hotspot")
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	rangeDate := c.Query("range")
	rangeDateInt, err := strconv.Atoi(rangeDate)
	if err != nil {
		rangeDateInt = 0
	}

	for i := rangeDateInt; i >= 0; i-- {
		now := time.Now().UTC()
		response.Labels = append(response.Labels, now.AddDate(0, 0, -i).Format("02 Jan"))
	}

	db := database.Instance()
	chain := db.Table("session_histories").Select("DATE_FORMAT(start_time, '%d %b') AS label, sum(duration)/count(id) AS data").Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND stop_time <= NOW()", rangeDateInt)
	rows, err := chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(start_time)").Rows()

	var total = 0
	var num = 0
	for rows.Next() {
		var label = ""
		var data = 0.0
		rows.Scan(&label, &data)

		response.Set0.Labels = append(response.Set0.Labels, label)
		response.Set0.Data = append(response.Set0.Data, int(math.Round(data)))

		total += int(math.Round(data))
		num++
	}
	var res = 0
	if num > 0 {
		res = total / num
	}
	response.Avg = append(response.Avg, res)

	c.JSON(http.StatusOK, response)
}

func GetHistorySMSYear(c *gin.Context) {
	var response GraphResponse
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Query("hotspot")
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	response.Labels = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

	db := database.Instance()
	chain := db.Table("hotspot_sms_counts").Select("DATE_FORMAT(sent, '%M') AS label, count(*) as data")
	rows, err := chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("label").Order("sent").Rows()

	for rows.Next() {
		var label = ""
		var data = 0.0
		rows.Scan(&label, &data)

		response.Set0.Labels = append(response.Set0.Labels, label)
		response.Set0.Data = append(response.Set0.Data, int(math.Round(data)))
	}

	c.JSON(http.StatusOK, response)
}

func GetHistorySMSHistory(c *gin.Context) {
	var response GraphResponse
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Query("hotspot")
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		hotspotIdInt = 0
	}

	rangeDate := c.Query("range")
	rangeDateInt, err := strconv.Atoi(rangeDate)
	if err != nil {
		rangeDateInt = 0
	}

	for i := rangeDateInt; i >= 0; i-- {
		now := time.Now().UTC()
		response.Labels = append(response.Labels, now.AddDate(0, 0, -i).Format("02 Jan"))
	}

	db := database.Instance()
	chain := db.Table("hotspot_sms_counts").Select("DATE_FORMAT(sent, '%d %b') AS label, count(*) AS data").Where("sent >= DATE_SUB(NOW(), INTERVAL ? DAY) AND sent <= NOW()", rangeDateInt)
	rows, err := chain.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Group("DATE(sent)").Rows()

	for rows.Next() {
		var label = ""
		var data = 0.0
		rows.Scan(&label, &data)

		response.Set0.Labels = append(response.Set0.Labels, label)
		response.Set0.Data = append(response.Set0.Data, int(math.Round(data)))
	}

	c.JSON(http.StatusOK, response)
}
