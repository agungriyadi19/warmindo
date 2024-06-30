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
	categoryAPI.Get("/:id", func(c *fiber.Ctx) error {
		return GetCategoryByID(c, dbConn)
	})
	categoryAPI.Post("/", func(c *fiber.Ctx) error {
		return CreateCategory(c, dbConn)
	})
	categoryAPI.Put("/:id", func(c *fiber.Ctx) error {
		return UpdateCategory(c, dbConn)
	})
	categoryAPI.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteCategory(c, dbConn)
	})
}

// Handlers untuk Category

func CreateCategory(c *fiber.Ctx, dbConn *sql.DB) error {
	var category db.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("INSERT INTO categories (name) VALUES ($1)", category.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
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

func GetCategoryByID(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var category db.Category

	err := dbConn.QueryRow("SELECT id, name, created_at, updated_at FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "category": category})
}

func UpdateCategory(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var category db.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("UPDATE categories SET name = $1, updated_at = NOW() WHERE id = $2", category.Name, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func DeleteCategory(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")

	_, err := dbConn.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
