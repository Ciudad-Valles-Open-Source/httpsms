package di

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedSystemUser ensures the system user required by the event queue exists.
//
// On first boot (no users in the DB), it creates the system user using the
// EVENTS_QUEUE_USER_ID and EVENTS_QUEUE_USER_API_KEY env vars. If those are
// empty, it auto-generates them, prints them to stdout, and updates the
// corresponding env vars in the current process so the rest of the app picks
// them up immediately.
//
// If users already exist, this function is a no-op.
//
// This mirrors the "seed" pattern used by frameworks like Phoenix (seeds.exs),
// Rails (db:seed), and Django (loaddata) where the first deploy auto-creates
// the essential admin/system records.
func (container *Container) SeedSystemUser(db *gorm.DB) {
	ctxLogger := container.logger.WithService("SeedSystemUser")

	// Check if any users exist
	var count int64
	if err := db.Model(&entities.User{}).Count(&count).Error; err != nil {
		ctxLogger.Error(fmt.Errorf("cannot count users: %w", err))
		return
	}

	if count > 0 {
		ctxLogger.Info("users already exist in the database, skipping system user seed")
		return
	}

	ctxLogger.Info("no users found — seeding system user...")

	// Resolve or generate User ID
	userID := strings.TrimSpace(os.Getenv("EVENTS_QUEUE_USER_ID"))
	if userID == "" {
		userID = uuid.New().String()
		os.Setenv("EVENTS_QUEUE_USER_ID", userID)
		ctxLogger.Info(fmt.Sprintf("auto-generated EVENTS_QUEUE_USER_ID: %s", userID))
	}

	// Resolve or generate API Key
	apiKey := strings.TrimSpace(os.Getenv("EVENTS_QUEUE_USER_API_KEY"))
	if apiKey == "" {
		generated, err := generateSeedAPIKey(64)
		if err != nil {
			ctxLogger.Error(fmt.Errorf("cannot generate API key for system user: %w", err))
			return
		}
		apiKey = "uk_" + generated
		os.Setenv("EVENTS_QUEUE_USER_API_KEY", apiKey)
		ctxLogger.Info(fmt.Sprintf("auto-generated EVENTS_QUEUE_USER_API_KEY: %s", apiKey))
	}

	// Resolve system user email
	email := strings.TrimSpace(os.Getenv("EVENTS_QUEUE_USER_EMAIL"))
	if email == "" {
		email = "system@httpsms.local"
	}

	now := time.Now().UTC()
	systemUser := &entities.User{
		ID:               entities.UserID(userID),
		Email:            email,
		APIKey:           apiKey,
		SubscriptionName: entities.SubscriptionNameFree,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := db.Create(systemUser).Error; err != nil {
		ctxLogger.Error(fmt.Errorf("cannot create system user: %w", err))
		return
	}

	// Print credentials prominently so operators can capture them
	fmt.Println("")
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              SYSTEM USER CREATED SUCCESSFULLY              ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  User ID : %-48s ║\n", userID)
	fmt.Printf("║  API Key : %-48s ║\n", truncateForDisplay(apiKey, 48))
	fmt.Printf("║  Email   : %-48s ║\n", email)
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║  Save these values in your environment variables:          ║")
	fmt.Println("║    EVENTS_QUEUE_USER_ID                                    ║")
	fmt.Println("║    EVENTS_QUEUE_USER_API_KEY                               ║")
	fmt.Println("║                                                            ║")
	fmt.Println("║  If auto-generated, set them before the next restart       ║")
	fmt.Println("║  to avoid creating a new system user.                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println("")

	ctxLogger.Info(fmt.Sprintf("system user [%s] seeded successfully", userID))
}

// generateSeedAPIKey produces a URL-safe random API key of length n.
func generateSeedAPIKey(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("cannot generate %d random bytes: %w", n, err)
	}
	encoded := base64.URLEncoding.EncodeToString(b)
	if len(encoded) < n {
		return encoded, nil
	}
	return encoded[:n], nil
}

// truncateForDisplay shortens a string for display inside the box, adding "…"
func truncateForDisplay(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}
