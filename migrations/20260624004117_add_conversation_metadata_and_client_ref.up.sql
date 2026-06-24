-- modify "conversations" table
ALTER TABLE "conversations" ADD COLUMN "client_id" uuid NOT NULL, ADD COLUMN "metadata" jsonb NOT NULL, ADD CONSTRAINT "fk_clients_conversations" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- create index "idx_conversations_client_id" to table: "conversations"
CREATE INDEX "idx_conversations_client_id" ON "conversations" ("client_id");
