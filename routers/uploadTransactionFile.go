package routers

import (
	"context"
	"log"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
)

func UploadTransactionFile(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {
	log.Println("Saving file...")
	var r models.ResposeAPI
	r.Status = 400

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))
	log.Printf("bucket name: %s\n", *bucket)

	r.Status = 200
	r.Message = "test."
	return r
}
