package routers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MartinMGomezVega/Tech_Challenge/bd"
	"github.com/MartinMGomezVega/Tech_Challenge/commons"
	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/gomail.v2"
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
	log.Printf("email: %s", account.AccountInfo.Email)
	log.Printf("Number of transactions: %v", len(account.Transactions))

	// Calculate the total number of transactions for each month.
	totalBalance := commons.CalculateTotalBalance(account.Transactions)
	log.Printf("Total balance: $%v", totalBalance)

	// Set up the email
	d := gomail.NewDialer("smtp.gmail.com", 587, "valkiria.jobs@gmail.com", "zzmp qkxj nmas kubm")

	bodyEmail := fmt.Sprintf("¡Hola %s!<br>", account.AccountInfo.Name)
	bodyEmail += "Ya está disponible el resumen de tu cuenta.<br><br>"
	bodyEmail += fmt.Sprintf("Saldo total: $%v<br>", totalBalance)

	// Implement a cycle to iterate over months
	transactionsByMonth := commons.CalculateTotalTransactionsByMonth(account.Transactions)
	for month, qtyTransactions := range transactionsByMonth {
		monthEsp, err := commons.GetMonthInSpanish(month)
		if err != nil {
			log.Println("Error getting the month in Spanish: ", err.Error())
			r.Status = 400
			r.Message = "Error DialAndSend: "
			return r
		}
		bodyEmail += fmt.Sprintf("Número de transacciones en %s: %d<br>", monthEsp, qtyTransactions)
	}

	// info extra
	bodyEmail += "<br><br><br>¡Si surgen preguntas, estamos aquí para ayudarte! <br>Explora nuestras Preguntas Frecuentes en <a href='http://www.storicard.com/preguntas-frecuentes'>www.storicard.com/preguntas-frecuentes</a> para obtener respuestas rápidas y útiles. ¡Tu tranquilidad es nuestra prioridad!<br>"
	bodyEmail += "Conoce más sobre nosotros en nuestro sitio web: <a href='http://www.storicard.com'>www.storicard.com</a><br><br>"

	// Add the image to the background of the email
	bodyEmail += `<br><img src="https://upload.wikimedia.org/wikipedia/commons/b/b0/Stori_Logo_2023.svg" alt="stori" style="width: 30%; height: auto;">`

	// Email subject
	subject := "Stori - Resumen"

	m := gomail.NewMessage()
	m.SetHeader("From", "valkiria.jobs@gmail.com")
	m.SetHeader("To", account.AccountInfo.Email) // Se le envía al usuario
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", bodyEmail)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error DialAndSend: ", err.Error())
		r.Status = 400
		r.Message = "Error DialAndSend: "
		return r
	}

	r.Status = 200
	r.Message = "Email successfully sent."
	return r

}
