-- create "c_users" table
CREATE TABLE "c_users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create "client_to_users" table
CREATE TABLE "client_to_users" (
  "client_id" uuid NOT NULL,
  "client_user_id" uuid NOT NULL,
  "external_id" character varying(255) NULL,
  CONSTRAINT "fk_client_to_users_client" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_client_to_users_client_user" FOREIGN KEY ("client_user_id") REFERENCES "c_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_client_to_users_client_user" to table: "client_to_users"
CREATE UNIQUE INDEX "idx_client_to_users_client_user" ON "client_to_users" ("client_id", "client_user_id");
-- create index "idx_client_to_users_external_id" to table: "client_to_users"
CREATE INDEX "idx_client_to_users_external_id" ON "client_to_users" ("external_id");
-- modify "client_user_conversations" table
ALTER TABLE "client_user_conversations" DROP CONSTRAINT "fk_client_users_client_user_conversations", ADD CONSTRAINT "fk_c_users_client_user_conversations" FOREIGN KEY ("client_user_id") REFERENCES "c_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- modify "conversation_messages" table
ALTER TABLE "conversation_messages" DROP CONSTRAINT "fk_client_users_conversation_messages", ADD CONSTRAINT "fk_c_users_conversation_messages" FOREIGN KEY ("client_user_id") REFERENCES "c_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- drop "client_users" table
DROP TABLE "client_users";
