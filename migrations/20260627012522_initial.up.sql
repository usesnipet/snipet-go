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
  CONSTRAINT "uni_api_keys_key" UNIQUE ("key")
);
-- create "bots" table
CREATE TABLE "bots" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "description" text NOT NULL,
  "configuration" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create "memories" table
CREATE TABLE "memories" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "type" character varying(255) NOT NULL,
  "is_default" boolean NOT NULL DEFAULT false,
  "provider" character varying(255) NOT NULL,
  "configuration" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create "bot_to_memories" table
CREATE TABLE "bot_to_memories" (
  "bot_id" uuid NOT NULL,
  "memory_id" uuid NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  CONSTRAINT "fk_bots_bot_to_memories" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_memories_bot_memories" FOREIGN KEY ("memory_id") REFERENCES "memories" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_bot_to_memories_bot_id" to table: "bot_to_memories"
CREATE INDEX "idx_bot_to_memories_bot_id" ON "bot_to_memories" ("bot_id");
-- create index "idx_bot_to_memories_memory_id" to table: "bot_to_memories"
CREATE INDEX "idx_bot_to_memories_memory_id" ON "bot_to_memories" ("memory_id");
-- create "clients" table
CREATE TABLE "clients" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "code" character(10) NOT NULL,
  "name" character varying(255) NOT NULL,
  "config" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_clients_code" UNIQUE ("code")
);
-- create "client_to_bots" table
CREATE TABLE "client_to_bots" (
  "client_id" uuid NOT NULL,
  "bot_id" uuid NOT NULL,
  CONSTRAINT "fk_bots_client_to_bots" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_clients_client_to_bots" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_client_bots_client_bot" to table: "client_to_bots"
CREATE UNIQUE INDEX "idx_client_bots_client_bot" ON "client_to_bots" ("client_id", "bot_id");
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
  CONSTRAINT "fk_clients_client_to_users" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_client_to_users" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_client_to_users_client_user" to table: "client_to_users"
CREATE UNIQUE INDEX "idx_client_to_users_client_user" ON "client_to_users" ("client_id", "user_id");
-- create index "idx_client_to_users_external_id" to table: "client_to_users"
CREATE INDEX "idx_client_to_users_external_id" ON "client_to_users" ("external_id");
-- create "sessions" table
CREATE TABLE "sessions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "client_id" uuid NOT NULL,
  "memory_id" uuid NOT NULL,
  "bot_id" uuid NOT NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_clients_sessions" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_bot" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_memory" FOREIGN KEY ("memory_id") REFERENCES "memories" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_sessions_bot_id" to table: "sessions"
CREATE INDEX "idx_sessions_bot_id" ON "sessions" ("bot_id");
-- create index "idx_sessions_client_id" to table: "sessions"
CREATE INDEX "idx_sessions_client_id" ON "sessions" ("client_id");
-- create index "idx_sessions_memory_id" to table: "sessions"
CREATE INDEX "idx_sessions_memory_id" ON "sessions" ("memory_id");
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
  CONSTRAINT "fk_sessions_user_to_sessions" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_user_to_sessions" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_user_to_sessions_session_id" to table: "user_to_sessions"
CREATE INDEX "idx_user_to_sessions_session_id" ON "user_to_sessions" ("session_id");
-- create index "idx_user_to_sessions_user_id" to table: "user_to_sessions"
CREATE INDEX "idx_user_to_sessions_user_id" ON "user_to_sessions" ("user_id");
