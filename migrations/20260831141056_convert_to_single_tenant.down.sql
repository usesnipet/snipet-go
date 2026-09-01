-- reverse: drop "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "email" character varying(255) NOT NULL,
  "password_hash" character varying(255) NULL,
  "picture" character varying(255) NULL,
  "is_admin" boolean NOT NULL DEFAULT false,
  "challenges" jsonb NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email")
);
-- reverse: drop "tokens" table
CREATE TABLE "tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "type" character varying(255) NOT NULL,
  "hash" character varying(255) NOT NULL,
  "user_id" uuid NOT NULL,
  "expires_at" timestamp NOT NULL,
  "revoked_at" timestamp NULL,
  "metadata" jsonb NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_tokens" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE UNIQUE INDEX "idx_tokens_hash" ON "tokens" ("hash");
CREATE INDEX "idx_tokens_type" ON "tokens" ("type");
CREATE INDEX "idx_tokens_user_id" ON "tokens" ("user_id");
-- reverse: drop "accounts" table
CREATE TABLE "accounts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "provider" character varying(255) NOT NULL,
  "external_id" character varying(255) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_accounts" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE UNIQUE INDEX "idx_accounts_provider_external_id" ON "accounts" ("provider", "external_id");
CREATE INDEX "idx_accounts_user_id" ON "accounts" ("user_id");
-- reverse: drop "tenants" table
CREATE TABLE "tenants" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "icon" character varying(255) NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_tenants_slug" UNIQUE ("slug")
);
-- reverse: drop "tenant_invitations" table
CREATE TABLE "tenant_invitations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tenant_id" uuid NOT NULL,
  "email" character varying(255) NOT NULL,
  "token" character varying(255) NOT NULL,
  "role" character varying(255) NOT NULL,
  "status" character varying(255) NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_tenant_invitations_token" UNIQUE ("token"),
  CONSTRAINT "fk_tenants_invitations" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_tenant_invitations_status" ON "tenant_invitations" ("status");
CREATE INDEX "idx_tenant_invitations_tenant_id" ON "tenant_invitations" ("tenant_id");
-- reverse: drop "members" table
CREATE TABLE "members" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "tenant_id" uuid NOT NULL,
  "role" character varying(255) NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tenants_members" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_members" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_members_tenant_id" ON "members" ("tenant_id");
CREATE UNIQUE INDEX "idx_members_user_tenant" ON "members" ("user_id", "tenant_id");
-- reverse: modify "sessions" table
ALTER TABLE "sessions" ADD COLUMN "tenant_id" uuid NOT NULL;
-- reverse: modify "llms" table
ALTER TABLE "llms" ADD COLUMN "tenant_id" uuid NOT NULL;
-- reverse: modify "knowledges" table
ALTER TABLE "knowledges" ADD COLUMN "tenant_id" uuid NOT NULL;
-- reverse: modify "knowledge_items" table
ALTER TABLE "knowledge_items" ADD COLUMN "tenant_id" uuid NOT NULL;
-- reverse: modify "knowledge_indices" table
ALTER TABLE "knowledge_indices" ADD COLUMN "tenant_id" uuid NOT NULL;
-- reverse: modify "indexed_knowledge_items" table
ALTER TABLE "indexed_knowledge_items" ADD COLUMN "tenant_id" uuid NOT NULL;
-- reverse: modify "executions" table
ALTER TABLE "executions" ADD COLUMN "tenant_id" uuid NOT NULL;
-- reverse: modify "execution_messages" table
ALTER TABLE "execution_messages" ADD COLUMN "tenant_id" uuid NOT NULL;
-- reverse: modify "apps" table
ALTER TABLE "apps" ADD COLUMN "tenant_id" uuid NOT NULL;
-- reverse: modify "api_keys" table
ALTER TABLE "api_keys" ADD COLUMN "tenant_id" uuid NOT NULL;
-- reverse: modify "agents" table
ALTER TABLE "agents" ADD COLUMN "tenant_id" uuid NOT NULL;
