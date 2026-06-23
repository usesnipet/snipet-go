-- reverse: create index "idx_conversation_messages_conversation_id" to table: "conversation_messages"
DROP INDEX "idx_conversation_messages_conversation_id";
-- reverse: create index "idx_conversation_messages_client_user_id" to table: "conversation_messages"
DROP INDEX "idx_conversation_messages_client_user_id";
-- reverse: create "conversation_messages" table
DROP TABLE "conversation_messages";
-- reverse: create index "idx_client_user_conversations_conversation_id" to table: "client_user_conversations"
DROP INDEX "idx_client_user_conversations_conversation_id";
-- reverse: create index "idx_client_user_conversations_client_user_id" to table: "client_user_conversations"
DROP INDEX "idx_client_user_conversations_client_user_id";
-- reverse: create "client_user_conversations" table
DROP TABLE "client_user_conversations";
-- reverse: create index "idx_conversations_memory_id" to table: "conversations"
DROP INDEX "idx_conversations_memory_id";
-- reverse: create index "idx_conversations_bot_id" to table: "conversations"
DROP INDEX "idx_conversations_bot_id";
-- reverse: create "conversations" table
DROP TABLE "conversations";
-- reverse: create "client_users" table
DROP TABLE "client_users";
-- reverse: create index "idx_client_bots_client_id" to table: "client_bots"
DROP INDEX "idx_client_bots_client_id";
-- reverse: create index "idx_client_bots_bot_id" to table: "client_bots"
DROP INDEX "idx_client_bots_bot_id";
-- reverse: create "client_bots" table
DROP TABLE "client_bots";
-- reverse: create "clients" table
DROP TABLE "clients";
-- reverse: create index "idx_bot_memories_memory_id" to table: "bot_memories"
DROP INDEX "idx_bot_memories_memory_id";
-- reverse: create index "idx_bot_memories_bot_id" to table: "bot_memories"
DROP INDEX "idx_bot_memories_bot_id";
-- reverse: create "bot_memories" table
DROP TABLE "bot_memories";
-- reverse: create "memories" table
DROP TABLE "memories";
-- reverse: create "bots" table
DROP TABLE "bots";
