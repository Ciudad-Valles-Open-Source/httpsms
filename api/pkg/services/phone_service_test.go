package services

import (
	"context"
	"testing"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type phoneCreationRepository struct {
	repositories.PhoneRepository
	saved entities.Phone
}

func (*phoneCreationRepository) Load(context.Context, entities.UserID, string) (*entities.Phone, error) {
	return nil, stacktrace.PropagateWithCodef(gorm.ErrRecordNotFound, repositories.ErrCodeNotFound, "phone not found")
}

func (repository *phoneCreationRepository) Save(_ context.Context, phone *entities.Phone) error {
	repository.saved = *phone
	return nil
}

func TestPhoneUpsertAppliesArchiveSettingOnFirstRegistration(t *testing.T) {
	enabled := true
	disabled := false
	for _, test := range []struct {
		name    string
		setting *bool
		want    bool
	}{
		{name: "enabled", setting: &enabled, want: true},
		{name: "disabled", setting: &disabled},
		{name: "omitted defaults off"},
	} {
		t.Run(test.name, func(t *testing.T) {
			repository := &phoneCreationRepository{}
			logger := &noopLogger{}
			tracer := telemetry.NewOtelLogger("test", logger)
			dispatcher := NewEventDispatcher(logger, tracer, nil, &phoneNotificationEventQueue{}, PushQueueConfig{})
			service := NewPhoneService(logger, tracer, repository, dispatcher)
			number, err := phonenumbers.Parse("+18005550199", phonenumbers.UNKNOWN_REGION)
			require.NoError(t, err)

			phone, err := service.Upsert(context.Background(), &PhoneUpsertParams{
				Source:                           "https://httpsms.test",
				UserID:                           entities.UserID("user-id"),
				PhoneNumber:                      number,
				SIM:                              entities.SIM1,
				ArchiveNewOutgoingMessageThreads: test.setting,
			})

			require.NoError(t, err)
			require.NotNil(t, phone)
			assert.Equal(t, test.want, phone.ArchiveNewOutgoingMessageThreads)
			assert.Equal(t, test.want, repository.saved.ArchiveNewOutgoingMessageThreads)
			assert.Equal(t, "+18005550199", repository.saved.PhoneNumber)
		})
	}
}

func TestPhoneUpdateArchiveNewOutgoingMessageThreads(t *testing.T) {
	enabled := true
	disabled := false
	tests := []struct {
		name    string
		initial bool
		setting *bool
		want    bool
	}{
		{name: "enable", setting: &enabled, want: true},
		{name: "disable", initial: true, setting: &disabled, want: false},
		{name: "omitted preserves enabled", initial: true, want: true},
		{name: "omitted defaults off", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			phone := &entities.Phone{ArchiveNewOutgoingMessageThreads: test.initial}
			params := PhoneUpsertParams{ArchiveNewOutgoingMessageThreads: test.setting}

			updated := (&PhoneService{}).update(phone, &params)

			assert.Equal(t, test.want, updated.ArchiveNewOutgoingMessageThreads)
		})
	}
}

func TestPhoneUpdateKeepsOmittedFcmTokenAndSchedule(t *testing.T) {
	token := "existing-token"
	newToken := "new-token"
	scheduleID := uuid.New()
	newScheduleID := uuid.New()
	reply := "Please text me"

	tests := []struct {
		name         string
		params       PhoneUpsertParams
		wantToken    *string
		wantSchedule *uuid.UUID
		wantReply    *string
	}{
		{
			name:         "omitted fields keep current values",
			params:       PhoneUpsertParams{SIM: entities.SIM1},
			wantToken:    &token,
			wantSchedule: &scheduleID,
			wantReply:    &reply,
		},
		{
			name:         "explicit fcm token updates it",
			params:       PhoneUpsertParams{SIM: entities.SIM1, FcmToken: &newToken},
			wantToken:    &newToken,
			wantSchedule: &scheduleID,
			wantReply:    &reply,
		},
		{
			name:         "explicit schedule replaces it",
			params:       PhoneUpsertParams{SIM: entities.SIM1, MessageSendScheduleID: &newScheduleID, UpdateMessageSendSchedule: true},
			wantToken:    &token,
			wantSchedule: &newScheduleID,
			wantReply:    &reply,
		},
		{
			name:      "explicit detach removes schedule",
			params:    PhoneUpsertParams{SIM: entities.SIM1, UpdateMessageSendSchedule: true},
			wantToken: &token,
			wantReply: &reply,
		},
		{
			name:         "empty auto reply removes it",
			params:       PhoneUpsertParams{SIM: entities.SIM1, MissedCallAutoReply: new(string)},
			wantToken:    &token,
			wantSchedule: &scheduleID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			currentToken := token
			currentSchedule := scheduleID
			currentReply := reply
			phone := &entities.Phone{
				FcmToken:              &currentToken,
				MessageSendScheduleID: &currentSchedule,
				MissedCallAutoReply:   &currentReply,
			}

			updated := (&PhoneService{}).update(phone, &test.params)

			assert.Equal(t, test.wantToken, updated.FcmToken)
			assert.Equal(t, test.wantSchedule, updated.MessageSendScheduleID)
			assert.Equal(t, test.wantReply, updated.MissedCallAutoReply)
		})
	}
}

func TestPhoneUpdateSetsFcmTokenOnPhoneWithoutToken(t *testing.T) {
	token := "first-token"
	phone := &entities.Phone{}

	updated := (&PhoneService{}).update(phone, &PhoneUpsertParams{SIM: entities.SIM1, FcmToken: &token})

	require.NotNil(t, updated.FcmToken)
	assert.Equal(t, token, *updated.FcmToken)
}
