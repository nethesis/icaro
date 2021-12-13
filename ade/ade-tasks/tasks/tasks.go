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
	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

func Init(action string, worker bool) {
	c := cron.New()

	// switch action and call methods
	switch action {

	case "send-surveys":
		c.AddFunc("@every 15m", sendSurveysActive)
		sendSurveysActive()
	case "send-surveys-expired":
		c.AddFunc("@every 15m", sendSurveysExpired)
		sendSurveysExpired()

	default:
		fmt.Println("Specify a valid action to execute, see -h option")
	}

	// if -w option is passed run as worker
	if worker {
		c.Run()
	}
}

type User struct {
	Id          int
	HotspotId   int
	Username    string
	Email       string
	AccountType string
	ValidUntil  time.Time
	ValidFrom   time.Time
	IsActive    bool
}

func sendSurveysActive() {
	var users []models.User

	db := database.Instance()
	
	db.Raw(`
		SELECT *
		FROM users
		WHERE survey_auth = 1
		  AND id IN (
		    SELECT user_id FROM ade_tokens
		    WHERE (
					feedback_sent_time != "0000-00-00 00:00:00" OR review_sent_time != "0000-00-00 00:00:00"
				) AND NOT (
					feedback_sent_time != "0000-00-00 00:00:00" AND review_sent_time != "0000-00-00 00:00:00"
				)
		  )
			AND hotspot_id IN (
				SELECT DISTINCT hotspot_id FROM hotspot_preferences
				WHERE ` + "`" + "key" + "`" + ` = "marketing_1_enabled" AND value = "true"
			)
			AND created >= NOW() - INTERVAL 1 YEAR

		UNION

		SELECT *
		FROM users
		WHERE survey_auth = 1
		  AND id NOT IN (
		    SELECT user_id FROM ade_tokens
		  )
			AND hotspot_id IN (
				SELECT DISTINCT hotspot_id FROM hotspot_preferences
				WHERE ` + "`" + "key" + "`" + ` = "marketing_1_enabled" AND value = "true"
			)
			AND created >= NOW() - INTERVAL 1 YEAR
	`).Scan(&users)

	usersList := make([]User, len(users))

	for i, u := range users {
		usersList[i].Id = u.Id
		usersList[i].HotspotId = u.HotspotId
		usersList[i].Username = u.Username
		usersList[i].Email = u.Email
		usersList[i].AccountType = u.AccountType
		usersList[i].ValidUntil = u.ValidUntil
		usersList[i].ValidFrom = u.ValidFrom
		usersList[i].IsActive = true
	}

	sendSurveys(usersList)
}

func sendSurveysExpired() {

	var users []models.UserHistory

	db := database.Instance()

	db.Raw(`
		SELECT *
		FROM user_histories
		WHERE survey_auth = 1
		  AND user_id IN (
		    SELECT user_id FROM ade_tokens
		    WHERE (
					feedback_sent_time != "0000-00-00 00:00:00" OR review_sent_time != "0000-00-00 00:00:00"
				) AND NOT (
					feedback_sent_time != "0000-00-00 00:00:00" AND review_sent_time != "0000-00-00 00:00:00"
				)
		  )
			AND hotspot_id IN (
				SELECT DISTINCT hotspot_id FROM hotspot_preferences
				WHERE ` + "`" + "key" + "`" + ` = "marketing_1_enabled" AND value = "true"
			)
			AND created >= NOW() - INTERVAL 1 YEAR

		UNION

		SELECT *
		FROM user_histories
		WHERE survey_auth = 1
		  AND user_id NOT IN (
		    SELECT user_id FROM ade_tokens
		  )
			AND hotspot_id IN (
				SELECT DISTINCT hotspot_id FROM hotspot_preferences
				WHERE ` + "`" + "key" + "`" + ` = "marketing_1_enabled" AND value = "true"
			)
			AND created >= NOW() - INTERVAL 1 YEAR
	`).Scan(&users)

	usersList := make([]User, len(users))

	for i, u := range users {
		usersList[i].Id = u.UserId
		usersList[i].HotspotId = u.HotspotId
		usersList[i].Username = u.Username
		usersList[i].Email = u.Email
		usersList[i].AccountType = u.AccountType
		usersList[i].ValidUntil = u.ValidUntil
		usersList[i].ValidFrom = u.ValidFrom
		usersList[i].IsActive = false
	}

	sendSurveys(usersList)
}

