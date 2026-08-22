-- modify "apps" table
ALTER TABLE "apps" ADD COLUMN "public" boolean NOT NULL DEFAULT false;
