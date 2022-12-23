package handler

import (
	"context"
	"encoding/json"
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

func Ping(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("ping message")

	const pongAction = "PONG"

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
		return events.APIGatewayProxyResponse{}, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("ws-messenger-table"),
		Item:      item,
	}

	_, err = svc.PutItemWithContext(ctx, input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := model.Response[model.PongResponsePayload]{
		Action:   pongAction,
		Response: model.PongResponsePayload{},
	}

	content, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	log.Printf("pong response: %s", string(content))

	return events.APIGatewayProxyResponse{
		Body:       string(content),
		StatusCode: http.StatusOK,
	}, nil
}
