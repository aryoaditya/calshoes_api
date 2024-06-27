package config

import (
	"calshoes_api/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSsl := os.Getenv("DB_SSL")

	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbHost, dbUser, dbPassword, dbName, dbPort, dbSsl)

	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})

	if err != nil {
		panic("connection failed to the database ")
	}
	DB = db
	fmt.Println("db conncted successfully")

	AutoMigrate(db)
}

func AutoMigrate(connection *gorm.DB) {
	connection.Debug().AutoMigrate(
		&models.CartItem{},
		&models.Cart{},
		&models.Category{},
		&models.Customer{},
		&models.OrderItem{},
		&models.Order{},
		&models.Payment{},
		&models.Product{})
}
