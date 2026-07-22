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

// errUserInfo answers with a machine-readable error code and no prose: the
// response is data-only and language-neutral, so the consumer (the My
// dashboard) maps the code to a localized message. The widget is a
// nice-to-have and the frontend stays silent on non-2xx.
func errUserInfo(c *gin.Context, status int, code string) {
	c.JSON(status, gin.H{"error": code})
}

// UserInfo returns aggregate statistics for the company (reseller) of the
// authenticated My user. The My dashboard widget calls it cross-origin with
// the user's Logto ID token as bearer: the token is verified against the
// same issuer used for the OIDC login, and the company is resolved through
// the organization claims exactly like the login flow (never by email).
//
// The response is data-only (raw counters + a link): the consumer composes
// its own widget from these, so there are no prebuilt labels/strings here.
func UserInfo(c *gin.Context) {
	config := configuration.Config

	if config.OIDC.Issuer == "" {
		errUserInfo(c, http.StatusInternalServerError, "oidc_not_configured")
		return
	}

	// Bearer token from the Authorization header
	authParts := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
	if len(authParts) != 2 || !strings.EqualFold(authParts[0], "Bearer") {
		errUserInfo(c, http.StatusUnauthorized, "invalid_authorization")
		return
	}

	provider, err := getOIDCProvider(config.OIDC.Issuer)
	if err != nil {
		log.Println("userinfo: OIDC provider initialization failed:", err)
		errUserInfo(c, http.StatusInternalServerError, "oidc_provider_unavailable")
		return
	}

	// The caller is My's own application, not our third-party client: the
	// audience is checked manually against the configured value
	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
	idToken, err := verifier.Verify(c.Request.Context(), authParts[1])
	if err != nil {
		log.Println("userinfo denied: invalid token:", err)
		errUserInfo(c, http.StatusUnauthorized, "invalid_token")
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
			errUserInfo(c, http.StatusUnauthorized, "invalid_audience")
			return
		}
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		errUserInfo(c, http.StatusUnauthorized, "invalid_token")
		return
	}

	// The company is the organization where the user holds the reseller
	// organization-role, mapped to the account through logto_org_id
	orgIDs := extractOrgIDsWithRole(claims, config.OIDC.OrgResellerRole)
	if len(orgIDs) == 0 {
		log.Printf("userinfo denied for %v: no organization with role %q (organization_roles=%v)",
			claims["email"], config.OIDC.OrgResellerRole, claims["organization_roles"])
		errUserInfo(c, http.StatusForbidden, "no_reseller_organization")
		return
	}

	db := database.Instance()
	var account models.Account
	db.Where("logto_org_id IN (?)", orgIDs).First(&account)
	if account.Id == 0 || account.Type != "reseller" {
		log.Printf("userinfo denied for %v: no company account linked to organizations %v", claims["email"], orgIDs)
		errUserInfo(c, http.StatusNotFound, "company_not_found")
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

	// Data-only response: raw counters + a link to the NethSpot frontend. The
	// consumer builds its own widget (labels/tones/per-row links) from these;
	// no prebuilt strings live here.
	base := strings.TrimRight(config.OIDC.FrontendURL, "/")
	c.JSON(http.StatusOK, gin.H{
		"company":  account.Name,
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
		"link": base + "/",
	})
}
