-- name: GetSend :one
SELECT * FROM send
WHERE id = $1 LIMIT 1;

-- name: GetSendByUrl :one
SELECT * FROM send
WHERE url = $1 LIMIT 1;

-- name: ListSend :many
SELECT * FROM send
ORDER BY updated_at;

-- name: CreateSend :one
INSERT INTO send (
  url, file, size, permanent
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: UpdateSend :one
UPDATE send
set url = $2,
    file = $3,
    size = $4,
    permanent = $5,
    updated_at = $6
WHERE id = $1
RETURNING *;

-- name: DeleteSend :exec
DELETE FROM send
WHERE id = $1;
