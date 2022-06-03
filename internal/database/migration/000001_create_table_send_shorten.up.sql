BEGIN;

CREATE TABLE "shorten" (
  "id" bigserial PRIMARY KEY,
  "url" varchar NOT NULL,
  "target" varchar NOT NULL,
  "permanent" boolean DEFAULT false NOT NULL,
  "created_at" timestamptz DEFAULT (now()) NOT NULL,
  "updated_at" timestamptz DEFAULT (now()) NOT NULL
);

CREATE TABLE "send" (
  "id" bigserial PRIMARY KEY,
  "url" varchar NOT NULL ,
  "file" varchar NOT NULL,
  "size" varchar NOT NULL,
  "permanent" boolean DEFAULT false NOT NULL,
  "created_at" timestamptz DEFAULT (now()) NOT NULL,
  "updated_at" timestamptz DEFAULT (now()) NOT NULL
);

COMMIT;
