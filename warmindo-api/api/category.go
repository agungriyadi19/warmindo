package api

import (
	"database/sql"
	"warmindo-api/db"

	"github.com/gofiber/fiber/v2"
)

func SetupCategoryRoutes(app *fiber.App, dbConn *sql.DB) {
	categoryAPI := app.Group("/api/categories")
	categoryAPI.Get("/", func(c *fiber.Ctx) error {
		return GetCategories(c, dbConn)
	})
}

func GetCategories(c *fiber.Ctx, dbConn *sql.DB) error {
	rows, err := dbConn.Query("SELECT id, name, created_at, updated_at FROM categories")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var categories []db.Category
	for rows.Next() {
		var category db.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		categories = append(categories, category)
	}

	return c.JSON(fiber.Map{"success": true, "categories": categories})
}
