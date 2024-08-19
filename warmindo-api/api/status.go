package api

import (
	"database/sql"
	"warmindo-api/db"

	"github.com/gofiber/fiber/v2"
)

func SetupStatusRoutes(app *fiber.App, dbConn *sql.DB) {
	statusAPI := app.Group("/api/statuses")
	statusAPI.Get("/", func(c *fiber.Ctx) error {
		return GetStatuses(c, dbConn)
	})
}

func GetStatuses(c *fiber.Ctx, dbConn *sql.DB) error {
	rows, err := dbConn.Query("SELECT id, name, created_at, updated_at FROM statuses")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var statuses []db.Status
	for rows.Next() {
		var status db.Status
		if err := rows.Scan(&status.ID, &status.Name, &status.CreatedAt, &status.UpdatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		statuses = append(statuses, status)
	}

	return c.JSON(fiber.Map{"success": true, "statuses": statuses})
}
