package routers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	// Obtener el nombre del archivo de la solicitud
	fileName := request.PathParameters["filename"] // Asumiendo que el parámetro en la URL es "filename"

	// Load Mexico's time zone
	location, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}
	now := time.Now().In(location) // Hora actual en la Ciudad de México

	// Formatear el nombre del archivo con la fecha y hora actual
	filename := fmt.Sprintf("files/%s_%s.csv", fileName, now.Format("02012006_030405PM"))
	r = AWSService.UploadFile(bucketName, filename, "/files/"+fileName)

	return r
}

func (awsSvc AWSService) UploadFile(bucketName string, bucketKey string, fileName string) models.ResposeAPI {
	var r models.ResposeAPI
	file, err := os.Open(fileName)
	if err != nil {
		log.Println(fmt.Errorf("failed to open file %q, %v", fileName, err))
		r.Status = 400
		r.Message = "Failed to open file."
	} else {
		defer file.Close()
		// Upload the file to S3.
		result, err := awsSvc.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(bucketKey),
			Body:   file,
		})
		if err != nil {
			log.Println(fmt.Errorf("failed to upload file, %v", err))
			r.Status = 400
			r.Message = "Failed to upload file."
		} else {
			r.Status = 200
			r.Message = "CSV file successfully uploaded."
		}
		log.Println(result)
	}

	return r
}
