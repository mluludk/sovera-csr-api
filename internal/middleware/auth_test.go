package middleware_test

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"sovera-core-api/internal/middleware"
)

func TestAuthenticateJWT_ValidToken(t *testing.T) {
	secret := "super_secret_jwt_key_enterprise"
	app := fiber.New()
	app.Get("/protected", middleware.AuthenticateJWT(secret), func(c *fiber.Ctx) error {
		orgID := c.Locals("org_id").(string)
		return c.SendString(orgID)
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    "usr_991823a",
		"org_id": "org_77123aa",
		"email":  "fundraiser@laznas.org",
		"role":   "FUNDRAISER",
		"exp":    time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestAuthenticateJWT_MissingToken(t *testing.T) {
	secret := "super_secret_jwt_key_enterprise"
	app := fiber.New()
	app.Get("/protected", middleware.AuthenticateJWT(secret), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", resp.StatusCode)
	}
}
