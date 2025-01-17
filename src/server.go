package main

import (
	scraper_control "genote-watcher/scraper-control"
	"log"
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func StartServer(commandChan chan scraper_control.Command, reponseChan chan scraper_control.Response) {

	app := fiber.New()

	app.Use(logger.New())

	defer func() {
		app.Shutdown()
	}()

	app.Static("/", "./client/dist")

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
