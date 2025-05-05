package domain

import "github.com/google/uuid"

type Profile struct {
	UserId         uuid.UUID `json:"user_id" validate:"required"`
	FollowerCount  int32     `json:"follow_count"`
	FollowingCount int32     `json:"following_count"`
	Private        bool      `json:"private" validate:"required"`
}
