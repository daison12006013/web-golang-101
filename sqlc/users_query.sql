-- name: UserExists :one
select exists(select 1 from users where id = $1 and deleted_at is null) as exists;

-- name: FindByEmail :one
select * from users where email_hash = $1 and deleted_at is null;

-- name: InsertUser :exec
insert into users (email, email_hash, password, first_name, last_name) values ($1, $2, $3, $4, $5);

-- name: UpdateVerifiedAt :exec
update users set email_verified_at = now() where email_hash = $1;
