package api

import (
	"database/sql"
	"fmt"
	"warmindo-api/db"
	"warmindo-api/middleware"
	"warmindo-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	db.User
	jwt.StandardClaims
}

// SetupUserRoutes sets up the routes for user management
func SetupUserRoutes(app *fiber.App, dbConn *sql.DB) {
	userAPI := app.Group("/api/users", middleware.AuthMiddleware(1))
	userAPI.Get("/", func(c *fiber.Ctx) error {
		return GetUsers(c, dbConn)
	})
	userAPI.Get("/:id", func(c *fiber.Ctx) error {
		return GetUserByID(c, dbConn)
	})
	userAPI.Post("/", func(c *fiber.Ctx) error {
		return CreateUser(c, dbConn)
	})
	userAPI.Put("/:id", func(c *fiber.Ctx) error {
		return UpdateUser(c, dbConn)
	})
	userAPI.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteUser(c, dbConn)
	})

	authAPI := app.Group("/api/auth")
	authAPI.Post("/login", func(c *fiber.Ctx) error {
		return Login(c, dbConn)
	})
}

// CreateUser handles creating a new user
func CreateUser(c *fiber.Ctx, dbConn *sql.DB) error {
	var user db.User
	if err := c.BodyParser(&user); err != nil {
		fmt.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	validationErrors := utils.ValidateUser(user)
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": validationErrors})
	}

	hashedPassword, err := utils.GetHash(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengenkripsi kata sandi"})
	}
	user.Password = hashedPassword

	_, err = dbConn.Exec("INSERT INTO staffs (email, password, name, username, role_id, phone, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())",
		user.Email, user.Password, user.Name, user.Username, user.RoleID, user.Phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memasukkan pengguna ke dalam database"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pengguna berhasil dibuat"})
}

// GetUsers retrieves all users
func GetUsers(c *fiber.Ctx, dbConn *sql.DB) error {
	rows, err := dbConn.Query("SELECT id, email, name, username, role_id, phone, created_at, updated_at FROM staffs")
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

// GetUserByID retrieves a single user by ID
func GetUserByID(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var user db.User

	err := dbConn.QueryRow("SELECT id, email, name, username, role_id, phone, created_at, updated_at FROM staffs WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.Name, &user.Username, &user.RoleID, &user.Phone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"success": true, "user": user})
}

// UpdateUser handles updating user details
func UpdateUser(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var user db.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	if err := user.HashPassword(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengenkripsi kata sandi"})
	}

	_, err := dbConn.Exec("UPDATE staffs SET email = $1, password = $2, name = $3, username = $4, role_id = $5, phone = $6, updated_at = NOW() WHERE id = $7",
		user.Email, user.Password, user.Name, user.Username, user.RoleID, user.Phone, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui pengguna di dalam database"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pengguna berhasil diperbarui"})
}

// DeleteUser handles deleting a user by ID
func DeleteUser(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")

	_, err := dbConn.Exec("DELETE FROM staffs WHERE id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus pengguna dari database"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pengguna berhasil dihapus"})
}

// Login handles user login and token generation
func Login(c *fiber.Ctx, dbConn *sql.DB) error {
	loginUser := &db.Login{}

	if err := c.BodyParser(loginUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "errors": []string{"Permintaan tidak valid"}})
	}

	user, err := db.GetUserByEmail(dbConn, loginUser.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "errors": []string{"Email belum terdaftar"}})
		}
		fmt.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "errors": []string{"Kesalahan database"}})
	}

	if !utils.ComparePassword(user.Password, loginUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "errors": []string{"Kredensial salah"}})
	}

	tokenValue, err := middleware.GenerateJWT(fmt.Sprintf("%d", user.ID), user.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "errors": []string{"Kesalahan dalam menghasilkan token"}})
	}

	return c.JSON(fiber.Map{"success": true, "user": user, "token": tokenValue})
}
