package models

import "github.com/aws/aws-lambda-go/events"

type ResposeAPI struct {
	Status     int
	Message    string
	CustomResp *events.APIGatewayProxyResponse
}
