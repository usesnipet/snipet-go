-- reverse: create index "idx_client_bots_client_bot" to table: "client_bots"
DROP INDEX "idx_client_bots_client_bot";
-- reverse: drop index "idx_client_bots_client_id" from table: "client_bots"
CREATE INDEX "idx_client_bots_client_id" ON "client_bots" ("client_id");
-- reverse: drop index "idx_client_bots_bot_id" from table: "client_bots"
CREATE INDEX "idx_client_bots_bot_id" ON "client_bots" ("bot_id");
