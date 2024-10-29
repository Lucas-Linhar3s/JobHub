package jwt

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Lucas-Linhar3s/JobHub/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWT represents the jwt configuration of the application.
type JWT struct {
	key []byte
}

// MyCustomClaims represents the custom claims of the jwt.
type MyCustomClaims struct {
	AuthProvider string
	UserId       string
	Email        string
	Role         string
	jwt.RegisteredClaims
}

// NewJwt returns a new JWT instance.
func NewJwt(config *config.Config) *JWT {
	key, err := json.Marshal(config.Security.Jwt.Key)
	if err != nil {
		panic(err)
	}
	return &JWT{key: key}
}

// GenToken generates a jwt token.
func (j *JWT) GenToken(oauthProvider, userId, email, role string, expiresAt time.Time) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		AuthProvider: oauthProvider,
		UserId:       userId,
		Email:        email,
		Role:         role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "github.com/Lucas-Linhar3s",
			Subject:   userId,
			ID:        uuid.New().String(),
		},
	})

	// Sign and get the complete encoded token as a string using the key
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (j *JWT) ParseToken(tokenString string) (*MyCustomClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if strings.TrimSpace(tokenString) == "" {
		return nil, errors.New("token is empty")
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
