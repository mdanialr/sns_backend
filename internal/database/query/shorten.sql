-- name: GetShorten :one
SELECT * FROM shorten
WHERE id = $1 LIMIT 1;

-- name: GetShortenByUrl :one
SELECT * FROM shorten
WHERE url = $1 LIMIT 1;

-- name: CreateShorten :one
INSERT INTO shorten (
  url, target, permanent
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: DeleteShorten :exec
DELETE FROM shorten
WHERE id = $1;
