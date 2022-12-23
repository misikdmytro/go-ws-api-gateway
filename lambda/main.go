package main

import (
	"context"
	"ws-messenger/internal/handler"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch event.RequestContext.EventType {
	case "CONNECT":
		return handler.Connect(ctx, event)
	case "DISCONNECT":
		return handler.Disconnect(ctx, event)
	default:
		return events.APIGatewayProxyResponse{Body: "no handler", StatusCode: 200}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
