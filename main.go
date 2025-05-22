package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-dax-go-v2/dax"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"test-aws/throttling"
	"time"
)

var daxEndpoint = "daxs://dax-cluster.4guxsh.dax-clusters.eu-west-2.amazonaws.com"
var ddEndpoint = "https://dynamodb.eu-west-2.amazonaws.com"
var daxClient, _ = daxClientNew(daxEndpoint, "eu-west-2")
var ddClient, _ = ddClientNew(ddEndpoint, "eu-west-2")
var tableThrottlingDd = throttling.TableThrottlingDD{
	TableName: "dax_test_table",
	Client:    ddClient,
}
var tableThrottlingDax = throttling.TableThrottlingDax{
	TableName: "dax_test_table",
	Client:    daxClient,
}

func ddClientNew(endpoint, region string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	cfg.Region = region
	cfg.BaseEndpoint = &endpoint

	var client = dynamodb.NewFromConfig(cfg)
	return client, err
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
	for i := 0; i < 10; i++ {
		t1 := time.Now()
		last, err := tableThrottlingDax.GetItem(context.Background(), "0123456789")
		if err != nil {
			log.Printf("Couldn't get item: %v\n", err)
		} else {
			log.Printf("Get item: %v microsec %v\n", last, time.Now().Sub(t1).Microseconds())
		}

		t1 = time.Now()
		newValue := time.Now().UnixNano()
		err = tableThrottlingDax.UpdateItem(context.Background(), "0123456789", newValue)
		if err != nil {
			log.Printf("Couldn't put item into table thr_req_to_dest_operator: %v\n", err)
		} else {
			log.Printf("Set item: %v microsec %v\n\n", newValue, time.Now().Sub(t1).Microseconds())
		}
	}

	log.Printf("-------------------------------------------------------------------------")

	for i := 0; i < 10; i++ {
		t1 := time.Now()
		last, err := tableThrottlingDd.GetItem(context.Background(), "0123456789")
		if err != nil {
			log.Printf("Couldn't get item: %v\n", err)
		} else {
			log.Printf("Get item: %v microsec %v\n", last, time.Now().Sub(t1).Microseconds())
		}

		t1 = time.Now()
		newValue := time.Now().UnixNano()
		err = tableThrottlingDd.UpdateItem(context.Background(), "0123456789", newValue)
		if err != nil {
			log.Printf("Couldn't put item into table thr_req_to_dest_operator: %v\n", err)
		} else {
			log.Printf("Set item: %v microsec %v\n\n", newValue, time.Now().Sub(t1).Microseconds())
		}
	}
}
