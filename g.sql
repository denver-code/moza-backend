BEGIN;

-- Update the balance
UPDATE bank_accounts 
SET balance = balance + 1000.00,
    last_activity = NOW()
WHERE id = 1;

-- Create a deposit transaction record
INSERT INTO transactions (
    from_account_id,
    to_account_id,
    amount,
    currency,
    description,
    type,
    status,
    reference,
    created_at,
    updated_at
)
VALUES (
    1,  -- from_account_id (same as recipient for deposit)
    1,  -- to_account_id
    1000.00,  -- amount
    'GBP',  -- currency
    'Initial deposit',  -- description
    'DEPOSIT',  -- type
    'COMPLETED',  -- status
    'TXN' || EXTRACT(EPOCH FROM NOW())::bigint,  -- reference
    NOW(),  -- created_at
    NOW()   -- updated_at
);

COMMIT;