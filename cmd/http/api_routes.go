package main

import (
	"net/http"

	"web-golang-101/internal/apikeys"
	"web-golang-101/internal/auth"
	ec "web-golang-101/pkg/errorcodes"
	"web-golang-101/pkg/utils"

	"github.com/go-chi/chi/v5"
)

func ApiRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Route"))
	})

	router.Group(publicApiRoutes())

	router.Group(requiredApiKeyRoutes())

	router.Group(requiredJwtAuthenticationRoutes())

	return router
}

func publicApiRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/register", registerHandler)
		r.Get("/verify-email/{token}", verifyEmailHandler)
		r.Post("/login", loginHandler)
		r.Post("/refresh-token", refreshTokenHandler)
	}
}

func requiredApiKeyRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(utils.RequireApiKeyMiddleware)
	}
}

func requiredJwtAuthenticationRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(utils.RequireJWTAuthMiddleware)
		r.Get("/api-keys", apiKeyListsHandler)
		r.Post("/api-keys", apiKeyGenerateHandler)
		r.Delete("/api-keys/{key}", apiKeyDeleteHandler)
	}
}

// apiKeyListsHandler
func apiKeyListsHandler(w http.ResponseWriter, r *http.Request) {
	logic := func(body []byte) (any, *ec.Error) {
		return apikeys.Lists(r.Context().Value(utils.ContextKey("UserID")).(string))
	}
	utils.CommonHandler(w, r, "API Key Lists", logic)
}

// apiKeyGenerateHandler
func apiKeyGenerateHandler(w http.ResponseWriter, r *http.Request) {
	logic := func(body []byte) (any, *ec.Error) {
		return apikeys.Generate(r.Context().Value(utils.ContextKey("UserID")).(string))
	}
	utils.CommonHandler(w, r, "API Key Generated", logic)
}

// apiKeyDeleteHandler
func apiKeyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	resp := utils.NewResponse(w)

	errc := apikeys.Delete(
		r.Context().Value(utils.ContextKey("UserID")).(string),
		chi.URLParam(r, "key"),
	)
	if errc != nil {
		utils.NewResponse(w).HandleErrorCode(errc)
	}

	resp.WriteSuccessResponse("API Key Deleted", nil)
}

// registerHandler
func registerHandler(w http.ResponseWriter, r *http.Request) {
	resp := utils.NewResponse(w)

	logic := func(body []byte) (any, *ec.Error) {
		return auth.Register(r, body)
	}

	data, errc := utils.CommonJsonHandler(r, logic)
	if errc != nil {
		if ok := resp.WriteValidationError(errc); ok {
			return
		}

		resp.HandleErrorCode(errc)
		return
	}

	resp.WriteSuccessResponse("Registered successfully", data)
}

// loginHandler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	resp := utils.NewResponse(w)
	logic := func(body []byte) (any, *ec.Error) {
		return auth.Login(body)
	}

	data, errc := utils.CommonJsonHandler(r, logic)
	if errc != nil {
		if ok := resp.WriteValidationError(errc); ok {
			return
		}

		resp.HandleErrorCode(errc)
		return
	}

	resp.WriteSuccessResponse("Logged in successfully", data)
}

// refreshTokenHandler
func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	logic := func(body []byte) (any, *ec.Error) {
		return auth.RefreshToken(r)
	}
	utils.CommonHandler(w, r, "Token refreshed successfully", logic)
}

// verifyEmailHandler
func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	logic := func(body []byte) (any, *ec.Error) {
		result, errc := auth.VerifyEmail(token)
		if !result {
			return nil, errc
		}
		return nil, nil
	}
	utils.CommonHandler(w, r, "Verified email successfully", logic)
}
