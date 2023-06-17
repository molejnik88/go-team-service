-- Create "teams" table
CREATE TABLE "public"."teams" (
  "uuid" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("uuid")
);
-- Create "team_members" table
CREATE TABLE "public"."team_members" (
  "uuid" uuid NOT NULL DEFAULT gen_random_uuid(),
  "team_uuid" uuid NOT NULL,
  "email" text NOT NULL,
  "is_admin" boolean NOT NULL DEFAULT false,
  "is_owner" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("uuid"),
  CONSTRAINT "team_uuid" FOREIGN KEY ("team_uuid") REFERENCES "public"."teams" ("uuid") ON UPDATE NO ACTION ON DELETE NO ACTION
);
