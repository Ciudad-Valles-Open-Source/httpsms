# Webhooks

Un webhook es una URL de **tu** aplicación a la que httpSMS hace un `POST` cada vez que ocurre algo: llega un SMS al teléfono, un mensaje se entrega, el teléfono se desconecta, etc.

¿Por qué usarlos? La alternativa sería preguntar a la API cada pocos segundos si hay novedades (polling), lo cual es lento y desperdicia recursos. Con webhooks te enteras al instante y solo cuando de verdad pasa algo.

## Qué necesitas para crear un webhook

| Campo                                     | Obligatorio               | Descripción                                                                                                                                                                            |
| ----------------------------------------- | ------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **URL de devolución de llamada** (`url`)  | Sí                        | La dirección de tu servidor que recibirá los `POST`. Máximo 255 caracteres. Si no indicas el esquema, se añade `https://`.                                                             |
| **Clave de firma** (`signing_key`)        | No, pero muy recomendable | Un secreto que solo conocen httpSMS y tu servidor. Se usa para firmar cada petición y que puedas comprobar que viene de tu servidor httpSMS y no de un tercero. Máximo 255 caracteres. |
| **Eventos** (`events`)                    | Sí                        | Los tipos de evento que quieres recibir (al menos uno). Consulta la [lista de eventos](/docs/webhooks-events).                                                                         |
| **Números de teléfono** (`phone_numbers`) | Sí                        | De qué teléfonos de tu cuenta quieres recibir eventos. Si tienes varios, puedes elegir solo algunos.                                                                                   |

Los eventos disponibles son:

| Evento                    | Cuándo se envía                                                     |
| ------------------------- | ------------------------------------------------------------------- |
| `message.phone.received`  | El teléfono recibió un SMS.                                         |
| `message.phone.sent`      | El teléfono envió un SMS.                                           |
| `message.phone.delivered` | El SMS llegó al teléfono del destinatario.                          |
| `message.send.failed`     | El teléfono no pudo enviar el SMS.                                  |
| `message.send.expired`    | El teléfono no envió el SMS a tiempo (por ejemplo, estaba apagado). |
| `message.call.missed`     | El teléfono recibió una llamada perdida.                            |
| `phone.heartbeat.offline` | El servidor lleva alrededor de una hora sin noticias del teléfono.  |
| `phone.heartbeat.online`  | El teléfono volvió a conectarse después de estar desconectado.      |

El detalle y el contenido de cada evento está en [Eventos de webhooks](/docs/webhooks-events).

> **Nota:** en esta versión autoalojada no hay límite de webhooks. Puedes crear tantos como necesites, por ejemplo uno para tu CRM y otro para tu sistema de alertas. Si varios webhooks están suscritos al mismo evento y teléfono, todos lo reciben.

## Crear un webhook desde el panel

1. Abre **Configuración** en `{{APP_URL}}/settings` y baja hasta la sección **Webhooks**.
2. Pulsa **Añadir webhook**.
3. Escribe la **URL de devolución de llamada**, por ejemplo `https://mi-app.ejemplo.com/webhooks/httpsms`.
4. Escribe una **Clave de firma**. Puedes generar una aleatoria con `openssl rand -hex 32`.
5. Elige los **Eventos** y los **Números de teléfono** que te interesan.
6. Pulsa **Guardar Webhook**.

Desde la misma tabla puedes editar o eliminar webhooks en cualquier momento.

## Crear un webhook con la API

También puedes gestionar webhooks desde código con estos endpoints (todos usan la cabecera `x-api-key`):

| Método y ruta                     | Para qué sirve                                      |
| --------------------------------- | --------------------------------------------------- |
| `GET /v1/webhooks`                | Listar tus webhooks.                                |
| `POST /v1/webhooks`               | Crear un webhook.                                   |
| `PUT /v1/webhooks/{webhookID}`    | Actualizar un webhook (se envían todos los campos). |
| `DELETE /v1/webhooks/{webhookID}` | Eliminar un webhook.                                |

