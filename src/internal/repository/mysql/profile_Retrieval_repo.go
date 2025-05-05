package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/RobsonDevCode/go-profile-service/src/internal/domain"
	profileInterfaces "github.com/RobsonDevCode/go-profile-service/src/internal/repository/interfaces"

	"github.com/google/uuid"
)

type ProfileRetrievalRepository struct {
	db *sql.DB
}

func NewProfileRetrievalRepository(db *sql.DB) profileInterfaces.ProfileRetrievalRepository {
	return &ProfileRetrievalRepository{
		db: db,
	}
}

func (r *ProfileRetrievalRepository) GetById(id uuid.UUID, ctx context.Context) (*domain.Profile, error) {

	rows := r.db.QueryRowContext(ctx, "SELECT * FROM profile WHERE userId = ? LIMIT 1", id)

	var profile domain.Profile
	if err := rows.Scan(&profile.UserId,
		&profile.FollowerCount, &profile.FollowingCount, &profile.Private); err != nil {
		return nil, fmt.Errorf("unable to scan row: %v", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %v", err)
	}

	return &profile, nil
}

func (r *ProfileRetrievalRepository) ProfileExits(id uuid.UUID, ctx context.Context) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT * FROM profile WHERE userId = ? LIMIT 1", id).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, which means profile doesn't exist
			return false, nil
		}

		return false, fmt.Errorf("checking profile existence: %v", err)
	}

	return true, nil
}
