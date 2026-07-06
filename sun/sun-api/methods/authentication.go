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
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/oauth2"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

// Global cleanup context and cancel function
var (
	cleanupCtx    context.Context
	cleanupCancel context.CancelFunc
)

// Cached OIDC provider: discovery and JWKS are fetched once and reused
// across logins instead of hitting the identity provider on every request
var (
	oidcProviderMu sync.Mutex
	oidcProvider   *oidc.Provider
)

// getOIDCProvider returns the cached OIDC provider, initializing it on
// first use. The provider is created with a background context because it
// outlives the request that triggered its initialization.
func getOIDCProvider(issuer string) (*oidc.Provider, error) {
	oidcProviderMu.Lock()
	defer oidcProviderMu.Unlock()

	if oidcProvider != nil {
		return oidcProvider, nil
	}

	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return nil, err
	}

	oidcProvider = provider
	return provider, nil
}

func Login(c *gin.Context) {
	var account models.Account
	var subscription models.Subscription

	var json models.Login
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	username := json.Username
	password := json.Password

	db := database.Instance()
	db.Where("username = ?", username).First(&account)

	if account.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No username found!"})
		return
	} else {
		// check password
		h := md5.New()
		h.Write([]byte(password))
		digest := fmt.Sprintf("%x", h.Sum(nil))

		if account.Password == digest {
			// create authorization token
			h := sha256.New()
			h.Write([]byte(time.Now().UTC().String() + username + password))
			token := fmt.Sprintf("%x", h.Sum(nil))

			// set expiration date
			expires := time.Now().UTC().AddDate(0, 0, configuration.Config.TokenExpiresDays)

			accessToken := models.AccessToken{
				AccountId:   account.Id,
				Token:       token,
				Role:        account.Type,
				Type:        "login",
				Expires:     expires,
				ACLs:        "full",
				Description: "",
			}

			db.Save(&accessToken)

			db.Set("gorm:auto_preload", true)
			if account.Type == "reseller" {
				db.Preload("SubscriptionPlan").Where("account_id = ?", account.Id).First(&subscription)
			} else {
				db.Preload("SubscriptionPlan").Where("account_id = ?", account.CreatorId).First(&subscription)
			}
			subscription.Expired = subscription.ValidUntil.Before(time.Now().UTC())

			c.JSON(http.StatusCreated, gin.H{
				"account_type": account.Type,
				"status":       "success",
				"token":        token,
				"expires":      expires.String(),
				"id":           account.Id,
				"subscription": subscription,
			})

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Password is invalid"})
		}
	}

}

func Logout(c *gin.Context) {
	token := c.GetHeader("Token")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
		return
	} else {
		// delete token
		utils.DeleteToken(token)

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

// Name of the cookie that binds the OIDC state to the browser that
// initiated the login flow (prevents login CSRF / session fixation)
const oidcStateCookieName = "icaro_oidc_state"

// OIDCStateData holds the per-login-attempt secrets generated at /login
// and needed again at /callback
type OIDCStateData struct {
	Expires      time.Time
	Nonce        string
	PKCEVerifier string
}

// OIDCStateStore manages OIDC state tokens with expiration
type OIDCStateStore struct {
	mu     sync.RWMutex
	states map[string]OIDCStateData
}

// OIDCCodeStore manages temporary codes for secure token exchange
type OIDCCodeStore struct {
	mu    sync.RWMutex
	codes map[string]OIDCCodeData
}

type OIDCCodeData struct {
	Token     string
	AccountID int
	Expires   time.Time
}

// Global state store instance
var stateStore = &OIDCStateStore{
	states: make(map[string]OIDCStateData),
}

// Global code store instance
var codeStore = &OIDCCodeStore{
	codes: make(map[string]OIDCCodeData),
}

// maxPendingOIDCStates bounds the memory used by pending login attempts:
// above this threshold new attempts are rejected until older ones expire
// (10 minutes) or complete, so spamming /auth/oidc/login cannot grow the
// store without limit
const maxPendingOIDCStates = 10000

// StoreState stores a state token with its associated per-attempt data.
// Returns false if the store is full.
func (s *OIDCStateStore) StoreState(state string, data OIDCStateData) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.states) >= maxPendingOIDCStates {
		// purge expired entries before giving up
		now := time.Now()
		for st, d := range s.states {
			if now.After(d.Expires) {
				delete(s.states, st)
			}
		}
		if len(s.states) >= maxPendingOIDCStates {
			return false
		}
	}

	s.states[state] = data
	return true
}

