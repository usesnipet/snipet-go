-- modify "agents" table
ALTER TABLE "agents" DROP COLUMN "tenant_id";
-- modify "api_keys" table
ALTER TABLE "api_keys" DROP COLUMN "tenant_id";
-- modify "apps" table
ALTER TABLE "apps" DROP COLUMN "tenant_id";
-- modify "execution_messages" table
ALTER TABLE "execution_messages" DROP COLUMN "tenant_id";
-- modify "executions" table
ALTER TABLE "executions" DROP COLUMN "tenant_id";
-- modify "indexed_knowledge_items" table
ALTER TABLE "indexed_knowledge_items" DROP COLUMN "tenant_id";
-- modify "knowledge_indices" table
ALTER TABLE "knowledge_indices" DROP COLUMN "tenant_id";
-- modify "knowledge_items" table
ALTER TABLE "knowledge_items" DROP COLUMN "tenant_id";
-- modify "knowledges" table
ALTER TABLE "knowledges" DROP COLUMN "tenant_id";
-- modify "llms" table
ALTER TABLE "llms" DROP COLUMN "tenant_id";
-- modify "sessions" table
ALTER TABLE "sessions" DROP COLUMN "tenant_id";
-- drop "members" table
DROP TABLE "members";
-- drop "tenant_invitations" table
DROP TABLE "tenant_invitations";
-- drop "tenants" table
DROP TABLE "tenants";
-- drop "accounts" table
DROP TABLE "accounts";
-- drop "tokens" table
DROP TABLE "tokens";
-- drop "users" table
DROP TABLE "users";
