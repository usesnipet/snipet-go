-- reverse: modify "execution_messages" table
ALTER TABLE "execution_messages" ADD COLUMN "tool_result" jsonb NULL, ADD COLUMN "tool_calls" jsonb NULL;
-- reverse: rename a column from "created_at" to "timestamp"
ALTER TABLE "execution_messages" RENAME COLUMN "timestamp" TO "created_at";
