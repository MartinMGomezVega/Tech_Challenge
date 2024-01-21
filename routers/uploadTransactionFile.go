package routers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"strings"
	"time"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type readSeeker struct {
	io.Reader
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func UploadTransactionFile(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {
	log.Println("Saving file...")
	var r models.ResposeAPI
	r.Status = 400

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))
	log.Printf("bucket name: %s\n", *bucket)

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary, ok := params["boundary"]
		if !ok {
			log.Println("Boundary not found in Content-Type header")
			r.Status = 400
			r.Message = "Boundary not found in Content-Type header."
			return r
		}

		mr := multipart.NewReader(strings.NewReader(request.Body), boundary)

		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading part: %v\n", err)
				r.Status = 500
				r.Message = err.Error()
				return r
			}

			log.Printf("Processing part: %s\n", p.FileName())

			if p.FileName() != "" {
				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, p); err != nil {
					log.Printf("Error copying part content: %v\n", err)
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				sess, err := session.NewSession(&aws.Config{
					Region: aws.String("us-east-1")},
				)

				if err != nil {
					log.Printf("Error creating AWS session: %v\n", err)
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				fileName := strings.TrimSuffix(p.FileName(), ".csv")
				location, err := time.LoadLocation("America/Mexico_City")
				if err != nil {
					log.Printf("Error loading time zone: %v\n", err)
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				now := time.Now().In(location)
				filename := fmt.Sprintf("transactions/%s_%s_%s.csv", fileName, now.Format("02012006"), now.Format("030405PM"))
				log.Printf("Name of the file with the transactions: %s\n", filename)

				uploader := s3manager.NewUploader(sess)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   &readSeeker{buf},
				})

				if err != nil {
					log.Printf("Error uploading to S3: %v\n", err)
					r.Status = 500
					r.Message = err.Error()
					return r
				}
			}
		}
	} else {
		r.Status = 400
		r.Message = "You must send a CSV file with the 'Content-Type' of type 'multipart/' in the Header."
		return r
	}

	r.Status = 200
	r.Message = "CSV file successfully uploaded."
	return r
}
