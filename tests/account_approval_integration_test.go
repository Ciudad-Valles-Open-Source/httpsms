package tests

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	httpsms "github.com/NdoleStudio/httpsms-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Account approval fixtures. The users and phone API keys below are inserted by
// seed.sql, and ADMIN_EMAILS in .env.test enables the approval gate.
const (
	approvalAdminEmail = "admin@httpsms.test" // must match ADMIN_EMAILS in .env.test

	approvalPendingUserID      = "approval-pending-user-id"
	approvalPendingAPIKey      = "approval-pending-api-key"
	approvalPendingPhoneAPIKey = "pk_approval-pending-phone-api-key"
	approvalPendingPhoneNumber = "+18005559001"

	approvalSuspendedUserID      = "approval-suspended-user-id"
	approvalSuspendedAPIKey      = "approval-suspended-api-key"
	approvalSuspendedPhoneAPIKey = "pk_approval-suspended-phone-api-key"
	approvalSuspendedPhoneNumber = "+18005559002"

	approvalPendingDeleteAPIKey = "approval-pending-delete-api-key"

	approvalAdminEmailUserID = "approval-admin-email-user-id"
	approvalAdminEmailAPIKey = "approval-admin-email-api-key"

	accountStatusActive    = "active"
	accountStatusPending   = "pending"
	accountStatusSuspended = "suspended"

	accountCodePending   = "account_pending"
	accountCodeSuspended = "account_suspended"
)

type approvalUserResponse struct {
	Data struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Status  string `json:"status"`
		IsAdmin bool   `json:"is_admin"`
	} `json:"data"`
}

type approvalGateErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Code string `json:"code"`
	} `json:"data"`
}

type approvalBlockedUser struct {
	name        string
	userID      string
	apiKey      string
	phoneAPIKey string
	phoneNumber string
	status      string
	code        string
}

func approvalBlockedUsers() []approvalBlockedUser {
	return []approvalBlockedUser{
		{
			name:        "pending",
			userID:      approvalPendingUserID,
			apiKey:      approvalPendingAPIKey,
			phoneAPIKey: approvalPendingPhoneAPIKey,
			phoneNumber: approvalPendingPhoneNumber,
			status:      accountStatusPending,
			code:        accountCodePending,
		},
		{
			name:        "suspended",
			userID:      approvalSuspendedUserID,
			apiKey:      approvalSuspendedAPIKey,
			phoneAPIKey: approvalSuspendedPhoneAPIKey,
			phoneNumber: approvalSuspendedPhoneNumber,
			status:      accountStatusSuspended,
			code:        accountCodeSuspended,
		},
	}
}

// approvalRequest sends a single request without the 401 retry loop used by
// requestJSON, because several assertions here expect a 401 on purpose.
func approvalRequest(
	ctx context.Context,
	t *testing.T,
	method string,
	path string,
	headers map[string]string,
	payload any,
) (int, []byte) {
	t.Helper()

	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		require.NoError(t, err)
		body = bytes.NewReader(encoded)
	}

	request, err := http.NewRequestWithContext(ctx, method, apiBaseURL+path, body)
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	return response.StatusCode, responseBody
}

func apiKeyHeader(apiKey string) map[string]string {
	return map[string]string{"x-api-key": apiKey}
}

func bearerHeader(token string) map[string]string {
	return map[string]string{"Authorization": "Bearer " + token}
}

func fetchApprovalUser(ctx context.Context, t *testing.T, headers map[string]string) approvalUserResponse {
	t.Helper()

	statusCode, body := approvalRequest(ctx, t, http.MethodGet, "/v1/users/me", headers, nil)
	require.Equal(t, http.StatusOK, statusCode, "GET /v1/users/me failed: %s", string(body))

	var response approvalUserResponse
	require.NoError(t, json.Unmarshal(body, &response), "response: %s", string(body))
	return response
}

func assertAccountGateError(t *testing.T, statusCode int, body []byte, expectedCode string) {
	t.Helper()

	require.Equal(t, http.StatusForbidden, statusCode, "expected the approval gate to reject the request: %s", string(body))

	var response approvalGateErrorResponse
	require.NoError(t, json.Unmarshal(body, &response), "response: %s", string(body))
	assert.Equal(t, "error", response.Status)
	assert.NotEmpty(t, response.Message)
	assert.Equal(t, expectedCode, response.Data.Code, "response: %s", string(body))
}

