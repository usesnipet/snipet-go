-- modify "clients" table
ALTER TABLE "clients" ADD COLUMN "code" character(10) NOT NULL, ADD COLUMN "webhook_url" text NULL, ADD CONSTRAINT "uni_clients_code" UNIQUE ("code");
