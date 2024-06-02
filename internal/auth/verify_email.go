package auth

import (
	"context"
	"errors"

	"github.com/daison12006013/web-golang-101/pkg/db"
	ec "github.com/daison12006013/web-golang-101/pkg/errorcodes"
	"github.com/daison12006013/web-golang-101/pkg/utils"
)

func VerifyEmail(token string) (bool, *ec.Error) {
	if token == "" {
		return false, ec.AsBadRequest(errors.New("token is required"))
	}

	emailHash, err := utils.NewCrypt().Decrypt(token)
	if err != nil {
		return false, ec.AsDefaultError(err)
	}

	conn, err := db.NewConnection()
	if err != nil {
		return false, ec.AsDatabaseConnection(err)
	}
	defer conn.DB.Close()

	q := conn.NewQuery()

	user, err := q.FindByEmail(context.Background(), emailHash)
	if err != nil {
		return false, ec.AsQueryError(err)
	}
	if user.EmailHash == "" {
		return false, ec.AsRecordNotFound(errors.New("verification token not found"))
	}
	if !user.EmailVerifiedAt.Time.IsZero() {
		return false, ec.AsConflict(errors.New("account is already verified"))
	}

	err = q.UpdateVerifiedAt(context.Background(), emailHash)
	if err != nil {
		return false, ec.AsQueryError(err)
	}

	return true, nil
}
