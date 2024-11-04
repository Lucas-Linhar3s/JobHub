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

// @Summary Register user
// @Description Register a new user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body application.UserRegisterReq true "User data to register"
// @Success 204 {string} string "No content"
// @Router /auth/ [post]
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

// @Summary Login with email and password
// @Description Login with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body application.UserRegisterReq true "User data to login"
// @Success 200 {object} application.SessionOut "User login response"
// @Router /auth/login [post]
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

// @Summary Redirect to login with oauth
// @Description Redirect to login with oauth
// @Tags auth
// @Accept  json
// @Produce  json
// @Param oauth_provider query string true "Oauth provider"
// @Router /auth/login [get]
func (h *AuthHandler) RedirectLoginOauth(ctx *gin.Context) {
	oatuhProvider := ctx.Request.URL.Query().Get("oauth_provider")
	h.app.RedirectLoginOauth(ctx, &oatuhProvider)
}

// @Summary Callback oauth
// @Description Callback oauth
// @Tags auth
// @Accept  json
// @Produce  json
// @Param code query string true "Code"
// @Param state query string true "State"
// @Success 200 {object} application.SessionOut "User login response"
// @Router /auth/login/callback [get]
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

	res, err := h.app.LoginOrRegisterUserOauth(ctx, result)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, http.StatusOK, res)
}
