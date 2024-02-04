package server

import (
	"context"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Use("/logs", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	s.App.Get("/", s.HelloWorldHandler)
	s.App.Post("/log/:channel", s.logHandler)
	s.App.Get("/logs/:channel", websocket.New(s.eventSocketHandler))
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) logHandler(c *fiber.Ctx) error {
	channel := c.Params("channel")

	var input struct {
		Log string `json:"log"`
	}

	if err := c.BodyParser(&input); err != nil {
		return err
	}

	call := s.Redis.Publish(c.Context(), channel, input.Log)
	_, err := call.Result()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Log sent",
	})
}

func (s *FiberServer) eventSocketHandler(c *websocket.Conn) {
	channel := c.Params("channel")
	sub := s.Redis.Subscribe(context.Background(), channel)

	defer sub.Close()

	ch := sub.Channel()
	for msg := range ch {
		c.WriteJSON(fiber.Map{
			"message": msg.Payload,
		})
	}
}
