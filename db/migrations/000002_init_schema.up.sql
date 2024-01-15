CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "lastname" varchar(50) NOT NULL,
  "firstname" varchar(50) NOT NULL,
  "hashed_password" varchar NOT NULL,
  "phone" varchar(11) UNIQUE NOT NULL,
  "address" varchar(300) NOT NULL,
  "email" varchar(200) UNIQUE NOT NULL,
  "is_admin" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "shops" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(100) UNIQUE NOT NULL,
  "phone" varchar(11) UNIQUE NOT NULL,
  "address" varchar(300) NOT NULL,
  "email" varchar(200) UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" text NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "images" varchar[] NOT NULL,
  "qty_aval" int NOT NULL,
  "shop_id" bigint NOT NULL,
  "shop_name" varchar NOT NULL,
  "category_id" bigint NOT NULL,
  "category_name" varchar NOT NULL,
  "sub_category_id" bigint NOT NULL,
  "sub_category_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "category" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sub_category" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "category_id" bigint NOT NULL,
  "category_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);


CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "items" jsonb NOT NULL,
  "transaction_ref" varchar UNIQUE NOT NULL,
  "pay_ref" varchar UNIQUE NOT NULL,
  "status" varchar NOT NULL ,
  "total_price" numeric(10,2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);


CREATE INDEX ON "shops" ("name");

CREATE INDEX ON "products" ("name");

CREATE INDEX ON "products" ("shop_id");

CREATE INDEX ON "products" ("shop_name");

CREATE INDEX ON "category" ("name");

CREATE INDEX ON "sub_category" ("name");

CREATE INDEX ON "sub_category" ("category_id");


COMMENT ON COLUMN "products"."description" IS 'description of the item';

COMMENT ON COLUMN "orders"."user_id" IS 'to know which user has an order';


ALTER TABLE "products" ADD FOREIGN KEY ("shop_id") REFERENCES "shops" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("shop_name") REFERENCES "shops" ("name");

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("category_name") REFERENCES "category" ("name");

ALTER TABLE "products" ADD FOREIGN KEY ("sub_category_id") REFERENCES "sub_category" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("sub_category_name") REFERENCES "sub_category" ("name");

ALTER TABLE "sub_category" ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");

ALTER TABLE "sub_category" ADD FOREIGN KEY ("category_name") REFERENCES "category" ("name");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");