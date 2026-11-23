# Eventos de webhooks

Esta página describe cada evento que tu servidor httpSMS puede enviar a tus [webhooks](/docs/webhooks-introduction): cuándo ocurre, qué datos trae y un ejemplo real del cuerpo de la petición.

## Estructura común

Todos los eventos llegan como un [CloudEvent](https://cloudevents.io/) en JSON. Lo que cambia de un evento a otro es el contenido de `data`:

```json
{
  "specversion": "1.0",
  "id": "identificador único del evento",
  "source": "ruta de la API que lo originó",
  "type": "tipo de evento, por ejemplo message.phone.received",
  "datacontenttype": "application/json",
  "time": "fecha de creación del evento (UTC)",
  "data": {}
}
```

El mismo tipo llega en la cabecera `X-Event-Type`, así que puedes decidir qué hacer sin leer el cuerpo.

Campos que se repiten en `data`:

| Campo        | Descripción                                                                                             |
| ------------ | ------------------------------------------------------------------------------------------------------- |
| `user_id`    | ID de tu usuario en httpSMS.                                                                            |
| `owner`      | El número de **tu** teléfono (la pasarela), en formato E.164.                                           |
| `contact`    | El número de la **otra** persona: el remitente en los SMS recibidos y el destinatario en los enviados.  |
| `sim`        | La SIM que se usó: `SIM1` o `SIM2`.                                                                     |
| `timestamp`  | Cuándo ocurrió el hecho (según el teléfono o el servidor, según el evento). Úsalo para ordenar eventos. |
| `encrypted`  | `true` si `content` está cifrado de extremo a extremo (ver [Primeros pasos](/docs/getting-started)).    |
| `request_id` | El `request_id` que enviaste al crear el mensaje, o `null` si no enviaste ninguno.                      |

Resumen de eventos:

| Evento                    | Cuándo se envía                          | ID del mensaje en |
| ------------------------- | ---------------------------------------- | ----------------- |
| `message.phone.received`  | El teléfono recibió un SMS.              | `data.message_id` |
| `message.phone.sent`      | El teléfono envió un SMS.                | `data.id`         |
| `message.phone.delivered` | El destinatario recibió el SMS.          | `data.id`         |
| `message.send.failed`     | El SMS no se pudo enviar.                | `data.id`         |
| `message.send.expired`    | El teléfono no envió el SMS a tiempo.    | `data.message_id` |
| `message.call.missed`     | El teléfono recibió una llamada perdida. | `data.message_id` |
| `phone.heartbeat.offline` | El teléfono dejó de dar señales de vida. | —                 |
| `phone.heartbeat.online`  | El teléfono volvió a conectarse.         | —                 |

> **Nota:** fíjate en que algunos eventos llaman `id` al ID del mensaje y otros `message_id`. No es un error de esta documentación: así lo envía la API.

## `message.phone.received`

Se envía cuando tu teléfono **recibe un SMS**. Es el evento que necesitas para responder mensajes automáticamente, crear tickets, guardar conversaciones en tu CRM, etc.

| Campo         | Descripción                                                                  |
| ------------- | ---------------------------------------------------------------------------- |
| `message_id`  | ID del mensaje en httpSMS.                                                   |
| `contact`     | Quién envió el SMS.                                                          |
| `owner`       | Tu número, el que recibió el SMS.                                            |
| `content`     | Texto del mensaje (cifrado si `encrypted` es `true`).                        |
| `attachments` | Lista de URLs de los archivos adjuntos si fue un MMS; lista vacía si no hay. |

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

> **Nota:** las URLs de `attachments` apuntan a tu propio servidor, con la forma `{{API_BASE_URL}}/v1/attachments/...`.

## `message.phone.sent`

Se envía cuando la app del teléfono **envía el SMS** a la red del operador. Que se haya enviado no garantiza que haya llegado; para eso está `message.phone.delivered`.

```json
{
  "specversion": "1.0",
  "id": "cf241cb9-0762-43d0-8b74-4f763f788c93",
  "source": "/v1/messages/ff313c14-17a3-4f74-bcb2-ca77213a64af/events",
  "type": "message.phone.sent",
  "datacontenttype": "application/json",
  "time": "2026-01-15T10:25:45.877240866Z",
  "data": {
    "id": "ff313c14-17a3-4f74-bcb2-ca77213a64af",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "request_id": "pedido-1234",
    "owner": "+5215512345678",
    "contact": "+5215587654321",
    "encrypted": false,
    "timestamp": "2026-01-15T10:25:45.166Z",
    "content": "Tu pedido 1234 ya va en camino",
    "sim": "SIM1"
  }
}
```

## `message.phone.delivered`

Se envía cuando el operador confirma que el SMS **llegó al teléfono del destinatario** (informe de entrega).

> **Nota:** no todos los operadores envían informes de entrega. Si el tuyo no lo hace, recibirás `message.phone.sent` pero nunca `message.phone.delivered`. No lo uses como única señal de éxito.

```json
{
  "specversion": "1.0",
  "id": "f4b86c6d-8c90-4ba7-a599-516a27f1a1e9",
  "source": "/v1/messages/5be8f09e-7007-4fe9-86b6-591d63fd38ad/events",
  "type": "message.phone.delivered",
  "datacontenttype": "application/json",
  "time": "2026-01-15T10:25:52.012355134Z",
  "data": {
    "id": "5be8f09e-7007-4fe9-86b6-591d63fd38ad",
    "owner": "+5215512345678",
    "contact": "+5215587654321",
    "request_id": null,
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "encrypted": false,
    "timestamp": "2026-01-15T10:25:51.262Z",
    "content": "Tu pedido 1234 ya va en camino",
    "sim": "SIM2"
  }
}
```

## `message.send.failed`

Se envía cuando el SMS **no se pudo enviar**. El motivo viene en `error_message`.

Valores habituales de `error_message`:

| Valor                                                                | Qué suele significar                                                                                                             |
| -------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `NO_SERVICE`                                                         | El teléfono no tenía cobertura.                                                                                                  |
| `RADIO_OFF`                                                          | La radio estaba apagada (por ejemplo, modo avión).                                                                               |
| `GENERIC_FAILURE`                                                    | Error genérico del operador o del teléfono; a menudo saldo insuficiente o número inválido.                                       |
| `LIMIT_EXCEEDED`                                                     | Android bloqueó el envío por mandar demasiados SMS seguidos. Ver [Controlar la velocidad de envío](/docs/control-sms-send-rate). |
| `NULL_PDU`                                                           | El mensaje no se pudo codificar.                                                                                                 |
| `UNKNOWN:<código>`                                                   | Android devolvió un código de error no reconocido.                                                                               |
| `Outgoing messages have been disabled on the mobile app`             | Los mensajes salientes están desactivados para esa SIM en la app.                                                                |
| `Outgoing message is encrypted but mobile app has no encryption key` | Enviaste `encrypted: true` pero el teléfono no tiene clave de encriptación configurada.                                          |

También puede contener errores del servidor, por ejemplo si no pudo avisar al teléfono mediante la notificación push.

```json
{
  "specversion": "1.0",
  "id": "d69c96fb-9f46-498b-811d-3adfb785e439",
  "source": "/v1/messages/508d783f-df33-4eaf-85e8-0fcab9958654/events",
  "type": "message.send.failed",
  "datacontenttype": "application/json",
  "time": "2026-01-15T10:30:58.979647773Z",
  "data": {
    "id": "508d783f-df33-4eaf-85e8-0fcab9958654",
    "error_message": "NO_SERVICE",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "owner": "+5215512345678",
    "request_id": null,
    "contact": "+5215587654321",
    "timestamp": "2026-01-15T10:30:58.218Z",
    "encrypted": false,
    "content": "Tu pedido 1234 ya va en camino",
    "sim": "SIM2"
  }
}
```

## `message.send.expired`

Se envía cuando el teléfono **no envió el SMS dentro del plazo de expiración** del teléfono (10 minutos por defecto, configurable por teléfono en `{{APP_URL}}/settings`). Suele ocurrir si el teléfono estaba apagado, sin internet o con la app dormida por el ahorro de batería.

httpSMS reintenta automáticamente los mensajes expirados hasta agotar los intentos configurados en el teléfono (2 por defecto). Por eso este evento puede llegar más de una vez para el mismo mensaje:

| Campo                | Descripción                                                                                          |
| -------------------- | ---------------------------------------------------------------------------------------------------- |
| `message_id`         | ID del mensaje.                                                                                      |
| `send_attempt_count` | Cuántos intentos se han hecho hasta ahora.                                                           |
| `is_final`           | `false`: httpSMS lo volverá a intentar. `true`: era el último intento y el mensaje ya no se enviará. |

```json
{
  "specversion": "1.0",
  "id": "8f5a9eaf-c495-4487-b11c-1824b8df6da2",
  "source": "/v1/messages/send",
  "type": "message.send.expired",
  "datacontenttype": "application/json",
  "time": "2026-01-15T10:45:43.461263458Z",
  "data": {
    "message_id": "ff313c14-17a3-4f74-bcb2-ca77213a64af",
    "owner": "+5215512345678",
    "send_attempt_count": 2,
    "is_final": true,
    "request_id": null,
    "contact": "+5215587654321",
    "encrypted": false,
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "timestamp": "2026-01-15T10:45:43.461254738Z",
    "content": "Tu pedido 1234 ya va en camino",
    "sim": "SIM1"
  }
}
```

> **Consejo:** si solo te interesa saber que un mensaje se perdió definitivamente, procesa únicamente los eventos con `is_final: true`.

## `message.call.missed`

Se envía cuando tu teléfono **recibe una llamada perdida**. Es útil, por ejemplo, para responder automáticamente con un SMS diciendo que ese número solo atiende por mensaje.

Para recibirlo, activa **Habilitar eventos de llamadas perdidas** para la SIM correspondiente en la configuración de la app Android.

```json
{
  "specversion": "1.0",
  "id": "f39d4b09-e2b7-42e2-b08f-d45820a18b05",
  "source": "/v1/messages/calls/missed",
  "type": "message.call.missed",
  "datacontenttype": "application/json",
  "time": "2026-01-15T11:02:35.3664741Z",
  "data": {
    "message_id": "682a47cd-adfd-4b5c-b4f2-65302746ac52",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "owner": "+5215512345678",
    "contact": "+5215599998888",
    "timestamp": "2026-01-15T11:02:34.527Z",
    "sim": "SIM1"
  }
}
```

## `phone.heartbeat.offline`

La app Android envía un "latido" (heartbeat) al servidor cada 15 minutos para decir "sigo aquí". Si el servidor pasa alrededor de **una hora** sin recibir ninguno (en concreto, más de 64 minutos), considera que el teléfono está desconectado y envía este evento.

Es la alerta que conviene vigilar: mientras el teléfono esté desconectado, tus mensajes se quedarán en cola o expirarán.

| Campo                      | Descripción                                |
| -------------------------- | ------------------------------------------ |
| `phone_id`                 | ID del teléfono en httpSMS.                |
| `monitor_id`               | ID del monitor de latidos de ese teléfono. |
| `owner`                    | El número del teléfono desconectado.       |
| `last_heartbeat_timestamp` | Cuándo se recibió el último latido.        |
| `timestamp`                | Cuándo se detectó la desconexión.          |

```json
{
  "specversion": "1.0",
  "id": "5a7e1d0c-2f4b-4c8e-9a61-3b2d7c9e8f10",
  "source": "/v1/heartbeats",
  "type": "phone.heartbeat.offline",
  "datacontenttype": "application/json",
  "time": "2026-01-15T12:20:23.461263458Z",
  "data": {
    "phone_id": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "last_heartbeat_timestamp": "2026-01-15T11:15:43.461254738Z",
    "timestamp": "2026-01-15T12:20:23.461254738Z",
    "monitor_id": "9f8e7d6c-5b4a-4392-8170-6e5d4c3b2a19",
    "owner": "+5215512345678"
  }
}
```

> **Nota:** el evento se envía una sola vez por desconexión, no cada hora. El siguiente aviso sobre ese teléfono será `phone.heartbeat.online` cuando vuelva.

## `phone.heartbeat.online`

Se envía cuando el servidor recibe un latido de un teléfono que estaba marcado como **desconectado**. Te sirve para cerrar la alerta que abriste con `phone.heartbeat.offline`. Tiene los mismos campos; aquí `last_heartbeat_timestamp` es el latido que acaba de llegar.

```json
{
  "specversion": "1.0",
  "id": "0c9b8a7f-6e5d-4c3b-a291-8f7e6d5c4b3a",
  "source": "/v1/heartbeats",
  "type": "phone.heartbeat.online",
  "datacontenttype": "application/json",
  "time": "2026-01-15T13:05:10.112233445Z",
  "data": {
    "phone_id": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "last_heartbeat_timestamp": "2026-01-15T13:05:09.871Z",
    "timestamp": "2026-01-15T13:05:10.112230112Z",
    "monitor_id": "9f8e7d6c-5b4a-4392-8170-6e5d4c3b2a19",
    "owner": "+5215512345678"
  }
}
```
