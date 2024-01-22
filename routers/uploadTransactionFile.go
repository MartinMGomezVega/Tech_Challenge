package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/MartinMGomezVega/Tech_Challenge/bd"
	"github.com/MartinMGomezVega/Tech_Challenge/commons"
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
	log.Println("bucket: " + *bucket)

	filename := "files/" + "20417027050" + ".csv"
	log.Println("filename: " + filename)

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
				log.Println("Uploaded Filename: " + uploadedFilename)

				// Get the filename quantile without the .csv extension
				cuil := commons.GetCuilFromFilename(uploadedFilename)

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

				log.Println("Before uploading to S3")
				uploader := s3manager.NewUploader(sess)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   &readSeeker{buf},
				})
				log.Println("After uploading to S3")

				if err != nil {
					log.Println("Error uploading the file to the bucket: " + *bucket)
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				// Get the email
				user, err := bd.GetUser(cuil)
				if err != nil {
					r.Status = 400
					r.Message = err.Error()
					return r
				}

				// Parsear el contenido del archivo CSV
				transactions, err := parseCSVContent(buf)
				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				// Crear un documento Account con la información de la cuenta y las transacciones
				account := models.Account{
					AccountInfo:  user.AccountInfo,
					Transactions: transactions,
				}

				// After uploading to S3, store in MongoDB
				_, status, err := bd.StoreInMongoDB(account)
				if err != nil {
					log.Println("Error storing in MongoDB:", err)
					r.Status = 500
					r.Message = err.Error()
					return r
				}

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

func parseCSVContent(reader io.Reader) ([]models.Transaction, error) {
	var transactions []models.Transaction

	// Create a CSV reader
	csvReader := csv.NewReader(reader)

	// Read the first line to get column names
	columns, err := csvReader.Read()
	if err != nil {
		log.Println("Error reading CSV header:", err)
		return nil, err
	}

	// Map column names to indexes
	columnIndex := map[string]int{
		"Id":          -1,
		"Date":        -1,
		"Transaction": -1,
	}

	for i, colName := range columns {
		columnIndex[colName] = i
	}

	// Read the remaining rows of the CSV
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error reading CSV:", err)
			return nil, err
		}

		// Get column values by name
		id, err := strconv.Atoi(record[columnIndex["Id"]])
		if err != nil {
			log.Println("Error parsing ID:", err)
			return nil, err
		}

		date := record[columnIndex["Date"]]
		transaction := record[columnIndex["Transaction"]]

		// parse the transaction and get the amount and the method
		amount, method, err := parseTransaction(transaction)
		if err != nil {
			log.Println("Error parsing Transaction:", err)
			return nil, err
		}

		// Create a transaction and add it to the list
		transactionObj := models.Transaction{
			ID:     id,
			Date:   date,
			Amount: amount,
			Method: method,
		}
		transactions = append(transactions, transactionObj)
	}

	return transactions, nil
}

// parseTransaction: parse the transaction and get the amount and the method
func parseTransaction(transaction string) (float64, string, error) {
	var amount float64
	var method string

	// Get the amount of the transaction
	if _, err := fmt.Sscanf(transaction, "%f", &amount); err != nil {
		log.Println("Error parsing Transaction:", err)
		return 0, "", err
	}

	// All amounts must be positive
	amount = math.Abs(amount)

	// Determine the method(Credit + o Debit -negativo)
	if strings.HasPrefix(transaction, "-") {
		method = "-"
	} else {
		method = "+"
	}

	return amount, method, nil
}
