-- modify "executions" table
ALTER TABLE "executions" ADD COLUMN "state_snapshot" jsonb NULL, ADD COLUMN "stream_messages" boolean NOT NULL DEFAULT true;
