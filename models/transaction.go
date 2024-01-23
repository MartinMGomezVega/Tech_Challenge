package models

// Transaction: represents the information of a transaction
type Transaction struct {
	ID     int     `json:"id" bson:"id"`
	Date   string  `json:"date" bson:"date"`
	Amount float64 `json:"amount" bson:"amount"`
	Method string  `json:"method" bson:"method"`
}

// AccountInfo: user account
type AccountInfo struct {
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
	Cuil    string `json:"cuil" bson:"cuil"`
	Email   string `json:"email" bson:"email"`
}

// Account: represents account information with a list of transactions.
type Account struct {
	AccountInfo  AccountInfo   `json:"accountInfo" bson:"accountInfo"`
	Transactions []Transaction `json:"transactions" bson:"transactions"`
}
