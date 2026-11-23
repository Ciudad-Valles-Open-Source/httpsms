# Tiempo real (Soketi)

La interfaz web se actualiza sola cuando llega, se envía o falla un SMS, sin recargar la página. Esto funciona con [Soketi](https://soketi.app/), un servidor WebSocket autohospedado compatible con el protocolo de Pusher.

> **Antes de configurarlo, lee la sección [⚠️ La `SOKETI_KEY` debe ser alfanumérica](#️-la-soketi_key-debe-ser-alfanumérica).** Es la causa más común de que el tiempo real no funcione.

---

## Cómo funciona

```
Teléfono Android ──▶ API (Go)
                      │  publica el evento (phone.updated, message.phone.received, …)
                      ▼
                 WebsocketListener ──HTTP──▶ Soketi (soketi:6001, red interna de Docker)
                                                │
Navegador ◀──WebSocket── nginx (web:3000, ruta /app/) ◀─┘
(pusher-js)
```

1. La API recibe un evento (por ejemplo, el teléfono reporta un SMS recibido).
2. El `WebsocketListener` de la API lo publica en Soketi, en el canal cuyo nombre es el **ID del usuario**.
3. El navegador está conectado por WebSocket a `wss://<tu-dominio>/app/<SOKETI_KEY>`. nginx (contenedor `web`) reenvía la ruta `/app/` a Soketi.
4. La web recibe el evento y recarga lo necesario:

| Evento                   | Qué se actualiza                                          |
| ------------------------ | --------------------------------------------------------- |
| `message.phone.received` | Lista de conversaciones y la conversación abierta         |
| `message.phone.sent`     | Lista de conversaciones y la conversación abierta         |
| `message.send.failed`    | Lista de conversaciones y la conversación abierta         |
| `message.call.missed`    | Lista de conversaciones y la conversación abierta         |
| `phone.updated`          | Lista de conversaciones, teléfonos, heartbeat y API keys |

La web abre **una sola conexión WebSocket** para toda la app (`web/app/composables/useRealtime.ts`).

---

## Variables de entorno

Las mismas claves deben llegar a los **tres** servicios: `api`, `soketi` y `web`. El `docker-compose.yml` ya las reparte a partir de las mismas variables.

| Variable         | Servicio(s)        | Descripción                                                                                   |
| ---------------- | ------------------ | --------------------------------------------------------------------------------------------- |
| `SOKETI_APP_ID`  | api, soketi        | ID de la aplicación en Soketi.                                                                |
| `SOKETI_KEY`     | api, soketi, web   | Clave pública. **Solo letras y números** (ver abajo).                                         |
| `SOKETI_SECRET`  | api, soketi        | Clave secreta. La API la usa para firmar los eventos. Si está vacía, la API no publica nada.  |
| `SOKETI_HOST`    | api                | Dónde la API encuentra a Soketi. Por defecto `soketi:6001` (red interna de Docker).            |
| `SOKETI_SECURE`  | api                | `false` cuando la API habla con Soketi por HTTP dentro de Docker (lo normal).                 |
| `SOKETI_CLUSTER` | api, web           | Solo se usa con Pusher en la nube. Con Soketi puede quedarse en `mt1`.                        |
| `USE_SOKETI`     | web                | `true` para conectar al mismo dominio de la web (`/app/`) en vez de a Pusher en la nube.      |

> **La web guarda `SOKETI_KEY` al compilar.** Si cambias la clave, hay que **recompilar el contenedor `web`** (`docker compose up -d --build web`, o volver a desplegar en Coolify). Reiniciarlo no basta.

---

## ⚠️ La `SOKETI_KEY` debe ser alfanumérica

**Problema:** el cliente de Pusher arma la URL del WebSocket pegando la clave tal cual:

```
wss://tu-dominio/app/<SOKETI_KEY>?protocol=7&client=js…
```

Si la clave se generó en **base64** (por ejemplo con `openssl rand -base64 48`), puede contener `/` y `+`:

- La `/` parte la ruta en dos (`/app/abc/def…`), Soketi no la reconoce y responde **404**.
- El `+` se interpreta como un espacio, así que la clave ya no coincide.

**Síntomas:**

- Nada se actualiza solo; hay que recargar la página para ver los mensajes nuevos.
- En los logs del contenedor `web` (nginx) aparece `GET /app/… 404` en lugar de `101`.
- En la consola del navegador aparece:
  `SOKETI_KEY must only contain letters, numbers, "-" or "_"…`
  La web muestra este aviso y no intenta conectar si la clave es inválida.

**Solución:** genera claves en hexadecimal:

```bash
openssl rand -hex 32
```

Usa una clave para `SOKETI_KEY` y otra para `SOKETI_SECRET` (el secreto técnicamente acepta cualquier carácter, pero en hex evitas problemas al copiarlo). Después:

1. Actualiza `SOKETI_KEY` y `SOKETI_SECRET` en tu `.env` local **y** en las variables de entorno de Coolify.
2. Recompila y reinicia: `docker compose up -d --build` (en Coolify: *Redeploy*).
3. Si tus usuarios tenían la página abierta, que la recarguen una vez.

`SOKETI_APP_ID` tampoco debería llevar `/` ni espacios; letras, números, `-` y `_` funcionan bien.

---

## Otro problema ya corregido: `PUSHER_SECRET`

La API solo activaba los listeners de WebSocket si existía la variable `PUSHER_SECRET`, que viene del proyecto original. Este fork usa `SOKETI_SECRET`, así que **la API nunca publicaba eventos**, aunque Soketi y nginx estuvieran bien.

Ya está corregido en `api/pkg/di/container.go`: ahora basta con `SOKETI_SECRET` (y `PUSHER_SECRET` sigue funcionando por compatibilidad).

---

## Cómo verificar que funciona

**1. La API publica eventos.** Al arrancar, los logs de la API deben mostrar:

```
registering listeners for listeners.WebsocketListener
creating *pusher.Client
```

Si en cambio ves `skipping websocket listeners because the SOKETI_SECRET env variable is not set`, falta `SOKETI_SECRET` en el servicio `api`.

**2. El navegador conecta.** En el navegador: DevTools → pestaña *Network* → filtro *WS*. Debe haber una conexión a `/app/<clave>` con estado **101 Switching Protocols** que permanece abierta. En los logs de `web` (nginx) debe aparecer `GET /app/… 101`.

**3. Prueba completa.** Con la app abierta en la lista de conversaciones, envía un SMS al teléfono Android: la conversación debe aparecer o moverse arriba en uno o dos segundos, sin recargar.

---

## En producción (Coolify)

- La web se sirve con HTTPS, así que el navegador usa `wss://` en el puerto 443 del mismo dominio. Traefik (Coolify) → nginx (`web`) → Soketi. No hace falta exponer el puerto 6001.
- El `nginx.conf` de `web` ya trae la configuración de WebSocket (`Upgrade`/`Connection`) para la ruta `/app/`.
- Tras cambiar cualquier `SOKETI_*`, vuelve a desplegar **todos** los servicios: la web necesita recompilarse y la API y Soketi necesitan reiniciarse con la misma clave.
