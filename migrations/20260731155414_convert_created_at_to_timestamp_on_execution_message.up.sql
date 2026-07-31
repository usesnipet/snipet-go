-- rename a column from "created_at" to "timestamp"
ALTER TABLE "execution_messages" RENAME COLUMN "created_at" TO "timestamp";
-- modify "execution_messages" table
ALTER TABLE "execution_messages" DROP COLUMN "tool_calls", DROP COLUMN "tool_result";
