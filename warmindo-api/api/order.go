package api

import (
	"database/sql"
	"fmt"
	"time"
	"warmindo-api/db"

	"github.com/gofiber/fiber/v2"
)

func SetupOrderRoutes(app *fiber.App, dbConn *sql.DB) {
	orderAPI := app.Group("/api/orders")
	orderAPI.Get("/", func(c *fiber.Ctx) error {
		return GetOrders(c, dbConn)
	})
	orderAPI.Get("/:order_code", func(c *fiber.Ctx) error {
		return GetOrdersByCode(c, dbConn)
	})
	orderAPI.Post("/", func(c *fiber.Ctx) error {
		return CreateOrder(c, dbConn)
	})
	orderAPI.Put("/:id", func(c *fiber.Ctx) error {
		return UpdateOrder(c, dbConn)
	})
	orderAPI.Patch("/status", func(c *fiber.Ctx) error {
		return UpdateOrderStatus(c, dbConn)
	})
	orderAPI.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteOrder(c, dbConn)
	})
}

// Handlers untuk Order
type CreateOrderRequest struct {
	TableNumber string `json:"table_number" validate:"required"`
	OrderCode   string `json:"order_code" validate:"required"`
	MenuID      int    `json:"menu_id" validate:"required,number"`
	Amount      int    `json:"amount" validate:"required,number"`
}

func CreateOrder(c *fiber.Ctx, dbConn *sql.DB) error {
	data := new(CreateOrderRequest)

	// Parse the request body
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Check if the order already exists
	var existingAmount int
	err := dbConn.QueryRow("SELECT amount FROM orders WHERE order_code = $1 AND menu_id = $2 AND status_id = 1", data.OrderCode, data.MenuID).Scan(&existingAmount)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if existingAmount > 0 {
		// Update the existing order by accumulating the amount
		newAmount := existingAmount + data.Amount
		_, err := dbConn.Exec(
			"UPDATE orders SET amount = $1, updated_at = NOW() WHERE order_code = $2 AND menu_id = $3 AND status_id = 1",
			newAmount, data.OrderCode, data.MenuID,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":      "Order updated successfully",
			"total_amount": newAmount,
		})
	}

	// Create a new order
	_, err = dbConn.Exec(
		"INSERT INTO orders (table_number, order_code, menu_id, amount, status_id, order_date) VALUES ($1, $2, $3, $4, $5, $6)",
		data.TableNumber, data.OrderCode, data.MenuID, data.Amount, 1, time.Now(),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
	})
}

func GetOrders(c *fiber.Ctx, dbConn *sql.DB) error {
	query := `
	SELECT o.id, o.amount, o.table_number, o.status_id, o.order_date, o.menu_id, o.order_code, 
           o.created_at, o.updated_at,
           s.name as status_name,
           m.name as menu_name, m.description as menu_description, m.price as menu_price,
           c.name as category_name,
           (o.amount * m.price) as total_price
    FROM orders o
    JOIN statuses s ON o.status_id = s.id
    JOIN menus m ON o.menu_id = m.id
    JOIN categories c ON m.category_id = c.id
	ORDER BY o.order_date DESC`

	rows, err := dbConn.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var orders []fiber.Map
	for rows.Next() {
		var order db.Order
		var statusName, menuName, menuDescription, categoryName string
		var menuPrice, totalPrice int

		if err := rows.Scan(&order.ID, &order.Amount, &order.TableNumber, &order.StatusID, &order.OrderDate, &order.MenuID, &order.OrderCode,
			&order.CreatedAt, &order.UpdatedAt, &statusName, &menuName, &menuDescription, &menuPrice, &categoryName, &totalPrice); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		orderMap := fiber.Map{
			"id":           order.ID,
			"amount":       order.Amount,
			"table_number": order.TableNumber,
			"status_id":    order.StatusID,
			"order_date":   order.OrderDate,
			"menu_id":      order.MenuID,
			"order_code":   order.OrderCode,
			"created_at":   order.CreatedAt,
			"updated_at":   order.UpdatedAt,
			"status_name":  statusName,
			"menu": fiber.Map{
				"name":          menuName,
				"description":   menuDescription,
				"price":         menuPrice,
				"category_name": categoryName,
			},

			"total_price": totalPrice,
		}

		orders = append(orders, orderMap)
	}

	return c.JSON(fiber.Map{"success": true, "orders": orders})
}

