package api

import (
	"database/sql"
	"encoding/base64"
	"io/ioutil"
	"strconv"
	"warmindo-api/db"

	"github.com/gofiber/fiber/v2"
)

func SetupMenuRoutes(app *fiber.App, dbConn *sql.DB) {
	menuAPI := app.Group("/api/menus")
	menuAPI.Get("/", func(c *fiber.Ctx) error {
		return GetMenus(c, dbConn)
	})
	menuAPI.Get("/:id", func(c *fiber.Ctx) error {
		return GetMenuByID(c, dbConn)
	})
	menuAPI.Post("/", func(c *fiber.Ctx) error {
		return CreateMenu(c, dbConn)
	})
	menuAPI.Put("/:id", func(c *fiber.Ctx) error {
		return UpdateMenu(c, dbConn)
	})
	menuAPI.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteMenu(c, dbConn)
	})
}

// Handlers untuk Menu

func CreateMenu(c *fiber.Ctx, dbConn *sql.DB) error {
	var menu db.Menu
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Extract fields from the form
	menu.Name = form.Value["name"][0]
	menu.Description = form.Value["description"][0]
	menu.Price, _ = strconv.Atoi(form.Value["price"][0])
	menu.CategoryID, _ = strconv.Atoi(form.Value["category_id"][0])

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil {
		fileContent, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		defer fileContent.Close()

		imageData, err := ioutil.ReadAll(fileContent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// Encode image to base64
		menu.Image = base64.StdEncoding.EncodeToString(imageData)
	}

	_, err = dbConn.Exec("INSERT INTO menus (name, image, description, price, category_id) VALUES ($1, $2, $3, $4, $5)",
		menu.Name, menu.Image, menu.Description, menu.Price, menu.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func GetMenus(c *fiber.Ctx, dbConn *sql.DB) error {
	rows, err := dbConn.Query("SELECT id, name, image, description, price, category_id, created_at, updated_at FROM menus WHERE deleted = false")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var menus []db.Menu
	for rows.Next() {
		var menu db.Menu
		if err := rows.Scan(&menu.ID, &menu.Name, &menu.Image, &menu.Description, &menu.Price, &menu.CategoryID, &menu.CreatedAt, &menu.UpdatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		menus = append(menus, menu)
	}

	return c.JSON(fiber.Map{"success": true, "menus": menus})
}

func GetMenuByID(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var menu db.Menu

	err := dbConn.QueryRow("SELECT id, name, image, description, price, category_id, created_at, updated_at FROM menus WHERE id = $1", id).Scan(&menu.ID, &menu.Name, &menu.Image, &menu.Description, &menu.Price, &menu.CategoryID, &menu.CreatedAt, &menu.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "menu": menu})
}

func UpdateMenu(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var menu db.Menu

	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Extract fields from the form
	menu.Name = form.Value["name"][0]
	menu.Description = form.Value["description"][0]
	menu.Price, _ = strconv.Atoi(form.Value["price"][0])
	menu.CategoryID, _ = strconv.Atoi(form.Value["category_id"][0])

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil {
		fileContent, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		defer fileContent.Close()

		imageData, err := ioutil.ReadAll(fileContent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// Encode image to base64
		menu.Image = base64.StdEncoding.EncodeToString(imageData)
	}

	_, err = dbConn.Exec("UPDATE menus SET name = $1, image = $2, description = $3, price = $4, category_id = $5, updated_at = NOW() WHERE id = $6",
		menu.Name, menu.Image, menu.Description, menu.Price, menu.CategoryID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func DeleteMenu(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")

	_, err := dbConn.Exec("UPDATE menus SET deleted = true WHERE id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
