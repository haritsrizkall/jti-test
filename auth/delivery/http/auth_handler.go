package http

import (
	"log"
	"net/http"
	"time"

	"github.com/haritsrizkall/jti-test/domain"
	"github.com/haritsrizkall/jti-test/utils"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) Register(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var request domain.RegisterRequest
	err := utils.DecodeBody(req, &request)
	if err != nil {
		utils.NewResponse(resp, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = h.authUsecase.Register(ctx, &request)
	if err != nil {
		utils.NewResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.NewResponse(resp, http.StatusOK, "Success", nil)
}

func (h *AuthHandler) Login(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var request domain.LoginRequest
	err := utils.DecodeBody(req, &request)
	if err != nil {
		utils.NewResponse(resp, http.StatusBadRequest, err.Error(), nil)
		return
	}

	loginResponse, err := h.authUsecase.Login(ctx, &request)
	if err != nil {
		utils.NewResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// set cookie
	utils.SetCookie(resp, "token", loginResponse.Token, time.Now().Add(24*time.Hour))

	utils.NewResponse(resp, http.StatusOK, "Success", loginResponse)
}

func (h *AuthHandler) Logout(resp http.ResponseWriter, req *http.Request) {
	utils.SetCookie(resp, "token", "", time.Now().Add(-1*time.Hour))

	utils.NewResponse(resp, http.StatusOK, "Success", nil)
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
	utils.SetCookie(resp, "token", token, time.Now().Add(24*time.Hour))

	http.Redirect(resp, req, "/input", http.StatusTemporaryRedirect)
}
