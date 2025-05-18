package handlers

import (
	userClient "github.com/RobsonDevCode/go-profile-service/src/internal/clients/user"
	"github.com/RobsonDevCode/go-profile-service/src/internal/services"
	"go.uber.org/zap"
)

type FollowerHandler struct {
	profileRetrievalService *services.ProfileRetrievalService
	followRetrievalService  *services.FollowerRetrievalService
	logger                  *zap.Logger
}

func NewFollowerHandler(profileRetrievalService *services.ProfileRetrievalService, followRetrievalService *services.FollowerRetrievalService
logger *zap.Logger
) *FollowerHandler{
	return &FollowerHandler{
		profileRetrievalService: profileRetrievalService,
		followRetrievalService: followRetrievalService,
		logger: logger,
	}
}


