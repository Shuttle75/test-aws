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
	setErr := client.Set(context.Background(), "1234", "5678", 0).Err()
	if setErr != nil {
		log.Printf("Couldn't set item: %v\n", setErr)
	}
	value, getErr := client.Get(context.Background(), "1234").Int64()
	if getErr != nil {
		log.Printf("Couldn't get item: %v\n", getErr)
	} else {
		log.Printf("Get item: %v\n", value)
	}
}
