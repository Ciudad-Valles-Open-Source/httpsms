# Primeros pasos

httpSMS convierte un teléfono Android en una **pasarela SMS**: tu aplicación hace una petición HTTP y el teléfono envía el SMS con su propia tarjeta SIM. También funciona al revés: cuando el teléfono recibe un SMS, httpSMS te lo reenvía a tu servidor.

Esta versión de httpSMS es **autoalojada y 100 % gratuita**: corre en tu propio servidor, no tiene planes de pago y no impone límites de mensajes, webhooks ni contactos. Lo único que limita cuántos SMS puedes enviar es tu operador y tu SIM.

En esta guía vas a conectar tu teléfono a tu servidor y enviar tu primer SMS en unos minutos.

## Cómo funciona

Antes de empezar conviene entender el recorrido de un mensaje, porque explica para qué sirve cada pieza que vas a configurar:

1. Tu aplicación llama a la API de tu servidor (`{{API_BASE_URL}}`) con tu clave API.
2. El servidor guarda el mensaje en una cola y avisa al teléfono mediante una notificación push.
3. La app de httpSMS en el teléfono descarga el mensaje y lo envía como SMS normal.
4. El teléfono informa al servidor si el SMS se envió, se entregó o falló.
5. Si configuraste [webhooks](/docs/webhooks-introduction), el servidor te avisa de cada uno de esos cambios en tiempo real.

## Antes de empezar

Necesitas:

- **Tu servidor httpSMS funcionando.** Si todavía no lo tienes, sigue la guía de autoalojamiento `docs/self-hosting_ES.md` del repositorio.
- **Una cuenta en el panel web** (`{{APP_URL}}`). Regístrate desde la página de inicio de sesión.
- **Un teléfono Android con SIM** que pueda enviar SMS. Será tu pasarela, así que conviene que esté siempre encendido, cargando y con conexión a internet.

## Paso 1: obtén tu clave API

Todas las peticiones a la API se autentican con una clave API que se envía en la cabecera HTTP `x-api-key`. Así el servidor sabe a qué cuenta pertenece cada petición.

1. Entra en el panel web y abre **Configuración** (`{{APP_URL}}/settings`).
2. En la sección **Clave API**, copia la clave con el botón de copiar.
3. Si vas a iniciar sesión en el teléfono, también puedes pulsar **Mostrar Código QR** para escanearla desde la app en lugar de escribirla.

Si una petición no lleva la clave, o la clave no es válida, la API responde `401 Unauthorized`:

```json
{
  "status": "error",
  "message": "You are not authorized to carry out this request.",
  "data": "Make sure your API key is set in the [x-api-key] header in the request"
}
```

> **Importante:** tu clave API da acceso completo a tu cuenta (enviar mensajes, leerlos, gestionar webhooks…). No la publiques ni la incluyas en código del navegador o de apps móviles. Si se filtra, usa **Rotar Clave API** en la misma página de configuración; la clave anterior dejará de funcionar.

> **Nota:** si tienes varios teléfonos, puedes crear una **clave API de teléfono** distinta para cada uno (empiezan por `pk_`). Solo sirven para que la app Android hable con el servidor y no permiten usar el resto de la API. Consulta [Claves API de teléfono](/docs/phone-api-keys).

## Paso 2: instala la app Android

1. En el teléfono, descarga el APK desde `{{APP_DOWNLOAD_URL}}` e instálalo. Android te pedirá permiso para instalar apps de origen desconocido.
2. Abre la app y concede el **permiso de SMS** (botón **Habilitar Permiso de SMS**).
3. Pulsa **Desactivar Optimización de Batería**. Sin esto, Android puede dormir la app y los mensajes se quedarán en cola sin enviarse.

> **Nota para quien administra el servidor:** la app recibe los avisos de "hay un mensaje nuevo" mediante Firebase Cloud Messaging de **tu** proyecto de Firebase. Por eso el APK debe compilarse con el `google-services.json` de tu proyecto (paso 8 de `docs/self-hosting_ES.md`). Un APK compilado para otro proyecto podrá iniciar sesión, pero no recibirá los mensajes a enviar.

## Paso 3: conecta el teléfono a tu servidor

La app no sabe dónde está tu servidor hasta que se lo dices. En la pantalla de inicio de sesión completa estos campos:

| Campo                         | Qué poner                                                                                                                        |
| ----------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| **API Key**                   | La clave del paso 1 (o una clave API de teléfono `pk_…`). Puedes pegarla o pulsar el icono de QR y escanear el código del panel. |
| **Número de teléfono (SIM1)** | El número de la SIM en formato internacional E.164, por ejemplo `+5215512345678`. La app intenta detectarlo sola.                |
| **Número de teléfono (SIM2)** | Solo aparece en teléfonos con doble SIM.                                                                                         |
| **URL del Servidor**          | La URL base de tu API: `{{API_BASE_URL}}` (sin `/v1` al final).                                                                  |

