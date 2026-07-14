-- modify "knowledge_items" table
ALTER TABLE "knowledge_items" ADD COLUMN "attributes" jsonb NULL;
-- rename "bots" table to "agents"
ALTER TABLE "bots" RENAME TO "agents";
-- rename "bot_to_knowledge_indices" table to "agent_to_knowledge_indices"
ALTER TABLE "bot_to_knowledge_indices" RENAME TO "agent_to_knowledge_indices";
ALTER TABLE "agent_to_knowledge_indices" RENAME COLUMN "bot_id" TO "agent_id";
ALTER TABLE "agent_to_knowledge_indices" RENAME CONSTRAINT "fk_bot_to_knowledge_indices_bot" TO "fk_agent_to_knowledge_indices_agent";
ALTER TABLE "agent_to_knowledge_indices" RENAME CONSTRAINT "fk_bot_to_knowledge_indices_index" TO "fk_agent_to_knowledge_indices_index";
-- rename "bot_to_knowledges" table to "agent_to_knowledges"
ALTER TABLE "bot_to_knowledges" RENAME TO "agent_to_knowledges";
ALTER TABLE "agent_to_knowledges" RENAME COLUMN "bot_id" TO "agent_id";
ALTER TABLE "agent_to_knowledges" RENAME CONSTRAINT "fk_bots_bot_to_knowledge" TO "fk_agents_agent_to_knowledge";
ALTER TABLE "agent_to_knowledges" RENAME CONSTRAINT "fk_bot_to_knowledges_knowledge" TO "fk_agent_to_knowledges_knowledge";
-- rename sessions.bot_id to agent_id
ALTER TABLE "sessions" DROP CONSTRAINT "fk_sessions_bot";
ALTER TABLE "sessions" RENAME COLUMN "bot_id" TO "agent_id";
ALTER TABLE "sessions" ADD CONSTRAINT "fk_sessions_agent" FOREIGN KEY ("agent_id") REFERENCES "agents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER INDEX "idx_sessions_bot_id" RENAME TO "idx_sessions_agent_id";
