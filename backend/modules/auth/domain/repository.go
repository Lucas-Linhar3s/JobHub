package domain

import (
	"errors"

	"github.com/Lucas-Linhar3s/JobHub/database"
	"github.com/Lucas-Linhar3s/JobHub/modules/auth/infrastructure"
	"github.com/Lucas-Linhar3s/JobHub/modules/auth/infrastructure/postgres"
	"github.com/Lucas-Linhar3s/JobHub/utils"
)

type repository struct {
	repo postgres.PGAuth
}

func newRepository(db *database.Database) *repository {
	return &repository{
		repo: postgres.PGAuth{
			Db: db,
		},
	}
}

// ConvertDomainToModelInfra implements IAuth.
func (r *repository) ConvertDomainToModelInfra(modelReq *AuthModel) (modelRes *infrastructure.AuthModel, err error) {
	if modelReq == nil {
		return nil, errors.New("modelReq is required")
	}

	if modelRes, err = utils.ConvertRequestToModel[infrastructure.AuthModel](modelReq); err != nil {
		return nil, err
	}

	return
}

// ConvertModelInfraToDomain implements IAuth.
func (r *repository) ConvertModelInfraToDomain(modelReq *infrastructure.AuthModel) (modelRes *AuthModel, err error) {
	if modelReq == nil {
		return nil, errors.New("modelReq is required")
	}

	if modelRes, err = utils.ConvertRequestToModel[AuthModel](modelReq); err != nil {
		return nil, err
	}

	return
}

// LoginWithEmailAndPassword implements IAuth.
func (r *repository) LoginWithEmailAndPassword(model *infrastructure.AuthModel) error {
	return r.repo.LoginWithEmailAndPassword(model)
}

// LoginWithOauth implements IAuth.
func (r *repository) LoginWithOauth(model *infrastructure.AuthModel) error {
	panic("unimplemented")
}

// RegisterUser implements IAuth.
func (r *repository) RegisterUser(model *infrastructure.AuthModel) error {
	return r.repo.RegisterUser(model)
}

// UpdateUser implements IAuth.
func (r *repository) UpdateUser(model *infrastructure.AuthModel) error {
	return r.repo.UpdateUser(model)
}

// VerifyEmail implements IAuth.
func (r *repository) VerifyEmail(email string) (bool, error) {
	return r.repo.VerifyEmail(email)
}

// VerifyRole implements IAuth.
func (r *repository) VerifyRole(model *infrastructure.AuthModel) error {
	return r.repo.VerifyRole(model)
}
