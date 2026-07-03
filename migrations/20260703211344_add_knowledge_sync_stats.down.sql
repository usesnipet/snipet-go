-- reverse: modify "knowledges" table
ALTER TABLE "knowledges" DROP COLUMN "sync_error", DROP COLUMN "sync_status", DROP COLUMN "last_synced_at";
-- reverse: create index "idx_knowledge_items_knowledge_external_id" to table: "knowledge_items"
DROP INDEX "idx_knowledge_items_knowledge_external_id";
-- reverse: drop index "idx_knowledge_items_knowledge_id" from table: "knowledge_items"
CREATE INDEX "idx_knowledge_items_knowledge_id" ON "knowledge_items" ("knowledge_id");
-- reverse: drop index "idx_knowledge_items_external_id" from table: "knowledge_items"
CREATE INDEX "idx_knowledge_items_external_id" ON "knowledge_items" ("external_id");
