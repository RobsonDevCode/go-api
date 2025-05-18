package handlers

import (
	"context"
	"net/http"
	"time"

	validator "github.com/RobsonDevCode/go-profile-service/src/internal/api/handlers/middleware"
	domain "github.com/RobsonDevCode/go-profile-service/src/internal/domain/models"
	"github.com/google/uuid"

	"github.com/RobsonDevCode/go-profile-service/src/internal/config"
	"github.com/RobsonDevCode/go-profile-service/src/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FollowerHandler struct {
	profileRetrievalService *services.ProfileRetrievalService
	followRetrievalService  *services.FollowerRetrievalService
	logger                  *zap.Logger
}

func NewFollowerHandler(profileRetrievalService *services.ProfileRetrievalService, followRetrievalService *services.FollowerRetrievalService,
	logger *zap.Logger,
) *FollowerHandler {
	return &FollowerHandler{
		profileRetrievalService: profileRetrievalService,
		followRetrievalService:  followRetrievalService,
		logger:                  logger,
	}
}

func (h *FollowerHandler) Register(router *gin.RouterGroup,
	config *config.Config, logger *zap.Logger) {
	followerGroup := router.Group("")
	followerGroup.Use(validator.JWTAuthMiddleWare(config, logger))
	{
		followerGroup.GET(":id")
	}
}

func (h *FollowerHandler) GetPaged(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing profile id param"})
		return
	}
	uuid, err := uuid.Parse(id)
	if err != nil {
		h.logger.Sugar().Errorf("unable to parse id to uuid, %w", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"validatoin error": "Invalid id sent",
		})
	}

	paginationOptions := domain.GetOptions(c)

	ctx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	pagedResult, err := h.followRetrievalService.GetPage(uuid, paginationOptions, ctx)
	if err != nil {
		h.logger.Sugar().Errorf("error getting page: %w", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong please try again later!",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result": pagedResult,
	})

	h.logger.Sugar().Infof("Get Page for %s succesfully completed", uuid)
}
