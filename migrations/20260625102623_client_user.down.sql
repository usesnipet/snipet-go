-- reverse: drop "client_users" table
CREATE TABLE "client_users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "anonymous" boolean NOT NULL DEFAULT false,
  "session_id" character varying(255) NOT NULL,
  "external_id" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- reverse: modify "conversation_messages" table
ALTER TABLE "conversation_messages" DROP CONSTRAINT "fk_c_users_conversation_messages", ADD CONSTRAINT "fk_client_users_conversation_messages" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- reverse: modify "client_user_conversations" table
ALTER TABLE "client_user_conversations" DROP CONSTRAINT "fk_c_users_client_user_conversations", ADD CONSTRAINT "fk_client_users_client_user_conversations" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- reverse: create index "idx_client_to_users_external_id" to table: "client_to_users"
DROP INDEX "idx_client_to_users_external_id";
-- reverse: create index "idx_client_to_users_client_user" to table: "client_to_users"
DROP INDEX "idx_client_to_users_client_user";
-- reverse: create "client_to_users" table
DROP TABLE "client_to_users";
-- reverse: create "c_users" table
DROP TABLE "c_users";
