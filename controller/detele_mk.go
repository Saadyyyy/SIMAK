package controller

import (
	"net/http"
	"simak/middleware"
	"simak/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DeleteMk(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		//mendapatkan id dari matakuliah
		deleteMkId := c.Param("id")

		var mk models.MataKuliah
		if db.First(&mk, deleteMkId); err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"Error": "Id Not found",
			})
		}

		if db.Delete(&mk); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"Error":   true,
				"Massage": "Cannot delete mata kuliah",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"Error":   true,
			"Massage": "Successfully to delete  matakuliah",
			"Id":      deleteMkId,
		})
	}
}
