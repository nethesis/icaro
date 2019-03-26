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
	"strconv"
	"time"

	"github.com/robfig/cron"

	"github.com/nethesis/icaro/ade/ade-api/utils"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

func Init(action string, worker bool) {
	c := cron.New()

	// switch action and call methods
	switch action {

	case "send-surveys":
		c.AddFunc("@every 1h", sendSurveys)
		sendSurveys()

	default:
		fmt.Println("Specify a valid action to execute, see -h option")
	}

	// if -w option is passed run as worker
	if worker {
		c.Run()
	}
}

func sendSurveys() {
	var users []models.User

	db := database.Instance()

	// get users with survey auth
	db.Where("survey_auth = 1").Find(&users)

	for _, u := range users {
		// get hotspot preferences
		var hotspot models.Hotspot
		db.Where("id = ?", u.HotspotId).Find(&hotspot)

		var hotspotName models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'captive_2_title'", u.HotspotId).Find(&hotspotName)

		var hotspotLogo models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'captive_3_logo'", u.HotspotId).Find(&hotspotLogo)

		var hotspotBg models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'captive_7_background'", u.HotspotId).Find(&hotspotBg)

		var marketingEnabled models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'marketing_1_enabled'", u.HotspotId).Find(&marketingEnabled)

		var marketingFirstEmail models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'marketing_3_first_email_enabled'", u.HotspotId).Find(&marketingFirstEmail)

		var marketingSecondEmail models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'marketing_4_second_email_enabled'", u.HotspotId).Find(&marketingSecondEmail)

		var marketingSMS models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'marketing_5_sms_enabled'", u.HotspotId).Find(&marketingSMS)

		var marketingFirstAfter models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'marketing_6_first_after'", u.HotspotId).Find(&marketingFirstAfter)

		var marketingSecondAfter models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'marketing_7_second_after'", u.HotspotId).Find(&marketingSecondAfter)

		var marketingSecondAfterDays models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'marketing_8_second_after_days'", u.HotspotId).Find(&marketingSecondAfterDays)

		var adeToken models.AdeToken
		db.Where("hotspot_id = ? AND user_id = ?", u.HotspotId, u.Id).First(&adeToken)

		var hotspotFeedbackBodyText models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'marketing_13_feedback_body_text'", u.HotspotId).Find(&hotspotFeedbackBodyText)

		var hotspotReviewBodyText models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'marketing_14_review_body_text'", u.HotspotId).Find(&hotspotReviewBodyText)

		// if token not exists create token
		if adeToken.Id == 0 {
			adeToken = utils.CreateAdeToken(u)
		}

		// check if marketing is enabled
		if marketingEnabled.Value == "true" {
			// check if first mail is enabled
			if marketingFirstEmail.Value == "true" {
				// check if time to send
				marketingFirstAfterInt, _ := strconv.Atoi(marketingFirstAfter.Value)
				if adeToken.FeedbackSentTime.IsZero() && u.Created.Add(time.Duration(marketingFirstAfterInt)*time.Hour).Before(time.Now()) {
					// send mail
					if len(u.Email) > 0 {
						utils.SendFeedBackMessageToUser(adeToken, u.Email, hotspotName.Value, hotspotLogo.Value, hotspotBg.Value, hotspot, hotspotFeedbackBodyText.Value)
					}

					// check if sms is enabled
					if marketingSMS.Value == "true" {
						if u.AccountType == "sms" {
							status := utils.SendSMS(adeToken, "Leave feedback", "feedbacks", u.Username, u.HotspotId)

							adeToken.FeedbackSentTime = time.Now()
							db.Save(&adeToken)

							if !status {
								fmt.Println("Feedback SMS not sent")
							}
						}
					}
				}

			}

			// check if second mail is enabled
			if marketingSecondEmail.Value == "true" {
				// check if time to send
				switch marketingSecondAfter.Value {
				case "days":
					marketingSecondAfterDaysInt, _ := strconv.Atoi(marketingSecondAfterDays.Value)
					if adeToken.ReviewSentTime.IsZero() && u.Created.AddDate(0, 0, marketingSecondAfterDaysInt).Before(time.Now()) {
						// send mail
						if len(u.Email) > 0 {
							utils.SendReviewMessageToUser(adeToken, u.Email, hotspotName.Value, hotspotLogo.Value, hotspotBg.Value, hotspot, hotspotReviewBodyText.Value)
						}

						// check if sms is enabled
						if marketingSMS.Value == "true" {
							if u.AccountType == "sms" {
								status := utils.SendSMS(adeToken, "Leave review", "reviews", u.Username, u.HotspotId)

								adeToken.ReviewSentTime = time.Now()
								db.Save(&adeToken)

								if !status {
									fmt.Println("Review SMS not sent")
								}
							}
						}
					}

				case "expiration":
					if adeToken.ReviewSentTime.IsZero() && u.ValidUntil.Before(time.Now()) {
						// send mail
						if len(u.Email) > 0 {
							utils.SendReviewMessageToUser(adeToken, u.Email, hotspotName.Value, hotspotLogo.Value, hotspotBg.Value, hotspot, hotspotReviewBodyText.Value)
						}

						// check if sms is enabled
						if marketingSMS.Value == "true" {
							if u.AccountType == "sms" {
								status := utils.SendSMS(adeToken, "Leave review", "reviews", u.Username, u.HotspotId)

								adeToken.ReviewSentTime = time.Now()
								db.Save(&adeToken)

								if !status {
									fmt.Println("Review SMS not sent")
								}
							}
						}
					}
				}

			}

		}
	}
}
