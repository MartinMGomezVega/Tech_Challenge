package commons

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
)

func ParseCSVContent(reader io.Reader) ([]models.Transaction, error) {
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

	return transactions, nil
}
