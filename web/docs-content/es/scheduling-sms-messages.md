# Programar mensajes SMS

Programar te permite decidir **cuándo** llega un mensaje: un recordatorio de cita una hora antes, una promoción a primera hora del lunes o simplemente evitar que tus clientes reciban SMS de madrugada.

Hay dos herramientas y puedes combinarlas:

| Herramienta                              | Se configura en… | Úsala cuando…                                                                          |
| ---------------------------------------- | ---------------- | -------------------------------------------------------------------------------------- |
| **Hora exacta** (`send_at` / `SendTime`) | Cada mensaje     | Un mensaje concreto debe salir en un momento preciso.                                  |
| **Horario de envío**                     | Cada teléfono    | Quieres que un teléfono solo envíe dentro de ciertas horas, p. ej. lun–vie 9:00–18:00. |

> **Nota:** un mensaje con hora exacta **ignora** el horario de envío y la [velocidad de envío](/docs/control-sms-send-rate) del teléfono: sale a la hora que indicaste. Los mensajes sin hora exacta respetan ambos.

## Programar un mensaje con la API

Añade el campo opcional `send_at` a `POST /v1/messages/send`:

```bash
curl -X POST '{{API_BASE_URL}}/v1/messages/send' \
  -H 'x-api-key: TU_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "from": "+18005550199",
    "to": "+18005550100",
    "content": "Hola, te recordamos tu cita de mañana a las 10:00.",
    "send_at": "2026-12-01T09:30:00-06:00"
  }'
```

Reglas de `send_at`:

| Regla              | Detalle                                                                                                                                                                               |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Formato            | [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) **con zona horaria**: `2026-12-01T09:30:00-06:00` o en UTC `2026-12-01T15:30:00Z`. Sin zona horaria la petición se rechaza. |
| Máximo             | Hasta **480 horas (20 días)** desde el momento de la petición. Más allá, la API responde `422` con un error en `send_at`.                                                             |
| Fecha en el pasado | El mensaje se envía de inmediato.                                                                                                                                                     |
| Contenido          | Hasta 2048 caracteres, igual que un envío normal.                                                                                                                                     |

La respuesta devuelve el mensaje con estado `pending` y el campo `scheduled_send_time` con la hora programada (en UTC). Guarda el `id` si quieres consultarlo o cancelarlo después.

> **Nota:** `POST /v1/messages/bulk-send` no acepta `send_at`. Para programar envíos a muchos destinatarios usa un archivo CSV o Excel (siguiente sección) o llama a `/v1/messages/send` una vez por destinatario.

## Programar mensajes desde un archivo CSV o Excel

En `{{APP_URL}}/bulk-messages` puedes subir un archivo con hasta **1000 mensajes**. Descarga una plantilla para empezar:

- CSV: `{{APP_URL}}/templates/httpsms-bulk.csv`
- Excel: `{{APP_URL}}/templates/httpsms-bulk.xlsx`

Columnas, en este orden:

| Columna                    | Obligatoria | Descripción                                                                   |
| -------------------------- | ----------- | ----------------------------------------------------------------------------- |
| `FromPhoneNumber`          | Sí          | Número de tu teléfono que enviará el SMS. Debe estar registrado en tu cuenta. |
| `ToPhoneNumber`            | Sí          | Número del destinatario.                                                      |
| `Content`                  | Sí          | Texto del mensaje, hasta 1024 caracteres.                                     |
| `SendTime(optional)`       | No          | Fecha y hora de envío. Déjala vacía para enviar en cuanto lo permita la cola. |
| `AttachmentURLs(optional)` | No          | URLs de adjuntos separadas por comas (máximo 10), para enviar como MMS.       |

Formatos aceptados en `SendTime(optional)`:

| Ejemplo                     | Cómo se interpreta                                   |
| --------------------------- | ---------------------------------------------------- |
| `2026-12-01T09:30:00`       | Hora local según la **zona horaria de tu perfil**.   |
| `2026-12-01 09:30:00`       | Igual que el anterior (con espacio en lugar de `T`). |
| `2026-12-01T09:30:00-06:00` | Usa la zona horaria incluida en el valor.            |

La zona horaria de tu perfil se elige en `{{APP_URL}}/settings`. Si no has elegido ninguna, se usa UTC.

> **Importante:** en archivos, `SendTime` admite como máximo **420 horas (17,5 días)** en el futuro, menos que las 480 horas de la API. Las filas que lo superen se rechazan con un error que indica el número de fila.

> **Consejo para Excel:** escribe la fecha como texto (por ejemplo, empezando la celda con un apóstrofo `'2026-12-01T09:30:00`). Si Excel la convierte a su propio formato de fecha, el archivo puede rechazarse por formato incorrecto.

También puedes subir el archivo con la API, en el campo `document` de un formulario:

```bash
curl -X POST '{{API_BASE_URL}}/v1/bulk-messages' \
  -H 'x-api-key: TU_API_KEY' \
  -F 'document=@mensajes.csv'
```

Las filas con `SendTime` salen a su hora exacta. Las filas sin `SendTime` se reparten según la velocidad de envío de cada teléfono remitente.

## Horarios de envío

Un horario de envío define **franjas horarias** en las que un teléfono puede enviar SMS. Los mensajes sin hora exacta que lleguen fuera de esas franjas esperan en la cola y salen cuando se abre la siguiente, respetando la velocidad de envío del teléfono.

Es útil para:

- Respetar el descanso de tus clientes (nada de SMS a las 3 de la mañana).
- Cumplir normativas locales sobre horarios de comunicaciones comerciales.
- Concentrar los envíos en horario laboral, cuando alguien puede atender respuestas.

### Crear un horario desde el panel