// forgedFirebaseIDToken builds a token that looks like a Firebase ID token for
// the admin email with email_verified=true, but is signed with a throwaway key
// that Google never issued. The API must not treat it as authenticated.
func forgedFirebaseIDToken(t *testing.T) string {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	now := time.Now().UTC()
	uid := "forged-" + uuid.NewString()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":            "https://securetoken.google.com/httpsms-test",
		"aud":            "httpsms-test",
		"sub":            uid,
		"user_id":        uid,
		"email":          approvalAdminEmail,
		"email_verified": true,
		"auth_time":      now.Add(-time.Minute).Unix(),
		"iat":            now.Add(-time.Minute).Unix(),
		"exp":            now.Add(time.Hour).Unix(),
	})
	token.Header["kid"] = "forged-" + uuid.NewString()

	signed, err := token.SignedString(privateKey)
	require.NoError(t, err)
	return signed
}

func TestAccountApproval_ActiveUserIsUnaffected(t *testing.T) {
	ctx := context.Background()

	user := fetchApprovalUser(ctx, t, apiKeyHeader(userAPIKey))
	assert.Equal(t, "test-user-id", user.Data.ID)
	assert.Equal(t, accountStatusActive, user.Data.Status)
	assert.False(t, user.Data.IsAdmin, "an API key must never report admin rights")

	// "Authorization: Bearer <api key>" is only accepted on the 3CX integration route,
	// so regular routes are exercised with the x-api-key header.
	statusCode, body := approvalRequest(ctx, t, http.MethodGet, "/v1/phones?skip=0&limit=10", apiKeyHeader(userAPIKey), nil)
	assert.Equal(t, http.StatusOK, statusCode, "active user should list phones: %s", string(body))
}

func TestAccountApproval_BlockedUserAPIKeyIsRejected(t *testing.T) {
	ctx := context.Background()

	for _, user := range approvalBlockedUsers() {
		t.Run(user.name, func(t *testing.T) {
			requests := []struct {
				name    string
				method  string
				path    string
				payload any
			}{
				{name: "list phones", method: http.MethodGet, path: "/v1/phones?skip=0&limit=10"},
				{name: "list phone api keys", method: http.MethodGet, path: "/v1/phone-api-keys?skip=0&limit=10"},
				{
					name:   "send message",
					method: http.MethodPost,
					path:   "/v1/messages/send",
					payload: map[string]any{
						"from":    user.phoneNumber,
						"to":      randomPhoneNumber(),
						"content": "blocked account " + uuid.NewString(),
					},
				},
				// Only GET and DELETE /v1/users/me are exempt from the gate.
				{name: "update profile", method: http.MethodPut, path: "/v1/users/me", payload: map[string]any{"timezone": "UTC"}},
				{name: "rotate api key", method: http.MethodDelete, path: fmt.Sprintf("/v1/users/%s/api-keys", user.userID)},
			}

			for _, request := range requests {
				t.Run(request.name+" with x-api-key", func(t *testing.T) {
					statusCode, body := approvalRequest(ctx, t, request.method, request.path, apiKeyHeader(user.apiKey), request.payload)
					assertAccountGateError(t, statusCode, body, user.code)
				})
			}

			// The 3CX integration is the only route that authenticates "Bearer <api key>".
			t.Run("3cx message with bearer api key", func(t *testing.T) {
				statusCode, body := approvalRequest(ctx, t, http.MethodPost, "/integration/3cx/messages", bearerHeader(user.apiKey), map[string]any{
					"from": "+18005550199", "to": "+18005550100", "text": "blocked",
				})
				assertAccountGateError(t, statusCode, body, user.code)
			})
		})
	}
}

func TestAccountApproval_BlockedUserCanReadOwnProfile(t *testing.T) {
	ctx := context.Background()

	for _, user := range approvalBlockedUsers() {
		t.Run(user.name, func(t *testing.T) {
			response := fetchApprovalUser(ctx, t, apiKeyHeader(user.apiKey))
			assert.Equal(t, user.userID, response.Data.ID)
			assert.Equal(t, user.status, response.Data.Status)
			assert.False(t, response.Data.IsAdmin)
		})
	}
}

