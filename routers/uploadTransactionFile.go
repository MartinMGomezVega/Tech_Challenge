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

	// Load Mexico's time zone
	// location, err := time.LoadLocation("America/Mexico_City")
	// if err != nil {
	// 	r.Status = 500
	// 	r.Message = err.Error()
	// 	return r
	// }
	// now := time.Now().In(location) // Mexico Time
	// filename := fmt.Sprintf("files/20417027050_%s_%s.csv", now.Format("02012006"), now.Format("030405PM"))
	r = AWSService.UploadFile(bucketName, "filename", "/files/20417027050.csv")

	return r
}

func (awsSvc AWSService) UploadFile(bucketName string, filename string, file string) models.ResposeAPI {
	var r models.ResposeAPI
	f, err := os.Open(file)
	if err != nil {
		log.Println(fmt.Errorf("failed to open file %q, %v", filename, err))
		r.Status = 400
		r.Message = "Failed to open file."
	} else {
		defer f.Close()
		// Upload the file to S3.
		result, err := awsSvc.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(filename),
			Body:   f,
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
