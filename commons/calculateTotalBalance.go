package commons

import "github.com/MartinMGomezVega/Tech_Challenge/models"

// CalculateTotalBalance: Calculates the total balance of debits and credits.
func CalculateTotalBalance(transactions []models.Transaction) float64 {
	totalBalance := 0.0

	for _, transaction := range transactions {
		if transaction.Method == "+" {
			totalBalance += transaction.Amount
		} else if transaction.Method == "-" {
			totalBalance -= transaction.Amount
		}
	}

	return totalBalance
}
