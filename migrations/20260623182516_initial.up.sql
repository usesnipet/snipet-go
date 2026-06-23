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
-- create "bot_memories" table
CREATE TABLE "bot_memories" (
  "bot_id" uuid NOT NULL,
  "memory_id" uuid NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  CONSTRAINT "fk_bots_bot_memories" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_memories_bot_memories" FOREIGN KEY ("memory_id") REFERENCES "memories" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_bot_memories_bot_id" to table: "bot_memories"
CREATE INDEX "idx_bot_memories_bot_id" ON "bot_memories" ("bot_id");
-- create index "idx_bot_memories_memory_id" to table: "bot_memories"
CREATE INDEX "idx_bot_memories_memory_id" ON "bot_memories" ("memory_id");
-- create "clients" table
CREATE TABLE "clients" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- create "client_bots" table
CREATE TABLE "client_bots" (
  "client_id" uuid NOT NULL,
  "bot_id" uuid NOT NULL,
  CONSTRAINT "fk_bots_client_bots" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_clients_client_bots" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_client_bots_bot_id" to table: "client_bots"
CREATE INDEX "idx_client_bots_bot_id" ON "client_bots" ("bot_id");
-- create index "idx_client_bots_client_id" to table: "client_bots"
CREATE INDEX "idx_client_bots_client_id" ON "client_bots" ("client_id");
-- create "client_users" table
CREATE TABLE "client_users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "anonymous" boolean NOT NULL DEFAULT false,
  "session_id" character varying(255) NOT NULL,
  "external_id" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- create "conversations" table
CREATE TABLE "conversations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "memory_id" uuid NOT NULL,
  "bot_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_conversations_bot" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_conversations_memory" FOREIGN KEY ("memory_id") REFERENCES "memories" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_conversations_bot_id" to table: "conversations"
CREATE INDEX "idx_conversations_bot_id" ON "conversations" ("bot_id");
-- create index "idx_conversations_memory_id" to table: "conversations"
CREATE INDEX "idx_conversations_memory_id" ON "conversations" ("memory_id");
-- create "client_user_conversations" table
CREATE TABLE "client_user_conversations" (
  "client_user_id" uuid NOT NULL,
  "conversation_id" uuid NOT NULL,
  CONSTRAINT "fk_client_users_client_user_conversations" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_conversations_client_user_conversations" FOREIGN KEY ("conversation_id") REFERENCES "conversations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_client_user_conversations_client_user_id" to table: "client_user_conversations"
CREATE INDEX "idx_client_user_conversations_client_user_id" ON "client_user_conversations" ("client_user_id");
-- create index "idx_client_user_conversations_conversation_id" to table: "client_user_conversations"
CREATE INDEX "idx_client_user_conversations_conversation_id" ON "client_user_conversations" ("conversation_id");
-- create "conversation_messages" table
CREATE TABLE "conversation_messages" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "client_user_id" uuid NOT NULL,
  "conversation_id" uuid NOT NULL,
  "role" character varying(255) NOT NULL,
  "parts" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_client_users_conversation_messages" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_conversations_conversation_messages" FOREIGN KEY ("conversation_id") REFERENCES "conversations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_conversation_messages_client_user_id" to table: "conversation_messages"
CREATE INDEX "idx_conversation_messages_client_user_id" ON "conversation_messages" ("client_user_id");
-- create index "idx_conversation_messages_conversation_id" to table: "conversation_messages"
CREATE INDEX "idx_conversation_messages_conversation_id" ON "conversation_messages" ("conversation_id");