// ValidateAndRemoveState checks if state is valid and removes it
func (s *OIDCStateStore) ValidateAndRemoveState(state string) (OIDCStateData, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.states[state]
	if !exists {
		return OIDCStateData{}, false
	}

	// Remove the state regardless of expiration (one-time use)
	delete(s.states, state)

	// Check if state has expired
	if time.Now().After(data.Expires) {
		return OIDCStateData{}, false
	}

	return data, true
}

// CleanupExpiredStates removes expired states (should be called periodically)
func (s *OIDCStateStore) CleanupExpiredStates() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for state, data := range s.states {
		if now.After(data.Expires) {
			delete(s.states, state)
		}
	}
}

// StoreCode stores a temporary code with token and account data
func (c *OIDCCodeStore) StoreCode(code string, token string, accountID int, expiration time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.codes[code] = OIDCCodeData{
		Token:     token,
		AccountID: accountID,
		Expires:   expiration,
	}
}

// ValidateAndRemoveCode checks if code is valid and removes it (one-time use)
func (c *OIDCCodeStore) ValidateAndRemoveCode(code string) (OIDCCodeData, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, exists := c.codes[code]
	if !exists {
		return OIDCCodeData{}, false
	}

	// Remove the code regardless of expiration (one-time use)
	delete(c.codes, code)

	// Check if code has expired
	if time.Now().After(data.Expires) {
		return OIDCCodeData{}, false
	}

	return data, true
}

// CleanupExpiredCodes removes expired codes
func (c *OIDCCodeStore) CleanupExpiredCodes() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for code, data := range c.codes {
		if now.After(data.Expires) {
			delete(c.codes, code)
		}
	}
}

// Initialize cleanup routine with graceful shutdown support
func init() {
	cleanupCtx, cleanupCancel = context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-cleanupCtx.Done():
				return
			case <-ticker.C:
				stateStore.CleanupExpiredStates()
				codeStore.CleanupExpiredCodes()
			}
		}
	}()
}

// StopCleanupRoutine stops the background cleanup goroutine gracefully
func StopCleanupRoutine() {
	if cleanupCancel != nil {
		cleanupCancel()
	}
}

