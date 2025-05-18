package middleware

import (
	domain "github.com/RobsonDevCode/go-profile-service/src/internal/domain/models"
	"github.com/go-playground/validator/v10"
)

func ValidateProfile(profile domain.Profile) error {
	validate := validator.New()
	return validate.Struct(profile)
}
