package postgres

import (
	"crm_system/config/crm_core"
	"crm_system/pkg/crm_core/logger"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Hello, World!")

	connectionStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		"localhost",
		5432,
		"postgres",
		//"12345",
		"auth_service_crm",
	)
	fmt.Print(connectionStr)
	db, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(db)
}

func ConnectDB(config *crm_core.Configuration, l *logger.Logger) *gorm.DB {
	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
	)
	print(connectionStr)
	db, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		l.Fatal(err)
	}
	return db
}
