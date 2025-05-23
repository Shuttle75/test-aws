package throttling

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
	"strconv"
)

type TableThrottlingDD struct {
	Client    *dynamodb.Client
	TableName string
}

func (table *TableThrottlingDD) UpdateItem(ctx context.Context, operator string, last int64) error {
	_, err := table.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(table.TableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: operator},
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":last": &types.AttributeValueMemberN{Value: strconv.FormatInt(last, 10)},
		},
		UpdateExpression: aws.String("set LastReq = :last"),
	})
	if err != nil {
		log.Printf("Couldn't update operator %v. Here's why: %v\n", operator, err)
	}
	return err
}

func (table *TableThrottlingDD) GetItem(ctx context.Context, operator string) (int64, error) {
	in := dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: operator},
		},
		TableName: aws.String(table.TableName),
	}

	var item struct {
		id      string
		LastReq int64
	}
	out, err := table.Client.GetItem(ctx, &in)
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", operator, err)
	} else {
		errMap := attributevalue.UnmarshalMap(out.Item, &item)
		if errMap != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}
	return item.LastReq, nil
}
