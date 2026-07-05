-- modify "indexed_knowledge_items" table
ALTER TABLE "indexed_knowledge_items" DROP COLUMN "version", ADD COLUMN "reason" text NULL;
-- modify "knowledge_items" table
ALTER TABLE "knowledge_items" ADD COLUMN "kinds" jsonb NULL;
-- modify "knowledges" table
ALTER TABLE "knowledges" DROP COLUMN "type";
