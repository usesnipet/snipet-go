-- reverse: drop "client_to_bots" table
CREATE TABLE "client_to_bots" (
  "client_id" uuid NOT NULL,
  "bot_id" uuid NOT NULL,
  PRIMARY KEY ("client_id", "bot_id"),
  CONSTRAINT "fk_bots_client_to_bots" FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_clients_client_to_bots" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
