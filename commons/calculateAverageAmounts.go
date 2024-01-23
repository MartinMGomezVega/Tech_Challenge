package commons

import (
	"fmt"
	"strings"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
)

// CalculateAverageAmounts: Calcula los importes medios de débito y crédito por mes.
func CalculateAverageAmounts(transactions []models.Transaction) map[string]float64 {
	averageAmounts := make(map[string]float64)

	// Mapas para almacenar la suma y el conteo de transacciones por tipo (Débito o Crédito) y mes
	sumDebit := make(map[string]float64)
	countDebit := make(map[string]int)
	sumCredit := make(map[string]float64)
	countCredit := make(map[string]int)

	// Procesar cada transacción
	for _, t := range transactions {
		// Obtener el mes de la fecha
		dateComponents := strings.Split(t.Date, "/")
		month := dateComponents[0]

		// Determinar si la transacción es débito o crédito según el campo Method
		isDebit := t.Method == "Debit"

		// Actualizar la suma y el conteo según el tipo de transacción y mes
		if isDebit {
			sumDebit[month] += t.Amount
			countDebit[month]++
		} else {
			sumCredit[month] += t.Amount
			countCredit[month]++
		}
	}

	// Calcular los importes medios por tipo y mes
	for month := range sumDebit {
		if countDebit[month] > 0 {
			averageAmounts[fmt.Sprintf("Debit_%s", month)] = sumDebit[month] / float64(countDebit[month])
		}

		if countCredit[month] > 0 {
			averageAmounts[fmt.Sprintf("Credit_%s", month)] = sumCredit[month] / float64(countCredit[month])
		}
	}

	return averageAmounts
}
