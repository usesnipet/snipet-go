-- create "knowledges" table
CREATE TABLE "knowledges" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "type" character varying(100) NOT NULL,
  "provider" character varying(100) NOT NULL,
  "configuration" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_knowledges_type" to table: "knowledges"
CREATE INDEX "idx_knowledges_type" ON "knowledges" ("type");
-- create "knowledge_indices" table
CREATE TABLE "knowledge_indices" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "driver" character varying(100) NOT NULL,
  "configuration" jsonb NOT NULL,
  "status" character varying(50) NULL DEFAULT 'ready',
  "knowledge_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_knowledges_indexes" FOREIGN KEY ("knowledge_id") REFERENCES "knowledges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_knowledge_indices_knowledge_id" to table: "knowledge_indices"
CREATE INDEX "idx_knowledge_indices_knowledge_id" ON "knowledge_indices" ("knowledge_id");
-- create "bot_to_knowledge_indices" table
CREATE TABLE "bot_to_knowledge_indices" (
  "bot_id" uuid NOT NULL,
  "index_id" uuid NOT NULL,
  PRIMARY KEY ("bot_id", "index_id"),
  CONSTRAINT "fk_bot_to_knowledge_indices_bot" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_bot_to_knowledge_indices_index" FOREIGN KEY ("index_id") REFERENCES "knowledge_indices" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "bot_to_knowledges" table
CREATE TABLE "bot_to_knowledges" (
  "bot_id" uuid NOT NULL,
  "knowledge_id" uuid NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("bot_id", "knowledge_id"),
  CONSTRAINT "fk_bot_to_knowledges_knowledge" FOREIGN KEY ("knowledge_id") REFERENCES "knowledges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_bots_bot_to_knowledge" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create "knowledge_items" table
CREATE TABLE "knowledge_items" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "external_id" character varying(255) NULL,
  "name" text NULL,
  "hash" character varying(128) NULL,
  "metadata" jsonb NULL,
  "last_modified" timestamptz NULL,
  "knowledge_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_knowledges_items" FOREIGN KEY ("knowledge_id") REFERENCES "knowledges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_knowledge_items_external_id" to table: "knowledge_items"
CREATE INDEX "idx_knowledge_items_external_id" ON "knowledge_items" ("external_id");
-- create index "idx_knowledge_items_hash" to table: "knowledge_items"
CREATE INDEX "idx_knowledge_items_hash" ON "knowledge_items" ("hash");
-- create index "idx_knowledge_items_knowledge_id" to table: "knowledge_items"
CREATE INDEX "idx_knowledge_items_knowledge_id" ON "knowledge_items" ("knowledge_id");
-- create "indexed_knowledge_items" table
CREATE TABLE "indexed_knowledge_items" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "status" character varying(20) NOT NULL DEFAULT 'pending',
  "version" bigint NULL DEFAULT 1,
  "hash" character varying(128) NULL,
  "indexed_at" timestamptz NULL,
  "last_error" text NULL,
  "metadata" jsonb NULL,
  "index_id" uuid NOT NULL,
  "knowledge_item_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_knowledge_indices_items" FOREIGN KEY ("index_id") REFERENCES "knowledge_indices" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_knowledge_items_indexes" FOREIGN KEY ("knowledge_item_id") REFERENCES "knowledge_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_indexed_knowledge_items_hash" to table: "indexed_knowledge_items"
CREATE INDEX "idx_indexed_knowledge_items_hash" ON "indexed_knowledge_items" ("hash");
-- create index "idx_indexed_knowledge_items_index_id" to table: "indexed_knowledge_items"
CREATE INDEX "idx_indexed_knowledge_items_index_id" ON "indexed_knowledge_items" ("index_id");
-- create index "idx_indexed_knowledge_items_knowledge_item_id" to table: "indexed_knowledge_items"
CREATE INDEX "idx_indexed_knowledge_items_knowledge_item_id" ON "indexed_knowledge_items" ("knowledge_item_id");
-- modify "sessions" table
ALTER TABLE "sessions" DROP COLUMN "memory_id";
-- drop "bot_to_memories" table
DROP TABLE "bot_to_memories";
-- drop "memories" table
DROP TABLE "memories";
