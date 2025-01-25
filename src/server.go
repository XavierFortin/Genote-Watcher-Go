package main

import (
	"embed"
	"genote-watcher/scrapers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed client/dist/*
var clientHtml embed.FS

func StartServer(scraper *scrapers.GenoteScraper) {

	app := fiber.New()

	app.Use(logger.New())

	defer func() {
		app.Shutdown()
	}()

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(clientHtml),
		PathPrefix: "client/dist",
	}))

	app.Use("/api", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		return c.Status(200).Next()
	})

	app.Get("/api/scraper/status", func(c *fiber.Ctx) error {
		status := scraper.GetStatus()
		return c.JSON(status)
	})

	app.Get("/api/logs", func(c *fiber.Ctx) error {
		file, err := os.ReadFile("log.txt")

		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(file))
	})

	app.Post("/api/scraper/start", func(c *fiber.Ctx) error {
		scraper.Start()
		return c.SendStatus(200)
	})

	app.Post("/api/scraper/stop", func(c *fiber.Ctx) error {
		scraper.Stop()
		return c.SendStatus(200)
	})

	app.Post("/api/scraper/force-start", func(c *fiber.Ctx) error {
		scraper.ScrapeOnce()
		return c.SendStatus(200)
	})

	app.Post("/api/scraper/change-interval", func(c *fiber.Ctx) error {
		type interval struct {
			Interval string `json:"interval"`
		}

		inter := interval{}
		if err := c.BodyParser(&inter); err != nil {
			c.Status(400).SendString(err.Error())
			return err
		}

		duration, err := time.ParseDuration(inter.Interval)

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid interval duration, the format must be like 1s, 1m, 1h, 1d or a combination of them")
		}

		scraper.SetInterval(duration)

		return c.SendStatus(200)
	})

	// Upgraded websocket request
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		mutex.Lock()
		clients[c] = true
		mutex.Unlock()

		defer func() {
			mutex.Lock()
			delete(clients, c)
			mutex.Unlock()
			c.Close()
		}()

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))

	log.Fatal(app.Listen(":4000"))
}
