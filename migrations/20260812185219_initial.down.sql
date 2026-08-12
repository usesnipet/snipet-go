-- reverse: create index "idx_tokens_user_id" to table: "tokens"
DROP INDEX "idx_tokens_user_id";
-- reverse: create index "idx_tokens_type" to table: "tokens"
DROP INDEX "idx_tokens_type";
-- reverse: create index "idx_tokens_hash" to table: "tokens"
DROP INDEX "idx_tokens_hash";
-- reverse: create "tokens" table
DROP TABLE "tokens";
-- reverse: create index "idx_tenant_invitations_tenant_id" to table: "tenant_invitations"
DROP INDEX "idx_tenant_invitations_tenant_id";
-- reverse: create index "idx_tenant_invitations_status" to table: "tenant_invitations"
DROP INDEX "idx_tenant_invitations_status";
-- reverse: create "tenant_invitations" table
DROP TABLE "tenant_invitations";
-- reverse: create index "idx_members_user_tenant" to table: "members"
DROP INDEX "idx_members_user_tenant";
-- reverse: create index "idx_members_tenant_id" to table: "members"
DROP INDEX "idx_members_tenant_id";
-- reverse: create "members" table
DROP TABLE "members";
-- reverse: create "tenants" table
DROP TABLE "tenants";
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
-- reverse: create "client_user_to_sessions" table
DROP TABLE "client_user_to_sessions";
-- reverse: create index "idx_sessions_client_id" to table: "sessions"
DROP INDEX "idx_sessions_client_id";
-- reverse: create index "idx_sessions_agent_id" to table: "sessions"
DROP INDEX "idx_sessions_agent_id";
-- reverse: create "sessions" table
DROP TABLE "sessions";
-- reverse: create index "idx_client_user_refresh_tokens_hash" to table: "client_user_refresh_tokens"
DROP INDEX "idx_client_user_refresh_tokens_hash";
-- reverse: create index "idx_client_user_refresh_tokens_client_user_id" to table: "client_user_refresh_tokens"
DROP INDEX "idx_client_user_refresh_tokens_client_user_id";
-- reverse: create "client_user_refresh_tokens" table
DROP TABLE "client_user_refresh_tokens";
-- reverse: create index "idx_client_to_client_users_external_id" to table: "client_to_client_users"
DROP INDEX "idx_client_to_client_users_external_id";
-- reverse: create "client_to_client_users" table
DROP TABLE "client_to_client_users";
-- reverse: create "clients" table
DROP TABLE "clients";
-- reverse: create "client_users" table
DROP TABLE "client_users";
-- reverse: create index "idx_agent_to_llms_agent_priority" to table: "agent_to_llms"
DROP INDEX "idx_agent_to_llms_agent_priority";
-- reverse: create "agent_to_llms" table
DROP TABLE "agent_to_llms";
-- reverse: create "llms" table
DROP TABLE "llms";
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
-- reverse: create index "idx_accounts_user_id" to table: "accounts"
DROP INDEX "idx_accounts_user_id";
-- reverse: create index "idx_accounts_provider_external_id" to table: "accounts"
DROP INDEX "idx_accounts_provider_external_id";
-- reverse: create "accounts" table
DROP TABLE "accounts";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: create "api_keys" table
DROP TABLE "api_keys";
