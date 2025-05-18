package mysql

import (
	"context"
	"database/sql"

	domain "github.com/RobsonDevCode/go-profile-service/src/internal/domain/models"
	profileInterfaces "github.com/RobsonDevCode/go-profile-service/src/internal/repository/interfaces"
	"go.uber.org/zap"
)

type ProfileWriterRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewWriterRetrievalRepository(db *sql.DB, logger *zap.Logger) profileInterfaces.ProfileWriterRepository {
	return &ProfileWriterRepository{
		db:     db,
		logger: logger,
	}
}

func (s *ProfileWriterRepository) Create(profile domain.Profile, ctx context.Context) error {
	query := `INSERT INTO profile(userId, followerCount, followingCount, private)
			  VALUES(?, ?, ?, ?)`

	s.logger.Sugar().Infof("writing to sql %s ", profile.UserId)

	result, err := s.db.Exec(query, profile.UserId, profile.FollowerCount, profile.FollowingCount, profile.Private)
	if err != nil {
		return err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil
	}

	s.logger.Sugar().Infof("Profile %v created", lastId)
	return nil
}
