package apikeys

import (
	"context"

	"web-golang-101/pkg/db"
	ec "web-golang-101/pkg/errorcodes"
	"web-golang-101/sqlc/queries"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type ApiKey struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Key       string    `json:"key"`
	CreatedAt null.Time `json:"created_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

func Lists(userID string) (*[]ApiKey, *ec.Error) {
	conn, err := db.NewConnection()
	if err != nil {
		return nil, ec.AsDatabaseConnection(err)
	}
	defer conn.DB.Close()
	q := conn.NewQuery()

	records, err := q.FindByUserId(context.Background(), userID)
	if err != nil {
		return nil, ec.AsQueryError(err)
	}

	apiKeys := make([]ApiKey, len(records))
	for i, record := range records {
		apiKeys[i] = ApiKey{
			ID:        record.ID,
			UserID:    record.UserID,
			Key:       record.Key,
			CreatedAt: null.NewTime(record.CreatedAt.Time, record.CreatedAt.Valid),
			DeletedAt: null.NewTime(record.DeletedAt.Time, record.DeletedAt.Valid),
		}
	}

	return &apiKeys, nil
}

func Generate(userID string) (*string, *ec.Error) {
	conn, err := db.NewConnection()
	if err != nil {
		return nil, ec.AsDatabaseConnection(err)
	}
	defer conn.DB.Close()

	q := conn.NewQuery()

	apiKey, err := q.InsertApiKey(
		context.Background(),
		queries.InsertApiKeyParams{
			UserID: userID,
			Key:    uuid.New().String(),
		})
	if err != nil {
		return nil, ec.AsQueryError(err)
	}

	return &apiKey, nil
}

func Delete(userID, key string) *ec.Error {
	conn, err := db.NewConnection()
	if err != nil {
		return ec.AsDatabaseConnection(err)
	}
	defer conn.DB.Close()

	q := conn.NewQuery()

	err = q.DeleteApiKey(
		context.Background(),
		queries.DeleteApiKeyParams{
			UserID: userID,
			Key:    key,
		})
	if err != nil {
		return ec.AsQueryError(err)
	}

	return nil
}
