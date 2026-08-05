-- create "api_keys" table
CREATE TABLE "api_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "key_id" character varying(255) NOT NULL,
  "key" text NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "expires_at" timestamp NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_api_keys_key" UNIQUE ("key"),
  CONSTRAINT "uni_api_keys_key_id" UNIQUE ("key_id")
);
-- create "agents" table
CREATE TABLE "agents" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "description" text NOT NULL,
  "instructions" text NOT NULL,
  "configuration" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create "knowledges" table
CREATE TABLE "knowledges" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "driver" character varying(100) NOT NULL,
  "configuration" jsonb NOT NULL,
  "last_synced_at" timestamp NULL,
  "sync_status" character varying(20) NULL DEFAULT NULL::character varying,
  "sync_error" text NULL,
  PRIMARY KEY ("id")
);
-- create "knowledge_indices" table
CREATE TABLE "knowledge_indices" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "driver" character varying(100) NOT NULL,
  "configuration" jsonb NOT NULL,
  "knowledge_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_knowledges_indexes" FOREIGN KEY ("knowledge_id") REFERENCES "knowledges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_knowledge_indices_knowledge_id" to table: "knowledge_indices"
CREATE INDEX "idx_knowledge_indices_knowledge_id" ON "knowledge_indices" ("knowledge_id");
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
-- create "clients" table
CREATE TABLE "clients" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "code" character(10) NOT NULL,
  "name" character varying(255) NOT NULL,
  "config" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_clients_code" UNIQUE ("code")
);
-- create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "picture" text NULL,
  "email" text NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create "client_to_users" table
CREATE TABLE "client_to_users" (
  "client_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "external_id" character varying(255) NULL,
  PRIMARY KEY ("client_id", "user_id"),
  CONSTRAINT "fk_clients_client_to_users" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_client_to_users" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_client_to_users_external_id" to table: "client_to_users"
CREATE INDEX "idx_client_to_users_external_id" ON "client_to_users" ("external_id");
-- create "sessions" table
CREATE TABLE "sessions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "client_id" uuid NOT NULL,
  "agent_id" uuid NOT NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_clients_sessions" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_agent" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_sessions_agent_id" to table: "sessions"
CREATE INDEX "idx_sessions_agent_id" ON "sessions" ("agent_id");
-- create index "idx_sessions_client_id" to table: "sessions"
CREATE INDEX "idx_sessions_client_id" ON "sessions" ("client_id");
-- create "executions" table
CREATE TABLE "executions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "session_id" uuid NULL,
  "agent_id" uuid NOT NULL,
  "status" character varying(50) NOT NULL,
  "turns" integer NOT NULL DEFAULT 0,
  "metadata" jsonb NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_executions_agent" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_executions" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_executions_agent_id" to table: "executions"
CREATE INDEX "idx_executions_agent_id" ON "executions" ("agent_id");
-- create index "idx_executions_session_id" to table: "executions"
CREATE INDEX "idx_executions_session_id" ON "executions" ("session_id");
-- create "execution_messages" table
CREATE TABLE "execution_messages" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "execution_id" uuid NOT NULL,
  "sequence" integer NOT NULL,
  "role" character varying(50) NOT NULL,
  "content" text NULL,
  "tool_calls" jsonb NULL,
  "tool_result" jsonb NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_executions_messages" FOREIGN KEY ("execution_id") REFERENCES "executions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_execution_messages_execution_id" to table: "execution_messages"
CREATE INDEX "idx_execution_messages_execution_id" ON "execution_messages" ("execution_id");
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
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_knowledges_items" FOREIGN KEY ("knowledge_id") REFERENCES "knowledges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_knowledge_items_hash" to table: "knowledge_items"
CREATE INDEX "idx_knowledge_items_hash" ON "knowledge_items" ("hash");
-- create index "idx_knowledge_items_knowledge_external_id" to table: "knowledge_items"
CREATE UNIQUE INDEX "idx_knowledge_items_knowledge_external_id" ON "knowledge_items" ("knowledge_id", "external_id");
-- create "indexed_knowledge_items" table
CREATE TABLE "indexed_knowledge_items" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "hash" character varying(128) NULL,
  "indexed_at" timestamptz NULL,
  "metadata" jsonb NULL,
  "status" character varying(20) NOT NULL DEFAULT 'pending',
  "reason" text NULL,
  "last_error" text NULL,
  "index_id" uuid NOT NULL,
  "knowledge_item_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_knowledge_indices_items" FOREIGN KEY ("index_id") REFERENCES "knowledge_indices" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_knowledge_items_indexes" FOREIGN KEY ("knowledge_item_id") REFERENCES "knowledge_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_indexed_knowledge_items_hash" to table: "indexed_knowledge_items"
CREATE INDEX "idx_indexed_knowledge_items_hash" ON "indexed_knowledge_items" ("hash");
-- create index "idx_indexed_knowledge_items_index_id" to table: "indexed_knowledge_items"
CREATE INDEX "idx_indexed_knowledge_items_index_id" ON "indexed_knowledge_items" ("index_id");
-- create index "idx_indexed_knowledge_items_knowledge_item_id" to table: "indexed_knowledge_items"
CREATE INDEX "idx_indexed_knowledge_items_knowledge_item_id" ON "indexed_knowledge_items" ("knowledge_item_id");
-- create "session_messages" table
CREATE TABLE "session_messages" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "session_id" uuid NOT NULL,
  "role" character varying(255) NOT NULL,
  "parts" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_sessions_session_messages" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_session_messages" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_session_messages_session_id" to table: "session_messages"
CREATE INDEX "idx_session_messages_session_id" ON "session_messages" ("session_id");
-- create index "idx_session_messages_user_id" to table: "session_messages"
CREATE INDEX "idx_session_messages_user_id" ON "session_messages" ("user_id");
-- create "user_to_sessions" table
CREATE TABLE "user_to_sessions" (
  "user_id" uuid NOT NULL,
  "session_id" uuid NOT NULL,
  PRIMARY KEY ("user_id", "session_id"),
  CONSTRAINT "fk_sessions_user_to_sessions" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_user_to_sessions" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
