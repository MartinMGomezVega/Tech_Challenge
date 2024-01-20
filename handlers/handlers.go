package handlers

import (
	"context"
	"fmt"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/MartinMGomezVega/Tech_Challenge/routers"
	"github.com/aws/aws-lambda-go/events"
)

func Handlers(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {

	fmt.Println("Processing: " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))

	var r models.ResposeAPI
	r.Status = 400

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
			return routers.CreateUser(ctx)

		case "uploadTransactionFile":
			return routers.UploadTransactionFile(ctx, request)

			// case "storeTransactionsInDB":
			// 	return routers.StoreTransactionsInDB(ctx)

			// case "sendEmail":
			// 	return routers.SendEmail(ctx)
		}

	case "GET":
		//
	case "PUT":
		//
	case "DELETE":
		//
	}

	r.Status = 400
	r.Message = "Method Invalid"
	return r
}

func validateAuthorization(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string) {
	path := ctx.Value(models.Key("path")).(string)
	if path == "uploadTransactionFile" || path == "sendEmail" || path == "storeTransactionsInDB" || path == "createUser" {
		return true, 200, ""
	}

	return false, 400, "The path entered is incorrect."
}
