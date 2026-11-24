package emails

import (
	"fmt"
	"strings"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/stacktrace"
	"github.com/go-hermes/hermes/v2"
)

func (factory *hermesUserEmailFactory) appName() string {
	if name := strings.TrimSpace(factory.config.AppName); name != "" {
		return name
	}
	return "httpSMS"
}

// AccountPendingApproval is the email sent to an administrator when a new account waits for approval
func (factory *hermesUserEmailFactory) AccountPendingApproval(adminEmail string, user *entities.User) (*Email, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				fmt.Sprintf("A new account was registered on %s with the email %s at %s.", factory.appName(), user.Email, user.CreatedAt.UTC().Format(time.RFC1123)),
				"The account can't use the API until an administrator approves it.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Review the account and approve or suspend it in the administration page.",
					Button: hermes.Button{
						Color:     "#329ef4",
						TextColor: "#FFFFFF",
						Text:      "REVIEW ACCOUNTS",
						Link:      factory.config.URL("/admin/users"),
					},
				},
			},
			Title:     "Hey,",
			Signature: "Cheers",
			Outros: []string{
				"You receive this email because your address is listed in ADMIN_EMAILS.",
			},
		},
	}

	return factory.generate(email, adminEmail, fmt.Sprintf("New %s account waiting for approval: %s", factory.appName(), user.Email))
}

// AccountApproved is the email sent to a user when an administrator approves their account
func (factory *hermesUserEmailFactory) AccountApproved(user *entities.User) (*Email, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				fmt.Sprintf("Good news: an administrator approved your %s account. You can now use the dashboard and the API.", factory.appName()),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Sign in to get started.",
					Button: hermes.Button{
						Color:     "#329ef4",
						TextColor: "#FFFFFF",
						Text:      "OPEN DASHBOARD",
						Link:      factory.config.URL("/"),
					},
				},
			},
			Title:     "Hey,",
			Signature: "Cheers",
		},
	}

	return factory.generate(email, user.Email, fmt.Sprintf("Your %s account has been approved", factory.appName()))
}

func (factory *hermesUserEmailFactory) generate(email hermes.Email, to string, subject string) (*Email, error) {
	html, err := factory.generator.GenerateHTML(email)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate html email")
	}

	text, err := factory.generator.GeneratePlainText(email)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate text email")
	}

	return &Email{
		ToEmail: to,
		Subject: subject,
		HTML:    html,
		Text:    text,
	}, nil
}
