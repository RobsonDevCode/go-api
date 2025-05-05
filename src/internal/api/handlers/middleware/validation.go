package middleware

import (
	"github.com/RobsonDevCode/go-profile-service/src/internal/domain"
	"github.com/go-playground/validator/v10"
)

func ValidateProfile(profile domain.Profile) error {
	validate := validator.New()
	return validate.Struct(profile)
}
