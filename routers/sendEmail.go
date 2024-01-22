package routers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/MartinMGomezVega/Tech_Challenge/bd"
	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/aws/aws-lambda-go/events"
)

// SendEmail: Send the email with the summary information of the costs to the user.
func SendEmail(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {
	log.Println("Sending the email...")
	var r models.ResposeAPI
	r.Status = 400

	// Get the number from the body of the JSON request
	var requestBody map[string]string
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		log.Println("Error parsing request body:", err)
		r.Status = 400
		r.Message = "Error parsing request body"
		return r
	}

	cuil, ok := requestBody["cuil"]
	if !ok {
		log.Println("Cuil not provided in the request body")
		r.Status = 400
		r.Message = "Cuil not provided in the request body"
		return r
	}
	log.Printf("Cuil received: %s", cuil)

	// Get transactions from MongoDB
	account, err := bd.GetAccountByCuil(cuil)
	if err != nil {
		log.Println("Error getting transactions from MongoDB:", err)
		r.Status = 500
		r.Message = err.Error()
		return r
	}
	log.Printf("User: %s %s", account.AccountInfo.Name, account.AccountInfo.Surname)
	log.Printf("Number of transactions: %v", len(account.Transactions))

	// Armar el resumen del email

	r.Status = 200
	r.Message = "Email successfully sent."
	return r

}
