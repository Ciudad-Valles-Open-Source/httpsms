// Package config holds small helpers that derive installation-wide settings
// from environment variables so every package resolves them the same way.
package config

import (
	"os"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

const (
	// DefaultAppSlug is used when APP_NAME is empty or has no usable characters.
	DefaultAppSlug = "httpsms"

	envAppName         = "APP_NAME"
	envJWTIssuer       = "JWT_ISSUER"
	envOTelServiceName = "OTEL_SERVICE_NAME"
)

// asciiFolds maps letters that do not decompose into "base letter + accent"
// under Unicode NFD to a reasonable ASCII spelling.
var asciiFolds = map[rune]string{
	'ß': "ss",
	'æ': "ae",
	'œ': "oe",
	'ø': "o",
	'ł': "l",
	'đ': "d",
	'ð': "d",
	'þ': "th",
	'ı': "i",
}

// Slug converts a human readable name into an identifier made only of
// [a-z0-9._-]:
//   - letters are lowercased and accents are stripped to ASCII ("Ñandú" → "nandu",
//     "ß" → "ss"); any other non-ASCII character is removed;
//   - whitespace becomes "_";
//   - every other character outside [a-z0-9._-] is removed;
//   - repeated "_" are collapsed into one;
//   - "_", "-" and "." are trimmed from both ends.
//
// The result may be empty (e.g. for "!!!"); callers decide the fallback.
func Slug(value string) string {
	var builder strings.Builder
	for _, r := range norm.NFD.String(strings.ToLower(value)) {
		switch {
		case unicode.Is(unicode.Mn, r):
			// combining accent left over by NFD: drop it
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '.', r == '-', r == '_':
			builder.WriteRune(r)
		case unicode.IsSpace(r):
			builder.WriteRune('_')
		default:
			if folded, ok := asciiFolds[r]; ok {
				builder.WriteString(folded)
			}
		}
	}

	slug := builder.String()
	for strings.Contains(slug, "__") {
		slug = strings.ReplaceAll(slug, "__", "_")
	}
	return strings.Trim(slug, "_-.")
}

// AppSlug returns the slug of APP_NAME, or DefaultAppSlug ("httpsms") when that is empty.
func AppSlug() string {
	if slug := Slug(os.Getenv(envAppName)); slug != "" {
		return slug
	}
	return DefaultAppSlug
}

// JWTIssuer returns the "iss" claim used in the JWTs this server signs
// (webhook requests and HTTP push notifications): JWT_ISSUER when set,
// otherwise AppSlug().
func JWTIssuer() string {
	if issuer := strings.TrimSpace(os.Getenv(envJWTIssuer)); issuer != "" {
		return issuer
	}
	return AppSlug()
}

// OTelServiceName returns the OpenTelemetry service name: OTEL_SERVICE_NAME
// when set, otherwise AppSlug().
func OTelServiceName() string {
	if name := strings.TrimSpace(os.Getenv(envOTelServiceName)); name != "" {
		return name
	}
	return AppSlug()
}
