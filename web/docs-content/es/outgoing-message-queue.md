# Cola de mensajes salientes

Tu teléfono Android no está conectado permanentemente al servidor esperando órdenes. Cuando envías un SMS (por la API, con `bulk-send` o subiendo un archivo), el mensaje entra primero en una **cola** en el servidor. El servidor decide en qué momento despertar al teléfono con una notificación push, y el teléfono entonces descarga el mensaje y lo envía.

Gracias a esta cola puedes enviar cientos de mensajes en una sola petición sin saturar el teléfono, respetar horarios y reintentar automáticamente cuando el teléfono no responde.

## El recorrido de un mensaje

1. **La API acepta el mensaje** y lo guarda con estado `pending`.
2. **El servidor calcula cuándo debe salir** (ver la siguiente sección) y lo deja en estado `scheduled`.
3. **Llegado el momento, envía una notificación push al teléfono.** Cada notificación cuenta como un intento de envío y arranca el temporizador de expiración.
4. **El teléfono descarga el mensaje** y este pasa a `sending`.
5. **El teléfono lo envía** y lo reporta como `sent`. Si el operador devuelve confirmación de entrega, pasa a `delivered`.
6. **Si algo sale mal:** si el teléfono reporta un error o no se le puede notificar, el mensaje queda en `failed`. Si el teléfono no confirma a tiempo, pasa a `expired` y, si le quedan intentos, se vuelve a programar.

### Estados de un mensaje saliente

| Estado      | Significado                                                                                                          |
| ----------- | -------------------------------------------------------------------------------------------------------------------- |
| `pending`   | El servidor recibió el mensaje y todavía no le ha asignado hora de envío.                                            |
| `scheduled` | Ya tiene hora asignada; espera su turno en la cola (por hora exacta, ritmo o horario de envío).                      |
| `sending`   | El teléfono ya lo descargó y lo está enviando.                                                                       |
| `sent`      | El teléfono lo envió a la red del operador.                                                                          |
| `delivered` | El operador confirmó que llegó al destinatario.                                                                      |
| `failed`    | El teléfono reportó un error, o el servidor no pudo notificar al teléfono (por ejemplo, porque no tiene token push). |
| `expired`   | El teléfono no confirmó el envío dentro del tiempo de expiración. Si quedan intentos, se reintenta.                  |

## Cómo decide la cola cuándo sale cada mensaje

El servidor toma la decisión con estas reglas, en orden:

```text
¿El mensaje tiene hora exacta (send_at o SendTime)?
├── Sí → sale a esa hora. Si ya pasó, sale de inmediato.
│        No se aplican ni la velocidad de envío ni el horario del teléfono.
└── No → 1. Si el teléfono tiene horario de envío y ahora está cerrado,
            se mueve al inicio de la siguiente franja.
         2. Se espera al menos (60 ÷ mensajes por minuto) segundos
            desde el mensaje anterior del mismo teléfono.
         3. Si el resultado cae fuera del horario, se mueve
            otra vez a la siguiente franja abierta.
```

### 1. Hora exacta

Si indicas `send_at` en `POST /v1/messages/send`, o la columna `SendTime(optional)` en un archivo, el mensaje sale exactamente a esa hora y **se salta** el resto de reglas. Es ideal para recordatorios y alertas que no pueden moverse. Encontrarás los formatos y límites en [Programar mensajes SMS](/docs/scheduling-sms-messages).

> **Nota:** como se saltan el ritmo, muchos mensajes programados para el mismo minuto salen todos juntos. Si son muchos, escalónalos.

### 2. Velocidad de envío

Los mensajes sin hora exacta se espacian según el ajuste **Mensajes Por Minuto** del teléfono remitente. Con el valor por defecto (10), sale uno cada 6 segundos. Cada número lleva su propio ritmo, así que los mensajes de un teléfono nunca retrasan a los de otro. Detalles en [Controlar la velocidad de envío de SMS](/docs/control-sms-send-rate).

En envíos masivos (`/v1/messages/bulk-send` o un archivo), el servidor además numera los mensajes de cada teléfono remitente y retrasa cada uno según su posición:

```text
retraso = posición × (60 ÷ mensajes_por_minuto)
```

| Posición del mensaje | Retraso con 10 por minuto |
| -------------------- | ------------------------- |
| 0 (el primero)       | 0 s                       |
| 1                    | 6 s                       |
| 2                    | 12 s                      |
| 9                    | 54 s                      |

Si un archivo usa varios números remitentes (columna `FromPhoneNumber`), cada uno tiene su propia numeración. Las filas con `SendTime` no cuentan en esa numeración.

Ejemplo de envío masivo por API (hasta 1000 destinatarios y 1024 caracteres):

