package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlug(t *testing.T) {
	cases := map[string]string{
		"httpSMS":              "httpsms",
		"Arquos SMS!":          "arquos_sms",
		"  Mi App  Ñandú ":     "mi_app_nandu",
		"Crème Brûlée":         "creme_brulee",
		"Straße":               "strasse",
		"api.example-sms":      "api.example-sms",
		"__My--App..":          "my--app",
		"a\t\nb":               "a_b",
		"SMS ✉️ Gateway":       "sms_gateway",
		"日本語":                  "",
		"!!!":                  "",
		"":                     "",
		"  ":                   "",
		"Gateway 2.0 (beta)":   "gateway_2.0_beta",
		"-._ Already_Slug _.-": "already_slug",
	}

	for input, expected := range cases {
		assert.Equal(t, expected, Slug(input), "Slug(%q)", input)
	}
}

func TestAppSlug(t *testing.T) {
	t.Setenv(envAppName, "Arquos SMS!")
	assert.Equal(t, "arquos_sms", AppSlug())

	t.Setenv(envAppName, "")
	assert.Equal(t, DefaultAppSlug, AppSlug())

	t.Setenv(envAppName, "日本語")
	assert.Equal(t, DefaultAppSlug, AppSlug())
}

func TestJWTIssuer(t *testing.T) {
	t.Setenv(envAppName, "Arquos SMS")
	t.Setenv(envJWTIssuer, "")
	assert.Equal(t, "arquos_sms", JWTIssuer())

	t.Setenv(envJWTIssuer, "  sms.example.com  ")
	assert.Equal(t, "sms.example.com", JWTIssuer())

	t.Setenv(envJWTIssuer, "   ")
	t.Setenv(envAppName, "")
	assert.Equal(t, "httpsms", JWTIssuer())
}

func TestOTelServiceName(t *testing.T) {
	t.Setenv(envAppName, "Mi App")
	t.Setenv(envOTelServiceName, "")
	assert.Equal(t, "mi_app", OTelServiceName())

	t.Setenv(envOTelServiceName, " sms-api ")
	assert.Equal(t, "sms-api", OTelServiceName())
}
