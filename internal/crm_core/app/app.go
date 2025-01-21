package app

import (
	"crm_system/config/crm_core"
	middleware2 "crm_system/internal/crm_core/controller/http/middleware"
	v1 "crm_system/internal/crm_core/controller/http/v1"
	debugRoute "crm_system/internal/crm_core/controller/http/v1/debug"
	_ "crm_system/internal/crm_core/docs"
	repoPkg "crm_system/internal/crm_core/repository"
	servicePkg "crm_system/internal/crm_core/service"
	"crm_system/internal/crm_core/transport"
	"crm_system/pkg/crm_core/cache"
	"crm_system/pkg/crm_core/httpserver/debug"
	"crm_system/pkg/crm_core/httpserver/public"
	"crm_system/pkg/crm_core/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(cfg *crm_core.Configuration) {
	l := logger.New(cfg.Gin.Mode)
	repo := repoPkg.New(cfg, l)

	//REDIS implementation
	redisClient, err := cache.NewRedisClient()
	fmt.Print("crm_system - Run - cache.NewRedisClient")
	if err != nil {
		fmt.Print("crm_system - Error - cache.NewRedisClient: %w", err)
		return
	}
	fmt.Print("crm_system - Run - cache.NewContactCache")

	contactCache := cache.NewContactCache(redisClient, 10*time.Minute)

	validateGrpcTransport := transport.NewValidateGrpcTransport(*cfg)

	service := servicePkg.New(cfg, repo, l)
	middleware := middleware2.New(repo, cfg, validateGrpcTransport)
	handler := gin.Default()
	handlerDebug := gin.Default()
	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8082"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	handler.GET("/rocket/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Print("crm_system - Run - http://localhost:8082/swagger/index.html")

	v1.NewRouter(handler, service, middleware, contactCache)
	debugRoute.NewDebugRouter(handlerDebug)
	httpServer := public.New(handler, cfg, public.Port(cfg.HTTP.Port))
	debugServer := debug.New(handlerDebug, cfg, debug.Port(cfg.HTTP.DebugPort))

	fmt.Print("crm_system - Run - http://localhost:8082")

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	fmt.Print("crm_system - Witing for signal - http://localhost:8082")

	select {
	case s := <-interrupt:
		l.Info("%s", "crm_system - Run - signal: "+s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("crm_system - Run - httpServer.Notify: %w", err))
	case err = <-debugServer.Notify():
		l.Error(fmt.Errorf("crm_system - Run - debugServer.Notify: %w", err))
	}

	print("crm_system - Run - Shutdown")

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("crm_system - Run - httpServer.Shutdown: %w", err))
	}
	err = debugServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("crm_system - Run - debugServer.Shutdown: %w", err))
	}
	err = redisClient.Close()
	fmt.Print("crm_system - Run - redisClient.Close : %w", err)
}
