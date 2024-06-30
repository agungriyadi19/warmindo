package api

import (
	"database/sql"
	"warmindo-api/db"

	"github.com/gofiber/fiber/v2"
)

func SetupRoleRoutes(app *fiber.App, dbConn *sql.DB) {
	roleAPI := app.Group("/api/roles")
	roleAPI.Get("/", func(c *fiber.Ctx) error {
		return GetRoles(c, dbConn)
	})
	roleAPI.Get("/:id", func(c *fiber.Ctx) error {
		return GetRoleByID(c, dbConn)
	})
	roleAPI.Post("/", func(c *fiber.Ctx) error {
		return CreateRole(c, dbConn)
	})
	roleAPI.Put("/:id", func(c *fiber.Ctx) error {
		return UpdateRole(c, dbConn)
	})
	roleAPI.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteRole(c, dbConn)
	})
}

// Handlers untuk Role

func CreateRole(c *fiber.Ctx, dbConn *sql.DB) error {
	var role db.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("INSERT INTO roles (name) VALUES ($1)", role.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func GetRoles(c *fiber.Ctx, dbConn *sql.DB) error {
	rows, err := dbConn.Query("SELECT id, name, created_at, updated_at FROM roles")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var roles []db.Role
	for rows.Next() {
		var role db.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		roles = append(roles, role)
	}

	return c.JSON(fiber.Map{"success": true, "roles": roles})
}

func GetRoleByID(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var role db.Role

	err := dbConn.QueryRow("SELECT id, name, created_at, updated_at FROM roles WHERE id = $1", id).Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "role": role})
}

func UpdateRole(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var role db.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("UPDATE roles SET name = $1, updated_at = NOW() WHERE id = $2", role.Name, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func DeleteRole(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")

	_, err := dbConn.Exec("DELETE FROM roles WHERE id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
