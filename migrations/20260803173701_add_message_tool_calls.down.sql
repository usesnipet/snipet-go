-- reverse: modify "execution_messages" table
ALTER TABLE "execution_messages" DROP COLUMN "tool_call_id", DROP COLUMN "tool_calls";
