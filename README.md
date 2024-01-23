# Tech_Challenge

## Descripción del Proyecto

Este proyecto consiste en una aplicación que procesa transacciones financieras y envía un resumen por correo electrónico a los usuarios. La aplicación utiliza AWS con API Gateway, Lambda y S3 para gestionar la infraestructura y ejecutar el código de procesamiento.

## Instrucciones de Ejecución

### Paso 1: Crear el usuario

- Para ejecutar el código, sigue estos pasos:

1. Accede a la API Gateway mediante Postman con el siguiente enlace: https://tnq2inp3d2.execute-api.us-east-1.amazonaws.com/prod/createUser
2. En el cuerpo de la solicitud, proporciona un JSON con la información del usuario:

```json
{
  "name": "Martin",
  "surname": "Gomez Vega",
  "email": "martin_gomezvega@hotmail.com",
  "cuil": "20417027050"
}
```

Este endpoint crea un usuario en la colección de MongoDB llamada "users". El CUIL se utiliza como identificador y se vincula al archivo que se cargará más adelante.

### Paso 2: Cargar Archivo de Transacciones

Para cargar el archivo de transacciones, sigue estos pasos:

1. Accede a la API Gateway mediante Postman utilizando el siguiente enlace: https://tnq2inp3d2.execute-api.us-east-1.amazonaws.com/prod/uploadTransactionFile

2. Selecciona la opción de `form-data` en el cuerpo de la solicitud para permitir la carga de archivos.

3. Adjunta el archivo CSV que contiene las transacciones. Asegúrate de que el archivo tenga el formato adecuado, con columnas como "Id", "Date", y "Transaction" (por ejemplo).

4. Envía la solicitud.

Este proceso cargará el archivo de transacciones en el sistema. Internamente, el sistema almacenará el archivo en el bucket S3 llamado "challenge" y registrará los detalles de las transacciones en la colección "transactions".
Una vez completado este paso, estarás listo para proceder al Paso 3 y obtener el resumen de transacciones.

### Paso 3: Envío del Resumen de Cuenta por Correo Electrónico

Para recibir el resumen de cuenta por correo electrónico, sigue estos pasos:

1. Accede a la API Gateway mediante Postman utilizando el siguiente enlace: https://tnq2inp3d2.execute-api.us-east-1.amazonaws.com/prod/sendEmail

2. En el cuerpo de la solicitud, proporciona un JSON con tu número de CUIL:

```json
{
  "cuil": "20417027050"
}
```

3. Envia la solicitud

Este proceso obtendrá la información del usuario asociada al cuil proporcionado. Luego, calculará el saldo total, el total de transacciones agrupadas por mes, el importe medio de débito y crédito. A continuación, armará un correo electrónico detallado y lo enviará desde la dirección de correo electrónico valkiria.jobs@gmail.com.
El correo electrónico contendrá información adicional para mejorar la visualización, incluyendo la imagen de Stori y el archivo que cargaste en el Paso 2. Utilice mi correo personal que cree para mi tesis de Ingeniería informática.

¡Con estos pasos, habrás completado el proceso! Si tienes alguna pregunta o encuentras algún problema, no dudes en contactarme.
