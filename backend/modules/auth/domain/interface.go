package domain

import "github.com/Lucas-Linhar3s/JobHub/modules/auth/infrastructure"

// IAuth is an interface that defines the methods that must be implemented by the Auth domain
type IAuth interface {
	VerifyEmail(email string) (bool, error)
	RegisterUser(model *infrastructure.AuthModel) error
	UpdateUser(model *infrastructure.AuthModel) error
	LoginWithEmailAndPassword(model *infrastructure.AuthModel) error
	LoginWithOauth(model *infrastructure.AuthModel) error
	VerifyRole(model *infrastructure.AuthModel) error
	ConvertModelInfraToDomain(modelReq *infrastructure.AuthModel) (modelRes *AuthModel, err error)
	ConvertDomainToModelInfra(modelReq *AuthModel) (modelRes *infrastructure.AuthModel, err error)
}
