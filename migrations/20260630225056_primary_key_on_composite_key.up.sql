-- modify "api_keys" table
ALTER TABLE "api_keys" ADD CONSTRAINT "uni_api_keys_key_id" UNIQUE ("key_id");
-- drop index "idx_bot_to_memories_bot_id" from table: "bot_to_memories"
DROP INDEX "idx_bot_to_memories_bot_id";
-- drop index "idx_bot_to_memories_memory_id" from table: "bot_to_memories"
DROP INDEX "idx_bot_to_memories_memory_id";
-- modify "bot_to_memories" table
ALTER TABLE "bot_to_memories" ADD PRIMARY KEY ("bot_id", "memory_id");
-- drop index "idx_client_bots_client_bot" from table: "client_to_bots"
DROP INDEX "idx_client_bots_client_bot";
-- modify "client_to_bots" table
ALTER TABLE "client_to_bots" ADD PRIMARY KEY ("client_id", "bot_id");
-- drop index "idx_client_to_users_client_user" from table: "client_to_users"
DROP INDEX "idx_client_to_users_client_user";
-- modify "client_to_users" table
ALTER TABLE "client_to_users" ADD PRIMARY KEY ("client_id", "user_id");
-- drop index "idx_user_to_sessions_session_id" from table: "user_to_sessions"
DROP INDEX "idx_user_to_sessions_session_id";
-- drop index "idx_user_to_sessions_user_id" from table: "user_to_sessions"
DROP INDEX "idx_user_to_sessions_user_id";
-- modify "user_to_sessions" table
ALTER TABLE "user_to_sessions" ADD PRIMARY KEY ("user_id", "session_id");
