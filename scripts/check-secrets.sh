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
# =============================================================================
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

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
done < <(git diff --cached --name-only --diff-filter=ACMR)

# ── 2. Secret patterns in added lines ───────────────────────────────────────
patterns=(
  'AIza[0-9A-Za-z_-]{35}'                                    # Google / Firebase API key
  '-----BEGIN ([A-Z]+ )?PRIVATE KEY-----'                    # PEM private key
  '"private_key"[[:space:]]*:[[:space:]]*"-----BEGIN'        # service account JSON
  'gh[pousr]_[0-9A-Za-z]{36,}'                               # GitHub token
  'github_pat_[0-9A-Za-z_]{40,}'                             # GitHub fine-grained token
  'xox[abposr]-[0-9A-Za-z-]{10,}'                            # Slack token
  '(sk|rk)_live_[0-9A-Za-z]{16,}'                            # Stripe live key
  'AKIA[0-9A-Z]{16}'                                         # AWS access key id
  'https://[0-9a-f]{32}@[a-z0-9.-]*ingest[a-z0-9.-]*sentry\.io' # Sentry DSN
  'xaat-[0-9a-f-]{20,}'                                      # Axiom token
  'eyJ[0-9A-Za-z_-]{10,}\.eyJ[0-9A-Za-z_-]{10,}\.[0-9A-Za-z_-]{10,}' # JWT
)

# Secret-looking variable names in .env-style files (KEY=value).
env_secret_var='^\+[[:space:]]*(export[[:space:]]+)?[A-Z0-9_]*(SECRET|PASSWORD|PASSWD|TOKEN|API_KEY|PRIVATE_KEY|CREDENTIALS|DSN)[A-Z0-9_]*[[:space:]]*=[[:space:]]*["'"'"']?[^"'"'"'[:space:]#]'
# Values that are clearly placeholders / local defaults.
placeholder='=[[:space:]]*["'"'"']?(\$\{|<|your|tu_|changeme|change_me|example|placeholder|xxx|test|app-(key|secret|id)|system|secret$|password$|dev$|local)'

current_file=''
while IFS= read -r line; do
  case "$line" in
    '+++ b/'*) current_file="${line#+++ b/}"; continue ;;
    '+++ '*|'---'*) continue ;;
    '+'*) ;;
    *) continue ;;
  esac
  printf '%s' "$line" | grep -q 'secret-scan: ignore' && continue

  for re in "${patterns[@]}"; do
    if printf '%s' "$line" | grep -Eq -- "$re"; then
      red "✖ Possible secret in $current_file:"
      printf '    %s\n' "$(printf '%s' "${line:1}" | cut -c1-80 | sed -E 's/([A-Za-z0-9_-]{6})[A-Za-z0-9_.-]{6,}/\1…/g')"
      found=1
      break
    fi
  done

  case "$current_file" in
    *.env|*.env.*|.env*|*/.env*)
      if printf '%s' "$line" | grep -Eq -- "$env_secret_var" &&
        ! printf '%s' "$line" | grep -Eiq -- "$placeholder"; then
        red "✖ Secret value in env file $current_file:"
        printf '    %s\n' "$(printf '%s' "${line:1}" | sed -E 's/=.*/=…/')"
        found=1
      fi
      ;;
  esac
done < <(git diff --cached -U0 --no-color --diff-filter=ACMR)

if [ "$found" -ne 0 ]; then
  echo
  red "Commit blocked: secrets must not be committed."
  echo "Put real values in .env (ignored by git), Coolify environment variables"
  echo "or GitHub Actions secrets. See docs/secrets.md (docs/secrets_ES.md)."
  echo "If it is a false positive, add the comment marker: secret-scan: ignore"
  exit 1
fi
