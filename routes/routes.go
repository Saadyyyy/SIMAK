package routes

import (
	"os"
	"simak/controller"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// e.Use(Logger())
	godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	secretKey := []byte(os.Getenv("SECRET_JWT"))

	e.POST("/signup", controller.Signup(db, secretKey))              // Register
	e.POST("/signin", controller.SignIn(db, secretKey))              //Login
	e.POST("/create-mk", controller.CreateMataKuliah(db, secretKey)) // create matakuliah
}
