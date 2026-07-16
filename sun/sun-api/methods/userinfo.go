/*
 * Copyright (C) 2026 Nethesis S.r.l.
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
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
)

// emptyUserInfo answers 200 with an empty widget: for the My dashboard a
// missing/invalid/unlinked identity simply means "nothing to show", and a
// non-2xx status would surface an error banner in the UI
func emptyUserInfo(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"status": "empty",
		"msg":    msg,
		"widget": gin.H{"items": []gin.H{}},
	})
}

// UserInfo returns aggregate statistics for the company (reseller) of the
// authenticated My user. The My dashboard widget calls it cross-origin with
// the user's Logto ID token as bearer: the token is verified against the
// same issuer used for the OIDC login, and the company is resolved through
// the organization claims exactly like the login flow (never by email).
func UserInfo(c *gin.Context) {
	config := configuration.Config

	if config.OIDC.Issuer == "" {
		emptyUserInfo(c, "OIDC not configured")
		return
	}

	// Bearer token from the Authorization header
	authParts := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
	if len(authParts) != 2 || !strings.EqualFold(authParts[0], "Bearer") {
		emptyUserInfo(c, "Missing or invalid authorization header")
		return
	}

	provider, err := getOIDCProvider(config.OIDC.Issuer)
	if err != nil {
		log.Println("userinfo: OIDC provider initialization failed:", err)
		emptyUserInfo(c, "Cannot retrieve OIDC provider")
		return
	}

	// The caller is My's own application, not our third-party client: the
	// audience is checked manually against the configured value
	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
	idToken, err := verifier.Verify(c.Request.Context(), authParts[1])
	if err != nil {
		log.Println("userinfo denied: invalid token:", err)
		emptyUserInfo(c, "Invalid token")
		return
	}

	if config.OIDC.UserinfoAudience != "" {
		audOk := false
		for _, aud := range idToken.Audience {
			if aud == config.OIDC.UserinfoAudience {
				audOk = true
				break
			}
		}
		if !audOk {
			log.Println("userinfo denied: audience mismatch:", idToken.Audience)
			emptyUserInfo(c, "Invalid token: audience mismatch")
			return
		}
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		emptyUserInfo(c, "Invalid token: claims extraction failed")
		return
	}

	// The company is the organization where the user holds the reseller
	// organization-role, mapped to the account through logto_org_id
	orgIDs := extractOrgIDsWithRole(claims, config.OIDC.OrgResellerRole)
	if len(orgIDs) == 0 {
		log.Printf("userinfo denied for %v: no organization with role %q (organization_roles=%v)",
			claims["email"], config.OIDC.OrgResellerRole, claims["organization_roles"])
		emptyUserInfo(c, "No reseller organization in token")
		return
	}

	db := database.Instance()
	var account models.Account
	db.Where("logto_org_id IN (?)", orgIDs).First(&account)
	if account.Id == 0 || account.Type != "reseller" {
		log.Printf("userinfo denied for %v: no company account linked to organizations %v", claims["email"], orgIDs)
		emptyUserInfo(c, "No company account linked to the organization")
		return
	}

	// Aggregate counters over the company's hotspots
	countByHotspot := func(table string) int {
		var n int
		db.Table(table).
			Joins("JOIN hotspots ON hotspots.id = "+table+".hotspot_id").
			Where("hotspots.account_id = ?", account.Id).
			Count(&n)
		return n
	}

	var hotspots, managers int
	db.Table("hotspots").Where("account_id = ?", account.Id).Count(&hotspots)
	db.Table("accounts").Where("creator_id = ? AND type IN (?)", account.Id, []string{"desk", "customer"}).Count(&managers)
	units := countByHotspot("units")
	guests := countByHotspot("users")
	devices := countByHotspot("devices")
	sessions := countByHotspot("sessions")

	var sms models.AccountSmsCount
	db.Where("account_id = ?", account.Id).First(&sms)
	smsRemaining := sms.SmsMaxCount - sms.SmsCount
	if smsRemaining < 0 {
		smsRemaining = 0
	}

	// SMS tone: warn the reseller when the remaining quota gets thin
	smsTone := "neutral"
	if sms.SmsMaxCount > 0 {
		switch {
		case smsRemaining*100 <= sms.SmsMaxCount*10:
			smsTone = "danger"
		case smsRemaining*100 <= sms.SmsMaxCount*25:
			smsTone = "warning"
		}
	}

	// Generic widget contract consumed by the My dashboard: My renders
	// `widget.items` without knowing anything about NethSpot
	base := strings.TrimRight(config.OIDC.FrontendURL, "/")
	widgetItems := []gin.H{
		{"label": "Hotspot", "value": hotspots, "tone": "neutral", "link": base + "/#/hotspots"},
		{"label": "Unit", "value": units, "tone": "neutral", "link": base + "/#/units"},
		{"label": "Utenti", "value": guests, "tone": "neutral", "link": base + "/#/users"},
		{"label": "SMS residui", "value": smsRemaining, "tone": smsTone, "link": base + "/"},
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"msg":     "Company " + account.Name + " authenticated",
		"company": account.Name,
		// Raw counters, for consumers that want to build their own view
		"hotspots": hotspots,
		"units":    units,
		"users":    guests,
		"devices":  devices,
		"sessions": sessions,
		"managers": managers,
		"sms": gin.H{
			"count":     sms.SmsCount,
			"max":       sms.SmsMaxCount,
			"remaining": smsRemaining,
		},
		// Generic contract rendered by the My dashboard
		"widget": gin.H{
			"items": widgetItems,
		},
	})
}
