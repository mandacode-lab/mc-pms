-- Create "namespaces" table
CREATE TABLE "public"."namespaces" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "description" character varying NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "projects" table
CREATE TABLE "public"."projects" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "description" character varying NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "namespace_projects" character varying NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "projects_namespaces_projects" FOREIGN KEY ("namespace_projects") REFERENCES "public"."namespaces" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
