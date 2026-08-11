-- create "client_users" table
CREATE TABLE "client_users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "picture" text NULL,
  "email" text NULL,
  "metadata" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create "client_to_client_users" table
CREATE TABLE "client_to_client_users" (
  "client_id" uuid NOT NULL,
  "client_user_id" uuid NOT NULL,
  "external_id" character varying(255) NULL,
  PRIMARY KEY ("client_id", "client_user_id"),
  CONSTRAINT "fk_client_users_client_to_client_users" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_clients_client_to_users" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_client_to_client_users_external_id" to table: "client_to_client_users"
CREATE INDEX "idx_client_to_client_users_external_id" ON "client_to_client_users" ("external_id");
-- create "client_user_to_sessions" table
CREATE TABLE "client_user_to_sessions" (
  "client_user_id" uuid NOT NULL,
  "session_id" uuid NOT NULL,
  PRIMARY KEY ("client_user_id", "session_id"),
  CONSTRAINT "fk_client_users_client_user_to_sessions" FOREIGN KEY ("client_user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sessions_client_user_to_sessions" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- modify "refresh_tokens" table
ALTER TABLE "refresh_tokens" DROP CONSTRAINT "fk_refresh_tokens_user", ADD CONSTRAINT "fk_refresh_tokens_user" FOREIGN KEY ("user_id") REFERENCES "client_users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- drop "client_to_users" table
DROP TABLE "client_to_users";
-- drop "user_to_sessions" table
DROP TABLE "user_to_sessions";
-- drop "users" table
DROP TABLE "users";
