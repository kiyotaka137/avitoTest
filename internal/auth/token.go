package auth

import (
	"errors"
	"os"
	"strings"
)

type Validator interface {
	Validate(token string) (Claims, error)
}

var globalValidator Validator

func SetValidator(v Validator) { globalValidator = v }

var ErrInvalidToken = errors.New("invalid token")

type DefaultEnvValidator struct{}

func (DefaultEnvValidator) Validate(token string) (Claims, error) {
	if token == "" {
		return Claims{}, ErrInvalidToken
	}
	admins := strings.Split(os.Getenv("ADMIN_TOKENS"), ",")
	for i := range admins {
		if strings.TrimSpace(admins[i]) == token {
			return Claims{Token: token, Role: "admin"}, nil
		}
	}
	return Claims{Token: token, Role: "user"}, nil
}

func ensureValidator() Validator {
	if globalValidator != nil {
		return globalValidator
	}
	return DefaultEnvValidator{}
}
