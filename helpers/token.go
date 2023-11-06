package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Token(c echo.Context) error {

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
	return c.JSON(http.StatusOK, map[string]string{"token": "tokenString"})
}
