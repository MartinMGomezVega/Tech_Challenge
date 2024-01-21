package routers

import (
	"context"
	"log"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-lambda-go/events"
)

func UploadTransactionFile(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {
	log.Println("UploadTransactionFile")
	var r models.ResposeAPI
	r.Status = 400

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))
	log.Println("bucketName: " + *bucket)

	r.Status = 200
	r.Message = "CSV file successfully uploaded."
	return r
}
