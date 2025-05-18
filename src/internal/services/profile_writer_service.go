package services

import (
	"context"
	"fmt"
	"sync"

	client "github.com/RobsonDevCode/go-profile-service/src/internal/clients/user"
	domain "github.com/RobsonDevCode/go-profile-service/src/internal/domain/models"
	profileInterfaces "github.com/RobsonDevCode/go-profile-service/src/internal/repository/interfaces"
	"go.uber.org/zap"
)

type ProfileWriterService struct {
	profileWriterRepo profileInterfaces.ProfileWriterRepository
	reader            ProfileRetrievalService
	userClient        *client.UserClient
	logger            *zap.Logger
}

func NewProfileWriterService(repo profileInterfaces.ProfileWriterRepository,
	reader ProfileRetrievalService,
	userClient *client.UserClient,
	logger zap.Logger) *ProfileWriterService {
	return &ProfileWriterService{
		profileWriterRepo: repo,
		reader:            reader,
		userClient:        userClient,
		logger:            &logger,
	}
}

func (s *ProfileWriterService) Create(profile domain.Profile, ctx context.Context) error {

	exists, err := s.reader.ProfileExists(profile.UserId, ctx)
	if err != nil {
		s.logger.Sugar().Errorf("error checking profile %w", err)
		return err
	}
	if exists {
		s.logger.Sugar().Errorf("profile, %s already exits", profile.UserId)
		return fmt.Errorf("profile, %s already exits", profile.UserId)
	}

	profileCh := make(chan domain.Result, 1)
	userCh := make(chan domain.Result, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		exists, err := s.reader.ProfileExists(profile.UserId, ctx)
		profileCh <- domain.Result{Exists: exists,
			Err: err}
	}()

	go func() {
		defer wg.Done()
		exists, err := s.userClient.UserExists(profile.UserId, ctx)
		userCh <- domain.Result{Exists: exists,
			Err: err}
	}()

	go func() {
		wg.Wait()
		close(profileCh)
		close(userCh)
	}()

	profileResult := <-profileCh
	userResult := <-userCh

	if profileResult.Err != nil {
		s.logger.Sugar().Errorf("failed to check profile , %w", profileResult.Err)
		return profileResult.Err
	}

	if userResult.Err != nil {
		s.logger.Sugar().Errorf("failed to check user exists, %w", userResult.Err)
		return userResult.Err
	}

	if profileResult.Exists {
		s.logger.Sugar().Errorf("profile, %s already exists", profile.UserId)
		return fmt.Errorf("profile, %s already exists", profile.UserId)
	}

	if userResult.Exists {
		s.logger.Sugar().Errorf("user, %s already exists", profile.UserId)
		return fmt.Errorf("user, %s already exists", profile.UserId)
	}

	if err := s.profileWriterRepo.Create(profile, ctx); err != nil {
		return err
	}

	return nil
}
