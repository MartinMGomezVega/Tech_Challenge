package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"strings"

	"github.com/MartinMGomezVega/Tech_Challenge/bd"
	"github.com/MartinMGomezVega/Tech_Challenge/commons"
	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"go.mongodb.org/mongo-driver/bson"

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

// UploadTransactionFile: Save the csv file in the s3 bucket of aws
func UploadTransactionFile(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {
	log.Println("Saving file...")
	var r models.ResposeAPI
	r.Status = 400

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		log.Println("Error in parsing content type.")
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		body, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			log.Println("Error processing multipart data.")
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
				uploadedFilename := p.FileName()
				log.Printf("Uploaded Filename: %s", uploadedFilename)

				filename := fmt.Sprintf("files/%s", uploadedFilename)
				log.Printf("filename: %s", filename)

				// Get the filename quantile without the .csv extension
				cuil := commons.GetCuilFromFilename(uploadedFilename)
				log.Printf("cuil: %s", cuil)

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
					log.Println("Error logging into aws.")
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				bucket := ctx.Value(models.Key("bucketName")).(string)
				log.Printf("bucket: %s", bucket)

				// Upload the csv file to Bucket S3
				log.Println("Before uploading to S3")
				uploader := s3manager.NewUploader(sess)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(filename),
					Body:   &readSeeker{buf},
				})

				if err != nil {
					log.Printf("Error uploading the file to the bucket: %s", bucket)
					r.Status = 500
					r.Message = err.Error()
					return r
				}
				log.Println("After uploading to S3")

				// Define the filter to search by account number
				filter := bson.M{"cuil": cuil}
				collection := "users"
				// Obtain user account information to store in the transaction collection
				user, err := bd.GetAccountByCuil(cuil, filter, collection)
				if err != nil {
					r.Status = 400
					r.Message = err.Error()
					return r
				}
				log.Printf("User: %s %s", user.AccountInfo.Name, user.AccountInfo.Surname)

				// Parse the contents of the CSV file
				transactions, err := commons.ParseCSVContent(buf)
				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}
				log.Printf("Number of transactions: %v", len(transactions))

				// Create an Account document with account and transaction information
				account := models.Account{
					AccountInfo:  user.AccountInfo,
					Transactions: transactions,
				}

				// After uploading to S3, store in MongoDB
				_, status, err := bd.StoreInMongoDB(account)
				if err != nil {
					log.Println("Error storing in MongoDB: ", err)
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				log.Printf("Upload status of the file to MongoDB: %v", status)
				if !status {
					r.Message = "Failed to insert user record."
					log.Println(r.Message)
					return r
				}
			}
		}

	} else {
		r.Status = 400
		r.Message = "You must send a csv with the 'Content-Type' of type 'multipart/' in the Header."
		return r
	}

	r.Status = 200
	r.Message = "csv upload OK"
	return r
}
