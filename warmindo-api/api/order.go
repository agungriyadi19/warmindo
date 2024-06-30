package api

import (
	"database/sql"
	"warmindo-api/db"

	"github.com/gofiber/fiber/v2"
)

func SetupOrderRoutes(app *fiber.App, dbConn *sql.DB) {
	orderAPI := app.Group("/api/orders")
	orderAPI.Get("/", func(c *fiber.Ctx) error {
		return GetOrders(c, dbConn)
	})
	orderAPI.Get("/:id", func(c *fiber.Ctx) error {
		return GetOrderByID(c, dbConn)
	})
	orderAPI.Post("/", func(c *fiber.Ctx) error {
		return CreateOrder(c, dbConn)
	})
	orderAPI.Put("/:id", func(c *fiber.Ctx) error {
		return UpdateOrder(c, dbConn)
	})
	orderAPI.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteOrder(c, dbConn)
	})
}

// Handlers untuk Order

func CreateOrder(c *fiber.Ctx, dbConn *sql.DB) error {
	var order db.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("INSERT INTO orders (amount, table_number, status_id, order_date, menu_id) VALUES ($1, $2, $3, $4, $5)",
		order.Amount, order.TableNumber, order.StatusID, order.OrderDate, order.MenuID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func GetOrders(c *fiber.Ctx, dbConn *sql.DB) error {
	rows, err := dbConn.Query("SELECT id, amount, table_number, status_id, order_date, menu_id, created_at, updated_at FROM orders")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var orders []db.Order
	for rows.Next() {
		var order db.Order
		if err := rows.Scan(&order.ID, &order.Amount, &order.TableNumber, &order.StatusID, &order.OrderDate, &order.MenuID, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		orders = append(orders, order)
	}

	return c.JSON(fiber.Map{"success": true, "orders": orders})
}

func GetOrderByID(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var order db.Order

	err := dbConn.QueryRow("SELECT id, amount, table_number, status_id, order_date, menu_id, created_at, updated_at FROM orders WHERE id = $1", id).Scan(&order.ID, &order.Amount, &order.TableNumber, &order.StatusID, &order.OrderDate, &order.MenuID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "order": order})
}

func UpdateOrder(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var order db.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("UPDATE orders SET amount = $1, table_number = $2, status_id = $3, order_date = $4, menu_id = $5, updated_at = NOW() WHERE id = $6",
		order.Amount, order.TableNumber, order.StatusID, order.OrderDate, order.MenuID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func DeleteOrder(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")

	_, err := dbConn.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
