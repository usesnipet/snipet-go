-- reverse: create index "idx_user_to_sessions_user_id" to table: "user_to_sessions"
DROP INDEX "idx_user_to_sessions_user_id";
-- reverse: create index "idx_user_to_sessions_session_id" to table: "user_to_sessions"
DROP INDEX "idx_user_to_sessions_session_id";
-- reverse: create "user_to_sessions" table
DROP TABLE "user_to_sessions";
-- reverse: create index "idx_session_messages_user_id" to table: "session_messages"
DROP INDEX "idx_session_messages_user_id";
-- reverse: create index "idx_session_messages_session_id" to table: "session_messages"
DROP INDEX "idx_session_messages_session_id";
-- reverse: create "session_messages" table
DROP TABLE "session_messages";
-- reverse: create index "idx_sessions_memory_id" to table: "sessions"
DROP INDEX "idx_sessions_memory_id";
-- reverse: create index "idx_sessions_client_id" to table: "sessions"
DROP INDEX "idx_sessions_client_id";
-- reverse: create index "idx_sessions_bot_id" to table: "sessions"
DROP INDEX "idx_sessions_bot_id";
-- reverse: create "sessions" table
DROP TABLE "sessions";
-- reverse: create index "idx_client_to_users_external_id" to table: "client_to_users"
DROP INDEX "idx_client_to_users_external_id";
-- reverse: create index "idx_client_to_users_client_user" to table: "client_to_users"
DROP INDEX "idx_client_to_users_client_user";
-- reverse: create "client_to_users" table
DROP TABLE "client_to_users";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: create index "idx_client_bots_client_bot" to table: "client_to_bots"
DROP INDEX "idx_client_bots_client_bot";
-- reverse: create "client_to_bots" table
DROP TABLE "client_to_bots";
-- reverse: create "clients" table
DROP TABLE "clients";
-- reverse: create index "idx_bot_to_memories_memory_id" to table: "bot_to_memories"
DROP INDEX "idx_bot_to_memories_memory_id";
-- reverse: create index "idx_bot_to_memories_bot_id" to table: "bot_to_memories"
DROP INDEX "idx_bot_to_memories_bot_id";
-- reverse: create "bot_to_memories" table
DROP TABLE "bot_to_memories";
-- reverse: create "memories" table
DROP TABLE "memories";
-- reverse: create "bots" table
DROP TABLE "bots";
-- reverse: create "api_keys" table
DROP TABLE "api_keys";
