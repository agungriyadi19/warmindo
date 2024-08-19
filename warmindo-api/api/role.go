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
