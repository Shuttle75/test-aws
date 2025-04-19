package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-dax-go-v2/dax"
	"github.com/redis/go-redis/v9"
	"log"
	"test-aws/throttling"
	"time"
)

//aws ssm start-session --target i-0a283e67790f4c8c7 --document-name AWS-StartPortForwardingSessionToRemoteHost --parameters '{\"host\":[\"redis-4guxsh.serverless.euw2.cache.amazonaws.com\"],\"portNumber\":[\"6379\"],\"localPortNumber\":[\"6379\"]}'

var client = redis.NewClient(&redis.Options{
	Addr:                  "master.redis.4guxsh.euw2.cache.amazonaws.com:6379",
	ContextTimeoutEnabled: false,
	DB:                    0,
})

var endpoint = "daxs://dax-cluster.4guxsh.dax-clusters.eu-west-2.amazonaws.com"
var daxClient, _ = daxClientNew(endpoint, "eu-west-2")
var tableThrottling = throttling.TableThrottling{
	TableName: "thr_req_to_dest_operator",
	Client:    daxClient,
}

func daxClientNew(endpoint, region string) (*dax.Dax, error) {
	cfg := dax.DefaultConfig()
	cfg.HostPorts = []string{endpoint}
	cfg.Region = region
	cfg.SkipHostnameVerification = false
	client, err := dax.New(cfg)
	if err != nil {
		panic(fmt.Errorf("unable to initialize client %v", err))
	}
	return client, err
}

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

	for i := 0; i < 100; i++ {
		t1 := time.Now()

		last, err := tableThrottling.GetItem(context.Background(), "0123456789")
		if err != nil {
			log.Printf("Couldn't get item: %v\n", err)
		} else {
			log.Printf("Get item: %v\n", last)
		}

		err = tableThrottling.UpdateItem(context.Background(), "0123456789", time.Now().UnixNano())
		if err != nil {
			log.Printf("Couldn't put item into table thr_req_to_dest_operator: %v\n", err)
		} else {
			log.Printf("Set item\n")
		}

		println("Milliseconds ", time.Now().Sub(t1).Milliseconds())
	}
}
