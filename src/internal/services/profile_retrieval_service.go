package services

import (
	"context"
	"fmt"
	"time"

	"github.com/RobsonDevCode/go-profile-service/src/internal/caching"
	domain "github.com/RobsonDevCode/go-profile-service/src/internal/domain/models"
	profileInterfaces "github.com/RobsonDevCode/go-profile-service/src/internal/repository/interfaces"
	"github.com/google/uuid"
)

type ProfileRetrievalService struct {
	profileRetrievalRepo profileInterfaces.ProfileRetrievalRepository
	cache                *caching.Cache
}

func NewProfileRetrievalService(repo profileInterfaces.ProfileRetrievalRepository, cache *caching.Cache) *ProfileRetrievalService {
	return &ProfileRetrievalService{
		profileRetrievalRepo: repo,
		cache:                cache,
	}
}

func (s *ProfileRetrievalService) GetById(id uuid.UUID, ctx context.Context) (domain.Profile, error) {
	key := fmt.Sprintf("profile-%s", id)

	result, err := s.cache.GetOrCreate(key, time.Minute*3, func() (interface{}, error) {
		if id == uuid.Nil {
			return domain.Profile{}, fmt.Errorf("argument error, id can't be null")
		}

		profile, err := s.profileRetrievalRepo.GetById(id, ctx)
		if err != nil {
			return domain.Profile{}, err
		}

		return *profile, nil
	})

	if err != nil {
		return domain.Profile{}, err
	}

	profile, ok := result.(domain.Profile)
	if !ok {
		return domain.Profile{}, fmt.Errorf("unexcpected response type")
	}

	return profile, nil
}

func (s *ProfileRetrievalService) ProfileExists(id uuid.UUID, ctx context.Context) (bool, error) {
	key := fmt.Sprintf("exists-%s", id)
	result, err := s.cache.GetOrCreate(key, time.Minute*5, func() (interface{}, error) {

		if id == uuid.Nil {
			return false, fmt.Errorf("argument error, user id can't be null")
		}

		exists, err := s.profileRetrievalRepo.ProfileExits(id, ctx)
		if err != nil {
			return false, err
		}

		return exists, nil
	})
	if err != nil {
		return false, nil
	}

	exists, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("unexpected response type")
	}

	return exists, nil
}
