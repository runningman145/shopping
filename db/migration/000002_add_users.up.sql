CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "changed_password_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00+00',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "role" text NOT NULL DEFAULT 'user'
);

ALTER TABLE "products" ADD COLUMN "user_id" BIGINT NOT NULL;

ALTER TABLE "products" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