func TestAccountApproval_BlockedUserPhoneAPIKeyIsRejected(t *testing.T) {
	ctx := context.Background()

	for _, user := range approvalBlockedUsers() {
		t.Run(user.name, func(t *testing.T) {
			statusCode, body := approvalRequest(ctx, t, http.MethodPost, "/v1/heartbeats", apiKeyHeader(user.phoneAPIKey), map[string]any{
				"phone_numbers": []string{user.phoneNumber},
				"charging":      true,
			})
			assertAccountGateError(t, statusCode, body, user.code)

			statusCode, body = approvalRequest(ctx, t, http.MethodPut, "/v1/phones/fcm-token", apiKeyHeader(user.phoneAPIKey), map[string]any{
				"phone_number": user.phoneNumber,
				"fcm_token":    "fcm-" + uuid.NewString(),
				"sim":          "SIM1",
			})
			assertAccountGateError(t, statusCode, body, user.code)

			statusCode, body = approvalRequest(
				ctx,
				t,
				http.MethodGet,
				"/v1/messages/outstanding?message_id="+uuid.NewString(),
				apiKeyHeader(user.phoneAPIKey),
				nil,
			)
			assertAccountGateError(t, statusCode, body, user.code)
		})
	}
}

func TestAccountApproval_PendingUserCanDeleteOwnAccount(t *testing.T) {
	ctx := context.Background()

	// A dedicated seeded user, so deleting it does not affect the other tests.
	statusCode, body := approvalRequest(ctx, t, http.MethodDelete, "/v1/users/me", apiKeyHeader(approvalPendingDeleteAPIKey), nil)
	assert.Equal(t, http.StatusNoContent, statusCode, "pending user should be able to delete their account: %s", string(body))
}

