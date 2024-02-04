package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

type LogEvent struct {
	Domain     string    `json:"domain"`
	AccessedAt time.Time `json:"accessed_at"`
	IsMyDomain bool      `json:"is_my_domain"`
}

func (s *FiberServer) logHandler(c *fiber.Ctx) error {
	channel := c.Params("channel")

	var input struct {
		Domain string `json:"domain"`
	}

	if err := c.BodyParser(&input); err != nil {
		return err
	}

	event := LogEvent{
		Domain:     input.Domain,
		AccessedAt: time.Now(),
		IsMyDomain: input.Domain == "mydomain.com",
	}

	json, err := json.Marshal(event)
	if err != nil {
		return err
	}

	call := s.Redis.Publish(c.Context(), channel, json)
	_, err = call.Result()
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
		var event LogEvent
		err := json.Unmarshal([]byte(msg.Payload), &event)
		if err != nil {
			fmt.Println("Error parsing message", err)
			continue
		}

		c.WriteJSON(event)
	}
}
