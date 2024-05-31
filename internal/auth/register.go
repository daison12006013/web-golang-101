package auth

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"web-golang-101/internal/mailer"
	"web-golang-101/pkg/db"
	ec "web-golang-101/pkg/errorcodes"
	"web-golang-101/pkg/utils"
	"web-golang-101/sqlc/queries"
)

type RegisterInput struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name,omitempty" validate:"required"`
	LastName        string `json:"last_name,omitempty" validate:"required"`
}

func Register(r *http.Request, body []byte) (*User, *ec.Error) {
	var input RegisterInput
	err := json.Unmarshal(body, &input)
	if err != nil {
		return nil, ec.AsDefaultError(err)
	}

	err = utils.Validator().Struct(input)
	if err != nil {
		return nil, ec.AsDefaultError(err)
	}

	c := utils.NewCrypt()
	email, err := c.Encrypt(input.Email)
	if err != nil {
		return nil, ec.AsDefaultError(err)
	}

	password, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, ec.AsDefaultError(err)
	}

	user := &User{
		Email:     email,
		EmailHash: utils.HashStr(input.Email),
		Password:  password,
		FirstName: sql.NullString{String: input.FirstName, Valid: input.FirstName != ""},
		LastName:  sql.NullString{String: input.LastName, Valid: input.LastName != ""},
	}

	conn, err := db.NewConnection()
	if err != nil {
		return nil, ec.AsDatabaseConnection(err)
	}
	defer conn.DB.Close()

	errEc := createUser(conn, user)
	if errEc != nil {
		return nil, errEc
	}

	go func() {
		if err := sendConfirmationEmail(r, user); err != nil {
			utils.Logger().Err(err).Msg("failed to send confirmation email")
		}
	}()

	return user, nil
}

func createUser(conn *db.DBC, user *User) *ec.Error {
	q := conn.NewQuery()

	// Start a new transaction
	tx, err := conn.DB.Begin()
	if err != nil {
		return ec.AsDatabaseConnection(err)
	}

	exists, err := q.UserExists(context.Background(), user.EmailHash)
	if err != nil {
		return ec.AsQueryError(err)
	}

	if exists {
		return ec.AsConflict(errors.New("email already exists"))
	}

	err = q.InsertUser(context.Background(), queries.InsertUserParams{
		Email:     user.Email,
		EmailHash: user.EmailHash,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})

	if err != nil {
		// If there is an error, rollback the transaction
		tx.Rollback()
		return ec.AsQueryError(err)
	}

	// If everything is fine, commit the transaction
	err = tx.Commit()
	if err != nil {
		return ec.AsDatabaseConnection(err)
	}

	return nil
}

type confirmationEmailData struct {
	Name             string
	VerificationLink string
}

//go:embed register_template.html
var tmplFS embed.FS

func sendConfirmationEmail(r *http.Request, user *User) error {
	name := fmt.Sprintf("%s %s", user.FirstName.String, user.LastName.String)
	to, err := user.DecryptEmail()
	if err != nil {
		return err
	}

	c := utils.NewCrypt()
	verifyEmail, err := c.Encrypt(user.EmailHash)
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFS(tmplFS, "register_template.html")
	if err != nil {
		return err
	}

	data := confirmationEmailData{
		Name:             name,
		VerificationLink: fmt.Sprintf("%s/verify-email/%s", utils.GetHost(r), verifyEmail),
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return err
	}

	htmlBody := tpl.String()

	m := mailer.New()
	resp, err := m.SendEmail(to, "Confirm your Registration", htmlBody)
	if err != nil {
		return err
	}
	utils.Logger().Debug().Msgf("Email sent %s", resp)

	return nil
}
