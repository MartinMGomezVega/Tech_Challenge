package commons

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/aws/aws-lambda-go/events"
)

// ParseCSVContent: reads the file from the request and returns the list of transactions.
func ParseCSVContent(request events.APIGatewayProxyRequest) ([]models.Transaction, error) {
	var transactions []models.Transaction

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		body, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			log.Println("Error processing multipart data.")
			return nil, err
		}

		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		p, err := mr.NextPart()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err != io.EOF {
			if p.FileName() != "" {
				// Create the buffer and copy the content
				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, p); err != nil {
					return nil, err
				}

				// Create a CSV reader for the buffer
				csvReader := csv.NewReader(buf)

				// Read the first line (header) to advance the reader
				_, err := csvReader.Read()
				if err != nil {
					log.Println("Error reading CSV header:", err)
					return nil, err
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
					id, err := strconv.Atoi(record[0])
					if err != nil {
						log.Println("Error parsing ID:", err)
						return nil, err
					}

					date := record[1]
					transaction := record[2]

					// parse the transaction and get the amount and the method
					amount, method, err := ParseTransaction(transaction)
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
			}
		}
	}

	return transactions, nil
}
