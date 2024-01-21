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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-lambda-go/events"
)

type readSeeker struct {
	io.Reader
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
func UploadTransactionFile(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {
	log.Println("UploadTransactionFile")
	var r models.ResposeAPI
	r.Status = 400

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))
	log.Println("bucketName: " + *bucket)

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		log.Println("1: ")

		mr := multipart.NewReader(strings.NewReader(request.Body), params["boundary"])
		p, err := mr.NextPart()
		if err != nil && err != io.EOF {
			log.Println("2: ")

			r.Status = 500
			r.Message = err.Error()
			return r
		}
		if err != io.EOF {
			if p.FileName() != "" {
				log.Println("3: ")

				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, p); err != nil {
					log.Println("4: ")

					r.Status = 500
					r.Message = err.Error()
					return r
				}

				sess, err := session.NewSession(&aws.Config{
					Region: aws.String("us-east-1")},
				)

				if err != nil {
					log.Println("5: ")

					r.Status = 500
					r.Message = err.Error()
					return r
				}

				fileName := strings.TrimSuffix(p.FileName(), ".csv")
				// Load Mexico's time zone
				location, err := time.LoadLocation("America/Mexico_City")
				if err != nil {
					log.Println("6: ")

					r.Status = 500
					r.Message = err.Error()
					return r
				}

				// Generate full filename with current date and time
				log.Println("7: ")

				now := time.Now().In(location) // Mexico Time
				filename := fmt.Sprintf("transactions/%s_%s_%s.csv", fileName, now.Format("02012006"), now.Format("030405PM"))
				fmt.Printf("Name of the file with the transactions: %s\n", filename)

				uploader := s3manager.NewUploader(sess)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   &readSeeker{buf},
				})

				if err != nil {
					log.Println("8: ")

					r.Status = 500
					r.Message = err.Error()
					return r
				}
			}
		} else {
			log.Println("9: ")

			r.Status = 400
			r.Message = "You must send a CSV file with the 'Content-Type' of type 'multipart/' in the Header."
			return r
		}
	}

	log.Println("10: fin ")

	r.Status = 200
	r.Message = "CSV file successfully uploaded."
	return r
}
