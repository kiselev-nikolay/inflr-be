package main

import (
	"fmt"
	"log"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v2"
)

const sk = "Kh4Hy=bKRZ^fkq!RE7P8cBx=KLAb#nU^4Es$7srGHdH8@g79q2"

var cfg = jwtware.Config{
	SigningKey: []byte(sk),
}

func main() {
	app := fiber.New()
	devMode := true
	if devMode {
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	} else {
		app.Use(compress.New(compress.Config{
			Level: compress.LevelBestCompression,
		}))
		app.Use(cache.New(cache.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.Query("refresh") == "true"
			},
			Expiration:   2 * time.Minute,
			CacheControl: true,
		}))
		app.Use(cors.New(cors.Config{
			AllowOrigins: "https://inflr.app",
		}))
		app.Use(limiter.New(limiter.Config{
			Max:      10,
			Duration: 5 * time.Second,
		}))
	}
	app.Post("/login", func(c *fiber.Ctx) error {
		user := c.FormValue("user")
		pass := c.FormValue("pass")
		if user != "john" || pass != "doe" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "John Doe"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		t, err := token.SignedString(cfg.SigningKey)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(fiber.Map{"token": t})
	})
	app.Use(jwtware.New(cfg))
	app.Get("/name", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		msg := fmt.Sprintf("Hello, %s 👋!", name)
		return c.SendString(msg)
	})

	log.Fatal(app.Listen(":3000"))
}
