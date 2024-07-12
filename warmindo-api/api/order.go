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

	_, err := dbConn.Exec("INSERT INTO orders (amount, table_number, status_id, order_date, menu_id, order_code) VALUES ($1, $2, $3, $4, $5, $6)",
		order.Amount, order.TableNumber, order.StatusID, order.OrderDate, order.MenuID, order.OrderCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func GetOrders(c *fiber.Ctx, dbConn *sql.DB) error {
	query := `
	SELECT o.id, o.amount, o.table_number, o.status_id, o.order_date, o.menu_id, o.order_code, 
           o.created_at, o.updated_at,
           s.name as status_name,
           m.name as menu_name, m.image as menu_image, m.description as menu_description, m.price as menu_price,
           c.name as category_name,
           (o.amount * m.price) as total_price
    FROM orders o
    JOIN statuses s ON o.status_id = s.id
    JOIN menus m ON o.menu_id = m.id
    JOIN categories c ON m.category_id = c.id`

	rows, err := dbConn.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var orders []fiber.Map
	for rows.Next() {
		var order db.Order
		var statusName, menuName, menuImage, menuDescription, categoryName string
		var menuPrice, totalPrice int

		if err := rows.Scan(&order.ID, &order.Amount, &order.TableNumber, &order.StatusID, &order.OrderDate, &order.MenuID, &order.OrderCode,
			&order.CreatedAt, &order.UpdatedAt, &statusName, &menuName, &menuImage, &menuDescription, &menuPrice, &categoryName, &totalPrice); err != nil {
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
				"image":         menuImage,
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

func GetOrderByID(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")

	query := `
	SELECT o.id, o.amount, o.table_number, o.status_id, o.order_date, o.menu_id, o.order_code, 
           o.created_at, o.updated_at,
           s.name as status_name,
           m.name as menu_name, m.image as menu_image, m.description as menu_description, m.price as menu_price,
           c.name as category_name,
           (o.amount * m.price) as total_price
    FROM orders o
    JOIN statuses s ON o.status_id = s.id
    JOIN menus m ON o.menu_id = m.id
    JOIN categories c ON m.category_id = c.id
    WHERE o.id = $1`

	var order db.Order
	var statusName, menuName, menuImage, menuDescription, categoryName string
	var menuPrice, totalPrice int

	err := dbConn.QueryRow(query, id).Scan(&order.ID, &order.Amount, &order.TableNumber, &order.StatusID, &order.OrderDate, &order.MenuID, &order.OrderCode,
		&order.CreatedAt, &order.UpdatedAt, &statusName, &menuName, &menuImage, &menuDescription, &menuPrice, &categoryName, &totalPrice)
	if err != nil {
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
			"image":         menuImage,
			"description":   menuDescription,
			"price":         menuPrice,
			"category_name": categoryName,
		},
		"total_price": totalPrice,
	}

	return c.JSON(fiber.Map{"success": true, "order": orderMap})
}

func UpdateOrder(c *fiber.Ctx, dbConn *sql.DB) error {
	id := c.Params("id")
	var order db.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := dbConn.Exec("UPDATE orders SET amount = $1, table_number = $2, status_id = $3, order_date = $4, menu_id = $5, order_code = $6, updated_at = NOW() WHERE id = $7",
		order.Amount, order.TableNumber, order.StatusID, order.OrderDate, order.MenuID, order.OrderCode, id)
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
