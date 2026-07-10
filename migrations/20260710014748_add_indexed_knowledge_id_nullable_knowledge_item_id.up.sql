-- modify "indexed_knowledge_items" table
ALTER TABLE "indexed_knowledge_items" ALTER COLUMN "knowledge_item_id" DROP NOT NULL;
-- modify "knowledge_items" table
ALTER TABLE "knowledge_items" DROP COLUMN "kinds", ADD COLUMN "kind" character varying(255) NULL;
