package config

import (
	"os"
	"sort"
	"strings"

	"github.com/NdoleStudio/httpsms/pkg/entities"
)

const envAdminEmails = "ADMIN_EMAILS"

// AdminPolicy decides who is an administrator and whether new accounts need
// an administrator's approval. It is built from the ADMIN_EMAILS environment
// variable: a comma-separated list of emails compared case-insensitively after
// trimming. When the list is empty the approval gate is disabled.
type AdminPolicy struct {
	emails map[string]struct{}
}

// NewAdminPolicy parses a comma-separated list of administrator emails.
func NewAdminPolicy(raw string) *AdminPolicy {
	emails := map[string]struct{}{}
	for _, part := range strings.Split(raw, ",") {
		if email := normalizeEmail(part); email != "" {
			emails[email] = struct{}{}
		}
	}
	return &AdminPolicy{emails: emails}
}

// AdminPolicyFromEnv creates an AdminPolicy from the ADMIN_EMAILS environment variable.
func AdminPolicyFromEnv() *AdminPolicy {
	return NewAdminPolicy(os.Getenv(envAdminEmails))
}

// Enabled is true when at least one administrator is configured, which turns
// on the account approval gate.
func (policy *AdminPolicy) Enabled() bool {
	return policy != nil && len(policy.emails) > 0
}

// Emails returns the configured administrator emails (normalized and sorted).
func (policy *AdminPolicy) Emails() []string {
	if policy == nil {
		return nil
	}
	result := make([]string, 0, len(policy.emails))
	for email := range policy.emails {
		result = append(result, email)
	}
	sort.Strings(result)
	return result
}

// IsAdminEmail checks if an email is in ADMIN_EMAILS. It says nothing about
// how a request was authenticated; use IsAdmin to authorize a request.
func (policy *AdminPolicy) IsAdminEmail(email string) bool {
	if !policy.Enabled() {
		return false
	}
	normalized := normalizeEmail(email)
	if normalized == "" {
		return false
	}
	_, ok := policy.emails[normalized]
	return ok
}

// IsAdmin checks if an authenticated request is made by an administrator.
// This is true ONLY for requests authenticated with a Firebase ID token whose
// email is verified and listed in ADMIN_EMAILS. API keys never grant admin.
func (policy *AdminPolicy) IsAdmin(authUser entities.AuthContext) bool {
	return policy.Enabled() &&
		!authUser.IsNoop() &&
		authUser.FirebaseIDToken &&
		authUser.EmailVerified &&
		authUser.PhoneAPIKeyID == nil &&
		policy.IsAdminEmail(authUser.Email)
}

// InitialStatus is the status of an account created on the first login of authUser.
func (policy *AdminPolicy) InitialStatus(authUser entities.AuthContext) entities.UserStatus {
	if !policy.Enabled() || policy.IsAdmin(authUser) {
		return entities.UserStatusActive
	}
	return entities.UserStatusPending
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
