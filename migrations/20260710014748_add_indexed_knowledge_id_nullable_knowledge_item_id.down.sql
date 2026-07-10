-- reverse: modify "knowledge_items" table
ALTER TABLE "knowledge_items" DROP COLUMN "kind", ADD COLUMN "kinds" jsonb NULL;
-- reverse: modify "indexed_knowledge_items" table
ALTER TABLE "indexed_knowledge_items" ALTER COLUMN "knowledge_item_id" SET NOT NULL;
