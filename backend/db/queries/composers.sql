-- name: GetComposer :one
SELECT * FROM composers
WHERE id = ? LIMIT 1;

-- name: GetAllComposers :many
SELECT * FROM composers
ORDER BY name;
