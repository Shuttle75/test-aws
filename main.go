package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-dax-go-v2/dax"
	"github.com/aws/aws-dax-go-v2/dax/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/redis/go-redis/v9"
	"log"
)

//aws ssm start-session --target i-0a283e67790f4c8c7 --document-name AWS-StartPortForwardingSessionToRemoteHost --parameters '{\"host\":[\"redis-4guxsh.serverless.euw2.cache.amazonaws.com\"],\"portNumber\":[\"6379\"],\"localPortNumber\":[\"6379\"]}'

var client = redis.NewClient(&redis.Options{
	Addr:                  "redis-4guxsh.serverless.euw2.cache.amazonaws.com:6379",
	ContextTimeoutEnabled: false,
	DB:                    0,
})

var endpoint = "daxs://dax-cluster.4guxsh.dax-clusters.eu-west-2.amazonaws.com"
var daxClient, _ = daxClientNew(endpoint, "eu-west-2")

func daxClientNew(endpoint, region string) (*dax.Dax, error) {
	cfg := dax.DefaultConfig()
	cfg.HostPorts = []string{endpoint}
	cfg.Region = region
	cfg.SkipHostnameVerification = false
	//cfg.ClientHealthCheckInterval = 30 * time.Second
	//cfg.IdleConnectionReapDelay = 30 * time.Second
	//cfg.MaxPendingConnectionsPerHost = 100
	cfg.LogLevel = utils.LogDebug
	//cfg.DialContext = func(ctx context.Context, network string, address string) (net.Conn, error) {
	//	// fmt.Println("Write your custom logic here")
	//	dialCon, err := dax.SecureDialContext(endpoint, cfg.SkipHostnameVerification)
	//	if err != nil {
	//		panic(fmt.Errorf("secure dialcontext creation failed %v", err))
	//	}
	//	return dialCon(ctx, network, address)
	//}
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

	in := dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"DestOperatorId": &types.AttributeValueMemberS{Value: "0123456789"},
		},
		TableName: aws.String("thr_req_to_dest_operator"),
	}
	var item struct {
		DestOperatorId string
		LastReq        int64
	}
	out, err := daxClient.GetItem(context.Background(), &in)
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", "0123456789", err)
	} else {
		errMap := attributevalue.UnmarshalMap(out.Item, &item)
		if errMap != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		} else {
			log.Printf("Get info about %v\n", item)
		}
	}
}
