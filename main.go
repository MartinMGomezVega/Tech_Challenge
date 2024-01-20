package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/MartinMGomezVega/Tech_Challenge/awsgo"
	"github.com/MartinMGomezVega/Tech_Challenge/bd"
	"github.com/MartinMGomezVega/Tech_Challenge/handlers"
	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/MartinMGomezVega/Tech_Challenge/secretmanager"
)

func ExecuteLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse
	awsgo.InitialiseAWS()

	if !ValidateParameter() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error in environment variables. must include 'SecretName', 'BucketName",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error in reading Secret " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	path := strings.Replace(request.PathParameters["techchallengego"], os.Getenv("UrlPrefix"), "", -1)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	if ctx.Value(models.Key("path")) != "uploadTransactionFile" {
		fmt.Println("Conectando a la DB...")
		// Chequeo Conexión a la BD o Conecto la BD
		err = bd.ConectBD(awsgo.Ctx)
		if err != nil {
			res = &events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Error connecting to DB: " + err.Error(),
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}
			return res, nil
		}
	}

	respAPI := handlers.Handlers(awsgo.Ctx, request)

	if respAPI.CustomResp == nil {
		headersResp := map[string]string{
			"Content-Type": "application/json",
		}
		res = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status,
			Body:       string(respAPI.Message),
			Headers:    headersResp,
		}

		return res, nil
	} else {
		return respAPI.CustomResp, nil
	}
}

func main() {
	lambda.Start(ExecuteLambda)
}

func ValidateParameter() bool {
	_, parameter := os.LookupEnv("SecretName")
	if !parameter {
		return parameter
	}

	_, parameter = os.LookupEnv("BucketName")
	if !parameter {
		return parameter
	}

	_, parameter = os.LookupEnv("UrlPrefix")
	if !parameter {
		return parameter
	}

	return parameter
}
