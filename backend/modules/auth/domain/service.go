package domain

import (
	"github.com/Lucas-Linhar3s/JobHub/database"
	"github.com/Lucas-Linhar3s/JobHub/modules/auth/infrastructure"
)

var (
	TokenRoleOauth = "Oauth2"
)

// Service is a struct that represent the service of auth
type Service struct {
	iAuth IAuth
}

// GetService is a function that returns a new Service struct
func GetService(repo IAuth) *Service {
	return &Service{
		iAuth: repo,
	}
}

// GetRepository is a function that returns a new repository struct
func GetRepository(db *database.Database) IAuth {
	return newRepository(db)
}

// RegisterUser is a function that registers a new user
func (s *Service) RegisterUser(model *AuthModel) error {
	var data *infrastructure.AuthModel
	var err error

	if data, err = s.iAuth.ConvertDomainToModelInfra(model); err != nil {
		return err
	}

	return s.iAuth.RegisterUser(data)
}

// UpdateUser is a function that updates a user
func (s *Service) UpdateUser(model *AuthModel) error {
	var data *infrastructure.AuthModel
	var err error

	if data, err = s.iAuth.ConvertDomainToModelInfra(model); err != nil {
		return err
	}

	return s.iAuth.UpdateUser(data)
}

// VerifyEmail is a function that verifies if an email exists
func (s *Service) VerifyEmail(email string) (bool, error) {
	return s.iAuth.VerifyEmail(email)
}

func (s *Service) VerifyRole(model *AuthModel) (*AuthModel, error) {
	var data *infrastructure.AuthModel
	var err error

	if data, err = s.iAuth.ConvertDomainToModelInfra(model); err != nil {
		return nil, err
	}

	if err = s.iAuth.VerifyRole(data); err != nil {
		return nil, err
	}

	return s.iAuth.ConvertModelInfraToDomain(data)
}

// LoginWithEmailAndPassword is a function that logs in with email and password
func (s *Service) LoginWithEmailAndPassword(model *AuthModel) (*AuthModel, error) {
	var data *infrastructure.AuthModel
	var err error

	if data, err = s.iAuth.ConvertDomainToModelInfra(model); err != nil {
		return nil, err
	}

	if err = s.iAuth.LoginWithEmailAndPassword(data); err != nil {
		return nil, err
	}

	return s.iAuth.ConvertModelInfraToDomain(data)
}
