package entities

// SecretHint returns a safe representation of a secret (API key, token) for
// logs and errors: only its last 4 characters, never the full value.
func SecretHint(secret string) string {
	if len(secret) <= 4 {
		return "****"
	}
	return "****" + secret[len(secret)-4:]
}
