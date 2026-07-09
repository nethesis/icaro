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
	"fmt"
	"html"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
)

// Device pairing lets a headless client (a NethSecurity unit) obtain a
// session token through the browser OIDC login without sharing an origin
// with the frontend:
//
//  1. the unit calls POST /auth/oidc/device/start and receives a public
//     pair_id (put in the browser URL) and a secret device_code (never
//     shown to the browser)
//  2. the unit opens the returned verification_url in the user's browser:
//     it is the normal /auth/oidc/login flow carrying ?pair=<pair_id>
//  3. the OIDC callback, instead of handing the session to the SPA,
//     attaches the freshly minted token to the pairing
//  4. the unit polls POST /auth/oidc/device/poll with the device_code and
//     receives the token once the user has completed the login
//
// The token is released only to the holder of the device_code, so knowing
// the pair_id (visible in the URL) is never enough to steal a session.

const (
	// how long a pairing stays valid waiting for the user to log in
	devicePairingTTL = 10 * time.Minute

	// suggested client poll interval, returned by /device/start
	devicePairingPollSeconds = 2

	// same anti-flood bound used for pending OIDC states
	maxPendingPairings = 10000
)

type pairingStatus int

const (
	pairingPending pairingStatus = iota
	pairingReady
	pairingFailed
)

// PairingData tracks one device pairing attempt from creation to the
// one-shot token pickup
type PairingData struct {
	DeviceCode  string
	UnitName    string
	Status      pairingStatus
	Expires     time.Time
	Token       string
	TokenExpiry time.Time
	AccountID   int
	AccountName string
	LoggedBy    string
	ErrorCode   string
}

// PairingStore manages pending device pairings, indexed both by pair_id
// (used by the browser flow) and by device_code (used by the polling unit)
type PairingStore struct {
	mu       sync.Mutex
	byPairID map[string]*PairingData
	byDevice map[string]string // device_code -> pair_id
}

var pairingStore = &PairingStore{
	byPairID: make(map[string]*PairingData),
	byDevice: make(map[string]string),
}

// StartPairing registers a new pending pairing. Returns false if the store
// is full.
func (p *PairingStore) StartPairing(pairID string, deviceCode string, unitName string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.byPairID) >= maxPendingPairings {
		now := time.Now()
		for id, data := range p.byPairID {
			if now.After(data.Expires) {
				delete(p.byDevice, data.DeviceCode)
				delete(p.byPairID, id)
			}
		}
		if len(p.byPairID) >= maxPendingPairings {
			return false
		}
	}

	p.byPairID[pairID] = &PairingData{
		DeviceCode: deviceCode,
		UnitName:   unitName,
		Status:     pairingPending,
		Expires:    time.Now().Add(devicePairingTTL),
	}
	p.byDevice[deviceCode] = pairID
	return true
}

// IsPending tells whether a pairing exists and is still waiting for a
// login: the OIDC login entry point uses it to reject unknown or already
// consumed pair ids before starting a flow
func (p *PairingStore) IsPending(pairID string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, exists := p.byPairID[pairID]
	return exists && data.Status == pairingPending && time.Now().Before(data.Expires)
}

// PendingUnitName returns the unit name declared for a pending pairing
// (empty if the pairing is unknown or no longer pending)
func (p *PairingStore) PendingUnitName(pairID string) string {
	p.mu.Lock()
	defer p.mu.Unlock()

	if data, exists := p.byPairID[pairID]; exists && data.Status == pairingPending {
		return data.UnitName
	}
	return ""
}

// Complete attaches the minted session token to a pending pairing and
// returns the unit name declared when the pairing was started
func (p *PairingStore) Complete(pairID string, token string, tokenExpiry time.Time, accountID int, accountName string, loggedBy string) (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, exists := p.byPairID[pairID]
	if !exists || data.Status != pairingPending || time.Now().After(data.Expires) {
		return "", false
	}

	data.Status = pairingReady
	data.Token = token
	data.TokenExpiry = tokenExpiry
	data.AccountID = accountID
	data.AccountName = accountName
	data.LoggedBy = loggedBy
	return data.UnitName, true
}

