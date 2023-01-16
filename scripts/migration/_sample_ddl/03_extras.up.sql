-- Table Definition
CREATE TABLE "public"."extras" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "locale" text,
    "admin_notes" text,
    "bio" text,
    "about" text,
    PRIMARY KEY ("id")
);
