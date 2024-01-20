package handlers

import (
	"context"
	"fmt"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/MartinMGomezVega/Tech_Challenge/routers"
	"github.com/aws/aws-lambda-go/events"
)

func Handlers(ctx context.Context, request events.APIGatewayProxyRequest) models.ResposeAPI {

	fmt.Println("Processing: " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))

	var r models.ResposeAPI
	r.Status = 400

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		case "uploadTransactionFile":
			return routers.UploadTransactionFile(ctx)

		case "sendEmail":
			return routers.SendEmail(ctx)

		case "storeTransactionsInDB":
			return routers.StoreTransactionsInDB(ctx)

		case "createUser":
			return routers.CreateUser(ctx)
		}

	case "GET":
		switch ctx.Value(models.Key("path")).(string) {
		// case "verperfil": // listo
		// 	return routers.VerPerfil(request)
		// case "leoTweets": // listo
		// 	return routers.LeoTweets(request)
		// case "consultaRelacion": // listo
		// 	return routers.ConsultaRelacion(request, claim)
		// case "listaUsuarios": // listo
		// 	return routers.ListaUsuarios(request, claim)
		// case "leoTweetsSeguidores": // listo
		// 	return routers.LeoTweetsSeguidores(request, claim)
		// case "obtenerAvatar": // listo
		// 	return routers.ObtenerImagen(ctx, "A", request, claim)
		// case "obtenerBanner": // listo
		// 	return routers.ObtenerImagen(ctx, "B", request, claim)
		}

	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {
		// case "modificarPerfil": // listo
		// 	return routers.ModificarPerfil(ctx, claim)
		}

	case "DELETE":
		//
	}

	r.Status = 400
	r.Message = "Method Invalid"
	return r
}
