package services

import (
	"context"
	"fmt"

	domain "github.com/RobsonDevCode/go-profile-service/src/internal/domain/models"
	followInterface "github.com/RobsonDevCode/go-profile-service/src/internal/repository/interfaces/follow"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FollowerRetrievalService struct {
	followerRetrievalRepo  followInterface.FollowerRetrivalRepository
	profileRetrivelService ProfileRetrievalService
	logger                 *zap.Logger
}

func NewFollowerRetrivalService(followerRepo followInterface.FollowerRetrivalRepository,
	profileService ProfileRetrievalService,
	logger zap.Logger) *FollowerRetrievalService {
	return &FollowerRetrievalService{
		followerRetrievalRepo:  followerRepo,
		profileRetrivelService: profileService,
		logger:                 &logger,
	}
}

func (s *FollowerRetrievalService) GetPage(id uuid.UUID, pageinationOptions domain.PageinationOptions, ctx context.Context) (domain.PagedResult[[]domain.User], error) {
	exists, err := s.profileRetrivelService.profileRetrievalRepo.ProfileExits(id, ctx)
	if err != nil {
		return domain.PagedResult[[]domain.User]{}, fmt.Errorf("error checking if profile exists: %w", err)
	}

	if !exists {
		return domain.PagedResult[[]domain.User]{}, fmt.Errorf("profile %s does not exist!", id)
	}

	result, err := s.followerRetrievalRepo.GetPage(id, pageinationOptions, ctx)
	if err != nil {
		return domain.PagedResult[[]domain.User]{}, fmt.Errorf("error getting page for %s, %w", id, err)
	}

	s.logger.Info("succesfully returned page")
	return *result, nil
}
