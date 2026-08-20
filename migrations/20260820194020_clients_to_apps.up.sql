-- modify "agent_to_knowledge_indices" table
ALTER TABLE "agent_to_knowledge_indices" DROP CONSTRAINT "fk_agent_to_knowledge_indices_agent", DROP CONSTRAINT "fk_agent_to_knowledge_indices_index", ADD CONSTRAINT "fk_agent_to_knowledge_indices_agent" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_agent_to_knowledge_indices_index" FOREIGN KEY ("index_id") REFERENCES "knowledge_indices" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- create "app_users" table
CREATE TABLE "app_users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "picture" text NULL,
  "email" text NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create "apps" table
CREATE TABLE "apps" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "code" character varying(10) NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NOT NULL,
  "status" character varying(255) NOT NULL,
  "last_verified_at" timestamp NULL,
  "key_id" character varying(255) NOT NULL,
  "key_hash" text NOT NULL,
  "auth_config" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_apps_code" UNIQUE ("code"),
  CONSTRAINT "uni_apps_key_hash" UNIQUE ("key_hash"),
  CONSTRAINT "uni_apps_key_id" UNIQUE ("key_id"),
  CONSTRAINT "fk_tenants_clients" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_apps_tenant_id" to table: "apps"
CREATE INDEX "idx_apps_tenant_id" ON "apps" ("tenant_id");
-- create "app_to_app_users" table
CREATE TABLE "app_to_app_users" (
  "app_id" uuid NOT NULL,
  "app_user_id" uuid NOT NULL,
  "external_id" character varying(255) NULL,
  PRIMARY KEY ("app_id", "app_user_id"),
  CONSTRAINT "fk_app_users_app_to_app_users" FOREIGN KEY ("app_user_id") REFERENCES "app_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_apps_app_to_users" FOREIGN KEY ("app_id") REFERENCES "apps" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_app_to_app_users_external_id" to table: "app_to_app_users"
CREATE INDEX "idx_app_to_app_users_external_id" ON "app_to_app_users" ("external_id");
-- create "app_user_refresh_tokens" table
CREATE TABLE "app_user_refresh_tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "hash" text NOT NULL,
  "app_user_id" uuid NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "revoked_at" timestamp NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_app_user_refresh_tokens_app_user" FOREIGN KEY ("app_user_id") REFERENCES "app_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_app_user_refresh_tokens_app_user_id" to table: "app_user_refresh_tokens"
CREATE INDEX "idx_app_user_refresh_tokens_app_user_id" ON "app_user_refresh_tokens" ("app_user_id");
-- create index "idx_app_user_refresh_tokens_hash" to table: "app_user_refresh_tokens"
CREATE UNIQUE INDEX "idx_app_user_refresh_tokens_hash" ON "app_user_refresh_tokens" ("hash");
-- drop index "idx_sessions_client_id" from table: "sessions"
DROP INDEX "idx_sessions_client_id";
-- rename a column from "client_id" to "app_id"
ALTER TABLE "sessions" RENAME COLUMN "client_id" TO "app_id";
-- modify "sessions" table
ALTER TABLE "sessions" DROP CONSTRAINT "fk_clients_sessions", ADD CONSTRAINT "fk_apps_sessions" FOREIGN KEY ("app_id") REFERENCES "apps" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- create index "idx_sessions_app_id" to table: "sessions"
CREATE INDEX "idx_sessions_app_id" ON "sessions" ("app_id");
-- create "app_user_to_sessions" table
CREATE TABLE "app_user_to_sessions" (
  "app_user_id" uuid NOT NULL,
  "session_id" uuid NOT NULL,
  PRIMARY KEY ("app_user_id", "session_id"),
  CONSTRAINT "fk_app_users_app_user_to_sessions" FOREIGN KEY ("app_user_id") REFERENCES "app_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_app_user_to_sessions" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- modify "indexed_knowledge_items" table
ALTER TABLE "indexed_knowledge_items" DROP CONSTRAINT "fk_knowledge_indices_items", DROP CONSTRAINT "fk_knowledge_items_indexes", ADD CONSTRAINT "fk_knowledge_indices_items" FOREIGN KEY ("index_id") REFERENCES "knowledge_indices" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_knowledge_items_indexes" FOREIGN KEY ("knowledge_item_id") REFERENCES "knowledge_items" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- drop "client_to_client_users" table
DROP TABLE "client_to_client_users";
-- drop "client_user_refresh_tokens" table
DROP TABLE "client_user_refresh_tokens";
-- drop "client_user_to_sessions" table
DROP TABLE "client_user_to_sessions";
-- drop "client_users" table
DROP TABLE "client_users";
-- drop "clients" table
DROP TABLE "clients";
