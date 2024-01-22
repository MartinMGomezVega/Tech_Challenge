package bd

import (
	"context"
	"time"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StoreInMongoDB: Stores an Account document in MongoDB
func StoreInMongoDB(account models.Account) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("transactions")

	result, err := col.InsertOne(ctx, account)
	if err != nil {
		return "", false, err
	}

	// Capture the id
	ObjID, _ := result.InsertedID.(primitive.ObjectID)

	return ObjID.String(), true, nil
}
