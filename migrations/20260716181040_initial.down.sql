-- reverse: create "user_to_sessions" table
DROP TABLE "user_to_sessions";
-- reverse: create index "idx_session_messages_user_id" to table: "session_messages"
DROP INDEX "idx_session_messages_user_id";
-- reverse: create index "idx_session_messages_session_id" to table: "session_messages"
DROP INDEX "idx_session_messages_session_id";
-- reverse: create "session_messages" table
DROP TABLE "session_messages";
-- reverse: create index "idx_indexed_knowledge_items_knowledge_item_id" to table: "indexed_knowledge_items"
DROP INDEX "idx_indexed_knowledge_items_knowledge_item_id";
-- reverse: create index "idx_indexed_knowledge_items_index_id" to table: "indexed_knowledge_items"
DROP INDEX "idx_indexed_knowledge_items_index_id";
-- reverse: create index "idx_indexed_knowledge_items_hash" to table: "indexed_knowledge_items"
DROP INDEX "idx_indexed_knowledge_items_hash";
-- reverse: create "indexed_knowledge_items" table
DROP TABLE "indexed_knowledge_items";
-- reverse: create index "idx_knowledge_items_knowledge_external_id" to table: "knowledge_items"
DROP INDEX "idx_knowledge_items_knowledge_external_id";
-- reverse: create index "idx_knowledge_items_hash" to table: "knowledge_items"
DROP INDEX "idx_knowledge_items_hash";
-- reverse: create "knowledge_items" table
DROP TABLE "knowledge_items";
-- reverse: create index "idx_execution_messages_execution_id" to table: "execution_messages"
DROP INDEX "idx_execution_messages_execution_id";
-- reverse: create "execution_messages" table
DROP TABLE "execution_messages";
-- reverse: create index "idx_executions_session_id" to table: "executions"
DROP INDEX "idx_executions_session_id";
-- reverse: create index "idx_executions_agent_id" to table: "executions"
DROP INDEX "idx_executions_agent_id";
-- reverse: create "executions" table
DROP TABLE "executions";
-- reverse: create index "idx_sessions_client_id" to table: "sessions"
DROP INDEX "idx_sessions_client_id";
-- reverse: create index "idx_sessions_agent_id" to table: "sessions"
DROP INDEX "idx_sessions_agent_id";
-- reverse: create "sessions" table
DROP TABLE "sessions";
-- reverse: create index "idx_client_to_users_external_id" to table: "client_to_users"
DROP INDEX "idx_client_to_users_external_id";
-- reverse: create "client_to_users" table
DROP TABLE "client_to_users";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: create "clients" table
DROP TABLE "clients";
-- reverse: create "agent_to_knowledges" table
DROP TABLE "agent_to_knowledges";
-- reverse: create "agent_to_knowledge_indices" table
DROP TABLE "agent_to_knowledge_indices";
-- reverse: create index "idx_knowledge_indices_knowledge_id" to table: "knowledge_indices"
DROP INDEX "idx_knowledge_indices_knowledge_id";
-- reverse: create "knowledge_indices" table
DROP TABLE "knowledge_indices";
-- reverse: create "knowledges" table
DROP TABLE "knowledges";
-- reverse: create "agents" table
DROP TABLE "agents";
-- reverse: create "api_keys" table
DROP TABLE "api_keys";
