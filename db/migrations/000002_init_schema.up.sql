CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "lastname" varchar(50) NOT NULL,
  "firstname" varchar(50) NOT NULL,
  "hashed_password" varchar NOT NULL,
  "phone" varchar(11) UNIQUE NOT NULL,
  "address" varchar(300) NOT NULL,
  "email" varchar(200) UNIQUE NOT NULL,
  "is_admin" bool DEFAULT (false),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "shops" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(100),
  "phone" varchar(11) UNIQUE NOT NULL,
  "address" varchar(300) NOT NULL,
  "email" varchar(200) UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" text,
  "price" float NOT NULL,
  "image" varchar NOT NULL,
  "qty_aval" int NOT NULL,
  "shop_id" bigint,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "carts" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigint,
  "qty_bought" int NOT NULL,
  "unit_price" float,
  "total_price" float NOT NULL,
  "user_id" bigint,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigint,
  "qty_bought" int NOT NULL,
  "unit_price" float,
  "total_price" float NOT NULL,
  "user_id" bigint,
  "session_id" bigint UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "invoice" (
  "id" bigserial PRIMARY KEY,
  "session_id" bigint,
  "order_cost" float NOT NULL,
  "shipping_cost" float NOT NULL,
  "invoice_no" bigserial NOT NULL,
  "user_id" bigint,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "shipping" (
  "id" bigserial PRIMARY KEY,
  "invoice_id" bigint,
  "courier_name" varchar NOT NULL,
  "eta" int NOT NULL,
  "time_left" timestamptz NOT NULL,
  "time_arrive" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "shops" ("name");

CREATE INDEX ON "products" ("name");

CREATE INDEX ON "products" ("shop_id");

CREATE INDEX ON "carts" ("user_id");

COMMENT ON COLUMN "products"."description" IS 'description of the item';

COMMENT ON COLUMN "carts"."user_id" IS 'to know which user has a cart';

COMMENT ON COLUMN "orders"."user_id" IS 'to know which user has an order';

COMMENT ON COLUMN "orders"."session_id" IS 'to track all orders';

ALTER TABLE "products" ADD FOREIGN KEY ("shop_id") REFERENCES "shops" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "invoice" ADD FOREIGN KEY ("session_id") REFERENCES "orders" ("session_id");

ALTER TABLE "invoice" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "shipping" ADD FOREIGN KEY ("invoice_id") REFERENCES "invoice" ("id");
