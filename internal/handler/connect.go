package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"ws-messenger/internal/helper"
	"ws-messenger/internal/model"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Connect(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("new connection. id: %s", event.RequestContext.ConnectionID)

	svc, err := helper.NewDynamoDB(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	item, err := attributevalue.MarshalMap(model.Connection{
		ConnectionID:   event.RequestContext.ConnectionID,
		ExpirationTime: int(time.Now().Add(5 * time.Minute).Unix()),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Item:      item,
	}

	svc.PutItem(ctx, input)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
