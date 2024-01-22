package bd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUser(cuil string) (models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("users")

	var user models.Account
	filter := bson.M{"cuil": cuil}
	err := col.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Printf("Error finding user in MongoDB for cuil %s: %v", cuil, err)
		return models.Account{}, fmt.Errorf("failed to find user: %v", err)
	}

	return user, nil
}
