-- modify "knowledge_indices" table
ALTER TABLE "knowledge_indices" DROP COLUMN "status";
-- rename a column from "provider" to "driver"
ALTER TABLE "knowledges" RENAME COLUMN "provider" TO "driver";
