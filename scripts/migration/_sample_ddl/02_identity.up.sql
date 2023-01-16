-- Table Definition
CREATE TABLE "public"."identity" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" text,
    "given_name" text,
    "family_name" text,
    "nick_name" text,
    "gender" text,
    "birthdate" text,
    "picture" text,
    PRIMARY KEY ("id")
);
