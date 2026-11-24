#!/usr/bin/env bash
# =============================================================================
# check-secrets.sh — blocks commits that contain secrets.
#
# Runs from the pre-commit hook (web/.husky/pre-commit). It scans ONLY what is
# staged for the commit (added lines + file names) and fails if it finds:
#   - files that must never be committed (keystores, private keys, Firebase /
#     GCP credentials, real .env files)
#   - well-known key formats (Google API keys, private key blocks, GitHub,
#     Slack, Stripe, AWS, Sentry DSNs, Axiom tokens, JWTs…)
#   - non-empty values for secret-looking variables in .env-style files
#     (…_SECRET, …_PASSWORD, …_TOKEN, …_API_KEY, …_CREDENTIALS, …_DSN)
#
# False positive? Add the marker  secret-scan: ignore  on that line (in a
# comment), or use a placeholder such as YOUR_API_KEY / changeme / <value>.
# Run manually:  bash scripts/check-secrets.sh
# In CI (scan every commit of a push/PR):
#   bash scripts/check-secrets.sh --range <base>..<head>
# =============================================================================
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

# Default: what is staged for the next commit. --range A..B: changes between
# two commits (used by CI).
diff_args=(--cached)
if [ "${1:-}" = "--range" ]; then
  [ -n "${2:-}" ] || { echo "usage: $0 [--range <base>..<head>]" >&2; exit 2; }
  diff_args=("$2")
fi

red() { printf '\033[31m%s\033[0m\n' "$*"; }
found=0

# ── 1. File names that must never be committed ──────────────────────────────
forbidden_files='(^|/)(\.env|\.env\.local|\.env\.[^/]*\.local|key\.properties|google-services\.json|GoogleService-Info\.plist)$|\.(jks|keystore|p12|pfx|pem|key)$|firebase-adminsdk[^/]*\.json$|service-account[^/]*\.json$|credentials[^/]*\.json$'

while IFS= read -r file; do
  [ -z "$file" ] && continue
  if printf '%s\n' "$file" | grep -Eq "$forbidden_files"; then
    red "✖ Sensitive file staged: $file"
    found=1
  fi
done < <(git diff "${diff_args[@]}" --name-only --diff-filter=ACMR)

# ── 2. Secret patterns in added lines ───────────────────────────────────────
# One Perl pass over the whole diff (fast even for huge diffs such as merges).
scan_content() {
  perl -ne '
    BEGIN {
      @patterns = (
        qr/AIza[0-9A-Za-z_-]{35}/,                                      # Google / Firebase API key
        qr/-----BEGIN (?:[A-Z]+ )?PRIVATE KEY-----/,                    # PEM private key
        qr/"private_key"\s*:\s*"-----BEGIN/,                            # service account JSON
        qr/gh[pousr]_[0-9A-Za-z]{36,}/,                                 # GitHub token
        qr/github_pat_[0-9A-Za-z_]{40,}/,                               # GitHub fine-grained token
        qr/xox[abposr]-[0-9A-Za-z-]{10,}/,                              # Slack token
        qr/(?:sk|rk)_live_[0-9A-Za-z]{16,}/,                            # Stripe live key
        qr/AKIA[0-9A-Z]{16}/,                                           # AWS access key id
        qr{https://[0-9a-f]{32}\@[a-z0-9.-]*ingest[a-z0-9.-]*sentry\.io}, # Sentry DSN
        qr/xaat-[0-9a-f-]{20,}/,                                        # Axiom token
        qr/eyJ[0-9A-Za-z_-]{10,}\.eyJ[0-9A-Za-z_-]{10,}\.[0-9A-Za-z_-]{10,}/, # JWT
      );
      # Secret-looking variable names in .env-style files (KEY=value) …
      $env_var = qr/^\+\s*(?:export\s+)?[A-Z0-9_]*(?:SECRET|PASSWORD|PASSWD|TOKEN|API_KEY|PRIVATE_KEY|CREDENTIALS|DSN)[A-Z0-9_]*\s*=\s*["\x27]?[^"\x27\s#]/;
      # … unless the value is clearly a placeholder / local default.
      $placeholder = qr/=\s*["\x27]?(?:\$\{|<|your|tu_|changeme|change_me|example|placeholder|xxx|test|app-(?:key|secret|id)|system|secret$|password$|dev$|local)/i;
      $bad = 0;
      sub mask { my $s = substr(shift, 0, 80); $s =~ s/([A-Za-z0-9_-]{6})[A-Za-z0-9_.-]{6,}/$1.../g; return $s; }
    }
    if (/^\+\+\+ b\/(.*)$/) { $file = $1; next; }
    next if /^(?:\+\+\+|---)/;
    next unless /^\+/;
    next if /secret-scan: ignore/;
    chomp;
    my $line = substr($_, 1);
    for my $re (@patterns) {
      if ($line =~ $re) {
        print "\e[31m[secret] Possible secret in $file:\e[0m\n    " . mask($line) . "\n";
        $bad = 1;
        last;
      }
    }
    if ($file =~ m{(?:^|/)\.env|\.env(?:\.|$)} && $_ =~ $env_var && $_ !~ $placeholder) {
      (my $shown = $line) =~ s/=.*/=.../;
      print "\e[31m[secret] Secret value in env file $file:\e[0m\n    $shown\n";
      $bad = 1;
    }
    END { exit $bad }
  '
}

if ! git diff "${diff_args[@]}" -U0 --no-color --diff-filter=ACMR | scan_content; then
  found=1
fi

if [ "$found" -ne 0 ]; then
  echo
  red "Blocked: secrets must not be committed."
  echo "Put real values in .env (ignored by git), Coolify environment variables"
  echo "or GitHub Actions secrets. See docs/secrets.md (docs/secrets_ES.md)."
  echo "If it is a false positive, add the comment marker: secret-scan: ignore"
  exit 1
fi
