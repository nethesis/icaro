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

package tasks

import (
	"fmt"
	"time"

	"github.com/robfig/cron"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

func Init(action string, worker bool) {
	c := cron.New()

	// switch action and call methods
	switch action {
	case "clean-tokens":
		c.AddFunc("@daily", cleanTokens)
		cleanTokens()

	case "store-sessions":
		c.AddFunc("@daily", storeSessions)
		storeSessions()

	case "clean-sessions":
		c.AddFunc("@every 2h", cleanSessions)
		cleanSessions()

	default:
		fmt.Println("Specify a valid action to execute, see -h option")
	}

	// if -w option is passed run as worker
	if worker {
		c.Run()
	}
}

func cleanTokens() {
	db := database.Instance()
	db.Where("expires < ?", time.Now().UTC()).Delete(models.AccessToken{})
}

func storeSessions() {
	var sessions []models.Session

	db := database.Instance()
	db.Find(&sessions)

	for _, s := range sessions {
		if !s.StopTime.IsZero() {
			// create session history model
			sessionHistory := models.SessionHistory{
				SessionId:   s.Id,
				UnitId:      s.UnitId,
				UnitMac:     s.UnitMac,
				HotspotId:   s.HotspotId,
				HotspotDesc: s.HotspotDesc,
				DeviceId:    s.DeviceId,
				DeviceMAC:   s.DeviceMAC,
				UserId:      s.UserId,
				Username:    s.Username,
				BytesUp:     s.BytesUp,
				BytesDown:   s.BytesDown,
				Duration:    s.Duration,
				AuthTime:    s.AuthTime,
				StartTime:   s.StartTime,
				UpdateTime:  s.UpdateTime,
				StopTime:    s.StopTime,
				SessionKey:  s.SessionKey,
			}
			// save to session_histories table
			db.Save(&sessionHistory)

			// delete from sessions table
			db.Delete(&s)
		}
	}

}

func cleanSessions() {
	var sessions []models.Session

	start := time.Now().UTC()
	maxSessionTimout := start.Add(time.Hour * -1)

	db := database.Instance()

	db.Where("update_time <= ? and stop_time = 0", maxSessionTimout).Find(&sessions)

	for _, s := range sessions {
		sessionHistory := models.SessionHistory{
			SessionId:   s.Id,
			UnitId:      s.UnitId,
			UnitMac:     s.UnitMac,
			HotspotId:   s.HotspotId,
			HotspotDesc: s.HotspotDesc,
			DeviceId:    s.DeviceId,
			DeviceMAC:   s.DeviceMAC,
			UserId:      s.UserId,
			Username:    s.Username,
			BytesUp:     s.BytesUp,
			BytesDown:   s.BytesDown,
			Duration:    s.Duration,
			AuthTime:    s.AuthTime,
			StartTime:   s.StartTime,
			UpdateTime:  s.UpdateTime,
			StopTime:    s.UpdateTime,
			SessionKey:  s.SessionKey,
		}
		// save to session_histories table
		db.Save(&sessionHistory)

		// delete from sessions table
		db.Delete(&s)
	}
}
