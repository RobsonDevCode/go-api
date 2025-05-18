package main

import (
	"log"

	"github.com/RobsonDevCode/go-profile-service/src/internal/api/handlers"
	"github.com/RobsonDevCode/go-profile-service/src/internal/caching"
	userClient "github.com/RobsonDevCode/go-profile-service/src/internal/clients/user"
	"github.com/RobsonDevCode/go-profile-service/src/internal/config"
	"github.com/RobsonDevCode/go-profile-service/src/internal/repository/mysql"
	"github.com/RobsonDevCode/go-profile-service/src/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
		return
	}

	database := mysql.NewUserDataBase(*config)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
		return
	}

	cache := &caching.Cache{}
	userClient, err := userClient.NewUserClient(*config, cache)
	if err != nil {
		logger.Sugar().Panicf("start up error, %w", err)
		return
	}

	profileRetrievalRepo := mysql.NewProfileRetrievalRepository(database)
	profileWriterRepo := mysql.NewWriterRetrievalRepository(database, logger)
	followRetrievalRepo := mysql.NewFollowerRetrivalRepository(database, logger)

	profileRetrievalService := services.NewProfileRetrievalService(profileRetrievalRepo, cache)
	profileWriterService := services.NewProfileWriterService(profileWriterRepo, *profileRetrievalService, userClient, *logger)
	followRetrievalService := services.NewFollowerRetrivalService(followRetrievalRepo, *profileRetrievalService, *logger)

	profileHandler := handlers.NewProfileHandler(profileRetrievalService, profileWriterService, userClient, logger)
	followerHandler := handlers.NewFollowerHandler(profileRetrievalService, followRetrievalService, logger)

	router := Setup(profileHandler, followerHandler, config, logger)

	if err := router.Run(":8080"); err != nil {
		logger.Sugar().Errorf("Failed to start server: %v", err)
	}

}

func Setup(profileHandler *handlers.ProfileHandler,
	followerHandler *handlers.FollowerHandler,
	config *config.Config, logger *zap.Logger) *gin.Engine {
	router := gin.Default()
	api := router.Group("profile/v1")
	{
		profileHandler.Register(api, config, logger)
		followerHandler.Register(api, config, logger)
	}

	return router
}
