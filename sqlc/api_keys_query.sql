-- name: ApiKeyExists :one
select exists(select 1 from api_keys where key = $1 and deleted_at is null) as exists;

-- name: FindByApiKey :one
select * from api_keys where key = $1 and deleted_at is null;

-- name: InsertApiKey :one
insert into api_keys (user_id, key) values ($1, $2) RETURNING key;

-- name: DeleteApiKey :exec
update api_keys set deleted_at = now() where user_id = $1 and key = $2 and deleted_at is null;

-- name: FindByUserId :many
select * from api_keys where user_id = $1 and deleted_at is null;