// Fail marks a pending pairing as failed so the polling unit can report
// the reason instead of timing out
func (p *PairingStore) Fail(pairID string, errorCode string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, exists := p.byPairID[pairID]
	if !exists || data.Status != pairingPending {
		return
	}

	data.Status = pairingFailed
	data.ErrorCode = errorCode
}

// Poll returns the current state of a pairing looked up by device_code.
// Terminal states (ready, failed) are one-shot: the entry is removed as
// it is returned.
func (p *PairingStore) Poll(deviceCode string) (PairingData, pairingStatus, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pairID, exists := p.byDevice[deviceCode]
	if !exists {
		return PairingData{}, pairingPending, false
	}
	data := p.byPairID[pairID]

	if time.Now().After(data.Expires) {
		delete(p.byDevice, deviceCode)
		delete(p.byPairID, pairID)
		return PairingData{}, pairingPending, false
	}

	if data.Status == pairingPending {
		return PairingData{}, pairingPending, true
	}

	// terminal state: hand it out once and forget the pairing
	delete(p.byDevice, deviceCode)
	delete(p.byPairID, pairID)
	return *data, data.Status, true
}

// CleanupExpiredPairings removes expired pairings (called by the periodic
// cleanup routine together with states and codes)
func (p *PairingStore) CleanupExpiredPairings() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for id, data := range p.byPairID {
		if now.After(data.Expires) {
			delete(p.byDevice, data.DeviceCode)
			delete(p.byPairID, id)
		}
	}
}

// oidcLoginURL derives the public /auth/oidc/login URL from the configured
// callback redirect URI (same host, sibling path)
func oidcLoginURL() string {
	return strings.TrimSuffix(configuration.Config.OIDC.RedirectURI, "/callback") + "/login"
}

// OIDCDeviceStart creates a new device pairing and returns the URL the
// unit must open in the user's browser plus the secret device_code used
// to poll for the token. The optional unit_name is echoed back to the
// user on the pairing result page.
func OIDCDeviceStart(c *gin.Context) {
	config := configuration.Config

	if err := validateOIDCConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "OIDC configuration error"})
		return
	}

	var request struct {
		UnitName string `json:"unit_name"`
	}
	// the body is optional: an empty or invalid one just means no unit name
	_ = c.ShouldBindJSON(&request)
	if len(request.UnitName) > 128 {
		request.UnitName = request.UnitName[:128]
	}

	deviceCode, err := generateSecureRandomString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate device code"})
		return
	}
	pairID, err := generateSecureRandomString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate pair id"})
		return
	}

	if !pairingStore.StartPairing(pairID, deviceCode, request.UnitName) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Too many pending pairing attempts, retry later"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_code":      deviceCode,
		"pair_id":          pairID,
		"verification_url": fmt.Sprintf("%s?pair=%s", oidcLoginURL(), pairID),
		"expires_in":       int(devicePairingTTL.Seconds()),
		"interval":         devicePairingPollSeconds,
	})
}

