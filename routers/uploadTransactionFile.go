package routers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type AWSService struct {
	S3Client *s3.Client
}

func UploadTransactionFile(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {
	log.Println("Saving file...")
	var r models.ResposeAPI

	bucketName := ctx.Value(models.Key("bucketName")).(string)
	log.Println("bucket: " + bucketName)

	// Ruta completa al archivo que deseas subir
	filePath := "C:/Users/Martin/Tech_Challenge/files/20417027050.csv"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Println("Error while loading the aws config: ", err)
		r.Status = 400
		r.Message = "Error while loading the aws config."
		return r
	}

	AWSService := AWSService{
		S3Client: s3.NewFromConfig(config),
	}

	// Utilizando el nombre del archivo como bucketKey
	r = AWSService.UploadFile(bucketName, "20417027050.csv", filePath)

	return r
}

func (awsSvc AWSService) UploadFile(bucketName string, bucketKey string, filePath string) models.ResposeAPI {
	var r models.ResposeAPI

	// Abrir el archivo
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to open file %q, %v", filePath, err))
		r.Status = 400
		r.Message = "Failed to open file."
		return r
	}
	defer file.Close()

	// Subir el contenido del archivo a S3
	result, err := awsSvc.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
		Body:   file,
	})
	if err != nil {
		fmt.Println(fmt.Errorf("failed to upload file, %v", err))
		r.Status = 400
		r.Message = "Failed to upload file."
		return r
	}

	r.Status = 200
	r.Message = "CSV file successfully uploaded."
	fmt.Println(result)

	return r
}
