-- reverse: rename a column from "provider" to "driver"
ALTER TABLE "knowledges" RENAME COLUMN "driver" TO "provider";
-- reverse: modify "knowledge_indices" table
ALTER TABLE "knowledge_indices" ADD COLUMN "status" character varying(50) NULL DEFAULT 'ready';
