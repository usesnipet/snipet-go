-- reverse: create index "idx_conversations_client_id" to table: "conversations"
DROP INDEX "idx_conversations_client_id";
-- reverse: modify "conversations" table
ALTER TABLE "conversations" DROP CONSTRAINT "fk_clients_conversations", DROP COLUMN "metadata", DROP COLUMN "client_id";
