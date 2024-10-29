package interfaces

import (
	"fmt"
	"net/http"

	"github.com/Lucas-Linhar3s/JobHub/modules/auth/application"
	v1 "github.com/Lucas-Linhar3s/JobHub/pkg/http/response/v1"
	"github.com/gin-gonic/gin"
)

// AuthHandler is a struct that represents the handler of auth
type AuthHandler struct {
	app *application.AuthApp
}

// NewAuthHandler is a function that returns a new AuthHandler struct
func NewAuthHandler(app *application.AuthApp) *AuthHandler {
	return &AuthHandler{
		app: app,
	}
}

// RegisterUser is a function that registers a new user
func (h *AuthHandler) RegisterUser(ctx *gin.Context) {
	var req application.UserRegisterReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if err := h.app.RegisterUser(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, http.StatusNoContent, nil)
}

// LoginWithEmailAndPassword is a function that logs in with email and password
func (h *AuthHandler) LoginWithEmailAndPassword(ctx *gin.Context) {
	var req application.UserRegisterReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}

	response, err := h.app.LoginWithEmailAndPassword(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, http.StatusOK, response)
}

// RedirectLoginOauth is a function that redirects to the login with oauth
func (h *AuthHandler) RedirectLoginOauth(ctx *gin.Context) {
	oatuhProvider := ctx.Request.URL.Query().Get("oauth_provider")
	h.app.RedirectLoginOauth(ctx, &oatuhProvider)
}

// Callback is a function that handles the callback of the oauth
func (h *AuthHandler) Callback(ctx *gin.Context) {
	fmt.Println(ctx.Request.URL)

	var req application.CalbackSSOReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}

	result, err := h.app.GetUserData(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	res, err := h.app.LoginOrRegisterUserOauth(result)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, http.StatusOK, res)
}
