ALTER TABLE projects DROP COLUMN IF EXISTS permission_scheme_id;
DROP TABLE IF EXISTS data_imports;
DROP TABLE IF EXISTS permission_scheme_grants;
DROP TABLE IF EXISTS permission_schemes;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS oauth_accounts;
DROP TABLE IF EXISTS oauth_states;