// generateSecureRandomString generates a cryptographically secure random
// URL-safe string (used for state, nonce and PKCE verifier)
func generateSecureRandomString() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// pkceChallengeS256 derives the S256 code challenge from a PKCE verifier
func pkceChallengeS256(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// createOAuth2Config creates the OAuth2 configuration based on the OIDC config
func createOAuth2Config(config configuration.Configuration, provider *oidc.Provider) oauth2.Config {
	scopes := []string{oidc.ScopeOpenID, "profile", "email"}

	// Add configured scopes if available, otherwise use defaults
	if len(config.OIDC.Scopes) > 0 {
		scopes = config.OIDC.Scopes
	} else {
		// Default scopes for Logto
		scopes = append(scopes, "roles", "urn:logto:scope:organizations", "urn:logto:scope:organization_roles")
	}

	return oauth2.Config{
		ClientID:     config.OIDC.ClientID,
		ClientSecret: config.OIDC.ClientSecret,
		RedirectURL:  config.OIDC.RedirectURI,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}
}

// hasRoleIgnoreCase reports whether the given role appears in roles
func hasRoleIgnoreCase(roles []string, want string) bool {
	for _, role := range roles {
		if strings.EqualFold(role, want) {
			return true
		}
	}
	return false
}

// extractOrgIDsWithRole returns the IDs of the organizations where the user
// holds the given role, from the Logto "organization_roles" claim
// (entries have the form "<organization_id>:<role_name>")
func extractOrgIDsWithRole(claims map[string]interface{}, adminRole string) []string {
	var orgIDs []string

	rolesClaim, ok := claims["organization_roles"].([]interface{})
	if !ok {
		return orgIDs
	}

	for _, entry := range rolesClaim {
		entryStr, ok := entry.(string)
		if !ok {
			continue
		}
		idx := strings.LastIndex(entryStr, ":")
		if idx <= 0 {
			continue
		}
		if strings.EqualFold(entryStr[idx+1:], adminRole) {
			orgIDs = append(orgIDs, entryStr[:idx])
		}
	}

	return orgIDs
}

// extractRolesFromClaims extracts roles from all possible claim fields.
// All fields are collected (not just the first non-empty one) so that
// e.g. Logto organization roles are considered even when the user also
// has unrelated global roles.
func extractRolesFromClaims(claims map[string]interface{}) []string {
	var roles []string

	// List of possible role claim fields
	roleFields := []string{"roles", "role", "groups", "organizations", "organization_roles"}

	for _, field := range roleFields {
		fieldValue, exists := claims[field]
		if !exists {
			continue
		}
		// Handle array of roles
		if roleSlice, ok := fieldValue.([]interface{}); ok {
			for _, role := range roleSlice {
				if roleStr, ok := role.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}
		// Handle single role as string
		if roleStr, ok := fieldValue.(string); ok {
			roles = append(roles, roleStr)
		}
	}

	return roles
}

// mapRoleToIcaro maps external roles to Icaro roles based on configuration.
// Mappings are evaluated in configuration order, so the first mapping wins:
// list the most privileged roles first (e.g. "super admin:admin" before
// "admin:reseller"). Returns empty string if no role is authorized.
func mapRoleToIcaro(externalRoles []string, config configuration.Configuration) string {
	// If no role mapping configured, deny access - no default mappings for security
	for _, mapping := range config.OIDC.RoleMapping {
		// Parse role mapping from configuration: "external:internal" format
		parts := strings.Split(mapping, ":")
		if len(parts) != 2 {
			continue
		}
		configRole := strings.TrimSpace(parts[0])
		icaroRole := strings.TrimSpace(parts[1])

		for _, externalRole := range externalRoles {
			if strings.EqualFold(externalRole, configRole) {
				return icaroRole
			}
			// Logto organization roles come as "<organization_id>:<role_name>":
			// match the role name part as well
			if idx := strings.LastIndex(externalRole, ":"); idx >= 0 {
				if strings.EqualFold(externalRole[idx+1:], configRole) {
					return icaroRole
				}
			}
		}
	}

	// Return empty string if no authorized role found
	return ""
}

// validateOIDCConfig validates that OIDC configuration is complete
func validateOIDCConfig(config configuration.Configuration) error {
	if config.OIDC.Issuer == "" {
		return fmt.Errorf("OIDC issuer not configured")
	}
	if config.OIDC.ClientID == "" {
		return fmt.Errorf("OIDC client ID not configured")
	}
	if config.OIDC.ClientSecret == "" {
		return fmt.Errorf("OIDC client secret not configured")
	}
	if config.OIDC.RedirectURI == "" {
		return fmt.Errorf("OIDC redirect URI not configured")
	}
	if config.OIDC.FrontendURL == "" {
		return fmt.Errorf("OIDC frontend URL not configured")
	}
	return nil
}

// myResellerResponse mirrors the relevant part of My's GET /resellers/:id
type myResellerResponse struct {
	Code int `json:"code"`
	Data struct {
		ID          string                 `json:"id"`
		LogtoID     *string                `json:"logto_id"`
		Name        string                 `json:"name"`
		CustomData  map[string]interface{} `json:"custom_data"`
		SuspendedAt *string                `json:"suspended_at"`
	} `json:"data"`
}

// provisionCompanyAccount creates the Icaro company account for a My
// reseller organization on its first OIDC login. The organization is
// verified against the My API (must exist as an active reseller) before
// anything is created. Returns the account, or a frontend error code.
func provisionCompanyAccount(orgID string) (*models.Account, string) {
	config := configuration.Config

	if config.OIDC.MyAPIURL == "" || config.OIDC.MyAPIKey == "" || config.OIDC.DefaultSubscriptionPlanID == 0 {
		// JIT provisioning not configured: only organizations already
		// linked to an account may log in
		log.Println("OIDC company login denied: organization", orgID, "not linked and JIT provisioning not configured")
		return nil, "account_not_found"
	}

	// Verify the organization on My: GET /resellers/:id looks up by
	// Logto organization ID and excludes deleted organizations
	req, err := http.NewRequest("GET", strings.TrimRight(config.OIDC.MyAPIURL, "/")+"/resellers/"+url.PathEscape(orgID), nil)
	if err != nil {
		log.Println("OIDC provisioning: request creation failed:", err)
		return nil, "provisioning_failed"
	}
	req.Header.Set("Authorization", "Bearer "+config.OIDC.MyAPIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("OIDC provisioning: My API unreachable:", err)
		return nil, "provisioning_failed"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// The organization exists on Logto but is not a reseller on My
		log.Println("OIDC company login denied: organization", orgID, "is not a reseller on My")
		return nil, "account_not_found"
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("OIDC provisioning: My API returned status", resp.StatusCode, "for organization", orgID)
		return nil, "provisioning_failed"
	}

	var myResp myResellerResponse
	if err := json.NewDecoder(resp.Body).Decode(&myResp); err != nil {
		log.Println("OIDC provisioning: invalid My API response:", err)
		return nil, "provisioning_failed"
	}

	if myResp.Data.SuspendedAt != nil {
		log.Println("OIDC company login denied: organization", orgID, "is suspended on My")
		return nil, "org_suspended"
	}

	// Organization contact email from custom data, when available
	orgEmail := ""
	for _, key := range []string{"email", "contactEmail", "contact_email"} {
		if v, ok := myResp.Data.CustomData[key].(string); ok && v != "" {
			orgEmail = v
			break
		}
	}

	name := myResp.Data.Name
	if name == "" {
		name = orgID
	}

	// JIT company accounts are created by the primary admin
	db := database.Instance()
	var admin models.Account
	db.Where("type = ?", "admin").Order("id ASC").First(&admin)

	logtoOrgID := orgID
	account := models.Account{
		CreatorId:  admin.Id,
		Uuid:       orgID, // stable unique identifier for accounts born from My organizations
		LogtoOrgId: &logtoOrgID,
		Type:       "reseller",
		Name:       name,
		Username:   orgID,
		Password:   "", // OIDC-only account, no password login
		Email:      orgEmail,
		Created:    time.Now().UTC(),
	}
	if err := db.Create(&account).Error; err != nil {
		log.Println("OIDC provisioning: account creation failed for organization", orgID, ":", err)
		return nil, "account_creation_failed"
	}

	// Create the subscription with the configured default plan and the
	// SMS accounting, mirroring the classic account creation
	var subscriptionPlan models.SubscriptionPlan
	db.Where("id = ?", config.OIDC.DefaultSubscriptionPlanID).First(&subscriptionPlan)
	if subscriptionPlan.ID == 0 {
		log.Println("OIDC provisioning: default subscription plan", config.OIDC.DefaultSubscriptionPlanID, "not found")
		return nil, "account_creation_failed"
	}

	subscription := models.Subscription{
		AccountID:          account.Id,
		SubscriptionPlanID: subscriptionPlan.ID,
		ValidFrom:          time.Now().UTC(),
		ValidUntil:         time.Now().UTC().AddDate(0, 0, subscriptionPlan.Period),
		Created:            time.Now().UTC(),
	}
	db.Save(&subscription)

	accountSMS := models.AccountSmsCount{
		AccountId:   account.Id,
		SmsMaxCount: subscriptionPlan.IncludedSMS,
	}
	db.Save(&accountSMS)

	log.Println("OIDC provisioning: created company account", account.Id, "for My organization", orgID, "("+name+")")
	return &account, ""
}

func GetOIDCConfig(c *gin.Context) {
	config := configuration.Config

	// Return only public configuration information
	response := gin.H{
		"enabled": config.OIDC.Issuer != "" && config.OIDC.ClientID != "",
	}

	// Only add provider name if OIDC is properly configured
	if config.OIDC.Issuer != "" && config.OIDC.ClientID != "" {
		response["provider_name"] = "My Nethesis"
	}

	c.JSON(http.StatusOK, response)
}

func OIDCLogin(c *gin.Context) {
	config := configuration.Config

	// Validate OIDC configuration
	if err := validateOIDCConfig(config); err != nil {
		log.Println("OIDC configuration error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "OIDC configuration error"})
		return
	}

	// Always use auto-discovery for endpoints (cached after first use)
	provider, err := getOIDCProvider(config.OIDC.Issuer)
	if err != nil {
		log.Println("OIDC provider initialization failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get OIDC provider"})
		return
	}

	oauth2Config := createOAuth2Config(config, provider)

	// Generate secure random state, nonce and PKCE verifier
	state, err := generateSecureRandomString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate state"})
		return
	}
	nonce, err := generateSecureRandomString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate nonce"})
		return
	}
	pkceVerifier, err := generateSecureRandomString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate PKCE verifier"})
		return
	}

	// Store state with expiration (valid for 10 minutes)
	stored := stateStore.StoreState(state, OIDCStateData{
		Expires:      time.Now().Add(10 * time.Minute),
		Nonce:        nonce,
		PKCEVerifier: pkceVerifier,
	})
	if !stored {
		log.Println("OIDC state store is full, rejecting login attempt")
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Too many pending login attempts, retry later"})
		return
	}

	// Bind the state to this browser with a short-lived HttpOnly cookie:
	// the callback only accepts a state that matches the cookie, so a state
	// minted for another browser cannot be replayed (login CSRF protection).
	// SameSite=Lax cookies are sent on the top-level GET navigation back
	// from the identity provider.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(oidcStateCookieName, state, 600, "/", "", strings.HasPrefix(config.OIDC.RedirectURI, "https://"), true)

	// Get authorization URL (with nonce and PKCE challenge)
	authURL := oauth2Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("nonce", nonce),
		oauth2.SetAuthURLParam("code_challenge", pkceChallengeS256(pkceVerifier)),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	// Redirect to authorization URL
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func OIDCCallback(c *gin.Context) {
	config := configuration.Config
	ctx := context.Background()

	// Validate OIDC configuration
	if err := validateOIDCConfig(config); err != nil {
		log.Println("OIDC configuration error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "OIDC configuration error"})
		return
	}

	// Get authorization code and state from query parameters
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=missing_code")
		return
	}

	if state == "" {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=missing_state")
		return
	}

	// The state must match the cookie set when this browser started the
	// flow: a valid state minted for a different browser is rejected
	// (login CSRF protection)
	stateCookie, err := c.Cookie(oidcStateCookieName)
	if err != nil || stateCookie != state {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=invalid_state")
		return
	}

	// Clear the state cookie: it is one-time use like the state itself
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(oidcStateCookieName, "", -1, "/", "", strings.HasPrefix(config.OIDC.RedirectURI, "https://"), true)

	// Validate state and retrieve the nonce and PKCE verifier of this attempt
	stateData, valid := stateStore.ValidateAndRemoveState(state)
	if !valid {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=invalid_state")
		return
	}

	// Initialize OIDC provider (cached after first use)
	provider, err := getOIDCProvider(config.OIDC.Issuer)
	if err != nil {
		log.Println("OIDC provider initialization failed:", err)
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=provider_init_failed")
		return
	}

	oauth2Config := createOAuth2Config(config, provider)

	// Exchange authorization code for token (PKCE verifier proves this is
	// the same client that started the flow)
	token, err := oauth2Config.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", stateData.PKCEVerifier))
	if err != nil {
		log.Println("OIDC token exchange failed:", err)
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=token_exchange_failed")
		return
	}

	// Extract and verify ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=no_id_token")
		return
	}

	// Create verifier
	oidcConfig := &oidc.Config{
		ClientID: config.OIDC.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	// Verify ID token
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Println("OIDC ID token verification failed:", err)
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=token_verification_failed")
		return
	}

	// The ID token must carry the nonce generated for this attempt:
	// binds the token to this login flow and prevents replay
	if idToken.Nonce != stateData.Nonce {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=invalid_nonce")
		return
	}

	// Extract claims
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=claims_extraction_failed")
		return
	}

	// Extract user information
	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)

	if sub == "" || email == "" {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=missing_user_info")
		return
	}

	// Extract and map global roles: the only mapping expected here is
	// "super admin:admin" (My super admins act as the Icaro admin); any
	// other user logs in as their company via the organization claims
	roles := extractRolesFromClaims(claims)
	icaroRole := mapRoleToIcaro(roles, config)

	// Handle account mapping based on role
	db := database.Instance()
	var account models.Account

	if icaroRole == "admin" {
		// If user should be admin, find the primary admin account (lowest ID)
		db.Where("type = ?", "admin").Order("id ASC").First(&account)

		if account.Id == 0 {
			// No admin account found, create error
			c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=no_admin_account_found")
			return
		}

		// Use the existing admin account - don't modify it
		// This allows super admin from Logto to act as the main admin
	} else {
		// Company login: on My the global role says what the user may do
		// ("Admin" of their company) and the organization-role encodes
		// the organization type ("Reseller"). Both are required; the
		// organization then maps to the Icaro company account through
		// accounts.logto_org_id (immutable link, never matched by email)
		if !hasRoleIgnoreCase(roles, config.OIDC.OrgAdminRole) {
			log.Printf("OIDC company login denied for %s: missing global role %q (roles=%v)",
				email, config.OIDC.OrgAdminRole, claims["roles"])
			c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=unauthorized_role")
			return
		}

		adminOrgIDs := extractOrgIDsWithRole(claims, config.OIDC.OrgResellerRole)
		if len(adminOrgIDs) == 0 {
			log.Printf("OIDC company login denied for %s: no organization with role %q (organization_roles=%v)",
				email, config.OIDC.OrgResellerRole, claims["organization_roles"])
			c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=unauthorized_role")
			return
		}

		db.Where("logto_org_id IN (?)", adminOrgIDs).First(&account)

		if account.Id == 0 {
			// First login of this organization: verify on My that it is
			// an active reseller, then create the linked company account
			// (the modern equivalent of the old "Get credentials" action)
			newAccount, errCode := provisionCompanyAccount(adminOrgIDs[0])
			if errCode != "" {
				c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error="+errCode)
				return
			}
			account = *newAccount
		}

		// The linked account must be a company (reseller) account
		if account.Type != "reseller" {
			c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=role_mismatch")
			return
		}
	}

	// Create authorization token for the session (256 random bits, same
	// length as the sha256-hex tokens issued by the password login)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=code_generation_failed")
		return
	}
	authToken := fmt.Sprintf("%x", tokenBytes)

	// Set expiration date
	expires := time.Now().UTC().AddDate(0, 0, configuration.Config.TokenExpiresDays)

	// The company account is shared by all org admins: record which My
	// user actually logged in for auditing
	accessToken := models.AccessToken{
		AccountId:   account.Id,
		Token:       authToken,
		Role:        account.Type,
		Type:        "oidc",
		Expires:     expires,
		ACLs:        "full",
		Description: fmt.Sprintf("OIDC login by %s (%s)", email, sub),
	}

	db.Save(&accessToken)

	// Generate a cryptographically secure temporary one-time code (expires in 2 minutes)
	codeBytes := make([]byte, 16)
	if _, err := rand.Read(codeBytes); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=code_generation_failed")
		return
	}
	tempCode := fmt.Sprintf("%x", codeBytes)
	codeExpiration := time.Now().Add(2 * time.Minute)

	// Store the code with token and account information
	codeStore.StoreCode(tempCode, authToken, account.Id, codeExpiration)

	// Redirect to frontend with temporary code instead of token
	callbackURL := fmt.Sprintf("%s/#/login/callback?code=%s",
		config.OIDC.FrontendURL,
		tempCode,
	)

	c.Redirect(http.StatusTemporaryRedirect, callbackURL)
}

