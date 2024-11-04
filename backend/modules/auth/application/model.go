package application

import "time"

// UserRegisterReq is a struct that contains the user register request
type UserRegisterReq struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required"`
}

// SessionOut is a struct that contains the session output
type SessionOut struct {
	UserID        *string    `json:"user_id"`
	AccessToken   *string    `json:"access_token"`
	DataExpiracao *time.Time `json:"data_expiracao"`
}

// CalbackSSOReq is a struct that contains the callback SSO request
type CalbackSSOReq struct {
	State         string `form:"state"`
	Code          string `form:"code"`
	OauthProvider string `form:"oauth_provider"`
}

// UserDataCallbackRes is a struct that contains the user data callback response
type UserDataCallbackRes struct {
	Email         *string `copier:"Email" validate:"required,email"`
	Picture       *string `copier:"Picture"`
	OauthProvider *string `copier:"OauthProvider"`
	OauthId       *string `copier:"OauthId"`
}
