# Claves API de teléfono

Una **clave API de teléfono** es una credencial pensada exclusivamente para la app de Android. En lugar de poner la clave API de tu cuenta (la que da acceso total a la API) dentro de cada teléfono, creas una clave distinta para cada dispositivo.

¿Por qué conviene hacerlo así?

- **Menos riesgo:** si un teléfono se pierde o te lo roban, solo eliminas su clave. La clave de tu cuenta, que usan tus integraciones y servidores, sigue funcionando sin cambios.
- **Permisos mínimos:** una clave de teléfono solo sirve para lo que la app necesita (latidos, mensajes recibidos, reportes de envío y entrega). No puede enviar SMS ni leer tu historial a través de la API.
- **Control por número:** cada clave sabe qué números de teléfono tiene asociados y rechaza peticiones para cualquier otro número.

En esta instalación puedes crear todas las claves de teléfono que necesites.

## Clave de la cuenta vs. clave de teléfono

|                      | Clave API de la cuenta                             | Clave API de teléfono                                      |
| -------------------- | -------------------------------------------------- | ---------------------------------------------------------- |
| Dónde se obtiene     | Página de configuración: `{{APP_URL}}/settings`    | Página de claves de teléfono: `{{APP_URL}}/phone-api-keys` |
| Formato              | Cadena aleatoria                                   | Empieza siempre con `pk_`                                  |
| Para qué sirve       | Toda la API: enviar mensajes, webhooks, contactos… | Solo los endpoints que usa la app de Android               |
| Cuántas puedes tener | Una por cuenta                                     | Las que quieras (por ejemplo, una por teléfono)            |
| Revocarla afecta a…  | Todas tus integraciones y teléfonos que la usen    | Solo a los teléfonos que inicien sesión con ella           |

> **Nota:** si quieres enviar mensajes, configurar webhooks o usar cualquier otro endpoint de la API (`{{API_DOCS_URL}}`), usa la clave de tu cuenta. Una clave `pk_` en esos endpoints recibe un error `401 Unauthorized`.

## Qué puede hacer una clave de teléfono

Una clave `pk_` solo es aceptada por estos endpoints, que son exactamente los que usa la app:

| Método | Endpoint                          | Uso en la app                                             |
| ------ | --------------------------------- | --------------------------------------------------------- |
| `PUT`  | `/v1/phones/fcm-token`            | Registrar el teléfono al iniciar sesión                   |
| `POST` | `/v1/heartbeats`                  | Enviar latidos para indicar que el teléfono está en línea |
| `GET`  | `/v1/messages/outstanding`        | Descargar el siguiente mensaje que debe enviar            |
| `POST` | `/v1/messages/{messageID}/events` | Reportar si un SMS se envió, falló o fue entregado        |
| `POST` | `/v1/messages/receive`            | Registrar un SMS recibido                                 |
| `POST` | `/v1/messages/calls/missed`       | Registrar una llamada perdida                             |

Además, para latidos, mensajes recibidos y reportes de envío, el servidor comprueba que el número de teléfono implicado esté asociado a la clave. Si no lo está, responde `401` con el mensaje _"You are not authorized to carry out the request for this phone number"_.

## Crear una clave de teléfono

### Desde el panel

1. Abre `{{APP_URL}}/phone-api-keys`.
2. Pulsa **Crear Clave API**.
3. Escribe un nombre que te ayude a reconocer el dispositivo (entre 1 y 60 caracteres), por ejemplo `Teléfono oficina`.
4. Pulsa **Crear Clave API de Teléfono**. La nueva clave aparece en la tabla.

### Con la API

Usa la clave de tu **cuenta** para crear claves de teléfono:

```bash
curl -X POST '{{API_BASE_URL}}/v1/phone-api-keys' \
  -H 'x-api-key: TU_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{ "name": "Teléfono oficina" }'
```

La respuesta incluye la clave completa en `data.api_key` (empieza con `pk_`), además de `phone_numbers` y `phone_ids`, que estarán vacíos hasta que un teléfono inicie sesión con ella.

## Iniciar sesión en la app de Android

Antes de empezar:

- Instala la versión más reciente de la app desde `{{APP_DOWNLOAD_URL}}`.
- El teléfono necesita Google Play Services: la app obtiene un token de notificaciones push (FCM) y sin él no puede iniciar sesión.

Pasos:

1. En `{{APP_URL}}/phone-api-keys`, pulsa **Ver** en la clave que quieres usar. Se abre una ventana con la clave en texto, un botón para copiarla y un código QR.
2. En la app, toca el icono de escáner QR junto al campo **API Key** y escanea el código. La clave se rellena sola (también puedes pegarla a mano).
3. Escribe el número de teléfono de la SIM en formato internacional E.164, por ejemplo `+18005550199`. Si el teléfono tiene dos SIM, la app te pedirá también el número de la SIM2.
4. En **URL del Servidor** escribe la dirección de tu API: `{{API_BASE_URL}}`.
5. Toca **Iniciar Sesión con API Key**.

