package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

//aws ssm start-session --target i-0a283e67790f4c8c7 --document-name AWS-StartPortForwardingSessionToRemoteHost --parameters '{\"host\":[\"redis-4guxsh.serverless.euw2.cache.amazonaws.com\"],\"portNumber\":[\"6379\"],\"localPortNumber\":[\"6379\"]}'

// var limiterRedis = throttling.NewLimiter(5)
var client = redis.NewClient(&redis.Options{
	Addr:                  "redis-4guxsh.serverless.euw2.cache.amazonaws.com:6379",
	ContextTimeoutEnabled: false,
	DB:                    0,
})

func main() {
	_ = client.Set(context.Background(), "1234", "5678", 0).Err()
	value, err := client.Get(context.Background(), "1234").Int64()
	if err != nil {
		log.Printf("Couldn't set item: %v\n", err)
	} else {
		log.Printf("Set item: %v\n", value)
	}
}
