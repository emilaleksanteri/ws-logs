package server

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type FiberServer struct {
	*fiber.App
	Redis *redis.Client
}

func New() *FiberServer {
	redisClient, err := createRedisClient("redis:6379")
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	server := &FiberServer{
		App:   fiber.New(),
		Redis: redisClient,
	}

	return server
}

func createRedisClient(address string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
