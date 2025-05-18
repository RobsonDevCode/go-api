package profileInterfaces

import (
	"context"

	domain "github.com/RobsonDevCode/go-profile-service/src/internal/domain/models"
	"github.com/google/uuid"
)

type ProfileRetrievalRepository interface {
	GetById(id uuid.UUID, ctx context.Context) (*domain.Profile, error)
	ProfileExits(id uuid.UUID, ctx context.Context) (bool, error)
}
