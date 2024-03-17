package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Default DB
	})
}

func WriteCredentials(client *redis.Client, email, password string) error {
	return client.Set(ctx, email, password, 0).Err()
}
func ReadCredentials(client *redis.Client, email string) (string, error) {
	return client.Get(ctx, email).Result()
}
func RemoveCredentials(client *redis.Client, email string) error {
	return client.Del(ctx, email).Err()
}
func CheckIfExists(client *redis.Client, key string) (bool, error) {
	_, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
