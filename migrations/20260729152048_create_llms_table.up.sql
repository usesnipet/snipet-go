-- create "llms" table
CREATE TABLE "llms" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "provider" character varying(255) NOT NULL,
  "configuration" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
