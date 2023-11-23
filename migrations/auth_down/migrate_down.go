package main

import (
	"crm_system/config/auth"
	entityRepo "crm_system/internal/auth/entity"
	repoPkg "crm_system/internal/auth/repository"
	"crm_system/pkg/auth/logger"
)

func main() {
	cfg := auth.NewConfig()
	l := logger.New(cfg.Gin.Mode)
	repo := repoPkg.New(cfg, l)
	repo.DB.Migrator().DropTable(&entityRepo.User{}, &entityRepo.Role{})
}
