package mysql

import (
	"context"
	"database/sql"
	"fmt"

	domain "github.com/RobsonDevCode/go-profile-service/src/internal/domain/models"
	followInterface "github.com/RobsonDevCode/go-profile-service/src/internal/repository/interfaces/follow"
	"github.com/google/uuid"
)

type FollowerRetrivalRepository struct {
	db *sql.DB
}

func NewFollowerRetrivalRepository(db *sql.DB) followInterface.FollowerRetrivalRepository {
	return &FollowerRetrivalRepository{
		db: db,
	}
}

func (r *FollowerRetrivalRepository) GetPage(id uuid.UUID, pageinationOptions domain.PageinationOptions, ctx context.Context) (*domain.PagedResult[[]domain.User], error) {

	offset := (pageinationOptions.Page - 1) * pageinationOptions.Size

	query := `SELECT * FROM follower WHERE userGuid = ? 
			  LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, id, pageinationOptions.Size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var user domain.User

		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	total, err := r.GetCount(id, ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting count: %w", err)
	}

	return &domain.PagedResult[[]domain.User]{
		Items: users,
		Page:  pageinationOptions.Page,
		Size:  pageinationOptions.Size,
		Total: total,
	}, nil
}

func (r *FollowerRetrivalRepository) GetCount(id uuid.UUID, ctx context.Context) (int, error) {

	query := `SELECT COUNT(userGuid) WHERE userGuid = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
