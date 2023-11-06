package controller

import (
	"errors"
	"net/http"
	"simak/middleware"
	"simak/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// function untuk register
func Signup(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var person models.Person
		if err := c.Bind(&person); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// Mengecek apakah Nim apakah sudah ada
		var existingperson models.Person
		result := db.Where("nim = ?", person.Nim).First(&existingperson)
		if result.Error == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Nim already exists"})
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check Nim"})
		}

		var pass models.Person
		result = db.Where("password = ?", person.Password).First(&pass)
		if result.Error == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Password already exites"})
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to check Password",
			})
		}

		var pn models.Person
		result = db.Where("phone_number = ?", person.PhoneNumber).First(&pn)
		if result.Error == nil {
			return c.JSON(http.StatusConflict, map[string]string{
				"Error": "Phone number already exites",
			})
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"Erorr": "Failed to check phone number",
			})
		}

		// // Meng-hash password dengan bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(person.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}

		// Menyimpan data pengguna ke database
		person.Password = string(hashedPassword)
		db.Create(&person)

		// Hapus password dari struct
		person.Password = ""

		// Generate JWT token

		// Menyertakan ID pengguna dalam response
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "User created successfully"})
	}
}

func SignIn(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		//cek apakah ada model person
		var person models.Person
		if err := c.Bind(&person); err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{
				"Error": err.Error(),
			})
		}

		// cek apakah nim sudah ada di database
		var nim models.Person
		result := db.Where("nim = ?", person.Nim).First(&nim)
		if result.Error == nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid nim or password"})
			}
		}

		// cek apakah nim dan password benar
		err := bcrypt.CompareHashAndPassword([]byte(nim.Password), []byte(person.Password))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"Error": "password",
			})
		}

		// Generate JWT token
		tokenString, err := middleware.GenerateToken(nim.Nim, secretKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"massage": "Succesful login",
			"nim":     nim.Nim,
			"token":   tokenString,
		})
	}
}
