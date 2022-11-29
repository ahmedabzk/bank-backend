-- name: CreateEntries :one
INSERT INTO entries(
    account_id,
    amount
)VALUES(
    $1, $2
)RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
LIMIT $1
OFFSET $2;

-- name: UpdateEntry :one
UPDATE entries
SET account_id = $1,
amount = $2
WHERE id = $1
RETURNING *;
