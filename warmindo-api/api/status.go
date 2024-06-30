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
	statusAPI.Get("/:id", func(c *fiber.Ctx) error {
		return GetStatusByID(c, dbConn)
	})
	statusAPI.Post("/", func(c *fiber.Ctx) error {
		return CreateStatus(c, dbConn)
	})
	statusAPI.Put("/:id", func(c *fiber.Ctx) error {
		return UpdateStatus(c, dbConn)
	})
	statusAPI.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteStatus(c, dbConn)
	})
}

// Handlers untuk Status

func CreateStatus(c *fiber.Ctx, dbConn *sql.DB) error {
	var status db.Status
	if err := c.BodyParser(&status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("INSERT INTO statuses (name) VALUES ($1)", status.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
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

func GetStatusByID(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var status db.Status

	err := dbConn.QueryRow("SELECT id, name, created_at, updated_at FROM statuses WHERE id = $1", id).Scan(&status.ID, &status.Name, &status.CreatedAt, &status.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "status": status})
}

func UpdateStatus(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var status db.Status
	if err := c.BodyParser(&status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("UPDATE statuses SET name = $1, updated_at = NOW() WHERE id = $2", status.Name, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func DeleteStatus(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")

	_, err := dbConn.Exec("DELETE FROM statuses WHERE id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
