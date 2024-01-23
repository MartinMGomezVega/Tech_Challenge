package commons

import (
	"github.com/MartinMGomezVega/Tech_Challenge/models"
)

// CalculateAverageAmounts: Calculates average debit and credit amounts per month.
func CalculateAverageAmounts(transactions []models.Transaction) map[string]float64 {
	averageAmounts := make(map[string]float64)

	// debit
	sumDebit := 0.0
	countDebit := 0

	// credit
	sumCredit := 0.0
	countCredit := 0

	// Process each transaction
	for _, t := range transactions {
		// Determine whether the transaction is a debit or credit based on the Method field
		if t.Method == "-" {
			sumDebit += t.Amount
			countDebit++
		} else if t.Method == "+" {
			sumCredit += t.Amount
			countCredit++
		}
	}

	// Calculating average debit and credit amounts
	if countDebit > 0 {
		averageAmounts["Debit"] = sumDebit / float64(countDebit)
	}

	if countCredit > 0 {
		averageAmounts["Credit"] = sumCredit / float64(countCredit)
	}

	return averageAmounts
}
