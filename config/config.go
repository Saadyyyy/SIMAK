package config

import (
	"log"
	"os"
	"simak/models"
	"simak/routes"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

func InitializeDatabase(config DatabaseConfig) (*gorm.DB, error) {
	// Konfigurasi koneksi database MySQL dengan GORM
	dsn := config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" + config.DBName + "?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate tabel pengguna
	db.AutoMigrate(&models.Person{})
	db.AutoMigrate(&models.MataKuliah{})

	return db, nil
}

func SetupRouter() *echo.Echo {
	// Inisialisasi database
	godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	dbConfig := DatabaseConfig{
		Host: os.Getenv("DBHOST"),
	}
	portStr := os.Getenv("DBPORT")
	dbConfig.Port, _ = strconv.Atoi(portStr)
	dbConfig.Username = os.Getenv("DBUSER")
	dbConfig.Password = os.Getenv("DBPASS")
	dbConfig.DBName = os.Getenv("DBNAME")

	db, err := InitializeDatabase(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Membuat instance router Echo
	router := echo.New()

	// Middleware untuk log
	router.Use(middleware.Logger())

	// Menggunakan routes yang telah dipisahkan
	routes.SetupRoutes(router, db)

	return router
}
