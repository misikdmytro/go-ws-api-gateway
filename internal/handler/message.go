package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"ws-messenger/internal/model"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func Message(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("new message")

	const messageAction = "MESSAGE"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	var request model.Request[model.MessageRequestPayload]
	if err := json.Unmarshal([]byte(event.Body), &request); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	filt := expression.Name("ConnectionID").NotEqual(expression.Value(event.RequestContext.ConnectionID))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String("ws-messenger-table"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	output, err := svc.ScanWithContext(ctx, input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	api := apigatewaymanagementapi.New(sess)
	for _, item := range output.Items {
		var conn model.Connection
		if err := dynamodbattribute.UnmarshalMap(item, &conn); err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		newMessage := model.Response[model.MessageRequestPayload]{
			Action: messageAction,
			Response: model.MessageRequestPayload{
				Message: request.Payload.Message,
			},
		}
		data, err := json.Marshal(newMessage)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		input := &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(conn.ConnectionID),
			Data:         data,
		}

		api.PostToConnectionWithContext(ctx, input)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
	}, nil
}
