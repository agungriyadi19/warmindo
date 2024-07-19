// api/routes.go

package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, dbConn *sql.DB) {
	// Set up user routes
	SetupUserRoutes(app, dbConn)

	// Set up menu routes
	SetupMenuRoutes(app, dbConn)

	// Set up order routes
	SetupOrderRoutes(app, dbConn)

	// Set up category routes
	SetupCategoryRoutes(app, dbConn)

	// Set up role routes
	SetupRoleRoutes(app, dbConn)

	// Set up status routes
	SetupStatusRoutes(app, dbConn)

	SetupSettingsRoutes(app, dbConn)
}
