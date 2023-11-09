package controller

import (
	"net/http"
	"simak/middleware"
	"simak/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateMataKuliah(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token dari header Authorization
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Authorization token is missing"})
		}

		// Memeriksa apakah header Authorization mengandung token Bearer
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token format. Use 'Bearer [token]'"})
		}

		// Ekstrak token dari header
		tokenString = tokenString[7:]

		nim, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			return c.JSON(http.StatusBadGateway, map[string]interface{}{
				"massage": "Invalid token",
			})
		}
		var person models.Person

		result := db.Where("nim = ?", nim).First(&person)
		if result.Error != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"massage": "Failed to fetch nim",
			})
		}
		if !person.IsAdmin {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"Error":   true,
				"massage": "Only admin can use this end point",
			})
		}

		var respon struct {
			Code     string `json:"code"`
			Sks      int    `json:"jumlah_sks"`
			Mk       string `json:"mata_kuliah"`
			Semester int16  `json:"semester"`
			Dosen    string `json:"dosen"`
		}

		if err := c.Bind(&respon); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"Error":   true,
				"Massage": err.Error(),
			})
		}

		NewCreateMk := models.MataKuliah{
			Code:     respon.Code,
			Mk:       respon.Mk,
			Sks:      respon.Sks,
			Semester: respon.Semester,
			Dosen:    respon.Dosen,
		}

		// if err := c.Bind(&NewCreateMk); err != nil {
		// 	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
		// 		"Error":   true,
		// 		"Massage": "Invalid connect to database",
		// 	})
		// }

		if err := db.Create(&NewCreateMk).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"Error":   true,
				"Massage": "Invalid create database",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"Error":   false,
			"massage": "successfully create mata kuliah",
			"Data":    NewCreateMk,
		})
	}
}
