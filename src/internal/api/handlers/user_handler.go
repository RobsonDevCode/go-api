package handlers

import (
	"context"
	"net/http"
	"time"

	validator "github.com/RobsonDevCode/go-profile-service/src/internal/api/handlers/middleware"
	client "github.com/RobsonDevCode/go-profile-service/src/internal/clients/user"
	"github.com/RobsonDevCode/go-profile-service/src/internal/domain"
	"github.com/RobsonDevCode/go-profile-service/src/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	readerService *services.ProfileRetrivelService
	writerService *services.ProfileWriterService
	userClient    *client.UserClient
	logger        *zap.Logger
}

func NewProfileHandler(readerService *services.ProfileRetrivelService,
	writerService *services.ProfileWriterService,
	userClient *client.UserClient,
	logger *zap.Logger) *ProfileHandler {
	return &ProfileHandler{
		readerService: readerService,
		writerService: writerService,
		userClient:    userClient,
		logger:        logger,
	}
}

func (h *ProfileHandler) Register(router *gin.RouterGroup) {
	profile := router.Group("")
	{
		profile.GET(":id", h.GetProfile)
		profile.POST("", h.CreateProfile)
	}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {

	h.logger.Info("Getting profile")

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing profile id paramter",
		})
		c.Abort()
		return
	}

	profileId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"validation error": "Invalid id sent",
		})
	}

	ctx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	profile, err := h.readerService.GetById(profileId, ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request Time Out"})
			return
		}
		if ctx.Err() == context.Canceled {
			c.AbortWithStatus(499)
			return
		} else {
			h.logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again later!"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": profile,
	})

	h.logger.Sugar().Infof("Profile: %v returned", profile)
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var profile domain.Profile
	if err := c.ShouldBindBodyWith(&profile, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validator.ValidateProfile(profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Attempting to create profile")

	ctx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	header := c.GetHeader("Authorization")
	jwt := header[7:]

	h.userClient.SetJwt(jwt)

	err := h.writerService.Create(profile, ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request Time Out"})
			return
		}
		if ctx.Err() == context.Canceled {
			c.AbortWithStatus(499)
			return
		} else {
			h.logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again later!"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "profile completed!"})

	h.logger.Sugar().Infof("Profil: %v created", profile)
}
