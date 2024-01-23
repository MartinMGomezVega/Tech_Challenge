package commons

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/MartinMGomezVega/Tech_Challenge/awsgo"
	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

// GetTransactionFile: obtains the file from the s3 bucket with the cuil
func GetTransactionFile(ctx context.Context, cuil string) (*bytes.Buffer, *models.ResposeAPI) {
	var r models.ResposeAPI
	r.Status = 400

	filename := fmt.Sprintf("files/%s.csv", cuil)
	log.Printf("filename: %s\n", filename)

	svc := s3.NewFromConfig(awsgo.Cfg)

	// Download the S3 file
	file, err := downloadFromS3(ctx, svc, filename)
	if err != nil {
		r.Status = 500
		r.Message = fmt.Sprintf("Error downloading csv file from S3: %s", err.Error())
		return nil, &r
	}

	return file, nil
}

// downloadFromS3: Function that downloads a file from AWS S3
func downloadFromS3(ctx context.Context, svc *s3.Client, filename string) (*bytes.Buffer, error) {
	bucket := ctx.Value(models.Key("bucketName")).(string)
	log.Printf("bucket: %s", *aws.String(bucket))

	obj, err := svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		log.Println("Error getting the s3 bucket.")
		return nil, err
	}
	defer obj.Body.Close()

	// Check if the file is CSV
	if obj.ContentType != nil && *obj.ContentType != "text/csv" {
		return nil, fmt.Errorf("the file is not of type csv")
	}

	file, err := io.ReadAll(obj.Body)
	if err != nil {
		log.Println("Error reading file.")
		return nil, err
	}
	buffer := bytes.NewBuffer(file)

	log.Println("The file was obtained correctly.")
	return buffer, nil
}
