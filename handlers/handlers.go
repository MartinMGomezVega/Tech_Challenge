package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/MartinMGomezVega/Tech_Challenge/routers"
	"github.com/aws/aws-lambda-go/events"
)

// Handlers: entry point for the AWS Lambda function
func Handlers(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {
	fmt.Println("Processing: " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))

	var r models.ResposeAPI
	r.Status = 400

	// Validate path
	isOk, statusCode, msg := validateAuthorization(ctx, request)
	if !isOk {
		r.Status = statusCode
		r.Message = msg
		return r
	}

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		case "createUser":
			log.Println("Accessing CreateUser")
			return routers.CreateUser(ctx)

		case "uploadTransactionFile":
			log.Println("Accessing UploadTransactionFile")
			return routers.UploadTransactionFile(ctx, request)

		case "sendEmail":
			return routers.SendEmail(ctx, request)
		}

	case "GET":
		//
	case "PUT":
		//
	case "DELETE":
		//
	}

	r.Status = 400
	r.Message = "Method Invalid."
	return r
}

// validateAuthorization: Validates the incoming path
func validateAuthorization(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string) {
	path := ctx.Value(models.Key("path")).(string)
	if path == "uploadTransactionFile" || path == "sendEmail" || path == "createUser" {
		return true, 200, ""
	}

	return false, 400, "The path entered is incorrect."
}
