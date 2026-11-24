-- Seed test data for integration tests
-- Run AFTER GORM has migrated the schema (i.e., after API starts)

-- Seeded users set status explicitly: ADMIN_EMAILS in .env.test enables the
-- account approval gate, and these users must stay active for the suite.

-- Test user
INSERT INTO users (id, email, api_key, timezone, subscription_name, status, created_at, updated_at)
VALUES (
    'test-user-id',
    'test@httpsms.com',
    'test-user-api-key',
    'UTC',
    'pro-monthly',
    'active',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- Test user for API key rotation tests (isolated to avoid mutating the shared test user)
INSERT INTO users (id, email, api_key, timezone, subscription_name, status, created_at, updated_at)
VALUES (
    'rotate-test-user-id',
    'rotate-test@httpsms.com',
    'rotate-test-api-key',
    'UTC',
    'pro-monthly',
    'active',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- System user (for event queue auth)
INSERT INTO users (id, email, api_key, timezone, subscription_name, status, created_at, updated_at)
VALUES (
    'system-user-id',
    'system@httpsms.com',
    'system-user-api-key',
    'UTC',
    'pro-monthly',
    'active',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- Account approval fixtures (see account_approval_integration_test.go).
-- ADMIN_EMAILS=admin@httpsms.test in .env.test turns the approval gate on.

-- Pending user: blocked everywhere except GET/DELETE /v1/users/me
INSERT INTO users (id, email, api_key, timezone, subscription_name, status, created_at, updated_at)
VALUES (
    'approval-pending-user-id',
    'approval-pending@httpsms.com',
    'approval-pending-api-key',
    'UTC',
    'pro-monthly',
    'pending',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- Suspended user: blocked everywhere except GET/DELETE /v1/users/me
INSERT INTO users (id, email, api_key, timezone, subscription_name, status, created_at, updated_at)
VALUES (
    'approval-suspended-user-id',
    'approval-suspended@httpsms.com',
    'approval-suspended-api-key',
    'UTC',
    'pro-monthly',
    'suspended',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- Pending user that deletes its own account (isolated from the user above)
INSERT INTO users (id, email, api_key, timezone, subscription_name, status, created_at, updated_at)
VALUES (
    'approval-pending-delete-user-id',
    'approval-pending-delete@httpsms.com',
    'approval-pending-delete-api-key',
    'UTC',
    'pro-monthly',
    'pending',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- Active user whose email is in ADMIN_EMAILS: its API key must never grant admin
INSERT INTO users (id, email, api_key, timezone, subscription_name, status, created_at, updated_at)
VALUES (
    'approval-admin-email-user-id',
    'admin@httpsms.test',
    'approval-admin-email-api-key',
    'UTC',
    'pro-monthly',
    'active',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- Phone API keys owned by the pending and suspended users
INSERT INTO phone_api_keys (id, name, user_id, user_email, phone_numbers, phone_ids, api_key, created_at, updated_at)
VALUES (
    'a0000000-0000-4000-8000-000000000001',
    'approval-pending-phone-key',
    'approval-pending-user-id',
    'approval-pending@httpsms.com',
    ARRAY['+18005559001'],
    ARRAY[]::TEXT[],
    'pk_approval-pending-phone-api-key',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO phone_api_keys (id, name, user_id, user_email, phone_numbers, phone_ids, api_key, created_at, updated_at)
VALUES (
    'a0000000-0000-4000-8000-000000000002',
    'approval-suspended-phone-key',
    'approval-suspended-user-id',
    'approval-suspended@httpsms.com',
    ARRAY['+18005559002'],
    ARRAY[]::TEXT[],
    'pk_approval-suspended-phone-api-key',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;
