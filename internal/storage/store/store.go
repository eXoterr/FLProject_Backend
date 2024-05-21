package store

import (
	"fmt"
	"log"
	"os"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/storage/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MustSetup(config config.PostgreSQL, doSync bool) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Moscow",
		config.Host, config.User, config.Password, config.Database, config.Port,
	)

	logg := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true,
		})

	driver, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logg,
	})

	if doSync {
		driver.AutoMigrate(
			&models.User{},
			&models.Client{},
			&models.Specialization{},
			&models.Worker{},
			&models.Category{},
			&models.Tag{},
			&models.Order{},
			&models.Token{},
		)
	}

	if err != nil {
		log.Fatalf("unable to setup storage: %s", err)
	}

	return driver
}
