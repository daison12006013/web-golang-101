// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: api_keys_query.sql

package queries

import (
	"context"
)

const apiKeyExists = `-- name: ApiKeyExists :one
select exists(select 1 from api_keys where key = $1 and deleted_at is null) as exists
`

func (q *Queries) ApiKeyExists(ctx context.Context, key string) (bool, error) {
	row := q.db.QueryRowContext(ctx, apiKeyExists, key)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const deleteApiKey = `-- name: DeleteApiKey :exec
update api_keys set deleted_at = now() where user_id = $1 and key = $2 and deleted_at is null
`

type DeleteApiKeyParams struct {
	UserID string
	Key    string
}

func (q *Queries) DeleteApiKey(ctx context.Context, arg DeleteApiKeyParams) error {
	_, err := q.db.ExecContext(ctx, deleteApiKey, arg.UserID, arg.Key)
	return err
}

const findByApiKey = `-- name: FindByApiKey :one
select id, user_id, key, created_at, deleted_at from api_keys where key = $1 and deleted_at is null
`

func (q *Queries) FindByApiKey(ctx context.Context, key string) (ApiKey, error) {
	row := q.db.QueryRowContext(ctx, findByApiKey, key)
	var i ApiKey
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Key,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const findByUserId = `-- name: FindByUserId :many
select id, user_id, key, created_at, deleted_at from api_keys where user_id = $1 and deleted_at is null
`

func (q *Queries) FindByUserId(ctx context.Context, userID string) ([]ApiKey, error) {
	rows, err := q.db.QueryContext(ctx, findByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ApiKey
	for rows.Next() {
		var i ApiKey
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Key,
			&i.CreatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertApiKey = `-- name: InsertApiKey :one
insert into api_keys (user_id, key) values ($1, $2) RETURNING key
`

type InsertApiKeyParams struct {
	UserID string
	Key    string
}

func (q *Queries) InsertApiKey(ctx context.Context, arg InsertApiKeyParams) (string, error) {
	row := q.db.QueryRowContext(ctx, insertApiKey, arg.UserID, arg.Key)
	var key string
	err := row.Scan(&key)
	return key, err
}