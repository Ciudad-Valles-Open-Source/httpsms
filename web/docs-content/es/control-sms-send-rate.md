# Controlar la velocidad de envío de SMS

Cuando mandas muchos mensajes de golpe, enviarlos todos al mismo tiempo suele salir mal:

- **Android limita los envíos.** Si una app intenta mandar alrededor de 30 SMS en un minuto, el sistema puede detenerla y pedir confirmación manual en la pantalla del teléfono.
- **Los operadores vigilan las ráfagas.** Muchos mensajes seguidos desde la misma SIM pueden marcarse como spam o hacer que bloqueen la línea.
- **Necesitas previsibilidad.** Con una velocidad fija sabes cuánto tardará una campaña en salir.

Por eso cada teléfono tiene un ajuste de **Mensajes Por Minuto**. El servidor reparte los envíos en el tiempo para que el teléfono nunca supere ese ritmo, sin importar cuántos mensajes recibas por la API.

## Cómo se calcula el ritmo

El servidor deja un intervalo fijo entre un mensaje y el siguiente del mismo teléfono:

```text
intervalo (segundos) = 60 ÷ mensajes_por_minuto
```

Con el valor por defecto de **10 mensajes por minuto**, el teléfono envía un SMS cada **6 segundos**:

| Mensaje | Sale aproximadamente a los… |
| ------- | --------------------------- |
| 1.º     | 0 s (de inmediato)          |
| 2.º     | 6 s                         |
| 3.º     | 12 s                        |
| 10.º    | 54 s                        |
| 100.º   | 9 min 54 s                  |

Cada teléfono (cada número) tiene su propio ritmo: los mensajes que salen de un teléfono no retrasan a los de otro.

## Ajustes de envío de cada teléfono

| Ajuste en el panel                | Campo de la API              | Por defecto        | Valores permitidos | Qué hace                                                                                                                                  |
| --------------------------------- | ---------------------------- | ------------------ | ------------------ | ----------------------------------------------------------------------------------------------------------------------------------------- |
| Mensajes Por Minuto               | `messages_per_minute`        | `10`               | 1 – 30             | Cuántos SMS como máximo envía el teléfono por minuto.                                                                                     |
| Máx. Intentos de Envío            | `max_send_attempts`          | `2`                | 1 – 5              | Número total de intentos por mensaje (el primero más los reintentos) si el teléfono no confirma el envío a tiempo.                        |
| Expiración del Mensaje (segundos) | `message_expiration_seconds` | `600` (10 minutos) | 60 – 3600          | Cuánto espera el servidor a que el teléfono confirme el envío antes de marcar el intento como expirado y, si quedan intentos, reintentar. |
| Horario de Envío                  | `message_send_schedule_id`   | ninguno            | ID de un horario   | Limita el envío a ciertas horas del día. Consulta [Programar mensajes SMS](/docs/scheduling-sms-messages).                                |

Los valores por defecto se aplican cuando un teléfono se registra por primera vez desde la app. Los reintentos y la expiración se explican en detalle en [Cola de mensajes salientes](/docs/outgoing-message-queue).

## Cambiar la velocidad desde el panel

1. Abre `{{APP_URL}}/settings` y baja hasta la sección **Teléfonos**. La columna **Tasa** muestra la velocidad actual de cada número.
2. Pulsa **Editar** en el teléfono que quieres ajustar.
3. Cambia **Mensajes Por Minuto** (y, si quieres, **Máx. Intentos de Envío** y **Expiración del Mensaje**).
4. Pulsa **Actualizar Teléfono**.

El nuevo ritmo se aplica a los mensajes que entren a la cola a partir de ese momento. Los que ya tenían una hora de envío asignada la conservan.

## Cambiar la velocidad con la API

El endpoint `PUT /v1/phones` actualiza el teléfono completo, así que lo más seguro es leer primero su configuración actual:

```bash
curl '{{API_BASE_URL}}/v1/phones?limit=20' \
  -H 'x-api-key: TU_API_KEY'
```

Después envía el objeto completo con los valores que quieres cambiar:

```bash
curl -X PUT '{{API_BASE_URL}}/v1/phones' \
  -H 'x-api-key: TU_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "phone_number": "+18005550199",
    "sim": "SIM1",
    "fcm_token": "TOKEN_ACTUAL_DEL_TELEFONO",
    "messages_per_minute": 6,
    "max_send_attempts": 2,
    "message_expiration_seconds": 600,
    "message_send_schedule_id": null
  }'
```

> **Importante:** copia `fcm_token`, `sim` y `message_send_schedule_id` tal como te los devolvió `GET /v1/phones`. Si omites `fcm_token`, el servidor borra el token y deja de poder avisar al teléfono; si omites `message_send_schedule_id`, el teléfono pierde su horario de envío.

Reglas de validación de la API:

- `messages_per_minute` acepta hasta `30`. Enviar `0` no cambia el valor actual.
- `max_send_attempts` acepta hasta `5`. Enviar `0` no cambia el valor actual.
- `message_expiration_seconds` debe estar entre `60` y `3600`, y es obligatorio cuando `max_send_attempts` es mayor que `0`.

## Qué mensajes respetan este límite

| Origen del mensaje                                                | ¿Respeta Mensajes Por Minuto?   |
| ----------------------------------------------------------------- | ------------------------------- |
| `POST /v1/messages/send` sin `send_at`                            | Sí                              |
| `POST /v1/messages/bulk-send`                                     | Sí                              |
| Archivo CSV o Excel en `{{APP_URL}}/bulk-messages` sin `SendTime` | Sí                              |
| Reintentos automáticos                                            | Sí                              |
| Cualquier mensaje con hora exacta (`send_at` o `SendTime`)        | **No**: sale a la hora indicada |

Los mensajes con hora exacta saltan la cola a propósito: si programas varios para la misma hora, saldrán juntos. Si necesitas que respeten el ritmo, escalónalos tú mismo (por ejemplo, con unos segundos de diferencia) o envíalos sin hora y usa un [horario de envío](/docs/scheduling-sms-messages).

## Cómo elegir un buen valor

- **Empieza con 5–10 por minuto** para campañas grandes. Es lo más fiable y lo menos propenso a filtros de spam.
- **Sube a 20–30 por minuto** solo si tu operador lo tolera y necesitas más rapidez.
- **Prefiere divisores de 60** (1, 2, 3, 4, 5, 6, 10, 12, 15, 20 o 30). El intervalo entre mensajes individuales se redondea hacia abajo a segundos enteros, así que con valores como 25 el teléfono puede llegar a enviar a ritmo de 30 por minuto.
- **¿Necesitas más capacidad?** Reparte el envío entre varios teléfonos o SIM: cada número tiene su propio límite.

> **Nota:** si en la columna **Tasa** ves **Ilimitado**, ese teléfono tiene el valor `0` guardado (por ejemplo, de una versión anterior) y no aplica espaciado. Edítalo y asígnale un valor entre 1 y 30.