```bash
curl -X POST '{{API_BASE_URL}}/v1/messages/bulk-send' \
  -H 'x-api-key: TU_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "from": "+18005550199",
    "to": ["+18005550100", "+18005550101", "+18005550102"],
    "content": "¡Hola! Este es un aviso para todos nuestros clientes."
  }'
```

Con 10 mensajes por minuto, el primero sale de inmediato, el segundo a los 6 segundos y el tercero a los 12.

### 3. Horario de envío

Si el teléfono tiene un horario de envío asignado, los mensajes sin hora exacta solo salen dentro de sus franjas. Lo que llegue fuera espera en estado `scheduled` hasta que se abra la siguiente franja, y a partir de ahí sale al ritmo configurado.

| Cuándo llega el mensaje             | Qué ocurre                                                    |
| ----------------------------------- | ------------------------------------------------------------- |
| Dentro de una franja                | Sale de inmediato, respetando la velocidad de envío.          |
| Antes de que abra la franja del día | Espera al inicio de la franja.                                |
| Después de la última franja del día | Espera a la siguiente franja disponible (puede ser otro día). |

Cómo crear horarios y asignarlos a un teléfono: [Programar mensajes SMS](/docs/scheduling-sms-messages).

> **Nota:** la hora de salida se calcula cuando el mensaje entra en la cola. Si después cambias la velocidad o el horario del teléfono, los mensajes ya programados conservan su hora.

## Expiración y reintentos

A veces el teléfono no responde: está sin batería, sin conexión o el sistema ha detenido la app. Para eso existen dos ajustes por teléfono:

| Ajuste                            | Por defecto        | Rango     | Qué controla                                                                                         |
| --------------------------------- | ------------------ | --------- | ---------------------------------------------------------------------------------------------------- |
| Expiración del Mensaje (segundos) | `600` (10 minutos) | 60 – 3600 | Tiempo que el servidor espera, desde que envía la notificación, a que el teléfono confirme el envío. |
| Máx. Intentos de Envío            | `2`                | 1 – 5     | Número **total** de intentos por mensaje, incluido el primero.                                       |

Cómo funciona:

1. El servidor envía la notificación y arranca el temporizador de expiración.
2. Si al terminar el tiempo el mensaje sigue en `pending`, `scheduled` o `sending`, pasa a `expired`.
3. Si todavía le quedan intentos, vuelve a la cola. Este reintento **sí respeta** la velocidad y el horario de envío, aunque el mensaje original tuviera hora exacta.
4. Cuando se agotan los intentos, el mensaje se queda en `expired` definitivamente.

Con los valores por defecto, un teléfono que no responde recibe dos intentos y el mensaje expira definitivamente unos 20 minutos después del primer aviso.

> **Nota:** la notificación push también caduca cuando se cumple el tiempo de expiración. Si el teléfono pasa más tiempo que ese desconectado, no recibirá ese aviso concreto, pero sí el del reintento.

> **Consejo:** no uses una expiración demasiado corta. Si el teléfono envía el SMS pero tarda en confirmarlo, un reintento podría producir un mensaje duplicado. Si la confirmación llega tarde, el mensaje igualmente pasa a `sent`.

## Seguir el estado de tus mensajes

- **Webhooks:** recibe en tu servidor los eventos `message.phone.sent`, `message.phone.delivered`, `message.send.failed` y `message.send.expired`. Este último se emite en cada expiración; su campo `is_final` vale `true` cuando ya no habrá más reintentos. Consulta [Eventos de webhooks](/docs/webhooks-events).
- **API:** `GET {{API_BASE_URL}}/v1/messages/{messageID}` devuelve el estado (`status`), los intentos realizados (`send_attempt_count`), el máximo permitido (`max_send_attempts`) y la hora exacta solicitada (`scheduled_send_time`), si la hay.
- **Panel:** el historial de `{{APP_URL}}/bulk-messages` muestra cuántos mensajes de cada carga están pendientes, programados, enviados, entregados, fallidos o expirados.

## Notas para instalaciones propias

Con la configuración por defecto de Docker (`EVENTS_QUEUE_TYPE=emulator`), los temporizadores de la cola viven en la memoria del proceso de la API. Si reinicias la API mientras hay mensajes esperando (programados para más tarde, envíos masivos en curso o reintentos pendientes), esos temporizadores se pierden y los mensajes afectados pueden quedarse en `pending` o `scheduled` sin salir.

Para evitarlo:

- Reinicia o actualiza la API cuando no haya envíos programados ni campañas en curso.
- Después de un reinicio, revisa los mensajes que sigan en `pending` o `scheduled` y vuelve a enviarlos si hace falta.
- Si necesitas una cola persistente, configura `EVENTS_QUEUE_TYPE=cloud_tasks` con Google Cloud Tasks.