// OIDCDevicePoll returns the pairing outcome to the unit holding the
// device_code: pending while the user logs in, then exactly once the
// token (ready) or the error (failed)
func OIDCDevicePoll(c *gin.Context) {
	var request struct {
		DeviceCode string `json:"device_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	data, status, found := pairingStore.Poll(request.DeviceCode)
	if !found {
		c.JSON(http.StatusOK, gin.H{"status": "expired"})
		return
	}

	switch status {
	case pairingReady:
		c.JSON(http.StatusOK, gin.H{
			"status":       "ready",
			"token":        data.Token,
			"expires":      data.TokenExpiry.Unix(),
			"account_id":   data.AccountID,
			"account_name": data.AccountName,
			"logged_by":    data.LoggedBy,
		})
	case pairingFailed:
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error":  data.ErrorCode,
		})
	default:
		c.JSON(http.StatusOK, gin.H{"status": "pending"})
	}
}

// pairingDetail is a label/value row shown on the pairing result page
type pairingDetail struct {
	Label string
	Value string
}

// pairingResultPage renders the page shown in the popup at the end of a
// pairing flow. It is styled after the NethSecurity UI (Tailwind gray
// palette, cyan primary) so the whole pairing feels like a single product;
// the outcome summary also lives in the unit UI that initiated the flow.
func pairingResultPage(c *gin.Context, success bool, title string, message string, details []pairingDetail) {
	icon := `<svg class="icon ok" viewBox="0 0 512 512" fill="currentColor"><path d="M256 512A256 256 0 1 0 256 0a256 256 0 1 0 0 512zM369 209L241 337c-9.4 9.4-24.6 9.4-33.9 0l-64-64c-9.4-9.4-9.4-24.6 0-33.9s24.6-9.4 33.9 0l47 47L335 175c9.4-9.4 24.6-9.4 33.9 0s9.4 24.6 0 33.9z"/></svg>`
	if !success {
		icon = `<svg class="icon err" viewBox="0 0 512 512" fill="currentColor"><path d="M256 512A256 256 0 1 0 256 0a256 256 0 1 0 0 512zm0-384c13.3 0 24 10.7 24 24V264c0 13.3-10.7 24-24 24s-24-10.7-24-24V152c0-13.3 10.7-24 24-24zM224 352a32 32 0 1 1 64 0 32 32 0 1 1 -64 0z"/></svg>`
	}

	rows := ""
	for _, d := range details {
		if d.Value == "" {
			continue
		}
		rows += fmt.Sprintf(`<div class="row"><span class="label">%s</span><span class="value">%s</span></div>`,
			html.EscapeString(d.Label), html.EscapeString(d.Value))
	}
	if rows != "" {
		rows = `<div class="details">` + rows + `</div>`
	}

	page := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600&display=swap" rel="stylesheet">
<style>
  :root {
    --gray-50: #f9fafb; --gray-100: #f3f4f6; --gray-200: #e5e7eb; --gray-400: #9ca3af;
    --gray-500: #6b7280; --gray-600: #4b5563; --gray-700: #374151; --gray-800: #1f2937;
    --gray-900: #111827; --gray-950: #030712;
    --green-500: #22c55e; --green-600: #16a34a; --red-500: #ef4444; --red-600: #dc2626;
  }
  * { box-sizing: border-box; }
  body {
    font-family: "Poppins", ui-sans-serif, system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    display: flex; align-items: center; justify-content: center; min-height: 100vh;
    margin: 0; padding: 1.5rem; background: var(--gray-50); color: var(--gray-700);
  }
  .card {
    background: #fff; border: 1px solid var(--gray-200); border-radius: .5rem;
    padding: 2rem 2.5rem; width: 100%%; max-width: 34rem;
    box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
    text-align: center;
  }
  .icon { width: 3rem; height: 3rem; margin-bottom: 1rem; }
  .icon.ok { color: var(--green-600); }
  .icon.err { color: var(--red-600); }
  h1 { font-size: 1.125rem; font-weight: 600; color: var(--gray-900); margin: 0 0 .5rem; }
  .message { font-size: .875rem; font-weight: 400; color: var(--gray-500); margin: 0 0 1.25rem; }
  .details { text-align: left; margin-bottom: 1.25rem; font-size: .875rem; }
  .row { display: flex; justify-content: space-between; gap: 1rem; padding: .75rem 0; }
  .row + .row { border-top: 1px solid var(--gray-200); }
  .label { color: var(--gray-500); white-space: nowrap; }
  .value { color: var(--gray-900); font-weight: 500; overflow-wrap: anywhere; text-align: right; }
  .hint { font-size: .875rem; font-weight: 400; color: var(--gray-500); margin: 0; }
  @media (prefers-color-scheme: dark) {
    body { background: var(--gray-950); color: var(--gray-200); }
    .card { background: var(--gray-900); border-color: var(--gray-700); }
    h1 { color: var(--gray-50); }
    .message, .hint, .label { color: var(--gray-400); }
    .row + .row { border-color: var(--gray-700); }
    .value { color: var(--gray-50); }
    .icon.ok { color: var(--green-500); }
    .icon.err { color: var(--red-500); }
  }
</style>
</head>
<body>
<div class="card">
%s
<h1>%s</h1>
<p class="message">%s</p>
%s
<p class="hint">You can close this window and go back to your unit.</p>
</div>
</body>
</html>`, html.EscapeString(title), icon, html.EscapeString(title), html.EscapeString(message), rows)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
}