func GetOrdersByCode(c *fiber.Ctx, dbConn *sql.DB) error {
	orderCode := c.Params("order_code")

	query := `
	SELECT o.id, o.amount, o.table_number, o.status_id, o.order_date, o.menu_id, o.order_code, 
           o.created_at, o.updated_at,
           s.name as status_name,
           m.name as menu_name, m.description as menu_description, m.price as menu_price,
           c.name as category_name,
           (o.amount * m.price) as total_price
    FROM orders o
    JOIN statuses s ON o.status_id = s.id
    JOIN menus m ON o.menu_id = m.id
    JOIN categories c ON m.category_id = c.id
    WHERE o.order_code = $1
	ORDER BY o.order_date DESC`

	rows, err := dbConn.Query(query, orderCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var orders []fiber.Map
	for rows.Next() {
		var order db.Order
		var statusName, menuName, menuDescription, categoryName string
		var menuPrice, totalPrice int

		if err := rows.Scan(&order.ID, &order.Amount, &order.TableNumber, &order.StatusID, &order.OrderDate, &order.MenuID, &order.OrderCode,
			&order.CreatedAt, &order.UpdatedAt, &statusName, &menuName, &menuDescription, &menuPrice, &categoryName, &totalPrice); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		orderMap := fiber.Map{
			"id":           order.ID,
			"amount":       order.Amount,
			"table_number": order.TableNumber,
			"status_id":    order.StatusID,
			"order_date":   order.OrderDate,
			"menu_id":      order.MenuID,
			"order_code":   order.OrderCode,
			"created_at":   order.CreatedAt,
			"updated_at":   order.UpdatedAt,
			"status_name":  statusName,
			"menu": fiber.Map{
				"name":          menuName,
				"description":   menuDescription,
				"price":         menuPrice,
				"category_name": categoryName,
			},
			"total_price": totalPrice,
		}

		orders = append(orders, orderMap)
	}

	return c.JSON(fiber.Map{"success": true, "orders": orders})
}

type UpdateOrderRequest struct {
	Amount      int    `json:"amount" validate:"required,number"`
	TableNumber string `json:"table_number"`
	OrderCode   string `json:"order_code"`
	MenuID      int    `json:"menu_id" validate:"number"`
}

func UpdateOrder(c *fiber.Ctx, dbConn *sql.DB) error {
	// Extract the order ID from the URL parameter
	id := c.Params("id")

	// Create a new instance of UpdateOrderRequest
	data := new(UpdateOrderRequest)

	// Parse the request body
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Update the order in the database
	_, err := dbConn.Exec(`
        UPDATE orders 
        SET amount = $1, table_number = $2, order_code = $3, menu_id = $4, updated_at = NOW() 
        WHERE id = $5
    `, data.Amount, data.TableNumber, data.OrderCode, data.MenuID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order updated successfully",
	})
}

func DeleteOrder(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")

	_, err := dbConn.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func UpdateOrderStatus(c *fiber.Ctx, dbConn *sql.DB) error {
	type UpdateStatusRequest struct {
		OrderCode string `json:"order_code"`
		ID        int    `json:"id"`
		StatusID  int    `json:"status_id"`
	}

	var request UpdateStatusRequest
	fmt.Print(request)
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if both order_code and id are provided
	if request.OrderCode == "" && request.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Order code or ID must be provided"})
	}

	var err error

	// Update the order status based on the provided order code or ID
	if request.ID != 0 {
		_, err = dbConn.Exec("UPDATE orders SET status_id = $1, updated_at = NOW() WHERE id = $2", request.StatusID, request.ID)
	} else if request.OrderCode != "" {
		_, err = dbConn.Exec("UPDATE orders SET status_id = $1, updated_at = NOW() WHERE order_code = $2", request.StatusID, request.OrderCode)
		if request.StatusID == 3 {
			_, err = dbConn.Exec("UPDATE customers SET active = false, end_date = NOW() WHERE order_code = $1", request.OrderCode)

		}
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
