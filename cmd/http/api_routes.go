package main

import (
	"net/http"

	"github.com/daison12006013/web-golang-101/internal/apikeys"
	"github.com/daison12006013/web-golang-101/internal/auth"
	ec "github.com/daison12006013/web-golang-101/pkg/errorcodes"
	"github.com/daison12006013/web-golang-101/pkg/utils"

	"github.com/go-chi/chi/v5"
)

func ApiRoutes() *chi.Mux {
	router := chi.NewRouter()

	utils.SetDefaultMiddlewares(router)

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
// @Summary API Key Lists
// @Description Get API Key Lists
// @Tags API Keys
// @Accept json
// @Produce json
// @Router /api-keys [get]
// @Security ApiKeyAuth
func apiKeyListsHandler(w http.ResponseWriter, r *http.Request) {
	logic := func(body []byte) (any, *ec.Error) {
		return apikeys.Lists(r.Context().Value(utils.ContextKey("UserID")).(string))
	}
	utils.CommonHandler(w, r, "API Key Lists", logic)
}

// apiKeyGenerateHandler
// @Summary API Key Generate
// @Description Generate API Key
// @Tags API Keys
// @Accept json
// @Produce json
// @Router /api-keys [post]
// @Security ApiKeyAuth
func apiKeyGenerateHandler(w http.ResponseWriter, r *http.Request) {
	logic := func(body []byte) (any, *ec.Error) {
		return apikeys.Generate(r.Context().Value(utils.ContextKey("UserID")).(string))
	}
	utils.CommonHandler(w, r, "API Key Generated", logic)
}

// apiKeyDeleteHandler
// @Summary API Key Delete
// @Description Delete API Key
// @Tags API Keys
// @Accept json
// @Produce json
// @Router /api-keys/{key} [delete]
// @Security ApiKeyAuth
func apiKeyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	resp := utils.NewResponse(w)

	errc := apikeys.Delete(
		r.Context().Value(utils.ContextKey("UserID")).(string),
		chi.URLParam(r, "key"),
	)
	if errc != nil {
		resp.HandleErrorCode(errc)
		return
	}

	resp.WriteSuccessResponse("API Key Deleted", nil)
}

// registerHandler
// @Summary Register
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Router /register [post]
// @Param body body auth.RegisterInput true "Register Input"
// @Success 200 {object} utils.Response
func registerHandler(w http.ResponseWriter, r *http.Request) {
	resp := utils.NewResponse(w)

	logic := func(body []byte) (any, *ec.Error) {
		return auth.Register(r, body)
	}

	data, errc := utils.CommonJsonHandler(r, logic)
	if errc != nil {
		resp.HandleErrorCode(errc)
		return
	}

	resp.WriteSuccessResponse("Registered successfully", data)
}

// loginHandler
// @Summary Login
// @Description Login to the system
// @Tags Auth
// @Accept json
// @Produce json
// @Router /login [post]
// @Param body body auth.LoginInput true "Login Input"
// @Success 200 {object} utils.Response
func loginHandler(w http.ResponseWriter, r *http.Request) {
	resp := utils.NewResponse(w)
	logic := func(body []byte) (any, *ec.Error) {
		return auth.Login(body)
	}

	data, errc := utils.CommonJsonHandler(r, logic)
	if errc != nil {
		resp.HandleErrorCode(errc)
		return
	}

	resp.WriteSuccessResponse("Logged in successfully", data)
}

// refreshTokenHandler
// @Summary Refresh Token
// @Description Refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Router /refresh-token [post]
// @Security BearerAuth
// @Success 200 {object} utils.Response
func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	logic := func(body []byte) (any, *ec.Error) {
		return auth.RefreshToken(r)
	}
	utils.CommonHandler(w, r, "Token refreshed successfully", logic)
}

// verifyEmailHandler
// @Summary Verify Email
// @Description Verify email
// @Tags Auth
// @Accept json
// @Produce json
// @Router /verify-email/{token} [get]
// @Param token path string true "Token"
// @Success 200 {object} utils.Response
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
