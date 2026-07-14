-- reverse: rename sessions.agent_id to bot_id
ALTER TABLE "sessions" DROP CONSTRAINT "fk_sessions_agent";
ALTER TABLE "sessions" RENAME COLUMN "agent_id" TO "bot_id";
ALTER INDEX "idx_sessions_agent_id" RENAME TO "idx_sessions_bot_id";
-- reverse: rename "agent_to_knowledges" table to "bot_to_knowledges"
ALTER TABLE "agent_to_knowledges" RENAME CONSTRAINT "fk_agents_agent_to_knowledge" TO "fk_bots_bot_to_knowledge";
ALTER TABLE "agent_to_knowledges" RENAME CONSTRAINT "fk_agent_to_knowledges_knowledge" TO "fk_bot_to_knowledges_knowledge";
ALTER TABLE "agent_to_knowledges" RENAME COLUMN "agent_id" TO "bot_id";
ALTER TABLE "agent_to_knowledges" RENAME TO "bot_to_knowledges";
-- reverse: rename "agent_to_knowledge_indices" table to "bot_to_knowledge_indices"
ALTER TABLE "agent_to_knowledge_indices" RENAME CONSTRAINT "fk_agent_to_knowledge_indices_agent" TO "fk_bot_to_knowledge_indices_bot";
ALTER TABLE "agent_to_knowledge_indices" RENAME CONSTRAINT "fk_agent_to_knowledge_indices_index" TO "fk_bot_to_knowledge_indices_index";
ALTER TABLE "agent_to_knowledge_indices" RENAME COLUMN "agent_id" TO "bot_id";
ALTER TABLE "agent_to_knowledge_indices" RENAME TO "bot_to_knowledge_indices";
-- reverse: rename "agents" table to "bots"
ALTER TABLE "agents" RENAME TO "bots";
-- restore sessions FK to bots
ALTER TABLE "sessions" ADD CONSTRAINT "fk_sessions_bot" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- reverse: modify "knowledge_items" table
ALTER TABLE "knowledge_items" DROP COLUMN "attributes";
