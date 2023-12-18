-- name: GetAllTransactions :many
SELECT * FROM exp_transactions;

-- name: GetTransactionByHash :one
SELECT * FROM exp_transactions WHERE hash = $1;

-- name: InsertTransaction :exec
INSERT INTO exp_transactions (
    identifier, 
    type, 
    status, 
    date, 
    processed_date, 
    original_amount, 
    original_currency, 
    charged_amount, 
    charged_currency, 
    description, 
    memo, 
    category, 
    account, 
    company_id, 
    hash
) VALUES 
( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
ON CONFLICT (hash) DO NOTHING;

-- name: GetUncategorizedDescriptions :many
SELECT DISTINCT(description) FROM exp_transactions
WHERE description not in (SELECT description FROM exp_desc_to_business);