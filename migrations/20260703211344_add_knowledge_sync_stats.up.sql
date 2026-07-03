-- drop index "idx_knowledge_items_external_id" from table: "knowledge_items"
DROP INDEX "idx_knowledge_items_external_id";
-- drop index "idx_knowledge_items_knowledge_id" from table: "knowledge_items"
DROP INDEX "idx_knowledge_items_knowledge_id";
-- create index "idx_knowledge_items_knowledge_external_id" to table: "knowledge_items"
CREATE UNIQUE INDEX "idx_knowledge_items_knowledge_external_id" ON "knowledge_items" ("knowledge_id", "external_id");
-- modify "knowledges" table
ALTER TABLE "knowledges" ADD COLUMN "last_synced_at" timestamp NULL, ADD COLUMN "sync_status" character varying(20) NULL DEFAULT NULL::character varying, ADD COLUMN "sync_error" text NULL;
