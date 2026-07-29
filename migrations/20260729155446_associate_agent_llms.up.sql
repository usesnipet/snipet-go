-- modify "agents" table
ALTER TABLE "agents" DROP COLUMN "configuration";
-- create "agent_to_llms" table
CREATE TABLE "agent_to_llms" (
  "agent_id" uuid NOT NULL,
  "llm_id" uuid NOT NULL,
  "priority" integer NOT NULL,
  PRIMARY KEY ("agent_id", "llm_id"),
  CONSTRAINT "fk_agent_to_llms_llm" FOREIGN KEY ("llm_id") REFERENCES "llms" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_agents_agent_to_ll_ms" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_agent_to_llms_agent_priority" to table: "agent_to_llms"
CREATE INDEX "idx_agent_to_llms_agent_priority" ON "agent_to_llms" ("agent_id", "priority");
