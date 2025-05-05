package profileInterfaces

import (
	"context"

	"github.com/RobsonDevCode/go-profile-service/src/internal/domain"
)

type ProfileWriterRepository interface {
	Create(profile domain.Profile, ctx context.Context) error
}