```bash
curl -X POST '{{API_BASE_URL}}/v1/webhooks' \
  -H 'x-api-key: TU_CLAVE_API' \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://mi-app.ejemplo.com/webhooks/httpsms",
    "signing_key": "una-clave-secreta-larga-y-aleatoria",
    "events": ["message.phone.received", "message.phone.delivered", "message.send.failed"],
    "phone_numbers": ["+5215512345678"]
  }'
```

Los números de `phone_numbers` deben ser teléfonos ya registrados en tu cuenta (es decir, con la app Android conectada); si no, la API responde `422`. Puedes ver todos los campos en la referencia de la API: `{{API_DOCS_URL}}`.

## Cómo es la petición que recibirás

Cada vez que ocurre un evento al que estás suscrito, tu servidor httpSMS hace un `POST` a tu URL.

### Cabeceras

| Cabecera        | Valor                                                                                                           |
| --------------- | --------------------------------------------------------------------------------------------------------------- |
| `Content-Type`  | Siempre `application/json`.                                                                                     |
| `X-Event-Type`  | El tipo de evento, por ejemplo `message.phone.received`. Útil para decidir qué hacer sin leer el cuerpo.        |
| `Authorization` | `Bearer <token JWT>`. **Solo se envía si el webhook tiene clave de firma.** Se explica en la siguiente sección. |

### Cuerpo

