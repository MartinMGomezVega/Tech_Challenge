package routers

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	fmt.Println("Saving file...")
	var r models.ResposeAPI
	r.Status = 400

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))

	// Get filename from the body of the request
	fmt.Println("request.Body:", request.Body)
	fileName := strings.TrimSuffix(request.Body, ".csv")

	// Load Mexico's time zone
	location, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	// Generate full filename with current date and time
	now := time.Now().In(location) // Mexico Time
	filename := fmt.Sprintf("transactions/%s_%s_%s.csv", fileName, now.Format("02012006"), now.Format("030405PM"))
	fmt.Printf("Name of the file with the transactions: %s\n", filename)

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(strings.NewReader(request.Body), params["boundary"])
		p, err := mr.NextPart()
		if err != nil && err != io.EOF {
			r.Status = 500
			r.Message = err.Error()
			return r
		}
		if err != io.EOF {
			if p.FileName() != "" {
				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, p); err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				sess, err := session.NewSession(&aws.Config{
					Region: aws.String("us-east-1")},
				)

				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				uploader := s3manager.NewUploader(sess)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   &readSeeker{buf},
				})

				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}
			}
		} else {
			r.Status = 400
			r.Message = "You must send a CSV file with the 'Content-Type' of type 'multipart/' in the Header."
			return r
		}
	}

	r.Status = 200
	r.Message = "CSV file successfully uploaded."
	return r
}
