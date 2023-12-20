CREATE TABLE IF NOT EXISTS exp_transaction_overrides (
	"id" bigserial PRIMARY KEY NOT NULL,
	"hash" text NOT NULL,

	"business_id" bigint,
	"charged_amount" double precision,

	"created_at" timestamp with time zone NOT NULL DEFAULT now(),

	FOREIGN KEY (hash) REFERENCES exp_transactions(hash),
	FOREIGN KEY (business_id) REFERENCES exp_businesses(id)
);