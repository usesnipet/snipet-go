-- modify "execution_messages" table
ALTER TABLE "execution_messages" ADD COLUMN "tool_calls" jsonb NULL, ADD COLUMN "tool_call_id" character varying(255) NULL;
