package api

import (
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Settings struct {
	Radius       int      `json:"radius"`
	WifiNetworks []string `json:"wifiNetworks"`
}

func GetSettings(c *fiber.Ctx) error {
	radius := os.Getenv("MAX_RADIUS")
	wifiNetworks := strings.Split(os.Getenv("ALLOWED_NETWORKS"), ",")
	return c.JSON(Settings{Radius: radius, WifiNetworks: wifiNetworks})
}

func SaveSettings(c *fiber.Ctx) error {
	var settings Settings
	if err := c.BodyParser(&settings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	os.Setenv("MAX_RADIUS", strconv.Itoa(settings.Radius))
	os.Setenv("ALLOWED_NETWORKS", strings.Join(settings.WifiNetworks, ","))

	return c.JSON(fiber.Map{"status": "success"})
}

func SetupRoutes(app *fiber.App) {
	app.Get("/api/settings", GetSettings)
	app.Post("/api/settings", SaveSettings)
}
