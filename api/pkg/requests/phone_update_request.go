package requests

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nyaruka/phonenumbers"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/services"
)

// PhoneUpsert is the payload for updating a phone
type PhoneUpsert struct {
	request
	MessagesPerMinute uint   `json:"messages_per_minute" example:"1" validate:"optional"`
	PhoneNumber       string `json:"phone_number" example:"+18005550199"`

	// MessageExpirationSeconds is the duration in seconds after sending a message when it is considered to be expired.
	MessageExpirationSeconds uint `json:"message_expiration_seconds" example:"12345" validate:"optional"`

	// MaxSendAttempts is the number of attempts when sending an SMS message to handle the case where the phone is offline.
	MaxSendAttempts uint `json:"max_send_attempts" example:"2" validate:"optional"`

	// FcmToken is either a Firebase registration token or a public HTTPS adapter callback URL.
	// When it is omitted, null or empty, the phone keeps its current token.
	FcmToken string `json:"fcm_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzd....." validate:"optional"`

	// MissedCallAutoReply is the SMS sent back when the phone misses a call.
	// When it is omitted the current value is kept; null or an empty string removes it.
	MissedCallAutoReply *string `json:"missed_call_auto_reply" example:"e.g. This phone cannot receive calls. Please send an SMS instead." validate:"optional"`

	// UnarchiveThread moves an archived thread back to the inbox when a new message is received on this phone.
	UnarchiveThread bool `json:"unarchive_thread" example:"false" validate:"optional"`

	// ArchiveNewOutgoingMessageThreads archives new conversations started by an outgoing message on this phone.
	ArchiveNewOutgoingMessageThreads bool `json:"archive_new_outgoing_message_threads" example:"false" validate:"optional"`

	// SIM is the SIM slot of the phone in case the phone has more than 1 SIM slot
	SIM string `json:"sim" example:"SIM1"`

	// MessageSendScheduleID is the ID of the send schedule attached to the phone.
	// When it is omitted the current schedule is kept; null or an empty string detaches it.
	MessageSendScheduleID string `json:"message_send_schedule_id,omitempty" example:"32343a19-da5e-4b1b-a767-3298a73703cb" validate:"optional"`
}

// Sanitize sets defaults to MessageOutstanding
func (input *PhoneUpsert) Sanitize() PhoneUpsert {
	input.FcmToken = strings.TrimSpace(input.FcmToken)
	input.PhoneNumber = input.sanitizeAddress(input.PhoneNumber)
	input.SIM = input.sanitizeSIM(input.SIM)
	if input.MissedCallAutoReply != nil {
		input.MissedCallAutoReply = input.sanitizeStringPointer(*input.MissedCallAutoReply)
	}
	return *input
}

// ToUpsertParams converts PhoneUpsert to services.PhoneUpsertParams.
// The body parameter is the raw JSON request body used to detect which fields were explicitly sent.
func (input *PhoneUpsert) ToUpsertParams(user entities.AuthContext, source string, body []byte) *services.PhoneUpsertParams {
	phone, _ := phonenumbers.Parse(input.PhoneNumber, phonenumbers.UNKNOWN_REGION)

	fields := make(map[string]json.RawMessage)
	_ = json.Unmarshal(body, &fields)

	var messagesPerMinute *uint
	if _, exists := fields["messages_per_minute"]; exists {
		messagesPerMinute = &input.MessagesPerMinute
	}

	// An omitted, null or empty fcm_token keeps the phone's current token. Clients like
	// the web dashboard send the whole phone back, so an empty value must not unregister it.
	var fcmToken *string
	if token := strings.TrimSpace(input.FcmToken); token != "" {
		if _, exists := fields["fcm_token"]; exists {
			fcmToken = &token
		}
	}

	// A non-nil empty string tells the service to remove the auto reply.
	missedCallAutoReply := input.MissedCallAutoReply
	if _, exists := fields["missed_call_auto_reply"]; exists && missedCallAutoReply == nil {
		empty := ""
		missedCallAutoReply = &empty
	}

	var timeout *time.Duration
	if _, exists := fields["message_expiration_seconds"]; exists {
		duration := time.Duration(input.MessageExpirationSeconds) * time.Second
		timeout = &duration
	}

	var maxSendAttempts *uint
	if _, exists := fields["max_send_attempts"]; exists {
		maxSendAttempts = &input.MaxSendAttempts
	}

	var unarchiveThread *bool
	if _, exists := fields["unarchive_thread"]; exists {
		unarchiveThread = &input.UnarchiveThread
	}

	var archiveNewOutgoingMessageThreads *bool
	if _, exists := fields["archive_new_outgoing_message_threads"]; exists {
		archiveNewOutgoingMessageThreads = &input.ArchiveNewOutgoingMessageThreads
	}

	// message_send_schedule_id: omitted keeps the current schedule, null or "" detaches it.
	var scheduleID *uuid.UUID
	_, updateSchedule := fields["message_send_schedule_id"]
	if updateSchedule {
		if parsed, err := uuid.Parse(strings.TrimSpace(input.MessageSendScheduleID)); err == nil {
			scheduleID = &parsed
		}
	}

	return &services.PhoneUpsertParams{
		Source:                           source,
		PhoneNumber:                      phone,
		MessagesPerMinute:                messagesPerMinute,
		MissedCallAutoReply:              missedCallAutoReply,
		MessageExpirationDuration:        timeout,
		MaxSendAttempts:                  maxSendAttempts,
		FcmToken:                         fcmToken,
		UserID:                           user.ID,
		SIM:                              entities.SIM(input.SIM),
		UnarchiveThread:                  unarchiveThread,
		ArchiveNewOutgoingMessageThreads: archiveNewOutgoingMessageThreads,
		MessageSendScheduleID:            scheduleID,
		UpdateMessageSendSchedule:        updateSchedule,
	}
}
