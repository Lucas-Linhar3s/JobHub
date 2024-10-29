package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Lucas-Linhar3s/JobHub/database"
	"github.com/Lucas-Linhar3s/JobHub/modules/auth/domain"
	"github.com/Lucas-Linhar3s/JobHub/pkg/config"
	v1 "github.com/Lucas-Linhar3s/JobHub/pkg/http/response/v1"
	"github.com/Lucas-Linhar3s/JobHub/pkg/jwt"
	"github.com/Lucas-Linhar3s/JobHub/pkg/log"
	"github.com/Lucas-Linhar3s/JobHub/services"
	"github.com/Lucas-Linhar3s/JobHub/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthApp represents the application service for auth.
type AuthApp struct {
	logger   *log.Logger
	validate *validator.Validate
	db       *database.Database
	config   *config.Config
	jwt      *jwt.JWT
}

// NewAuthApp creates a new AuthApp.
func NewAuthApp(
	logger *log.Logger,
	db *database.Database,
	config *config.Config,
	jwt *jwt.JWT,
) *AuthApp {
	return &AuthApp{
		logger:   logger,
		validate: validator.New(),
		db:       db,
		config:   config,
		jwt:      jwt,
	}
}

// RegisterUser registers a new user.
func (app *AuthApp) RegisterUser(ctx *gin.Context, req *UserRegisterReq) error {
	const msg = "Error while registering user"

	tx, err := app.db.NewTransaction()
	if err != nil {
		app.logger.Error(msg, zap.Error(err))
		return err
	}
	defer tx.Rollback()

	var (
		service = domain.GetService(domain.GetRepository(tx))
		data    = new(domain.AuthModel)
	)

	if err = app.validate.Struct(req); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return err
	}

	if data, err = utils.ConvertRequestToModel[domain.AuthModel](req); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return err
	}

	if exist, err := service.VerifyEmail(*data.Email); err != nil {
		app.logger.Error(msg+"VerifyEmail", zap.Error(err))
		return err
	} else if exist {
		return v1.ErrEmailAlreadyUse
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		app.logger.Error(msg, zap.Error(err))
		return err
	}
	data.PasswordHash = utils.GetStringPointer(string(bytes))

	if err = service.RegisterUser(data); err != nil {
		app.logger.Error(msg+"RegisterUser", zap.Error(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		app.logger.Error(msg+"Commit", zap.Error(err))
		return err
	}

	return nil
}

// LoginWithEmailAndPassword logs in with email and password.
func (app *AuthApp) LoginWithEmailAndPassword(ctx *gin.Context, req *UserRegisterReq) (*SessionOut, error) {
	var msg = "Error while logging in with email and password"

	tx, err := app.db.NewTransaction()
	if err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()

	var (
		service = domain.GetService(domain.GetRepository(tx))
		data    = new(domain.AuthModel)
		res     = new(SessionOut)
		role    = "user"
	)

	if err = app.validate.Struct(req); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	}

	if data, err = utils.ConvertRequestToModel[domain.AuthModel](req); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	}

	if exist, err := service.VerifyEmail(*data.Email); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	} else if !exist {
		return nil, v1.ErrEmailNotExists
	}

	if data, err = service.LoginWithEmailAndPassword(data); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*data.PasswordHash), []byte(*req.Password)); err != nil {
		return nil, v1.ErrInvalidPassword
	}

	if data.Role != nil {
		role = *data.Role
	}

	now := time.Now()
	expiresAt := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+app.config.Security.Jwt.ExpiresAt, now.Second(), now.Nanosecond(), now.Location())

	res.UserID = data.ID
	res.DataExpiracao = &expiresAt

	if res.AccessToken, err = app.jwt.GenToken("", *data.ID, *data.Email, role, expiresAt); err != nil {
		return nil, err
	}

	return res, nil
}

// LoginOauth logs in with oauth.
func (app *AuthApp) RedirectLoginOauth(ctx *gin.Context, oauthProvider *string) {
	now := time.Now()
	expiresAt := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+5, now.Second(), now.Nanosecond(), now.Location())

	emailDomain := fmt.Sprintf("%s.com", *oauthProvider)

	token, err := app.jwt.GenToken(*oauthProvider, uuid.NewString(), emailDomain, "Oauth2", expiresAt)
	if err != nil {
		app.logger.Error("Error while generating token", zap.Error(err))
		return
	}

	var (
		sso = utils.GetOauth2Config(*oauthProvider, app.config)
		url = sso.AuthCodeURL(*token)
	)

	utils.RedirectRequest(ctx.Writer, ctx.Request, url)
}

// GetUserData gets the user data.
func (app *AuthApp) GetUserData(ctx *gin.Context, req CalbackSSOReq) (*UserDataCallbackRes, error) {
	claims, err := app.jwt.ParseToken(req.State)
	if err != nil {
		app.logger.Error("Error while parsing token", zap.Error(err))
		return nil, err
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		app.logger.Error("Token expired")
		return nil, err
	}

	if claims.Role != domain.TokenRoleOauth {
		app.logger.Error("Role not allowed")
		return nil, err
	}

	data, err := utils.GetOauth2Config(claims.AuthProvider, app.config).Exchange(context.Background(), req.Code)
	if err != nil {
		app.logger.Error("Error while exchanging code", zap.Error(err))
		return nil, err
	}

	response, err := services.GetUserDataRequest(&claims.AuthProvider, data, app.logger)
	if err != nil {
		app.logger.Error("Error while getting user data", zap.Error(err))
		return nil, err
	}
	res := new(UserDataCallbackRes)
	if res, err = utils.ConvertRequestToModel[UserDataCallbackRes](response); err != nil {
		app.logger.Error("Error while converting request to model", zap.Error(err))
	}
	res.OauthProvider = &claims.AuthProvider

	return res, nil
}

// LoginOrRegisterUserOauth logs in or registers a user with oauth.
func (app *AuthApp) LoginOrRegisterUserOauth(req *UserDataCallbackRes) (*SessionOut, error) {
	const msg = "Error while logging in or registering user with oauth"

	tx, err := app.db.NewTransaction()
	if err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	}
	// defer tx.Rollback()

	var (
		data    *domain.AuthModel
		service = domain.GetService(domain.GetRepository(tx))
		res     = new(SessionOut)
		role    = "user"
	)

	if err = app.validate.Struct(req); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	}

	if data, err = utils.ConvertRequestToModel[domain.AuthModel](req); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	}

	if exist, err := service.VerifyEmail(*data.Email); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	} else if !exist {
		if err = service.RegisterUser(data); err != nil {
			app.logger.Error(msg, zap.Error(err))
			return nil, err
		}
	} else if exist {
		if err = service.UpdateUser(data); err != nil {
			app.logger.Error(msg, zap.Error(err))
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	}

	if data, err = service.VerifyRole(data); err != nil {
		app.logger.Error(msg, zap.Error(err))
		return nil, err
	}

	if data.Role != nil {
		role = *data.Role
	}

	now := time.Now()
	expiresAt := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+app.config.Security.Jwt.ExpiresAt, now.Second(), now.Nanosecond(), now.Location())

	res.UserID = utils.GetStringPointer(*data.OauthId)
	res.DataExpiracao = &expiresAt

	if res.AccessToken, err = app.jwt.GenToken(*data.OauthProvider, *data.OauthId, *data.Email, role, expiresAt); err != nil {
		return nil, err
	}

	return res, nil
}
