package auth

import (
	"database/sql"
	"encoding/json"

	"web-golang-101/pkg/utils"
)

// User a clone of queries.User{} from sqlc/queries/models.go
type User struct {
	ID              string         `json:"id"`
	FirstName       sql.NullString `json:"first_name"`
	LastName        sql.NullString `json:"last_name"`
	Email           string         `json:"email"`
	EmailHash       string         `json:"-"`
	EmailVerifiedAt sql.NullTime   `json:"email_verified_at"`
	Password        string         `json:"-"`
	CreatedAt       sql.NullTime   `json:"created_at"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
	DeletedAt       sql.NullTime   `json:"deleted_at"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	decryptedEmail, err := u.DecryptEmail()
	if err != nil {
		return nil, err
	}

	type Alias User
	return json.Marshal(&struct {
		*Alias
		Email string `json:"email"`
	}{
		Alias: (*Alias)(u),
		Email: decryptedEmail,
	})
}

func (u *User) DecryptEmail() (string, error) {
	c := utils.NewCrypt()
	return c.Decrypt(u.Email)
}
