package api

import (
	"database/sql"
	"warmindo-api/db"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, dbConn *sql.DB) {
	userAPI := app.Group("/api/users")
	userAPI.Get("/", func(c *fiber.Ctx) error {
		return GetUsers(c, dbConn) // Replace dbConn with your db connection
	})
	userAPI.Get("/:id", func(c *fiber.Ctx) error {
		return GetUserByID(c, dbConn) // Replace dbConn with your db connection
	})
	userAPI.Post("/", func(c *fiber.Ctx) error {
		return CreateUser(c, dbConn) // Replace dbConn with your db connection
	})
	userAPI.Put("/:id", func(c *fiber.Ctx) error {
		return UpdateUser(c, dbConn) // Replace dbConn with your db connection
	})
	userAPI.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteUser(c, dbConn) // Replace dbConn with your db connection
	})
}

// Handlers untuk User

func CreateUser(c *fiber.Ctx, dbConn *sql.DB) error {
	var user db.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := user.HashPassword(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("INSERT INTO users (email, password, name, username, role_id, phone) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Email, user.Password, user.Name, user.Username, user.RoleID, user.Phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func GetUsers(c *fiber.Ctx, dbConn *sql.DB) error {
	rows, err := dbConn.Query("SELECT id, email, name, username, role_id, phone, created_at, updated_at FROM users")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var users []db.User
	for rows.Next() {
		var user db.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Username, &user.RoleID, &user.Phone, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		users = append(users, user)
	}

	return c.JSON(fiber.Map{"success": true, "users": users})
}

func GetUserByID(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var user db.User

	err := dbConn.QueryRow("SELECT id, email, name, username, role_id, phone, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.Name, &user.Username, &user.RoleID, &user.Phone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "user": user})
}

func UpdateUser(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var user db.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := user.HashPassword(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("UPDATE users SET email = $1, password = $2, name = $3, username = $4, role_id = $5, phone = $6, updated_at = NOW() WHERE id = $7",
		user.Email, user.Password, user.Name, user.Username, user.RoleID, user.Phone, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func DeleteUser(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")

	_, err := dbConn.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
