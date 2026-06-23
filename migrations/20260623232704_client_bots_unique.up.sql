-- drop index "idx_client_bots_bot_id" from table: "client_bots"
DROP INDEX "idx_client_bots_bot_id";
-- drop index "idx_client_bots_client_id" from table: "client_bots"
DROP INDEX "idx_client_bots_client_id";
-- create index "idx_client_bots_client_bot" to table: "client_bots"
CREATE UNIQUE INDEX "idx_client_bots_client_bot" ON "client_bots" ("client_id", "bot_id");
