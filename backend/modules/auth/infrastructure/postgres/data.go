package postgres

import (
	"database/sql"

	"github.com/Lucas-Linhar3s/JobHub/database"
	"github.com/Lucas-Linhar3s/JobHub/modules/auth/infrastructure"
	"github.com/Masterminds/squirrel"
)

// PGAuth is a struct that represent the postgres auth
type PGAuth struct {
	Db *database.Database
}

// RegisterUser is a function that registers a new user
func (pg *PGAuth) RegisterUser(model *infrastructure.AuthModel) error {
	if err := pg.Db.Builder.
		Insert("jobhub.users").
		Columns("email", "password_hash", "oauth_provider", "oauth_id", "picture_url").
		Values(model.Email, model.PasswordHash, model.OauthProvider, *model.OauthId, model.Picture).
		Suffix("RETURNING id").
		Scan(new(string)); err != nil {
		return err
	}

	return nil
}

func (pg *PGAuth) UpdateUser(model *infrastructure.AuthModel) error {
	if err := pg.Db.Builder.
		Update("jobhub.users").
		Set("oauth_provider", model.OauthProvider).
		Set("oauth_id", model.OauthId).
		Set("picture_url", model.Picture).
		Where(
			squirrel.Eq{
				"email": model.Email,
			},
		).
		Suffix("RETURNING id").
		Scan(new(string)); err != nil {
		return err
	}

	return nil
}

// VerifyEmail is a function that verifies if an email exists
func (pg *PGAuth) VerifyEmail(email string) (bool, error) {
	var exist bool
	if err := pg.Db.Builder.
		Select("COUNT(*) > 0").
		From("jobhub.users").
		Where("email = ?", email).
		Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}

// LoginWithEmailAndPassword is a function that logs in with email and password
func (pg *PGAuth) LoginWithEmailAndPassword(model *infrastructure.AuthModel) error {
	if err := pg.Db.Builder.
		Select("id", "password_hash").
		From("jobhub.users").
		Where(
			squirrel.Eq{
				"email": model.Email,
			},
		).
		Scan(
			&model.ID,
			&model.PasswordHash,
		); err != nil {
		return err
	}

	return nil
}

// VerifyRole is a function that verifies the role of a user
func (pg *PGAuth) VerifyRole(model *infrastructure.AuthModel) error {
	if err := pg.Db.Builder.
		Select("role").
		From("jobhub.company_users CPU").
		Join("jobhub.users US ON CPU.user_id = US.id").
		Where(squirrel.Eq{
			"US.oauth_id": model.OauthId,
			"US.email":    model.Email,
		}).
		Scan(&model.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	return nil
}
