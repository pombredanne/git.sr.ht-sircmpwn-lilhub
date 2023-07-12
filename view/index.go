package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /
func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
