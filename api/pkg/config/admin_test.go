package config

import (
	"testing"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func firebaseAuth(email string, verified bool) entities.AuthContext {
	return entities.AuthContext{ID: "firebase-uid", Email: email, EmailVerified: verified, FirebaseIDToken: true}
}

func TestNewAdminPolicy_ParsesTrimsAndLowercases(t *testing.T) {
	policy := NewAdminPolicy(" Admin@Example.com , ,other@example.com,")

	assert.True(t, policy.Enabled())
	assert.Equal(t, []string{"admin@example.com", "other@example.com"}, policy.Emails())
	assert.True(t, policy.IsAdminEmail("ADMIN@example.COM "))
	assert.False(t, policy.IsAdminEmail("admin@example.co"))
	assert.False(t, policy.IsAdminEmail(""))
}

func TestNewAdminPolicy_EmptyDisablesGate(t *testing.T) {
	for _, raw := range []string{"", " ", " , ,"} {
		policy := NewAdminPolicy(raw)
		assert.False(t, policy.Enabled())
		assert.False(t, policy.IsAdmin(firebaseAuth("admin@example.com", true)))
		assert.Equal(t, entities.UserStatusActive, policy.InitialStatus(firebaseAuth("user@example.com", false)))
	}

	var nilPolicy *AdminPolicy
	assert.False(t, nilPolicy.Enabled())
	assert.False(t, nilPolicy.IsAdmin(firebaseAuth("admin@example.com", true)))
}

func TestAdminPolicy_IsAdmin(t *testing.T) {
	policy := NewAdminPolicy("admin@example.com")
	phoneKeyID := uuid.New()

	tests := []struct {
		name     string
		auth     entities.AuthContext
		expected bool
	}{
		{"verified firebase admin", firebaseAuth("admin@example.com", true), true},
		{"verified firebase admin, different case", firebaseAuth("ADMIN@Example.COM", true), true},
		{"unverified firebase admin", firebaseAuth("admin@example.com", false), false},
		{"verified firebase non admin", firebaseAuth("user@example.com", true), false},
		{"user api key of admin", entities.AuthContext{ID: "firebase-uid", Email: "admin@example.com"}, false},
		{"user api key with verified flag but no id token", entities.AuthContext{ID: "firebase-uid", Email: "admin@example.com", EmailVerified: true}, false},
		{"phone api key of admin", entities.AuthContext{ID: "firebase-uid", Email: "admin@example.com", EmailVerified: true, FirebaseIDToken: true, PhoneAPIKeyID: &phoneKeyID}, false},
		{"missing id", entities.AuthContext{Email: "admin@example.com", EmailVerified: true, FirebaseIDToken: true}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, policy.IsAdmin(test.auth))
		})
	}
}

func TestAdminPolicy_InitialStatus(t *testing.T) {
	policy := NewAdminPolicy("admin@example.com")

	assert.Equal(t, entities.UserStatusActive, policy.InitialStatus(firebaseAuth("admin@example.com", true)))
	assert.Equal(t, entities.UserStatusPending, policy.InitialStatus(firebaseAuth("admin@example.com", false)))
	assert.Equal(t, entities.UserStatusPending, policy.InitialStatus(firebaseAuth("user@example.com", true)))
	assert.Equal(t, entities.UserStatusPending, policy.InitialStatus(entities.AuthContext{ID: "id", Email: "admin@example.com"}))
}
