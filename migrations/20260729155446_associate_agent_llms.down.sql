-- reverse: create index "idx_agent_to_llms_agent_priority" to table: "agent_to_llms"
DROP INDEX "idx_agent_to_llms_agent_priority";
-- reverse: create "agent_to_llms" table
DROP TABLE "agent_to_llms";
-- reverse: modify "agents" table
ALTER TABLE "agents" ADD COLUMN "configuration" jsonb NOT NULL;
