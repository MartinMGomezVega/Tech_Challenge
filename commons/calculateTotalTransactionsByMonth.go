package commons

import (
	"log"
	"strings"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
)

// CalculateTotalTransactionsByMonth: Calculate the total number of transactions for each month.
func CalculateTotalTransactionsByMonth(transactions []models.Transaction) map[string]int {
	// Create a map to store total transactions per month
	transactionsByMonth := make(map[string]int)

	// Iterate on transactions
	for _, transaction := range transactions {
		// Parse the transaction date
		dateComponents := strings.Split(transaction.Date, "/")
		if len(dateComponents) != 2 {
			log.Printf("Invalid date format in transaction: %v", transaction)
			continue
		}

		month := dateComponents[0]

		// Increment transaction counter for the corresponding month
		transactionsByMonth[month]++
	}

	return transactionsByMonth
}
