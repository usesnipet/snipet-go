-- create "api_keys" table
CREATE TABLE "api_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "key_id" character varying(255) NOT NULL,
  "key" text NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "expires_at" timestamp NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_api_keys_key" UNIQUE ("key")
);
