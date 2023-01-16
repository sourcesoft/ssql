-- Table Definition
CREATE TABLE "public"."user" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "username" text,
    "email" text,
    "email_verified" bool NOT NULL DEFAULT false,
    "active" bool NOT NULL DEFAULT false,
    "created_at" bigint NOT NULL,
    "updated_at" bigint NOT NULL,
    "deleted_at" bigint,
    PRIMARY KEY ("id")
);
