package commons

import (
	"fmt"
	"log"
	"math"
	"strings"
)

// ParseTransaction: parse the transaction and get the amount and the method
func ParseTransaction(transaction string) (float64, string, error) {
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
