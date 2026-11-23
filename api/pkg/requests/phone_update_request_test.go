package requests

import (
	"encoding/json"
	"testing"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPhoneUpsertArchiveNewOutgoingMessageThreadsPresence(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		present bool
		want    bool
	}{
		{name: "enabled", body: `{"phone_number":"+18005550199","archive_new_outgoing_message_threads":true}`, present: true, want: true},
		{name: "disabled", body: `{"phone_number":"+18005550199","archive_new_outgoing_message_threads":false}`, present: true, want: false},
		{name: "omitted by older client", body: `{"phone_number":"+18005550199"}`, present: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var input PhoneUpsert
			require.NoError(t, json.Unmarshal([]byte(test.body), &input))

			params := input.ToUpsertParams(entities.AuthContext{}, "test", []byte(test.body))

			if !test.present {
				assert.Nil(t, params.ArchiveNewOutgoingMessageThreads)
				return
			}
			require.NotNil(t, params.ArchiveNewOutgoingMessageThreads)
			assert.Equal(t, test.want, *params.ArchiveNewOutgoingMessageThreads)
		})
	}
}

func TestPhoneUpsertFcmTokenPresence(t *testing.T) {
	tests := []struct {
		name string
		body string
		want *string
	}{
		{name: "omitted keeps token", body: `{"phone_number":"+18005550199","sim":"SIM1"}`},
		{name: "null keeps token", body: `{"phone_number":"+18005550199","sim":"SIM1","fcm_token":null}`},
		{name: "empty keeps token", body: `{"phone_number":"+18005550199","sim":"SIM1","fcm_token":"  "}`},
		{name: "explicit token updates", body: `{"phone_number":"+18005550199","sim":"SIM1","fcm_token":" new-token "}`, want: ptr("new-token")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var input PhoneUpsert
			require.NoError(t, json.Unmarshal([]byte(test.body), &input))
			input.Sanitize()

			params := input.ToUpsertParams(entities.AuthContext{}, "test", []byte(test.body))

			assert.Equal(t, test.want, params.FcmToken)
		})
	}
}

func TestPhoneUpsertMessageSendSchedulePresence(t *testing.T) {
	scheduleID := "32343a19-da5e-4b1b-a767-3298a73703cb"
	tests := []struct {
		name       string
		body       string
		wantUpdate bool
		wantID     *string
	}{
		{name: "omitted keeps schedule", body: `{"phone_number":"+18005550199","sim":"SIM1"}`},
		{name: "null detaches schedule", body: `{"phone_number":"+18005550199","sim":"SIM1","message_send_schedule_id":null}`, wantUpdate: true},
		{name: "empty string detaches schedule", body: `{"phone_number":"+18005550199","sim":"SIM1","message_send_schedule_id":""}`, wantUpdate: true},
		{name: "explicit schedule attaches", body: `{"phone_number":"+18005550199","sim":"SIM1","message_send_schedule_id":"` + scheduleID + `"}`, wantUpdate: true, wantID: &scheduleID},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var input PhoneUpsert
			require.NoError(t, json.Unmarshal([]byte(test.body), &input))
			input.Sanitize()

			params := input.ToUpsertParams(entities.AuthContext{}, "test", []byte(test.body))

			assert.Equal(t, test.wantUpdate, params.UpdateMessageSendSchedule)
			if test.wantID == nil {
				assert.Nil(t, params.MessageSendScheduleID)
				return
			}
			require.NotNil(t, params.MessageSendScheduleID)
			assert.Equal(t, *test.wantID, params.MessageSendScheduleID.String())
		})
	}
}

func TestPhoneUpsertMissedCallAutoReplyPresence(t *testing.T) {
	tests := []struct {
		name string
		body string
		want *string
	}{
		{name: "omitted keeps reply", body: `{"phone_number":"+18005550199","sim":"SIM1"}`},
		{name: "null clears reply", body: `{"phone_number":"+18005550199","sim":"SIM1","missed_call_auto_reply":null}`, want: ptr("")},
		{name: "empty clears reply", body: `{"phone_number":"+18005550199","sim":"SIM1","missed_call_auto_reply":"  "}`, want: ptr("")},
		{name: "explicit reply updates", body: `{"phone_number":"+18005550199","sim":"SIM1","missed_call_auto_reply":" Text me "}`, want: ptr("Text me")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var input PhoneUpsert
			require.NoError(t, json.Unmarshal([]byte(test.body), &input))
			input.Sanitize()

			params := input.ToUpsertParams(entities.AuthContext{}, "test", []byte(test.body))

			assert.Equal(t, test.want, params.MissedCallAutoReply)
		})
	}
}

func ptr(value string) *string {
	return &value
}
