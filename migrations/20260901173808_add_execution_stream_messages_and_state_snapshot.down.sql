-- reverse: modify "executions" table
ALTER TABLE "executions" DROP COLUMN "stream_messages", DROP COLUMN "state_snapshot";