> **Importante:** el código QR solo contiene la clave, no la dirección del servidor. Revisa siempre el campo **URL del Servidor**: si trae otra dirección, reemplázala por `{{API_BASE_URL}}`. La app solo acepta direcciones `https://`.

Al iniciar sesión, la app registra cada número en el servidor y este lo asocia automáticamente a la clave. En unos segundos verás el número en la columna **Números de Teléfono** de la tabla (la página se actualiza sola).

> **Nota:** la app también acepta la clave de tu cuenta para iniciar sesión, pero no es recomendable: cualquiera con acceso al teléfono tendría acceso completo a tu API.

## Cómo se asocian los números a una clave

- **Una clave puede tener varios números.** Por ejemplo, las dos SIM de un mismo teléfono quedan asociadas a la misma clave. También puedes compartir una clave entre varios teléfonos, aunque lo más seguro es usar una por dispositivo.
- **Un número solo puede estar en una clave a la vez.** Si inicias sesión con un número usando una clave diferente, el servidor lo quita de la clave anterior y lo pasa a la nueva automáticamente.
- **Si eliminas un teléfono** de tu cuenta, su número se quita también de la clave.

## Quitar un número de una clave

Úsalo cuando quieras que un teléfono deje de operar con esa clave sin borrar la clave completa.

1. En `{{APP_URL}}/phone-api-keys`, pulsa **Eliminar** junto al número.
2. Confirma con **Eliminar Teléfono de la clave**.

Desde ese momento, el servidor rechazará los latidos, mensajes recibidos y reportes de ese número hechos con esa clave. Para que el teléfono vuelva a funcionar, **cierra sesión en la app e inicia sesión de nuevo** con una clave válida.

Con la API necesitas el ID de la clave y el ID del teléfono (lo obtienes con `GET {{API_BASE_URL}}/v1/phones`):

```bash
curl -X DELETE '{{API_BASE_URL}}/v1/phone-api-keys/ID_DE_LA_CLAVE/phones/ID_DEL_TELEFONO' \
  -H 'x-api-key: TU_API_KEY'
```

## Eliminar una clave

1. En `{{APP_URL}}/phone-api-keys`, pulsa **Eliminar** en la columna **Acciones**.
2. Confirma con **Eliminar Clave API**.

La eliminación es permanente. Todos los teléfonos que usaban esa clave dejan de poder comunicarse con el servidor: **cierra sesión en la app e inicia sesión con una clave nueva** en cada uno.

```bash
curl -X DELETE '{{API_BASE_URL}}/v1/phone-api-keys/ID_DE_LA_CLAVE' \
  -H 'x-api-key: TU_API_KEY'
```

## Listar tus claves con la API

```bash
curl '{{API_BASE_URL}}/v1/phone-api-keys?limit=100' \
  -H 'x-api-key: TU_API_KEY'
```

| Parámetro | Valor por defecto | Rango    | Descripción                                     |
| --------- | ----------------- | -------- | ----------------------------------------------- |
| `limit`   | `1`               | 1 – 100  | Cuántas claves devolver                         |
| `skip`    | `0`               | 0 o más  | Cuántas claves saltar (paginación)              |
| `query`   | —                 | máx. 100 | Filtra las claves cuyo nombre contenga el texto |

> **Nota:** si no envías `limit`, la API devuelve **solo una** clave. Indica siempre `limit` para ver la lista completa.

## Solución de problemas

| Síntoma                                                          | Causa probable                                              | Qué hacer                                          |
| ---------------------------------------------------------------- | ----------------------------------------------------------- | -------------------------------------------------- |
| La app muestra _"Cannot validate the API key"_ al iniciar sesión | La clave está mal copiada o fue eliminada                   | Vuelve a escanear el QR o crea una clave nueva     |
| La app dice que la URL del servidor _"must be HTTPS"_            | La dirección empieza con `http://`                          | Usa `{{API_BASE_URL}}` con `https://`              |
| La app no encuentra el token FCM                                 | El teléfono no tiene Google Play Services                   | Instala o actualiza Google Play Services           |
| Respuestas `401` para un número después de usar **Eliminar**     | El número ya no está asociado a la clave                    | Cierra sesión e inicia sesión de nuevo en la app   |
| `401` al intentar enviar un SMS con una clave `pk_`              | Las claves de teléfono no pueden enviar mensajes por la API | Usa la clave de tu cuenta (`{{APP_URL}}/settings`) |
