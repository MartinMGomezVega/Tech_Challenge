package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"strings"

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

	var filename string

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))

	filename = "files/" + "stori_logo" + ".png"

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		body, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
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
		}

	} else {
		r.Status = 400
		r.Message = "Debe enviar una imagen con el 'Content-Type' de tipo 'multipart/' en el Header"
		return r
	}

	r.Status = 200
	r.Message = "Image Upload OK"
	return r
}