Pulsa **Iniciar Sesión con API Key**. Si todo va bien, el teléfono quedará registrado y aparecerá en la sección **Teléfonos** de `{{APP_URL}}/settings`.

> **Importante:** la app solo acepta URLs que empiecen por `https://`. Si tu servidor está en una red local sin certificado, ponlo detrás de un proxy con HTTPS (por ejemplo Caddy, Traefik o un túnel) antes de conectar el teléfono.

> **Nota:** puedes cambiar la URL del servidor más adelante desde **Configuración de la App** dentro de la propia app, por ejemplo si migras a otro dominio.

## Paso 4: envía tu primer SMS

Para enviar un SMS, haz una petición `POST` a `{{API_BASE_URL}}/v1/messages/send` con un cuerpo JSON.

### Parámetros

| Campo         | Tipo           | Obligatorio | Descripción                                                                                                                            |
| ------------- | -------------- | ----------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| `from`        | texto          | Sí          | Número del teléfono que enviará el SMS, en formato E.164. Debe ser un teléfono registrado en tu cuenta.                                |
| `to`          | texto          | Sí          | Número del destinatario.                                                                                                               |
| `content`     | texto          | Sí          | Texto del mensaje, **hasta 2048 caracteres**. Si supera lo que cabe en un SMS, el teléfono lo envía como mensaje multiparte.           |
| `request_id`  | texto          | No          | Un identificador tuyo (máx. 255 caracteres). Te llega de vuelta en los webhooks, útil para relacionar el mensaje con tu base de datos. |
| `send_at`     | fecha ISO 8601 | No          | Programa el envío para más adelante, hasta 20 días (480 horas). Ver [Programar SMS](/docs/scheduling-sms-messages).                    |
| `encrypted`   | booleano       | No          | Indica que `content` va cifrado de extremo a extremo. Ver la sección "Cifrado de extremo a extremo" más abajo.                         |
| `attachments` | lista de URLs  | No          | Hasta 10 URLs de archivos. Si las incluyes, el mensaje se envía como MMS.                                                              |

### Ejemplos

Sustituye `TU_CLAVE_API` por tu clave y los números por los tuyos.

```bash
curl -X POST '{{API_BASE_URL}}/v1/messages/send' \
  -H 'x-api-key: TU_CLAVE_API' \
  -H 'Content-Type: application/json' \
  -d '{
    "from": "+5215512345678",
    "to": "+5215587654321",
    "content": "Hola, este es mi primer SMS con httpSMS"
  }'
```

```javascript
const response = await fetch('{{API_BASE_URL}}/v1/messages/send', {
  method: 'POST',
  headers: {
    'x-api-key': process.env.HTTPSMS_API_KEY,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    from: '+5215512345678',
    to: '+5215587654321',
    content: 'Hola, este es mi primer SMS con httpSMS',
  }),
})

console.log(await response.json())
```

```python
import os
import requests

response = requests.post(
    "{{API_BASE_URL}}/v1/messages/send",
    headers={"x-api-key": os.environ["HTTPSMS_API_KEY"]},
    json={
        "from": "+5215512345678",
        "to": "+5215587654321",
        "content": "Hola, este es mi primer SMS con httpSMS",
    },
    timeout=10,
)

print(response.status_code, response.json())
```

```php
<?php
// Requiere Guzzle: composer require guzzlehttp/guzzle
$client = new GuzzleHttp\Client();

$response = $client->post('{{API_BASE_URL}}/v1/messages/send', [
    'headers' => ['x-api-key' => getenv('HTTPSMS_API_KEY')],
    'json' => [
        'from' => '+5215512345678',
        'to' => '+5215587654321',
        'content' => 'Hola, este es mi primer SMS con httpSMS',
    ],
]);

echo $response->getBody();
```

Si usas Go, el SDK oficial `github.com/NdoleStudio/httpsms-go` funciona con tu servidor indicándole la URL base:

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/NdoleStudio/httpsms-go"
)