1. Abre `{{APP_URL}}/settings` y baja hasta la sección de horarios de envío.
2. Pulsa **Crear Horario de Envío**.
3. Escribe un nombre, por ejemplo `Horario comercial`.
4. Elige la **Zona Horaria** en la que se interpretan las horas (por defecto, la de tu perfil).
5. Activa los días en los que se puede enviar y define para cada uno una o varias franjas con hora de **inicio** y **fin**. Un horario nuevo viene con lunes a viernes de 09:00 a 17:00.
6. Pulsa **Guardar Horario**.

### Asignar el horario a un teléfono

1. En `{{APP_URL}}/settings`, sección **Teléfonos**, pulsa **Editar** en el teléfono.
2. Elige el horario en el campo **Horario de Envío**.
3. Pulsa **Actualizar Teléfono**.

Cada teléfono puede tener un horario como máximo, y un mismo horario puede usarse en todos los teléfonos que quieras. Para quitarlo, borra el campo **Horario de Envío** y actualiza el teléfono.

### Reglas de un horario

| Campo                    | Valores permitidos                                     | Descripción                                                                                         |
| ------------------------ | ------------------------------------------------------ | --------------------------------------------------------------------------------------------------- |
| `name`                   | 2 – 100 caracteres                                     | Nombre para identificarlo.                                                                          |
| `timezone`               | Zona horaria IANA válida, p. ej. `America/Mexico_City` | Las franjas se interpretan en esta zona; los cambios de horario de verano se aplican solos.         |
| `windows`                | Al menos 1 franja                                      | Lista de franjas horarias.                                                                          |
| `windows[].day_of_week`  | `0` – `6`                                              | Día de la semana: `0` = domingo, `1` = lunes … `6` = sábado.                                        |
| `windows[].start_minute` | `0` – `1439`                                           | Inicio, en minutos desde la medianoche (09:00 = `540`).                                             |
| `windows[].end_minute`   | `1` – `1440`                                           | Fin, en minutos desde la medianoche (17:00 = `1020`, 24:00 = `1440`). Debe ser mayor que el inicio. |

Además:

- Cada día admite **hasta 6 franjas**, y no pueden solaparse entre sí.
- Una franja no puede cruzar la medianoche. Para enviar de 22:00 a 02:00, crea dos franjas: `22:00–24:00` en un día y `00:00–02:00` en el siguiente.
- La hora de fin no está incluida: con una franja de 09:00 a 17:00, un mensaje que llega exactamente a las 17:00 espera a la siguiente franja.

### Gestionar horarios con la API

Crear un horario de lunes a viernes, de 9:00 a 18:00 en Ciudad de México:

```bash
curl -X POST '{{API_BASE_URL}}/v1/send-schedules' \
  -H 'x-api-key: TU_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Horario comercial",
    "timezone": "America/Mexico_City",
    "windows": [
      { "day_of_week": 1, "start_minute": 540, "end_minute": 1080 },
      { "day_of_week": 2, "start_minute": 540, "end_minute": 1080 },
      { "day_of_week": 3, "start_minute": 540, "end_minute": 1080 },
      { "day_of_week": 4, "start_minute": 540, "end_minute": 1080 },
      { "day_of_week": 5, "start_minute": 540, "end_minute": 1080 }
    ]
  }'
```

| Acción           | Petición                                                          |
| ---------------- | ----------------------------------------------------------------- |
| Listar horarios  | `GET /v1/send-schedules`                                          |
| Crear un horario | `POST /v1/send-schedules`                                         |
| Actualizar       | `PUT /v1/send-schedules/{scheduleID}` (mismo cuerpo que al crear) |
| Eliminar         | `DELETE /v1/send-schedules/{scheduleID}`                          |

Para asignar el horario a un teléfono por API, envía su `id` en `message_send_schedule_id` con `PUT /v1/phones` (consulta las precauciones en [Controlar la velocidad de envío](/docs/control-sms-send-rate)).

> **Nota:** al eliminar un horario, se desvincula automáticamente de todos los teléfonos que lo usaban, que vuelven a enviar a cualquier hora.

### Cómo decide el servidor la hora de envío

| Momento en que llega el mensaje     | Qué ocurre                                                                 |
| ----------------------------------- | -------------------------------------------------------------------------- |
| Dentro de una franja                | Sale de inmediato (respetando la velocidad de envío).                      |
| Antes de la primera franja del día  | Espera al inicio de esa franja.                                            |
| Después de la última franja del día | Espera a la siguiente franja disponible, aunque sea otro día de la semana. |

El horario se evalúa cuando el mensaje entra en la cola. Si después cambias el horario, los mensajes que ya tenían una hora asignada la conservan; los nuevos usan el horario actualizado.

## Ver y cancelar mensajes programados

- En `{{APP_URL}}/bulk-messages`, el historial muestra cuántos mensajes de cada carga siguen en estado **Programados**.
- Con la API puedes consultar un mensaje con `GET {{API_BASE_URL}}/v1/messages/{messageID}`: el campo `scheduled_send_time` indica la hora exacta solicitada y `status` su estado actual.

Para cancelar un mensaje que todavía no ha salido (estado `pending` o `scheduled`), elimínalo:

```bash
curl -X DELETE '{{API_BASE_URL}}/v1/messages/ID_DEL_MENSAJE' \
  -H 'x-api-key: TU_API_KEY'
```

El teléfono ya no podrá descargarlo, así que no se enviará. Ten en cuenta que el mensaje también desaparece de tu historial.

Para entender todo el recorrido de un mensaje, desde que llega a la API hasta que se entrega, consulta [Cola de mensajes salientes](/docs/outgoing-message-queue).
