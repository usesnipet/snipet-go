-- reverse: drop "clients" table
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
CREATE INDEX "idx_clients_tenant_id" ON "clients" ("tenant_id");
-- reverse: drop "client_users" table
CREATE TABLE "client_users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "picture" text NULL,
  "email" text NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- reverse: drop "client_user_to_sessions" table
CREATE TABLE "client_user_to_sessions" (
  "client_user_id" uuid NOT NULL,
  "session_id" uuid NOT NULL,
  PRIMARY KEY ("client_user_id", "session_id"),
  CONSTRAINT "fk_client_users_client_user_to_sessions" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_client_user_to_sessions" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- reverse: drop "client_user_refresh_tokens" table
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
CREATE INDEX "idx_client_user_refresh_tokens_client_user_id" ON "client_user_refresh_tokens" ("client_user_id");
CREATE UNIQUE INDEX "idx_client_user_refresh_tokens_hash" ON "client_user_refresh_tokens" ("hash");
-- reverse: drop "client_to_client_users" table
CREATE TABLE "client_to_client_users" (
  "client_id" uuid NOT NULL,
  "client_user_id" uuid NOT NULL,
  "external_id" character varying(255) NULL,
  PRIMARY KEY ("client_id", "client_user_id"),
  CONSTRAINT "fk_client_users_client_to_client_users" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_clients_client_to_users" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_client_to_client_users_external_id" ON "client_to_client_users" ("external_id");
-- reverse: modify "indexed_knowledge_items" table
ALTER TABLE "indexed_knowledge_items" DROP CONSTRAINT "fk_knowledge_items_indexes", DROP CONSTRAINT "fk_knowledge_indices_items", ADD CONSTRAINT "fk_knowledge_items_indexes" FOREIGN KEY ("knowledge_item_id") REFERENCES "knowledge_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_knowledge_indices_items" FOREIGN KEY ("index_id") REFERENCES "knowledge_indices" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- reverse: create "app_user_to_sessions" table
DROP TABLE "app_user_to_sessions";
-- reverse: create index "idx_sessions_app_id" to table: "sessions"
DROP INDEX "idx_sessions_app_id";
-- reverse: modify "sessions" table
ALTER TABLE "sessions" DROP CONSTRAINT "fk_apps_sessions", ADD CONSTRAINT "fk_clients_sessions" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- reverse: rename a column from "client_id" to "app_id"
ALTER TABLE "sessions" RENAME COLUMN "app_id" TO "client_id";
-- reverse: drop index "idx_sessions_client_id" from table: "sessions"
CREATE INDEX "idx_sessions_client_id" ON "sessions" ("client_id");
-- reverse: create index "idx_app_user_refresh_tokens_hash" to table: "app_user_refresh_tokens"
DROP INDEX "idx_app_user_refresh_tokens_hash";
-- reverse: create index "idx_app_user_refresh_tokens_app_user_id" to table: "app_user_refresh_tokens"
DROP INDEX "idx_app_user_refresh_tokens_app_user_id";
-- reverse: create "app_user_refresh_tokens" table
DROP TABLE "app_user_refresh_tokens";
-- reverse: create index "idx_app_to_app_users_external_id" to table: "app_to_app_users"
DROP INDEX "idx_app_to_app_users_external_id";
-- reverse: create "app_to_app_users" table
DROP TABLE "app_to_app_users";
-- reverse: create index "idx_apps_tenant_id" to table: "apps"
DROP INDEX "idx_apps_tenant_id";
-- reverse: create "apps" table
DROP TABLE "apps";
-- reverse: create "app_users" table
DROP TABLE "app_users";
-- reverse: modify "agent_to_knowledge_indices" table
ALTER TABLE "agent_to_knowledge_indices" DROP CONSTRAINT "fk_agent_to_knowledge_indices_index", DROP CONSTRAINT "fk_agent_to_knowledge_indices_agent", ADD CONSTRAINT "fk_agent_to_knowledge_indices_index" FOREIGN KEY ("index_id") REFERENCES "knowledge_indices" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_agent_to_knowledge_indices_agent" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
