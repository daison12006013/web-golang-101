-- +goose Up
-- +goose StatementBegin
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ language 'plpgsql';

    CREATE TABLE users (
        id VARCHAR(36) DEFAULT uuid_generate_v4() PRIMARY KEY,
        first_name VARCHAR(255) NULL DEFAULT NULL,
        last_name VARCHAR(255) NULL DEFAULT NULL,
        email TEXT NOT NULL,
        email_hash VARCHAR(255) UNIQUE NOT NULL,
        email_verified_at TIMESTAMP(3) NULL DEFAULT NULL,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP(3) WITH TIME ZONE NULL DEFAULT NULL
    );

    CREATE INDEX idx_users_email_hash ON users(email_hash);
    CREATE INDEX idx_users_deleted_at ON users(deleted_at);

    CREATE TRIGGER update_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    DROP TRIGGER update_timestamp ON users;
    DROP TABLE users;
    DROP FUNCTION update_updated_at_column;
    DROP EXTENSION "uuid-ossp";
-- +goose StatementEnd
