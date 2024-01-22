package commons

import (
	"fmt"
)

// GetMonthInSpanish: converts the number of the month to its English name.
func GetMonthInSpanish(month string) (string, error) {
	// Convert the number of the month to an integer
	var monthNumber int
	_, err := fmt.Sscanf(month, "%d", &monthNumber)
	if err != nil {
		return "", fmt.Errorf("the number of the month c: %v", err)
	}

	// Validate that the month number is in the correct range.
	if monthNumber < 1 || monthNumber > 12 {
		return "", fmt.Errorf("month number out of range: %d", monthNumber)
	}

	// Create a slice of Spanish month names
	monthNames := []string{
		"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
		"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
	}

	// Get the name of the month in English
	translatedMonth := monthNames[monthNumber-1]

	return translatedMonth, nil
}
