-- reverse: modify "user_to_sessions" table
ALTER TABLE "user_to_sessions" DROP CONSTRAINT "user_to_sessions_pkey";
-- reverse: drop index "idx_user_to_sessions_user_id" from table: "user_to_sessions"
CREATE INDEX "idx_user_to_sessions_user_id" ON "user_to_sessions" ("user_id");
-- reverse: drop index "idx_user_to_sessions_session_id" from table: "user_to_sessions"
CREATE INDEX "idx_user_to_sessions_session_id" ON "user_to_sessions" ("session_id");
-- reverse: modify "client_to_users" table
ALTER TABLE "client_to_users" DROP CONSTRAINT "client_to_users_pkey";
-- reverse: drop index "idx_client_to_users_client_user" from table: "client_to_users"
CREATE UNIQUE INDEX "idx_client_to_users_client_user" ON "client_to_users" ("client_id", "user_id");
-- reverse: modify "client_to_bots" table
ALTER TABLE "client_to_bots" DROP CONSTRAINT "client_to_bots_pkey";
-- reverse: drop index "idx_client_bots_client_bot" from table: "client_to_bots"
CREATE UNIQUE INDEX "idx_client_bots_client_bot" ON "client_to_bots" ("client_id", "bot_id");
-- reverse: modify "bot_to_memories" table
ALTER TABLE "bot_to_memories" DROP CONSTRAINT "bot_to_memories_pkey";
-- reverse: drop index "idx_bot_to_memories_memory_id" from table: "bot_to_memories"
CREATE INDEX "idx_bot_to_memories_memory_id" ON "bot_to_memories" ("memory_id");
-- reverse: drop index "idx_bot_to_memories_bot_id" from table: "bot_to_memories"
CREATE INDEX "idx_bot_to_memories_bot_id" ON "bot_to_memories" ("bot_id");
-- reverse: modify "api_keys" table
ALTER TABLE "api_keys" DROP CONSTRAINT "uni_api_keys_key_id";
