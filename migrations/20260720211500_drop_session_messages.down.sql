-- reverse: drop "session_messages" table
CREATE TABLE "session_messages" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "session_id" uuid NOT NULL,
  "role" character varying(255) NOT NULL,
  "parts" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_sessions_session_messages" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_session_messages" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_session_messages_session_id" ON "session_messages" ("session_id");
CREATE INDEX "idx_session_messages_user_id" ON "session_messages" ("user_id");
