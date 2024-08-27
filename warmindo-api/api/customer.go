package api

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"time"
	"warmindo-api/db"

	"github.com/gofiber/fiber/v2"
)

const EarthRadiusKm = 6371 // Earth's radius in kilometers

// SetupCustomerRoutes sets up the routes for customer-related operations
func SetupCustomerRoutes(app *fiber.App, dbConn *sql.DB) {
	customerAPI := app.Group("/api/customer")

	customerAPI.Post("/check-radius", func(c *fiber.Ctx) error {
		return CheckLocationAndGenerateOrderCode(c, dbConn)
	})

	customerAPI.Post("/check-code", func(c *fiber.Ctx) error {
		return CheckOrderCodeAndTableNumber(c, dbConn)
	})
	// New endpoint to check if a table has an active order
	customerAPI.Post("/check-active", func(c *fiber.Ctx) error {
		return CheckActiveOrder(c, dbConn)
	})
}

// CheckActiveOrder handles checking if a table has an active order
func CheckActiveOrder(c *fiber.Ctx, dbConn *sql.DB) error {
	var req struct {
		TableNumber int    `json:"table_number"`
		OrderCode   string `json:"order_code"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "error": "Body permintaan tidak valid"})
	}

	// Check if the active field for the table_number is true
	var active bool
	err := dbConn.QueryRow("SELECT active FROM customers WHERE order_code = $1 AND table_number = $2 ORDER BY start_date DESC LIMIT 1", req.TableNumber).Scan(&active)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "error": "Gagal memeriksa status aktif"})
	}

	return c.JSON(fiber.Map{"status": true, "table_number": req.TableNumber, "active": active})
}

// isWithinRadius calculates if a location is within a specified radius using the Haversine formula
func isWithinRadius(lat1, lon1, lat2, lon2, radius float64) bool {
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := EarthRadiusKm * c
	return distance <= radius
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Generate a unique order code
func generateUniqueOrderCode(dbConn *sql.DB) (string, error) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const digits = "0123456789"
	rand.Seed(time.Now().UnixNano())

	for {
		orderCode := fmt.Sprintf("%c%c%c%c-%d%d%d%d",
			letters[rand.Intn(len(letters))],
			letters[rand.Intn(len(letters))],
			letters[rand.Intn(len(letters))],
			letters[rand.Intn(len(letters))],
			digits[rand.Intn(len(digits))],
			digits[rand.Intn(len(digits))],
			digits[rand.Intn(len(digits))],
			digits[rand.Intn(len(digits))],
		)

		// Check if the order code already exists
		var count int
		err := dbConn.QueryRow("SELECT COUNT(*) FROM customers WHERE order_code = $1", orderCode).Scan(&count)
		if err != nil {
			return "", err
		}
		if count == 0 {
			return orderCode, nil
		}
	}
}

// CheckLocationAndGenerateOrderCode handles checking the location and generating an order code
func CheckLocationAndGenerateOrderCode(c *fiber.Ctx, dbConn *sql.DB) error {
	var req struct {
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		TableNumber int     `json:"table_number"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "error": "Body permintaan tidak valid"})
	}

	// Get settings from the database
	var settings db.Settings
	err := dbConn.QueryRow("SELECT latitude, longitude, radius FROM settings WHERE id = 1").Scan(&settings.Latitude, &settings.Longitude, &settings.Radius)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "error": "Gagal mengambil pengaturan"})
	}

	// Check if the location is within the radius
	if !isWithinRadius(req.Latitude, req.Longitude, settings.Latitude, settings.Longitude, settings.Radius) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": false, "error": "Lokasi berada di luar radius"})
	}

	// Check if the active field for the table_number is true
	var active bool
	err = dbConn.QueryRow("SELECT active FROM customers WHERE table_number = $1 ORDER BY start_date DESC LIMIT 1", req.TableNumber).Scan(&active)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "error": "Gagal memeriksa status aktif"})
	}

	if active {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": false, "error": "Meja sudah memiliki kode pesanan aktif"})
	}

	// Generate a unique order code
	orderCode, err := generateUniqueOrderCode(dbConn)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "error": "Gagal membuat kode pesanan"})
	}

	// Insert into the database
	_, err = dbConn.Exec("INSERT INTO customers (order_code, table_number, start_date, active) VALUES ($1, $2, NOW(), true)", orderCode, req.TableNumber)
	if err != nil {
		fmt.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "error": "Gagal menyimpan kode pesanan"})
	}

	return c.JSON(fiber.Map{"status": true, "order_code": orderCode})
}

// CheckOrderCodeAndTableNumber handles checking the status of an order code and table number
func CheckOrderCodeAndTableNumber(c *fiber.Ctx, dbConn *sql.DB) error {
	var req struct {
		OrderCode   string `json:"order_code"`
		TableNumber int    `json:"table_number"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "error": "Body permintaan tidak valid"})
	}

	// Check if the provided order code and table number are valid and get the end date
	var endDate sql.NullTime
	err := dbConn.QueryRow("SELECT end_date FROM customers WHERE order_code = $1 AND table_number = $2", req.OrderCode, req.TableNumber).Scan(&endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "error": "Kode pesanan atau nomor meja tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "error": "Gagal memeriksa kode pesanan"})
	}

	// Determine the status based on end_date
	orderStatus := !endDate.Valid // true if end_date is empty, false otherwise

	return c.JSON(fiber.Map{"status": true, "order_status": orderStatus, "table_number": req.TableNumber, "end_date": endDate.Time})
}