func main() {
	client := httpsms.New(
		httpsms.WithBaseURL("{{API_BASE_URL}}"),
		httpsms.WithAPIKey(os.Getenv("HTTPSMS_API_KEY")),
	)

	message, _, err := client.Messages.Send(context.Background(), &httpsms.MessageSendParams{
		From:    "+5215512345678",
		To:      "+5215587654321",
		Content: "Hola, este es mi primer SMS con httpSMS",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(message.Data.ID)
}
```

### Respuesta

Una respuesta `200` significa que el mensaje **entró en la cola**, no que ya se haya enviado. El envío real ocurre unos segundos después, cuando el teléfono lo recoge:

```json
{
  "status": "success",
  "message": "message added to queue",
  "data": {
    "id": "32343a19-da5e-4b1b-a767-3298a73703cb",
    "request_id": null,
    "owner": "+5215512345678",
    "contact": "+5215587654321",
    "content": "Hola, este es mi primer SMS con httpSMS",
    "encrypted": false,
    "type": "mobile-terminated",
    "status": "pending",
    "sim": "SIM1",
    "send_attempt_count": 0,
    "max_send_attempts": 2,
    "created_at": "2026-01-15T10:20:30.123456Z"
  }
}
```

(La respuesta real incluye más campos con fechas, como `sent_at` o `delivered_at`, que se rellenan a medida que avanza el mensaje.)

Guarda el `id`: con él puedes consultar el mensaje en `GET {{API_BASE_URL}}/v1/messages/{id}`. El campo `status` irá pasando por `pending` → `sending` → `sent` → `delivered`, o terminará en `failed` o `expired` si algo sale mal.

### Errores habituales

| Código | Qué significa                                                  | Qué hacer                                                                                                                               |
| ------ | -------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------- |
| `400`  | El cuerpo no es JSON válido.                                   | Revisa el JSON y la cabecera `Content-Type: application/json`.                                                                          |
| `401`  | Falta la clave API o no es válida.                             | Revisa la cabecera `x-api-key`.                                                                                                         |
| `422`  | Algún campo no pasó la validación; el detalle viene en `data`. | Lo más común: `from` no es un teléfono registrado en tu cuenta, un número no está en formato E.164, o `content` supera 2048 caracteres. |
| `500`  | Error interno del servidor.                                    | Revisa los logs de tu API.                                                                                                              |

## Paso 5: recibe SMS

Cuando alguien escribe al número de tu teléfono, la app reenvía el SMS a tu servidor y aparece en el panel web. Para que llegue también a **tu aplicación**, tienes dos opciones:

- **Webhooks (recomendado):** el servidor te hace un `POST` en tiempo real con el evento `message.phone.received`. Ver [Webhooks](/docs/webhooks-introduction).
- **Consultar la API:** lista la conversación entre dos números con `GET {{API_BASE_URL}}/v1/messages?owner=+5215512345678&contact=+5215587654321`.

## Cifrado de extremo a extremo

Si no quieres que el contenido de los mensajes quede legible en tu base de datos (por ejemplo, porque el servidor lo administra otra persona), puedes cifrarlo. El cifrado y descifrado ocurre solo en tu aplicación y en el teléfono; el servidor únicamente guarda y reenvía texto cifrado.

1. En la app Android, abre **Configuración de la App** y escribe una **Clave de Encriptación** (una contraseña que solo conocen tu aplicación y el teléfono).
2. Para **enviar** mensajes cifrados: cifra `content` en tu aplicación y envía `"encrypted": true`. El teléfono lo descifra justo antes de mandar el SMS.
3. Para **recibir** mensajes cifrados: activa **Encriptar mensajes recibidos** en la app. Los eventos `message.phone.received` llegarán con `"encrypted": true` y el `content` cifrado.

El formato es AES-256 en modo CFB: la clave AES es el SHA-256 de tu clave de encriptación, se genera un IV aleatorio de 16 bytes, y el resultado es `base64(IV + texto cifrado)`. Ejemplo en Python (requiere `pip install cryptography`):

```python
import base64
import hashlib
import os

from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes


def _key(encryption_key: str) -> bytes:
    return hashlib.sha256(encryption_key.encode()).digest()


def encrypt(encryption_key: str, text: str) -> str:
    iv = os.urandom(16)
    encryptor = Cipher(algorithms.AES(_key(encryption_key)), modes.CFB(iv)).encryptor()
    data = encryptor.update(text.encode()) + encryptor.finalize()
    return base64.b64encode(iv + data).decode()


def decrypt(encryption_key: str, payload: str) -> str:
    raw = base64.b64decode(payload)
    decryptor = Cipher(algorithms.AES(_key(encryption_key)), modes.CFB(raw[:16])).decryptor()
    return (decryptor.update(raw[16:]) + decryptor.finalize()).decode()
```

> **Nota:** el límite de 2048 caracteres se aplica al texto **ya cifrado**, que es más largo que el original (IV + Base64 añaden aproximadamente un tercio). Si el teléfono recibe un mensaje cifrado sin tener clave configurada, lo marcará como fallido.

## Siguientes pasos

- [Webhooks](/docs/webhooks-introduction): recibe eventos en tu servidor en tiempo real.
- [Eventos de webhooks](/docs/webhooks-events): lista completa de eventos y sus datos.
- [Controlar la velocidad de envío](/docs/control-sms-send-rate): evita que tu operador bloquee la SIM por enviar demasiado rápido.
- [Cola de mensajes salientes](/docs/outgoing-message-queue): qué pasa con los mensajes mientras esperan.
- Referencia completa de la API (Swagger): `{{API_DOCS_URL}}`.
