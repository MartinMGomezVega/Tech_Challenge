package routers

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/MartinMGomezVega/Tech_Challenge/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func StoreTransactionsInDB(ctx context.Context) models.ResposeAPI {
	var r models.ResposeAPI
	r.Status = 400

	folderName := "files/"
	bucketName := aws.String(ctx.Value(models.Key("bucketName")).(string))

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	s3Client := s3.New(sess)

	// Listar objetos en el bucket y carpeta especificados
	input := &s3.ListObjectsV2Input{
		Bucket: bucketName,
		Prefix: aws.String(folderName),
	}

	result, err := s3Client.ListObjectsV2(input)
	if err != nil {
		r.Status = 400
		r.Message = "Error listing s3 bucket files."
		fmt.Println(r.Message)
		return r
	}

	// Iterar sobre los objetos y leer su contenido
	for _, item := range result.Contents {
		// Obtener el nombre del objeto
		objectKey := aws.StringValue(item.Key)

		// Leer el contenido del objeto
		content, err := readObjectContent(s3Client, aws.StringValue(bucketName), objectKey)
		if err != nil {
			fmt.Printf("Error reading object %s: %v\n", objectKey, err)
			continue
		}

		// Manejar el contenido leído según tus necesidades
		fmt.Printf("Content of %s:\n%s\n", objectKey, content)
	}

	r.Status = 200
	r.Message = "The update of the database was a success."
	fmt.Println(r.Message)
	return r
}

func readObjectContent(s3Client *s3.S3, bucketName, objectKey string) (string, error) {
	// Crear una solicitud para obtener el objeto
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	// Obtener el objeto
	result, err := s3Client.GetObject(input)
	if err != nil {
		return "", err
	}
	defer result.Body.Close()

	// Leer el contenido del objeto
	content, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
