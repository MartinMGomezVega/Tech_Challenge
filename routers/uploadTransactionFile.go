package routers

import (
	"context"
	"log"

	"github.com/MartinMGomezVega/Tech_Challenge/models"

	"github.com/aws/aws-lambda-go/events"
)

func UploadTransactionFile(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {
	log.Println("Saving file...")
	var r models.ResposeAPI

	bucketName := ctx.Value(models.Key("bucketName")).(string)
	log.Println("bucket: " + bucketName)

	S3.Region = "us-east-1"
	S3.NewSession(S3.Region)
	S3.Ls()
	S3.Upload("27426626956.csv", bucketName, "subido.csv")
	S3.GenerateUrl(bucketName, "subido.csv")

	return r
}
