package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	// "github.com/ptilotta/twittor/jwt"
	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/MartinMGomezVega/Tech_Challenge/routers"
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
			// case "tweet": // listo
			// 	return routers.GraboTweet(ctx, claim)
			// case "altaRelacion": // listo
			// 	return routers.AltaRelacion(ctx, request, claim)
			// case "subirAvatar": // listo
			// 	return routers.UploadImage(ctx, "A", request, claim)
			// case "subirBanner": // listo
			// 	return routers.UploadImage(ctx, "B", request, claim)
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
