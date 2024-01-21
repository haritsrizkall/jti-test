package http

import (
	"log"
	"net/http"

	"github.com/haritsrizkall/jti-test/domain"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) LoginWithGoogle(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	url := h.authUsecase.LoginWithGoogle(ctx)
	http.Redirect(resp, req, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) LoginWithGoogleCallback(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	code := req.URL.Query().Get("code")
	token, err := h.authUsecase.LoginWithGoogleCallback(ctx, code)
	if err != nil {
		log.Println(err)
		http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
		return
	}

	// set cookie
	resp.Header().Set("Set-Cookie", "token="+token+"; Path=/; HttpOnly")

	http.Redirect(resp, req, "/input", http.StatusTemporaryRedirect)
}
