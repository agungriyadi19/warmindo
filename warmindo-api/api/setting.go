package api

import (
	"database/sql"
	"warmindo-api/db"

	"github.com/gofiber/fiber/v2"
)

// SetupSettingsRoutes sets up the routes for settings
func SetupSettingsRoutes(app *fiber.App, dbConn *sql.DB) {
	settingsAPI := app.Group("/api/settings")
	settingsAPI.Get("/", func(c *fiber.Ctx) error {
		return GetSettings(c, dbConn)
	})
	settingsAPI.Post("/", func(c *fiber.Ctx) error {
		return CreateOrUpdateSettings(c, dbConn)
	})
}

// Handlers for Settings

func CreateOrUpdateSettings(c *fiber.Ctx, dbConn *sql.DB) error {
	var settings db.Settings
	if err := c.BodyParser(&settings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Assuming a single row for settings, use UPSERT
	query := `
		INSERT INTO settings (title, subtitle, seats, latitude, longitude, radius)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
		title = EXCLUDED.title,
		subtitle = EXCLUDED.subtitle,
		seats = EXCLUDED.seats,
		latitude = EXCLUDED.latitude,
		longitude = EXCLUDED.longitude,
		radius = EXCLUDED.radius;
	`

	_, err := dbConn.Exec(query, settings.Title, settings.Subtitle, settings.Seats, settings.Latitude, settings.Longitude, settings.Radius)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func GetSettings(c *fiber.Ctx, dbConn *sql.DB) error {
	var settings db.Settings

	err := dbConn.QueryRow("SELECT title, subtitle, seats, latitude, longitude, radius FROM settings WHERE id = 1").Scan(&settings.Title, &settings.Subtitle, &settings.Seats, &settings.Latitude, &settings.Longitude, &settings.Radius)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "settings": settings})
}
