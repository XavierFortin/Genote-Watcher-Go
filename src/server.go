package main

import (
	"embed"
	scraper_control "genote-watcher/scraper-control"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed client/dist/*
var clientHtml embed.FS

func StartServer(commandChan chan scraper_control.Command, reponseChan chan scraper_control.Response) {

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
		commandChan <- scraper_control.Command{Action: scraper_control.Status}
		response := <-reponseChan

		return c.JSON(response)
	})

	app.Get("/api/logs", func(c *fiber.Ctx) error {
		file, err := os.ReadFile("log.txt")

		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(file))
	})

	app.Post("/api/scraper/start", func(c *fiber.Ctx) error {
		log.Println("Scraper started successfully")
		commandChan <- scraper_control.Command{Action: scraper_control.Start}
		return c.SendStatus(200)
	})

	app.Post("/api/scraper/stop", func(c *fiber.Ctx) error {
		log.Println("Stopped scraper successfully")
		commandChan <- scraper_control.Command{Action: scraper_control.Stop}
		return c.SendStatus(200)
	})

	app.Post("/api/scraper/force-start", func(c *fiber.Ctx) error {
		log.Println("Scraper force started once")
		commandChan <- scraper_control.Command{Action: scraper_control.ForceStartOnce}
		return c.SendStatus(200)
	})

	app.Post("/api/scraper/restart", func(c *fiber.Ctx) error {
		log.Println("Scraper restarted")
		commandChan <- scraper_control.Command{Action: scraper_control.Restart}
		return c.SendStatus(200)
	})

	app.Post("/api/scraper/change-interval", func(c *fiber.Ctx) error {
		log.Println("Scraper interval changed")

		type interval struct {
			Interval string `json:"interval"`
		}

		inter := interval{}
		if err := c.BodyParser(inter); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		commandChan <- scraper_control.Command{Action: scraper_control.ChangeInterval, Data: inter}

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
