-- create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "email" character varying(255) NOT NULL,
  "password_hash" character varying(255) NULL,
  "picture" character varying(255) NULL,
  "is_admin" boolean NOT NULL DEFAULT false,
  "challenges" jsonb NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email")
);
-- create "accounts" table
CREATE TABLE "accounts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "provider" character varying(255) NOT NULL,
  "external_id" character varying(255) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_accounts" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_accounts_provider_external_id" to table: "accounts"
CREATE UNIQUE INDEX "idx_accounts_provider_external_id" ON "accounts" ("provider", "external_id");
-- create index "idx_accounts_user_id" to table: "accounts"
CREATE INDEX "idx_accounts_user_id" ON "accounts" ("user_id");
-- create "tenants" table
CREATE TABLE "tenants" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "icon" character varying(255) NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_tenants_slug" UNIQUE ("slug")
);
-- create "agents" table
CREATE TABLE "agents" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NOT NULL,
  "instructions" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tenants_agents" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_agents_tenant_id" to table: "agents"
CREATE INDEX "idx_agents_tenant_id" ON "agents" ("tenant_id");
-- create "knowledges" table
CREATE TABLE "knowledges" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "driver" character varying(100) NOT NULL,
  "configuration" jsonb NOT NULL,
  "last_synced_at" timestamp NULL,
  "sync_status" character varying(20) NULL DEFAULT NULL::character varying,
  "sync_error" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tenants_knowledges" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_knowledges_tenant_id" to table: "knowledges"
CREATE INDEX "idx_knowledges_tenant_id" ON "knowledges" ("tenant_id");
-- create "knowledge_indices" table
CREATE TABLE "knowledge_indices" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "driver" character varying(100) NOT NULL,
  "configuration" jsonb NOT NULL,
  "knowledge_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_knowledge_indices_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_knowledges_indexes" FOREIGN KEY ("knowledge_id") REFERENCES "knowledges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_knowledge_indices_knowledge_id" to table: "knowledge_indices"
