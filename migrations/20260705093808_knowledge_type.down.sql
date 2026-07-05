-- reverse: modify "knowledges" table
ALTER TABLE "knowledges" ADD COLUMN "type" character varying(100) NOT NULL;
-- reverse: modify "knowledge_items" table
ALTER TABLE "knowledge_items" DROP COLUMN "kinds";
-- reverse: modify "indexed_knowledge_items" table
ALTER TABLE "indexed_knowledge_items" DROP COLUMN "reason", ADD COLUMN "version" bigint NULL DEFAULT 1;
