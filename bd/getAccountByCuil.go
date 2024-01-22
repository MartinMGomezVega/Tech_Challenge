package bd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAccountByCuil: Get full account by account number.
func GetAccountByCuil(cuil string) (models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("transactions")

	// Define el filtro para buscar por cuil
	filter := bson.M{"accountInfo.cuil": cuil}

	var account models.Account
	err := col.FindOne(ctx, filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No account was found for the cuil: %s", cuil)
			return models.Account{}, fmt.Errorf("no se encontró niNo account was found for the cuil:nguna cuenta para el cuil %s", cuil)
		}
		log.Printf("Error looking up the account for the %s quantile in MongoDB: %v", cuil, err)
		return models.Account{}, fmt.Errorf("error looking up the account for the %s quantile in MongoDB: %v", cuil, err)
	}

	return account, nil
}
