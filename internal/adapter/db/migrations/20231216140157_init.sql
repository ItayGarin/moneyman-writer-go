CREATE TABLE IF NOT EXISTS "exp_businesses" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"name" text NOT NULL,
	"category" text NOT NULL,
	"sub_category" text NOT NULL,
	CONSTRAINT "exp_businesses_name_unique" UNIQUE("name")
);

CREATE TABLE IF NOT EXISTS "exp_desc_to_business" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"description" text NOT NULL,
	"business_name" text NOT NULL,
	CONSTRAINT "exp_desc_to_business_description_unique" UNIQUE("description")
);

CREATE TABLE IF NOT EXISTS "exp_transactions" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"identifier" text NOT NULL,
	"type" text NOT NULL,
	"status" text NOT NULL,
	"date" timestamp with time zone NOT NULL,
	"processed_date" timestamp with time zone NOT NULL,
	"original_amount" double precision NOT NULL,
	"original_currency" text NOT NULL,
	"charged_amount" double precision NOT NULL,
	"charged_currency" text NOT NULL,
	"description" text NOT NULL,
	"memo" text NOT NULL,
	"category" text NOT NULL,
	"account" text NOT NULL,
	"company_id" text NOT NULL,
	"hash" text NOT NULL,
	CONSTRAINT "exp_transactions_hash_unique" UNIQUE("hash")
);

DO $$ BEGIN
 ALTER TABLE "exp_desc_to_business" ADD CONSTRAINT "exp_desc_to_business_business_name_exp_businesses_name_fk" FOREIGN KEY ("business_name") REFERENCES "exp_businesses"("name") ON DELETE no action ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
