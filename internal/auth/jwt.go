package auth

import (
	"errors"
	"time"

	ec "web-golang-101/pkg/errorcodes"
	"web-golang-101/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string) (string, *ec.Error) {
	// Generate access token
	expDuration, err := time.ParseDuration(utils.GetEnvWithDefault("JWT_EXP", "1h"))
	if err != nil {
		return "", ec.AsDefaultError(err)
	}
	expTime := time.Now().Add(expDuration)
	claims := &utils.Claims{
		UserID: userID,
		For:    "access_token",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tknStr, err := tkn.SignedString(utils.GetJWTKey())
	if err != nil {
		return "", ec.AsDefaultError(err)
	}

	return tknStr, nil
}

func GenerateRefreshToken(userID string) (string, *ec.Error) {
	// Generate refresh token
	refreshExpDuration, err := time.ParseDuration(utils.GetEnvWithDefault("JWT_REFRESH_EXP", "168h"))
	if err != nil {
		return "", ec.AsDefaultError(err)
	}
	refreshExpTime := time.Now().Add(refreshExpDuration)
	refreshClaims := &utils.Claims{
		UserID: userID,
		For:    "refresh_token",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpTime),
		},
	}

	refreshTkn := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTknStr, err := refreshTkn.SignedString(utils.GetJWTKey())
	if err != nil {
		return "", ec.AsDefaultError(err)
	}

	return refreshTknStr, nil
}

func VerifyRefreshToken(refreshTknStr string) (string, *ec.Error) {
	// Verify the refresh token
	claims := &utils.Claims{}
	tkn, err := jwt.ParseWithClaims(refreshTknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.GetJWTKey(), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", ec.AsBadRequest(errors.New("refresh token is invalid"))
		}
		return "", ec.AsDefaultError(err)
	}

	if !tkn.Valid {
		return "", ec.AsBadRequest(errors.New("refresh token is invalid"))
	}

	if claims.For != "refresh_token" {
		return "", ec.AsBadRequest(errors.New("claim `for` is not 'refresh_token'"))
	}

	return claims.UserID, nil
}