// OIDCExchange exchanges a temporary code for authentication token
func OIDCExchange(c *gin.Context) {
	// Get the temporary code from request
	var request struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	// Validate and retrieve the code data
	codeData, valid := codeStore.ValidateAndRemoveCode(request.Code)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired code"})
		return
	}

	// Get account information
	db := database.Instance()
	var account models.Account
	if err := db.Where("id = ?", codeData.AccountID).First(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Account not found"})
		return
	}

	// Get token expiration from database
	var accessToken models.AccessToken
	if err := db.Where("token = ?", codeData.Token).First(&accessToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token not found"})
		return
	}

	// Load the subscription like the password login does
	var subscription models.Subscription
	db.Set("gorm:auto_preload", true)
	if account.Type == "reseller" {
		db.Preload("SubscriptionPlan").Where("account_id = ?", account.Id).First(&subscription)
	} else {
		db.Preload("SubscriptionPlan").Where("account_id = ?", account.CreatorId).First(&subscription)
	}
	subscription.Expired = subscription.ValidUntil.Before(time.Now().UTC())

	// Return the token information
	c.JSON(http.StatusOK, gin.H{
		"token":        codeData.Token,
		"expires":      accessToken.Expires.Unix(),
		"id":           account.Id,
		"account_type": account.Type,
		"subscription": subscription,
	})
}
