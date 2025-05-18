package profileInterfaces

import (
	"context"

	domain "github.com/RobsonDevCode/go-profile-service/src/internal/domain/models"
)

type ProfileWriterRepository interface {
	Create(profile domain.Profile, ctx context.Context) error
}
