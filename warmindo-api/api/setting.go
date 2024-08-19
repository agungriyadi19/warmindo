package api

import (
	"database/sql"
	"warmindo-api/db"

	"warmindo-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupSettingsRoutes sets up the routes for settings with authentication middleware
func SetupSettingsRoutes(app *fiber.App, dbConn *sql.DB) {
	settingsAPI := app.Group("/api/settings", middleware.AuthMiddleware(1))
	settingsAPI.Get("/", func(c *fiber.Ctx) error {
		return GetSettings(c, dbConn)
	})
	settingsAPI.Put("/", func(c *fiber.Ctx) error {
		return CreateOrUpdateSettings(c, dbConn) // Updated to use PUT method for editing
	})
}

// CreateOrUpdateSettings handles updating settings with a fixed id of 1
func CreateOrUpdateSettings(c *fiber.Ctx, dbConn *sql.DB) error {
	var settings db.Settings
	if err := c.BodyParser(&settings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Update settings with id = 1
	query := `
		UPDATE settings
		SET total_table = $1,
			latitude = $2,
			longitude = $3,
			radius = $4
		WHERE id = 1;
	`

	res, err := dbConn.Exec(query, settings.TotalTable, settings.Latitude, settings.Longitude, settings.Radius)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if any row was updated
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Settings with id 1 not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

// GetSettings handles retrieving settings
func GetSettings(c *fiber.Ctx, dbConn *sql.DB) error {
	var settings db.Settings

	err := dbConn.QueryRow("SELECT total_table, latitude, longitude, radius FROM settings WHERE id = 1").Scan(&settings.TotalTable, &settings.Latitude, &settings.Longitude, &settings.Radius)
	if err != nil {
		if err == sql.ErrNoRows {
			// Optionally handle the case where no settings are found
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Settings not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "settings": settings})
}
