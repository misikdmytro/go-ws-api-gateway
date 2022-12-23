package helper

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

func NewAPIGatewayManagementAPI(ctx context.Context) (*apigatewaymanagementapi.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           os.Getenv("API_GATEWAY_ENDPOINT"),
			SigningRegion: os.Getenv("AWS_REGION"),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		return nil, err
	}

	return apigatewaymanagementapi.NewFromConfig(cfg), nil
}
