package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func GetJWTKey() []byte {
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		panic("JWT_KEY is not set")
	}
	return []byte(jwtKey)
}

type Claims struct {
	UserID string `json:"user_id"`
	For    string `json:"for"`
	jwt.RegisteredClaims
}

func RequireJWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing auth token", http.StatusUnauthorized)
			return
		}

		resp := NewResponse(w)

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			resp.WriteErrorResponse("Invalid access token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(bearerToken[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return GetJWTKey(), nil
		})

		if err != nil {
			resp.WriteErrorResponse("Invalid access token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			if claims.For != "access_token" {
				resp.WriteErrorResponse("Invalid access token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextKey("UserID"), claims.UserID)
			r = r.WithContext(ctx)
		} else {
			resp.WriteErrorResponse("Invalid access token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
