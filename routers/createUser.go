package routers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MartinMGomezVega/Tech_Challenge/bd"
	"github.com/MartinMGomezVega/Tech_Challenge/models"
)

func CreateUser(ctx context.Context) models.ResposeAPI {
	var u models.User
	var r models.ResposeAPI
	r.Status = 400

	fmt.Println("Creating the user...")

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &u)
	if err != nil {
		r.Message = err.Error()
		fmt.Println(r.Message)
		return r
	}

	if len(u.Email) == 0 {
		r.Message = "User email is required."
		fmt.Println(r.Message)
		return r
	}
	if len(u.Cuil) == 0 {
		r.Message = "The cuil is required."
		fmt.Println(r.Message)
		return r
	}

	_, encontrado, _ := bd.CheckUserExists(u.Email)
	if encontrado {
		r.Message = "The user already exists in the DB."
		fmt.Println(r.Message)
		return r
	}

	_, status, err := bd.InsertUser(u)
	if err != nil {
		r.Message = "An error occurred while trying to perform user registration. Error: " + err.Error()
		fmt.Println(r.Message)
		return r
	}

	if !status {
		r.Message = "Failed to insert user record."
		fmt.Println(r.Message)
		return r
	}

	r.Status = 200
	r.Message = "Registro OK"
	fmt.Println(r.Message)
	return r

}