El cuerpo es un evento en formato [CloudEvents 1.0](https://cloudevents.io/) serializado como JSON. Los campos comunes a todos los eventos son:

| Campo             | Descripción                                                                                                |
| ----------------- | ---------------------------------------------------------------------------------------------------------- |
| `id`              | Identificador único del evento (UUID). Úsalo para detectar duplicados.                                     |
| `type`            | El tipo de evento (igual que la cabecera `X-Event-Type`).                                                  |
| `source`          | La ruta de la API que originó el evento. Es informativa; no dependas de su valor.                          |
| `specversion`     | Siempre `1.0`.                                                                                             |
| `time`            | Momento en que se creó el evento (UTC).                                                                    |
| `datacontenttype` | Siempre `application/json`.                                                                                |
| `data`            | Los datos específicos del evento. Cambian según el tipo; ver [Eventos de webhooks](/docs/webhooks-events). |

Ejemplo de un SMS recibido:

```json
{
  "specversion": "1.0",
  "id": "f4aed1d3-ab4f-42b9-b9dd-9fc7182f7197",
  "source": "/v1/messages/receive",
  "type": "message.phone.received",
  "datacontenttype": "application/json",
  "time": "2026-01-15T10:20:31.524331882Z",
  "data": {
    "message_id": "0b0123bb-ef2e-468f-908a-c026d51636aa",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "owner": "+5215512345678",
    "encrypted": false,
    "contact": "+5215587654321",
    "timestamp": "2026-01-15T10:20:30.814Z",
    "content": "Hola, ¿a qué hora abren mañana?",
    "sim": "SIM1",
    "attachments": []
  }
}
```

## Verificar que la petición es auténtica

Tu URL de webhook es pública, así que cualquiera podría enviarle peticiones falsas. Para evitarlo, configura una clave de firma: httpSMS enviará en cada petición un **token JWT firmado con HMAC-SHA256 (`HS256`)** usando esa clave. Si la firma es válida, la petición viene de tu servidor httpSMS.

El token contiene estos datos (claims):

| Claim | Valor                                                                                                                                 |
| ----- | ------------------------------------------------------------------------------------------------------------------------------------- |
| `iss` | El emisor configurado en tu servidor: el valor de `JWT_ISSUER` o, si no está definido, un slug de `APP_NAME` (por defecto `httpsms`). |
| `aud` | La URL del webhook, tal como está guardada (por ejemplo `https://mi-app.ejemplo.com/webhooks/httpsms`).                               |
| `sub` | El ID de tu usuario en httpSMS.                                                                                                       |
| `iat` | Momento de creación del token.                                                                                                        |
| `nbf` | 10 minutos antes de la creación (margen para relojes desincronizados).                                                                |
| `exp` | 10 minutos después de la creación.                                                                                                    |

> **Sobre `iss`:** el emisor es el valor que el administrador del servidor define en la variable de entorno `JWT_ISSUER`. Si está vacía, la API lo deriva de `APP_NAME`: minúsculas, los espacios pasan a `_`, se quitan los acentos y solo se conservan `a-z`, `0-9`, `.`, `_` y `-` (por ejemplo `httpSMS` → `httpsms`, `Arquos SMS` → `arquos_sms`). Si lo compruebas, usa el valor configurado en **tu** servidor.

Para validarlo:

1. Lee la cabecera `Authorization` y quita el prefijo `Bearer `.
2. Verifica el token con una librería JWT, usando tu clave de firma y aceptando **solo** el algoritmo `HS256`.
3. Opcionalmente, comprueba que `aud` sea tu URL y que `iss` coincida con el emisor configurado en tu servidor.
4. Si algo falla, responde `401` y descarta la petición.

> **Importante:** el token firma la **identidad del remitente**, no el contenido del cuerpo. Usa siempre HTTPS en tu URL para que nadie pueda leer ni modificar la petición por el camino.

### Ejemplo con Node.js (Express)

```javascript
// npm install express jsonwebtoken
const express = require('express')
const jwt = require('jsonwebtoken')

const app = express()
app.use(express.json())

const SIGNING_KEY = process.env.HTTPSMS_WEBHOOK_SIGNING_KEY
// Debe coincidir con JWT_ISSUER de tu servidor (por defecto, slug de APP_NAME, p. ej. httpsms)
const ISSUER = process.env.HTTPSMS_JWT_ISSUER || 'httpsms'

app.post('/webhooks/httpsms', (req, res) => {
  const token = (req.get('Authorization') || '').replace(/^Bearer /, '')

  try {
    jwt.verify(token, SIGNING_KEY, {
      algorithms: ['HS256'],
      issuer: ISSUER,
    })
  } catch (error) {
    return res.status(401).send('firma no válida')
  }

  // Responde rápido y procesa después
  res.sendStatus(200)

  const event = req.body
  if (event.type === 'message.phone.received') {
    console.log(`SMS de ${event.data.contact}: ${event.data.content}`)
  }
})

app.listen(3000)
```

### Ejemplo con Python (Flask)

```python
# pip install flask pyjwt
import os

import jwt
from flask import Flask, abort, request

app = Flask(__name__)
SIGNING_KEY = os.environ["HTTPSMS_WEBHOOK_SIGNING_KEY"]
# Debe coincidir con JWT_ISSUER de tu servidor (por defecto, slug de APP_NAME, p. ej. httpsms)
ISSUER = os.environ.get("HTTPSMS_JWT_ISSUER", "httpsms")


@app.post("/webhooks/httpsms")
def httpsms_webhook():
    token = request.headers.get("Authorization", "").removeprefix("Bearer ")
    try:
        jwt.decode(
            token,
            SIGNING_KEY,
            algorithms=["HS256"],
            issuer=ISSUER,
            options={"verify_aud": False},
        )
    except jwt.InvalidTokenError:
        abort(401)

    event = request.get_json()
    if event["type"] == "message.phone.received":
        print(f"SMS de {event['data']['contact']}: {event['data']['content']}")

    return "", 200
```

### Ejemplo con Dart (Shelf)

```dart
// dart pub add shelf shelf_router dart_jsonwebtoken
import 'dart:convert';
import 'dart:io';

import 'package:dart_jsonwebtoken/dart_jsonwebtoken.dart';
import 'package:shelf/shelf.dart';
import 'package:shelf/shelf_io.dart' as io;
import 'package:shelf_router/shelf_router.dart';

final signingKey = Platform.environment['HTTPSMS_WEBHOOK_SIGNING_KEY']!;
// Debe coincidir con JWT_ISSUER de tu servidor (por defecto: slug de APP_NAME, p. ej. httpsms)
final issuer = Platform.environment['HTTPSMS_JWT_ISSUER'] ?? 'httpsms';

Future<Response> httpsmsWebhook(Request request) async {
  final token = (request.headers['authorization'] ?? '')
      .replaceFirst(RegExp('^Bearer '), '');

  try {
    final jwt = JWT.verify(token, SecretKey(signingKey), issuer: issuer);
    if (jwt.header?['alg'] != 'HS256') {
      return Response.unauthorized('invalid signature');
    }
  } on JWTException {
    return Response.unauthorized('invalid signature');
  }

  final event = jsonDecode(await request.readAsString()) as Map<String, dynamic>;
  if (event['type'] == 'message.phone.received') {
    final data = event['data'] as Map<String, dynamic>;
    print('SMS de ${data['contact']}: ${data['content']}');
  }

  return Response.ok('');
}

Future<void> main() async {
  final router = Router()..post('/webhooks/httpsms', httpsmsWebhook);
  await io.serve(router.call, InternetAddress.anyIPv4, 3000);
}
```

> **Nota:** en los ejemplos no se valida `aud` para simplificar. Si quieres validarlo, pasa tu URL exacta (`audience` en Node, `audience=` en PyJWT, `audience:` en Dart) y quita `verify_aud: False`.

## Respuestas, tiempos de espera y reintentos

Tu servidor debe responder con un código **2xx** (por ejemplo `200`) para indicar que recibió el evento. Esto es lo que hace httpSMS con cada entrega:

- **Tiempo de espera:** cada intento tiene un máximo de **5 segundos**. Si tu servidor tarda más, el intento cuenta como fallido.
- **Reintentos:** si el primer intento falla (error de red, tiempo agotado o cualquier código **400 o superior**, incluidos los 4xx), httpSMS lo reintenta **una sola vez** tras una pausa breve. En total, como máximo **2 intentos** por evento y webhook.
- **Si ambos intentos fallan,** el evento se descarta (no hay más reintentos). Si tu servidor tiene configurado el envío de correos y tienes activado **Correos de webhooks y discord** en la configuración, recibirás un correo avisándote del fallo (como máximo uno por hora y por teléfono, para no llenarte la bandeja).

Recomendaciones:

- **Responde rápido.** Devuelve `200` en cuanto valides la firma y procesa el evento en segundo plano (una cola, un job, etc.).
- **Sé idempotente.** Si tu servidor tardó en responder pero sí procesó el evento, el reintento te lo enviará otra vez. Guarda el `id` del evento y descarta los que ya procesaste.
- **No dependas del orden.** Los eventos pueden llegar desordenados (por ejemplo, `delivered` antes que `sent`). Usa el campo `timestamp` de `data` para ordenarlos.
- **Evita redirecciones.** Configura la URL final; una redirección puede hacer que la petición pierda el método `POST` o el cuerpo.
- **Si el evento es crítico,** no dependas solo del webhook: puedes consultar el estado del mensaje en `GET {{API_BASE_URL}}/v1/messages/{id}`.

## Enviar SMS recibidos a Discord

Si la URL del webhook empieza por `https://discord.com/api/webhooks/` (la URL de un webhook de canal de Discord), el evento `message.phone.received` no se envía en formato CloudEvents, sino como un mensaje de Discord con el remitente, el destinatario, el texto y el ID del mensaje. Así puedes ver los SMS entrantes en un canal sin escribir código. El resto de eventos se envían en formato CloudEvents normal, que Discord no entiende, así que para Discord suscríbete solo a `message.phone.received`.

## Probar tus webhooks

- Tu servidor httpSMS es quien hace la petición, así que la URL debe ser accesible **desde el servidor**, no desde tu navegador. Si desarrollas en tu ordenador, usa un túnel (por ejemplo `cloudflared` o `ngrok`) para exponer tu puerto local.
- Para ver exactamente qué llega, puedes apuntar el webhook temporalmente a un servicio de inspección de peticiones.
- La forma más fácil de disparar un evento es enviar un SMS desde otro teléfono al número de tu pasarela (`message.phone.received`) o enviar un SMS por la API (`message.phone.sent` y `message.phone.delivered`).
- Si no llega nada, revisa los logs de la API: cada envío de webhook registra la URL, el evento y el código de respuesta.
