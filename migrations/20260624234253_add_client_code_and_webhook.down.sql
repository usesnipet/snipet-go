-- reverse: modify "clients" table
ALTER TABLE "clients" DROP CONSTRAINT "uni_clients_code", DROP COLUMN "webhook_url", DROP COLUMN "code";
