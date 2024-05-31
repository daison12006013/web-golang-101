-- +goose Up
-- +goose StatementBegin
CREATE TABLE api_keys (
    id VARCHAR(36) DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id),
    key VARCHAR(36) NOT NULL,
    created_at TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3) WITH TIME ZONE NULL DEFAULT NULL
);

CREATE INDEX idx_api_keys_key ON api_keys(key);
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE api_keys;
-- +goose StatementEnd
