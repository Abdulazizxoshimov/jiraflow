ALTER TABLE oauth_accounts
    DROP COLUMN IF EXISTS refresh_token,
    DROP COLUMN IF EXISTS token_expiry;