func TestAccountApproval_SystemUserCanDispatchEvents(t *testing.T) {
	ctx := context.Background()

	// The internal system user (EVENTS_QUEUE_USER_ID) must never be stopped by the
	// approval gate, otherwise every asynchronous listener breaks.
	dispatchInternalEvent(ctx, t, map[string]any{
		"specversion":     "1.0",
		"id":              uuid.NewString(),
		"source":          "/tests/account-approval",
		"type":            "integration.test.account-approval",
		"time":            time.Now().UTC().Format(time.RFC3339),
		"datacontenttype": "application/json",
		"data":            map[string]any{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})

	// A pending user cannot reach the internal events endpoint either.
	statusCode, body := approvalRequest(ctx, t, http.MethodPost, "/v1/events", apiKeyHeader(approvalPendingAPIKey), map[string]any{
		"specversion": "1.0",
		"id":          uuid.NewString(),
		"source":      "/tests/account-approval",
		"type":        "integration.test.account-approval",
	})
	assertAccountGateError(t, statusCode, body, accountCodePending)
}

func TestAdminEndpoints_RejectUnauthenticatedRequests(t *testing.T) {
	ctx := context.Background()

	for _, request := range adminEndpointRequests() {
		t.Run(request.name, func(t *testing.T) {
			statusCode, body := approvalRequest(ctx, t, request.method, request.path, nil, request.payload)
			assert.Equal(t, http.StatusUnauthorized, statusCode, "response: %s", string(body))
		})
	}
}

func TestAdminEndpoints_RejectAPIKeys(t *testing.T) {
	ctx := context.Background()

	callers := map[string]map[string]string{
		"active user x-api-key":         apiKeyHeader(userAPIKey),
		"active user bearer api key":    bearerHeader(userAPIKey),
		"admin email x-api-key":         apiKeyHeader(approvalAdminEmailAPIKey),
		"admin email bearer api key":    bearerHeader(approvalAdminEmailAPIKey),
		"pending user x-api-key":        apiKeyHeader(approvalPendingAPIKey),
		"suspended user bearer api key": bearerHeader(approvalSuspendedAPIKey),
	}

	for callerName, headers := range callers {
		for _, request := range adminEndpointRequests() {
			t.Run(callerName+" "+request.name, func(t *testing.T) {
				statusCode, body := approvalRequest(ctx, t, request.method, request.path, headers, request.payload)
				if strings.Contains(callerName, "bearer") {
					// "Bearer <api key>" does not authenticate on these routes at all.
					assert.Contains(t, []int{http.StatusUnauthorized, http.StatusForbidden}, statusCode, "response: %s", string(body))
					return
				}
				assert.Equal(t, http.StatusForbidden, statusCode, "response: %s", string(body))
			})
		}
	}

	// The rejected calls above must not have changed anything.
	assert.Equal(t, accountStatusPending, fetchApprovalUser(ctx, t, apiKeyHeader(approvalPendingAPIKey)).Data.Status)
	assert.Equal(t, accountStatusSuspended, fetchApprovalUser(ctx, t, apiKeyHeader(approvalSuspendedAPIKey)).Data.Status)
	statusCode, body := approvalRequest(ctx, t, http.MethodGet, "/v1/phones?skip=0&limit=10", apiKeyHeader(approvalPendingAPIKey), nil)
	assertAccountGateError(t, statusCode, body, accountCodePending)
}

func TestAdminEndpoints_AdminEmailAPIKeyIsNotAdmin(t *testing.T) {
	ctx := context.Background()

	// The seeded user's email is listed in ADMIN_EMAILS, but an API key is not a
	// verified Firebase identity, so it must not be reported as admin.
	user := fetchApprovalUser(ctx, t, apiKeyHeader(approvalAdminEmailAPIKey))
	assert.Equal(t, approvalAdminEmailUserID, user.Data.ID)
	assert.Equal(t, approvalAdminEmail, user.Data.Email)
	assert.Equal(t, accountStatusActive, user.Data.Status)
	assert.False(t, user.Data.IsAdmin, "an API key must never grant admin rights")
}

func TestAdminEndpoints_RejectPhoneAPIKeys(t *testing.T) {
	ctx := context.Background()

	apiKeyResponse, response, err := newAPIClient().PhoneAPIKeys.Store(ctx, &httpsms.PhoneAPIKeyStoreParams{
		Name: "approval-admin-test-key-" + uuid.NewString(),
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.HTTPResponse.StatusCode, "phone api key store failed")
	phoneAPIKey := apiKeyResponse.Data.APIKey
	require.NotEmpty(t, phoneAPIKey)

	// Phone API keys only authenticate phone routes, so the admin routes see the
	// request as unauthenticated (401) or forbidden (403); never as allowed.
	for _, request := range adminEndpointRequests() {
		t.Run(request.name, func(t *testing.T) {
			statusCode, body := approvalRequest(ctx, t, request.method, request.path, apiKeyHeader(phoneAPIKey), request.payload)
			assert.Contains(t, []int{http.StatusUnauthorized, http.StatusForbidden}, statusCode, "response: %s", string(body))
		})
	}
}

func TestAdminEndpoints_RejectForgedFirebaseToken(t *testing.T) {
	ctx := context.Background()
	headers := bearerHeader(forgedFirebaseIDToken(t))

	for _, request := range adminEndpointRequests() {
		t.Run(request.name, func(t *testing.T) {
			statusCode, body := approvalRequest(ctx, t, request.method, request.path, headers, request.payload)
			assert.Equal(t, http.StatusUnauthorized, statusCode, "response: %s", string(body))
		})
	}

	statusCode, body := approvalRequest(ctx, t, http.MethodGet, "/v1/users/me", headers, nil)
	assert.Equal(t, http.StatusUnauthorized, statusCode, "forged token must not authenticate: %s", string(body))
}

type adminEndpointRequest struct {
	name    string
	method  string
	path    string
	payload any
}

func adminEndpointRequests() []adminEndpointRequest {
	return []adminEndpointRequest{
		{name: "list users", method: http.MethodGet, path: "/v1/admin/users?status=pending&skip=0&limit=20"},
		{
			name:    "approve user",
			method:  http.MethodPut,
			path:    fmt.Sprintf("/v1/admin/users/%s/status", approvalPendingUserID),
			payload: map[string]any{"status": accountStatusActive},
		},
		{
			name:    "reactivate user",
			method:  http.MethodPut,
			path:    fmt.Sprintf("/v1/admin/users/%s/status", approvalSuspendedUserID),
			payload: map[string]any{"status": accountStatusActive},
		},
		{name: "delete user", method: http.MethodDelete, path: fmt.Sprintf("/v1/admin/users/%s", approvalSuspendedUserID)},
	}
}
