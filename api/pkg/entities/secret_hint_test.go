package entities

import "testing"

func TestSecretHint(t *testing.T) {
	cases := map[string]string{
		"":                      "****",
		"abc":                   "****",
		"pk_1234567890abcdWXYZ": "****WXYZ",
	}
	for secret, want := range cases {
		if got := SecretHint(secret); got != want {
			t.Errorf("SecretHint(%q) = %q, want %q", secret, got, want)
		}
	}
}
