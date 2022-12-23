package handler

import (
	"context"
	"log"
	"net/http"
	"time"
	"ws-messenger/internal/model"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func Connect(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("new connection. id: %s", event.RequestContext.ConnectionID)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	conn := model.Connection{
		ConnectionID:   event.RequestContext.ConnectionID,
		ExpirationTime: int(time.Now().Add(5 * time.Minute).Unix()),
	}

	item, err := dynamodbattribute.MarshalMap(conn)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "unknown error",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("ws-messenger-table"),
		Item:      item,
	}

	svc.PutItemWithContext(ctx, input)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