CREATE INDEX "idx_knowledge_indices_knowledge_id" ON "knowledge_indices" ("knowledge_id");
-- create index "idx_knowledge_indices_tenant_id" to table: "knowledge_indices"
CREATE INDEX "idx_knowledge_indices_tenant_id" ON "knowledge_indices" ("tenant_id");
-- create "agent_to_knowledge_indices" table
CREATE TABLE "agent_to_knowledge_indices" (
  "agent_id" uuid NOT NULL,
  "index_id" uuid NOT NULL,
  PRIMARY KEY ("agent_id", "index_id"),
  CONSTRAINT "fk_agent_to_knowledge_indices_agent" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_agent_to_knowledge_indices_index" FOREIGN KEY ("index_id") REFERENCES "knowledge_indices" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "agent_to_knowledges" table
CREATE TABLE "agent_to_knowledges" (
  "agent_id" uuid NOT NULL,
  "knowledge_id" uuid NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("agent_id", "knowledge_id"),
  CONSTRAINT "fk_agent_to_knowledges_knowledge" FOREIGN KEY ("knowledge_id") REFERENCES "knowledges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_agents_agent_to_knowledge" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create "llms" table
CREATE TABLE "llms" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "provider" character varying(255) NOT NULL,
  "configuration" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tenants_ll_ms" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_llms_tenant_id" to table: "llms"
CREATE INDEX "idx_llms_tenant_id" ON "llms" ("tenant_id");
-- create "agent_to_llms" table
CREATE TABLE "agent_to_llms" (
  "agent_id" uuid NOT NULL,
  "llm_id" uuid NOT NULL,
  "priority" integer NOT NULL,
  PRIMARY KEY ("agent_id", "llm_id"),
  CONSTRAINT "fk_agent_to_llms_llm" FOREIGN KEY ("llm_id") REFERENCES "llms" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_agents_agent_to_ll_ms" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_agent_to_llms_agent_priority" to table: "agent_to_llms"
CREATE INDEX "idx_agent_to_llms_agent_priority" ON "agent_to_llms" ("agent_id", "priority");
-- create "api_keys" table
CREATE TABLE "api_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "key_id" character varying(255) NOT NULL,
  "key" text NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "expires_at" timestamp NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_api_keys_key" UNIQUE ("key"),
  CONSTRAINT "uni_api_keys_key_id" UNIQUE ("key_id"),
  CONSTRAINT "fk_tenants_api_keys" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_api_keys_tenant_id" to table: "api_keys"
CREATE INDEX "idx_api_keys_tenant_id" ON "api_keys" ("tenant_id");
-- create "client_users" table
CREATE TABLE "client_users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "picture" text NULL,
  "email" text NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create "clients" table
CREATE TABLE "clients" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "code" character varying(10) NOT NULL,
  "name" character varying(255) NOT NULL,
  "config" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_clients_code" UNIQUE ("code"),
  CONSTRAINT "fk_tenants_clients" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_clients_tenant_id" to table: "clients"
CREATE INDEX "idx_clients_tenant_id" ON "clients" ("tenant_id");
-- create "client_to_client_users" table
CREATE TABLE "client_to_client_users" (
  "client_id" uuid NOT NULL,
  "client_user_id" uuid NOT NULL,
  "external_id" character varying(255) NULL,
  PRIMARY KEY ("client_id", "client_user_id"),
  CONSTRAINT "fk_client_users_client_to_client_users" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_clients_client_to_users" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_client_to_client_users_external_id" to table: "client_to_client_users"
CREATE INDEX "idx_client_to_client_users_external_id" ON "client_to_client_users" ("external_id");
-- create "client_user_refresh_tokens" table
CREATE TABLE "client_user_refresh_tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "hash" text NOT NULL,
  "client_user_id" uuid NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "revoked_at" timestamp NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_client_user_refresh_tokens_client_user" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_client_user_refresh_tokens_client_user_id" to table: "client_user_refresh_tokens"
CREATE INDEX "idx_client_user_refresh_tokens_client_user_id" ON "client_user_refresh_tokens" ("client_user_id");
-- create index "idx_client_user_refresh_tokens_hash" to table: "client_user_refresh_tokens"
CREATE UNIQUE INDEX "idx_client_user_refresh_tokens_hash" ON "client_user_refresh_tokens" ("hash");
-- create "sessions" table
CREATE TABLE "sessions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "client_id" uuid NOT NULL,
  "agent_id" uuid NOT NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_clients_sessions" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_agent" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_sessions_agent_id" to table: "sessions"
CREATE INDEX "idx_sessions_agent_id" ON "sessions" ("agent_id");
-- create index "idx_sessions_client_id" to table: "sessions"
CREATE INDEX "idx_sessions_client_id" ON "sessions" ("client_id");
-- create index "idx_sessions_tenant_id" to table: "sessions"
CREATE INDEX "idx_sessions_tenant_id" ON "sessions" ("tenant_id");
-- create "client_user_to_sessions" table
CREATE TABLE "client_user_to_sessions" (
  "client_user_id" uuid NOT NULL,
  "session_id" uuid NOT NULL,
  PRIMARY KEY ("client_user_id", "session_id"),
  CONSTRAINT "fk_client_users_client_user_to_sessions" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_client_user_to_sessions" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create "executions" table
CREATE TABLE "executions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "session_id" uuid NULL,
  "agent_id" uuid NOT NULL,
  "status" character varying(50) NOT NULL,
  "error_message" text NOT NULL,
  "turns" integer NOT NULL DEFAULT 0,
  "metadata" jsonb NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_executions_agent" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_executions_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_executions" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_executions_agent_id" to table: "executions"
CREATE INDEX "idx_executions_agent_id" ON "executions" ("agent_id");
-- create index "idx_executions_session_id" to table: "executions"
CREATE INDEX "idx_executions_session_id" ON "executions" ("session_id");
-- create index "idx_executions_tenant_id" to table: "executions"
CREATE INDEX "idx_executions_tenant_id" ON "executions" ("tenant_id");
-- create "execution_messages" table
CREATE TABLE "execution_messages" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "role" character varying(50) NOT NULL,
  "sequence" integer NOT NULL,
  "content" text NULL,
  "timestamp" timestamp NOT NULL DEFAULT now(),
  "tool_calls" jsonb NULL,
  "tool_call_id" character varying(255) NULL,
  "tenant_id" uuid NOT NULL,
  "execution_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_execution_messages_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_executions_messages" FOREIGN KEY ("execution_id") REFERENCES "executions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_execution_messages_execution_id" to table: "execution_messages"
CREATE INDEX "idx_execution_messages_execution_id" ON "execution_messages" ("execution_id");
-- create index "idx_execution_messages_tenant_id" to table: "execution_messages"
CREATE INDEX "idx_execution_messages_tenant_id" ON "execution_messages" ("tenant_id");
-- create "knowledge_items" table
CREATE TABLE "knowledge_items" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "external_id" character varying(255) NULL,
  "name" text NULL,
  "hash" character varying(128) NULL,
  "metadata" jsonb NULL,
  "attributes" jsonb NULL,
  "kind" character varying(255) NULL,
  "last_modified" timestamptz NULL,
  "knowledge_id" uuid NOT NULL,
  "tenant_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_knowledge_items_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_knowledges_items" FOREIGN KEY ("knowledge_id") REFERENCES "knowledges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_knowledge_items_hash" to table: "knowledge_items"
CREATE INDEX "idx_knowledge_items_hash" ON "knowledge_items" ("hash");
-- create index "idx_knowledge_items_knowledge_external_id" to table: "knowledge_items"
CREATE UNIQUE INDEX "idx_knowledge_items_knowledge_external_id" ON "knowledge_items" ("knowledge_id", "external_id");
-- create index "idx_knowledge_items_tenant_id" to table: "knowledge_items"
CREATE INDEX "idx_knowledge_items_tenant_id" ON "knowledge_items" ("tenant_id");
-- create "indexed_knowledge_items" table
CREATE TABLE "indexed_knowledge_items" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "hash" character varying(128) NULL,
  "indexed_at" timestamptz NULL,
  "metadata" jsonb NULL,
  "status" character varying(20) NOT NULL DEFAULT 'pending',
  "reason" text NULL,
  "last_error" text NULL,
  "tenant_id" uuid NOT NULL,
  "index_id" uuid NOT NULL,
  "knowledge_item_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_indexed_knowledge_items_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_knowledge_indices_items" FOREIGN KEY ("index_id") REFERENCES "knowledge_indices" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_knowledge_items_indexes" FOREIGN KEY ("knowledge_item_id") REFERENCES "knowledge_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_indexed_knowledge_items_hash" to table: "indexed_knowledge_items"
CREATE INDEX "idx_indexed_knowledge_items_hash" ON "indexed_knowledge_items" ("hash");
-- create index "idx_indexed_knowledge_items_index_id" to table: "indexed_knowledge_items"
CREATE INDEX "idx_indexed_knowledge_items_index_id" ON "indexed_knowledge_items" ("index_id");
-- create index "idx_indexed_knowledge_items_knowledge_item_id" to table: "indexed_knowledge_items"
CREATE INDEX "idx_indexed_knowledge_items_knowledge_item_id" ON "indexed_knowledge_items" ("knowledge_item_id");
-- create index "idx_indexed_knowledge_items_tenant_id" to table: "indexed_knowledge_items"
CREATE INDEX "idx_indexed_knowledge_items_tenant_id" ON "indexed_knowledge_items" ("tenant_id");
-- create "members" table
CREATE TABLE "members" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "tenant_id" uuid NOT NULL,
  "role" character varying(255) NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tenants_members" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_members" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_members_tenant_id" to table: "members"
CREATE INDEX "idx_members_tenant_id" ON "members" ("tenant_id");
-- create index "idx_members_user_tenant" to table: "members"
CREATE UNIQUE INDEX "idx_members_user_tenant" ON "members" ("user_id", "tenant_id");
-- create "tenant_invitations" table
CREATE TABLE "tenant_invitations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "email" character varying(255) NOT NULL,
  "token" character varying(255) NOT NULL,
  "role" character varying(255) NOT NULL,
  "status" character varying(255) NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_tenant_invitations_token" UNIQUE ("token"),
  CONSTRAINT "fk_tenants_invitations" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_tenant_invitations_status" to table: "tenant_invitations"
CREATE INDEX "idx_tenant_invitations_status" ON "tenant_invitations" ("status");
-- create index "idx_tenant_invitations_tenant_id" to table: "tenant_invitations"
CREATE INDEX "idx_tenant_invitations_tenant_id" ON "tenant_invitations" ("tenant_id");
-- create "tokens" table
CREATE TABLE "tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "type" character varying(255) NOT NULL,
  "hash" character varying(255) NOT NULL,
  "user_id" uuid NOT NULL,
  "expires_at" timestamp NOT NULL,
  "revoked_at" timestamp NULL,
  "metadata" jsonb NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_tokens" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_tokens_hash" to table: "tokens"
CREATE UNIQUE INDEX "idx_tokens_hash" ON "tokens" ("hash");
-- create index "idx_tokens_type" to table: "tokens"
CREATE INDEX "idx_tokens_type" ON "tokens" ("type");
-- create index "idx_tokens_user_id" to table: "tokens"
CREATE INDEX "idx_tokens_user_id" ON "tokens" ("user_id");
