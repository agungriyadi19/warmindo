package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// GenerateJWT generates a JWT token
func GenerateJWT(userID string, roleID int) (string, error) {
	claims := jwt.MapClaims{
		"sub":     userID,
		"role_id": roleID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

// JWTProtected returns a middleware handler for JWT authentication
func AuthMiddleware(requiredRole int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No token provided"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token
		secretKey := os.Getenv("JWT_KEY")
		if secretKey == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server configuration error"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		roleID, ok := claims["role_id"].(float64)
		if !ok || int(roleID) != requiredRole { // Check if role_id is 1 (Admin)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}

		return c.Next()
	}
}
