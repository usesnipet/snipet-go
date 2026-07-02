-- reverse: drop "memories" table
CREATE TABLE "memories" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "type" character varying(255) NOT NULL,
  "is_default" boolean NOT NULL DEFAULT false,
  "provider" character varying(255) NOT NULL,
  "configuration" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- reverse: drop "bot_to_memories" table
CREATE TABLE "bot_to_memories" (
  "bot_id" uuid NOT NULL,
  "memory_id" uuid NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("bot_id", "memory_id"),
  CONSTRAINT "fk_bots_bot_to_memories" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_memories_bot_memories" FOREIGN KEY ("memory_id") REFERENCES "memories" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- reverse: modify "sessions" table
ALTER TABLE "sessions" ADD COLUMN "memory_id" uuid NOT NULL;
-- reverse: create index "idx_indexed_knowledge_items_knowledge_item_id" to table: "indexed_knowledge_items"
DROP INDEX "idx_indexed_knowledge_items_knowledge_item_id";
-- reverse: create index "idx_indexed_knowledge_items_index_id" to table: "indexed_knowledge_items"
DROP INDEX "idx_indexed_knowledge_items_index_id";
-- reverse: create index "idx_indexed_knowledge_items_hash" to table: "indexed_knowledge_items"
DROP INDEX "idx_indexed_knowledge_items_hash";
-- reverse: create "indexed_knowledge_items" table
DROP TABLE "indexed_knowledge_items";
-- reverse: create index "idx_knowledge_items_knowledge_id" to table: "knowledge_items"
DROP INDEX "idx_knowledge_items_knowledge_id";
-- reverse: create index "idx_knowledge_items_hash" to table: "knowledge_items"
DROP INDEX "idx_knowledge_items_hash";
-- reverse: create index "idx_knowledge_items_external_id" to table: "knowledge_items"
DROP INDEX "idx_knowledge_items_external_id";
-- reverse: create "knowledge_items" table
DROP TABLE "knowledge_items";
-- reverse: create "bot_to_knowledges" table
DROP TABLE "bot_to_knowledges";
-- reverse: create "bot_to_knowledge_indices" table
DROP TABLE "bot_to_knowledge_indices";
-- reverse: create index "idx_knowledge_indices_knowledge_id" to table: "knowledge_indices"
DROP INDEX "idx_knowledge_indices_knowledge_id";
-- reverse: create "knowledge_indices" table
DROP TABLE "knowledge_indices";
-- reverse: create index "idx_knowledges_type" to table: "knowledges"
DROP INDEX "idx_knowledges_type";
-- reverse: create "knowledges" table
DROP TABLE "knowledges";
