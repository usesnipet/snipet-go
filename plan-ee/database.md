# Tenant
- ID (uuid)
- name (varchar(255))
- slug (unique varchar(255))
- icon (varchar(255))
- created_at (timestamp)
- updated_at (timestamp)

# User
- ID (uuid)
- name (varchar(255))
- email (varchar(255))
- password_hash (varchar(255) or null)
- picture (varchar(255) or null)
- is_admin (boolean)
- challenges (jsonb) - array of challenges (active_account, change_password, etc.)
- created_at (timestamp)
- updated_at (timestamp)

# Account
- ID (uuid)
- user_id (uuid - foreign key to User)
- provider (varchar(255)) (google, github, etc.)
- external_id (varchar(255))
- created_at (timestamp)
- updated_at (timestamp)

# Token
- ID (uuid)
- type (varchar(255)) (refresh, activate_account, reset_password, etc.)
- hash (varchar(255))
- user_id (uuid - foreign key to User)
- expires_at (timestamp)
- created_at (timestamp)
- revoked_at (timestamp or null)
- metadata (jsonb)

# TenantInvitation
- ID (uuid)
- tenant_id (uuid - foreign key to Tenant)
- email (varchar(255))
- token (varchar(255))
- role (varchar(255)) (admin, user, etc.)
- status (varchar(255)) (pending, accepted, declined)
- expires_at (timestamp)
- created_at (timestamp)
- updated_at (timestamp)

# Member
- ID (uuid)
- user_id (uuid - foreign key to User)
- tenant_id (uuid - foreign key to Tenant)
- role (varchar(255)) (admin, user, etc.)
- is_active (boolean)
- created_at (timestamp)
- updated_at (timestamp)
- unique constraint: (user_id, tenant_id)