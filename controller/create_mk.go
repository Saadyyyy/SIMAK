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
		if person.IsAdmin != person.IsAdmin {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"Error":   true,
				"massage": "Only admin can use this end point",
			})
		}

		result := db.Where("nim = ?", nim).First(&person)
		if result.Error != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"massage": "Failed to fetch nim",
			})
		}

		// newMatakuliah := models.MataKuliah{}

		var mataKuliah struct {
			Mk       string `json:"mata_kuliah"`
			Semester int16  `json:"semester"`
			Dosen    string `json:"dosen"`
		}

		if err := c.Bind(&mataKuliah); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"Eror": err.Error(),
			})
		}

		newMatakuliah := models.MataKuliah{
			Mk:       mataKuliah.Mk,
			Semester: mataKuliah.Semester,
			Dosen:    mataKuliah.Dosen,
		}

		if err := db.Create(&newMatakuliah).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"Error": "Failed create matakuliah",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"Error":   false,
			"massage": "successfully create mata kuliah",
			"Data":    newMatakuliah,
		})
	}
}
