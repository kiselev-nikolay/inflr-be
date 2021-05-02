package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
)

const (
	ModeDev        = iota
	ModeProduction = iota
)

func GetApp(mode int) *fiber.App {
	app := fiber.New()

	switch mode {
	case ModeDev:
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	case ModeProduction:
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
	app.Use(authware.NewAuthware(&authware.Config{
		Key: []byte("Kh4Hy=bKRZ^fkq!RE7P8cBx=KLAb#nU^4Es$7srGHdH8@g79q2"),
	}))
	return app
}
