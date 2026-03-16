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
	"fmt"
	"net/http"
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

// OIDCStateStore manages OIDC state tokens with expiration
type OIDCStateStore struct {
	mu     sync.RWMutex
	states map[string]time.Time
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
	states: make(map[string]time.Time),
}

// Global code store instance
var codeStore = &OIDCCodeStore{
	codes: make(map[string]OIDCCodeData),
}

// StoreState stores a state token with expiration time
func (s *OIDCStateStore) StoreState(state string, expiration time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[state] = expiration
}

// ValidateAndRemoveState checks if state is valid and removes it
func (s *OIDCStateStore) ValidateAndRemoveState(state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiration, exists := s.states[state]
	if !exists {
		return false
	}

	// Remove the state regardless of expiration (one-time use)
	delete(s.states, state)

	// Check if state has expired
	if time.Now().After(expiration) {
		return false
	}

	return true
}

// CleanupExpiredStates removes expired states (should be called periodically)
func (s *OIDCStateStore) CleanupExpiredStates() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for state, expiration := range s.states {
		if now.After(expiration) {
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

// generateSecureState generates a cryptographically secure random state string
func generateSecureState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
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

// extractRolesFromClaims extracts roles from various possible claim fields
func extractRolesFromClaims(claims map[string]interface{}) []string {
	var roles []string

	// List of possible role claim fields
	roleFields := []string{"roles", "role", "groups", "organizations", "organization_roles"}

	for _, field := range roleFields {
		if fieldValue, exists := claims[field]; exists {
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

		// If we found roles, stop looking
		if len(roles) > 0 {
			break
		}
	}

	return roles
}

// mapRoleToIcaro maps external roles to Icaro roles based on configuration
// Returns empty string if role is not authorized for OIDC login
func mapRoleToIcaro(externalRoles []string, config configuration.Configuration) string {
	// Parse role mapping from configuration: "external:internal" format
	roleMapping := make(map[string]string)

	for _, mapping := range config.OIDC.RoleMapping {
		parts := strings.Split(mapping, ":")
		if len(parts) == 2 {
			roleMapping[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// If no role mapping configured, deny access - no default mappings for security
	if len(roleMapping) == 0 {
		return ""
	}

	// Check for role mappings based on configured roles only
	for _, externalRole := range externalRoles {
		for configRole, icaroRole := range roleMapping {
			if strings.EqualFold(externalRole, configRole) {
				return icaroRole
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "OIDC configuration error", "error": err.Error()})
		return
	}

	ctx := context.Background()

	// Always use auto-discovery for endpoints
	provider, err := oidc.NewProvider(ctx, config.OIDC.Issuer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get OIDC provider", "error": err.Error()})
		return
	}

	oauth2Config := createOAuth2Config(config, provider)

	// Generate secure random state
	state, err := generateSecureState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate state", "error": err.Error()})
		return
	}

	// Store state with expiration (valid for 10 minutes)
	stateStore.StoreState(state, time.Now().Add(10*time.Minute))

	// Get authorization URL
	authURL := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Redirect to authorization URL
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func OIDCCallback(c *gin.Context) {
	config := configuration.Config
	ctx := context.Background()

	// Validate OIDC configuration
	if err := validateOIDCConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "OIDC configuration error", "error": err.Error()})
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

	// Validate state
	if !stateStore.ValidateAndRemoveState(state) {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=invalid_state")
		return
	}

	// Initialize OIDC provider
	provider, err := oidc.NewProvider(ctx, config.OIDC.Issuer)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=provider_init_failed")
		return
	}

	oauth2Config := createOAuth2Config(config, provider)

	// Exchange authorization code for token
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
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
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=token_verification_failed")
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
	name, _ := claims["name"].(string)

	if name == "" {
		name, _ = claims["preferred_username"].(string)
	}
	if name == "" {
		name = email
	}

	if sub == "" || email == "" {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=missing_user_info")
		return
	}

	// Extract and map roles
	roles := extractRolesFromClaims(claims)
	icaroRole := mapRoleToIcaro(roles, config)

	// Block access if role is not authorized for OIDC login
	if icaroRole == "" {
		c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=unauthorized_role")
		return
	}

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
		// For non-admin roles, look for existing account by email first, then Logto sub
		db.Where("email = ?", email).First(&account)

		if account.Id == 0 {
			// If not found by email, try by Logto sub
			db.Where("username = ? OR uuid = ?", sub, sub).First(&account)
		}

		if account.Id == 0 {
			// Create new account for non-admin users if none found
			account = models.Account{
				Username: sub,
				Password: "", // No password for OIDC accounts
				Name:     name,
				Email:    email,
				Type:     icaroRole,
				Uuid:     sub, // Use Logto subject ID as UUID for traceability
				Created:  time.Now().UTC(),
			}

			// Create the account
			if err := db.Create(&account).Error; err != nil {
				c.Redirect(http.StatusTemporaryRedirect, config.OIDC.FrontendURL+"/?error=account_creation_failed")
				return
			}
		}
		// If user exists, do not update - use existing account as is
	}

	// Create authorization token for the session
	h := sha256.New()
	h.Write([]byte(time.Now().UTC().String() + sub + rawIDToken))
	authToken := fmt.Sprintf("%x", h.Sum(nil))

	// Set expiration date
	expires := time.Now().UTC().AddDate(0, 0, configuration.Config.TokenExpiresDays)

	accessToken := models.AccessToken{
		AccountId:   account.Id,
		Token:       authToken,
		Role:        account.Type,
		Type:        "oidc",
		Expires:     expires,
		ACLs:        "full",
		Description: "OIDC Login",
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

	// Return the token information
	c.JSON(http.StatusOK, gin.H{
		"token":        codeData.Token,
		"expires":      accessToken.Expires.Unix(),
		"id":           account.Id,
		"account_type": account.Type,
	})
}
