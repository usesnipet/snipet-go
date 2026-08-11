-- reverse: drop "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "picture" text NULL,
  "email" text NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- reverse: drop "user_to_sessions" table
CREATE TABLE "user_to_sessions" (
  "user_id" uuid NOT NULL,
  "session_id" uuid NOT NULL,
  PRIMARY KEY ("user_id", "session_id"),
  CONSTRAINT "fk_sessions_user_to_sessions" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_user_to_sessions" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- reverse: drop "client_to_users" table
CREATE TABLE "client_to_users" (
  "client_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "external_id" character varying(255) NULL,
  PRIMARY KEY ("client_id", "user_id"),
  CONSTRAINT "fk_clients_client_to_users" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_client_to_users" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_client_to_users_external_id" ON "client_to_users" ("external_id");
-- reverse: modify "refresh_tokens" table
ALTER TABLE "refresh_tokens" DROP CONSTRAINT "fk_refresh_tokens_user", ADD CONSTRAINT "fk_refresh_tokens_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- reverse: create "client_user_to_sessions" table
DROP TABLE "client_user_to_sessions";
-- reverse: create index "idx_client_to_client_users_external_id" to table: "client_to_client_users"
DROP INDEX "idx_client_to_client_users_external_id";
-- reverse: create "client_to_client_users" table
DROP TABLE "client_to_client_users";
-- reverse: create "client_users" table
DROP TABLE "client_users";
