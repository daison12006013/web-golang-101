package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"web-golang-101/pkg/db"
	ec "web-golang-101/pkg/errorcodes"
	"web-golang-101/pkg/utils"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Login(body []byte) (*TokenResponse, *ec.Error) {
	var input LoginInput
	err := json.Unmarshal(body, &input)
	if err != nil {
		return nil, ec.AsDefaultError(err)
	}

	err = utils.Validator().Struct(input)
	if err != nil {
		return nil, ec.AsDefaultError(err)
	}

	conn, err := db.NewConnection()
	if err != nil {
		return nil, ec.AsDatabaseConnection(err)
	}
	defer conn.DB.Close()

	q := conn.NewQuery()
	user, err := q.FindByEmail(context.Background(), utils.HashStr(input.Email))
	if err != nil {
		return nil, ec.AsQueryError(err)
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		return nil, ec.AsRecordNotFound(errors.New("email or password is invalid"))
	}

	tokenString, errc := GenerateToken(user.ID)
	if errc != nil {
		return nil, errc
	}

	refreshTokenString, errc := GenerateRefreshToken(user.ID)
	if errc != nil {
		return nil, errc
	}

	return &TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func RefreshToken(r *http.Request) (*TokenResponse, *ec.Error) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		return nil, ec.AsBadRequest(errors.New("authorization header is required"))
	}

	refreshToken := strings.TrimPrefix(bearer, "Bearer ")

	userID, errc := VerifyRefreshToken(refreshToken)
	if errc != nil {
		return nil, errc
	}

	tokenString, errc := GenerateToken(userID)
	if errc != nil {
		return nil, errc
	}

	refreshTokenString, errc := GenerateRefreshToken(userID)
	if errc != nil {
		return nil, errc
	}

	return &TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