func sendSurveys(users []User) {

	db := database.Instance()

	for _, u := range users {
		// get hotspot preferences
		var hotspot models.Hotspot
		db.Where("id = ?", u.HotspotId).Find(&hotspot)

		var hotspotName models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'captive_2_title'", u.HotspotId).Find(&hotspotName)

		var hotspotBg models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'captive_7_background'", u.HotspotId).Find(&hotspotBg)

		var hotspotContainerBgColor models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'captive_82_container_bg_color'", u.HotspotId).Find(&hotspotContainerBgColor)

		var hotspotTitleColor models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'captive_83_title_color'", u.HotspotId).Find(&hotspotTitleColor)

		var hotspotTextColor models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'captive_84_text_color'", u.HotspotId).Find(&hotspotTextColor)

		var hotspotTextStyle models.HotspotPreference
		db.Where("hotspot_id = ? AND `key` = 'captive_85_text_style'", u.HotspotId).Find(&hotspotTextStyle)

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
		if adeToken.Id == 0 && u.IsActive && marketingEnabled.Value == "true" {
			var user models.User
			db.Where("id = ?", u.Id).First(&user)
			adeToken = utils.CreateAdeToken(user)
		}

		if adeToken.Id != 0 {
			// check if marketing is enabled
			if marketingEnabled.Value == "true" {
				// check if first mail is enabled
				if marketingFirstEmail.Value == "true" {
					// check if time to send
					marketingFirstAfterInt, _ := strconv.Atoi(marketingFirstAfter.Value)
					if adeToken.FeedbackSentTime.IsZero() && u.ValidFrom.Add(time.Duration(marketingFirstAfterInt)*time.Minute).Before(time.Now()) {
						// send mail
						if len(u.Email) > 0 {
							utils.SendFeedBackMessageToUser(adeToken, u.Email, hotspotName.Value, hotspotBg.Value, hotspotContainerBgColor.Value, hotspotTitleColor.Value, hotspotTextColor.Value, hotspotTextStyle.Value, hotspot, hotspotFeedbackBodyText.Value)
						}

						// check if sms is enabled
						if marketingSMS.Value == "true" {
							if u.AccountType == "sms" {
								status := utils.SendSMS(adeToken, configuration.Config.Survey.FeedbackSmsBodyText, "feedbacks", u.Username, u.HotspotId)

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
						if adeToken.ReviewSentTime.IsZero() && u.ValidFrom.AddDate(0, 0, marketingSecondAfterDaysInt).Before(time.Now()) {
							// send mail
							if len(u.Email) > 0 {
								utils.SendReviewMessageToUser(adeToken, u.Email, hotspotName.Value, hotspotBg.Value, hotspotContainerBgColor.Value, hotspotTitleColor.Value, hotspotTextColor.Value, hotspotTextStyle.Value, hotspot, hotspotReviewBodyText.Value)
							}

							// check if sms is enabled
							if marketingSMS.Value == "true" {
								if u.AccountType == "sms" {
									status := utils.SendSMS(adeToken, configuration.Config.Survey.ReviewSmsBodyText, "reviews", u.Username, u.HotspotId)

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
								utils.SendReviewMessageToUser(adeToken, u.Email, hotspotName.Value, hotspotBg.Value, hotspotContainerBgColor.Value, hotspotTitleColor.Value, hotspotTextColor.Value, hotspotTextStyle.Value, hotspot, hotspotReviewBodyText.Value)
							}

							// check if sms is enabled
							if marketingSMS.Value == "true" {
								if u.AccountType == "sms" {
									status := utils.SendSMS(adeToken, configuration.Config.Survey.ReviewSmsBodyText, "reviews", u.Username, u.HotspotId)

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
}
