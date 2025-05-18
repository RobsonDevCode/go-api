package followInterface

import (
	"context"

	domain "github.com/RobsonDevCode/go-profile-service/src/internal/domain/models"
	"github.com/google/uuid"
)

type FollowerRetrivalRepository interface {
	GetPage(id uuid.UUID, pageinationOptions domain.PageinationOptions, ctx context.Context) (*domain.PagedResult[[]domain.User], error)
}
