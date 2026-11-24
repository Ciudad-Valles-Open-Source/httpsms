# Aprobación de cuentas y administradores

> English version: [admin.md](./admin.md)

Por defecto, cualquier persona que pueda abrir tu instancia de httpSMS puede crear una cuenta con Firebase y empezar a usarla: generar claves de API, registrar teléfonos y enviar SMS a través de ellos. En una instalación privada eso casi nunca es lo que quieres.

La **aprobación de cuentas** resuelve esto: las cuentas nuevas quedan **pendientes** hasta que un administrador las aprueba desde el panel. Mientras tanto no pueden hacer nada más que ver su propio perfil o borrar su cuenta.

---

## Índice

- [Cómo funciona](#cómo-funciona)
- [Configurar `ADMIN_EMAILS`](#configurar-admin_emails)
- [Por qué el correo del administrador debe estar verificado](#por-qué-el-correo-del-administrador-debe-estar-verificado)
- [Administrar usuarios desde el panel](#administrar-usuarios-desde-el-panel)
- [Qué ven los usuarios pendientes o suspendidos](#qué-ven-los-usuarios-pendientes-o-suspendidos)
- [Efecto en claves de API y teléfonos](#efecto-en-claves-de-api-y-teléfonos)
- [Al actualizar desde una versión anterior](#al-actualizar-desde-una-versión-anterior)
- [Configuración recomendada de Firebase](#configuración-recomendada-de-firebase)
- [Desactivar la aprobación](#desactivar-la-aprobación)
- [Notas de seguridad](#notas-de-seguridad)

---

## Cómo funciona

Cada usuario tiene un **estado**:

| Estado | Significado |
| --- | --- |
| `active` | Cuenta aprobada. Funciona con normalidad. |
| `pending` | Cuenta nueva que espera aprobación. No puede usar la API. |
| `suspended` | Cuenta bloqueada por un administrador. No puede usar la API. |

La aprobación se activa en cuanto defines la variable `ADMIN_EMAILS` en la API con al menos un correo. A partir de ese momento:

1. Toda cuenta nueva se crea como `pending`, salvo que su correo esté en `ADMIN_EMAILS` (y verificado): esa se crea como `active`.
2. Cualquier petición autenticada de un usuario `pending` o `suspended` recibe `403` con un código que explica el motivo:

   ```json
   {
     "status": "error",
     "message": "...",
     "data": { "code": "account_pending" }
   }
   ```

   El código es `account_pending` o `account_suspended`.
3. Solo hay dos excepciones para esos usuarios: `GET /v1/users/me` (para que el panel sepa en qué estado están) y `DELETE /v1/users/me` (para que puedan borrar su cuenta si quieren).
4. Los administradores y el usuario interno del sistema (`EVENTS_QUEUE_USER_ID`, el que procesa la cola de eventos) nunca quedan bloqueados.

Si `ADMIN_EMAILS` está vacía, la instancia se comporta exactamente como antes: no hay aprobación y todas las cuentas funcionan.

Un usuario es **administrador** solo cuando se cumplen las dos condiciones:

- La petición viene del panel web, autenticada con un token de sesión de Firebase (no con una clave de API).
- El correo de esa sesión está en `ADMIN_EMAILS` **y Firebase lo marca como verificado** (`email_verified`).

---

## Configurar `ADMIN_EMAILS`

`ADMIN_EMAILS` es una lista de correos separados por comas. No distingue mayúsculas de minúsculas. Es una variable **de la API** (en tiempo de ejecución); el frontend web no la necesita.

```env
ADMIN_EMAILS=tu-correo@ejemplo.com,otra-persona@ejemplo.com
```

### Local (`docker compose`)

1. Abre el archivo `.env` en la raíz del repositorio (la copia de `.env.example`).
2. Agrega la línea `ADMIN_EMAILS=...` con tu correo.
3. Reinicia la API para que lea el nuevo valor:

   ```bash
   docker compose up -d api
   ```

### Coolify

1. Abre el recurso de httpSMS en Coolify y entra a **Environment Variables**.
2. Agrega `ADMIN_EMAILS` con tu correo (o varios separados por comas). Es una variable normal de ejecución, **no** una *Build Variable*.
3. Guarda y vuelve a desplegar (o reinicia) el servicio `api`.

Si despliegas la API como aplicación independiente en Coolify, agrega la variable en esa aplicación.

> Usa exactamente el mismo correo con el que inicias sesión en el panel (el de tu cuenta de Google, por ejemplo).

---

## Por qué el correo del administrador debe estar verificado

Firebase permite registrarse con correo y contraseña **sin demostrar que el correo es tuyo**. Si la API confiara solo en el correo, cualquiera podría registrarse con tu dirección (por ejemplo `tu-correo@ejemplo.com`) antes que tú, o en paralelo, y obtendría permisos de administrador automáticamente.

Por eso la API exige que Firebase marque el correo como verificado. Una dirección verificada significa que quien inició sesión controla ese buzón (o la cuenta de Google asociada), así que no basta con escribir tu correo en un formulario.

En la práctica:

- **Inicio de sesión con Google:** Google ya entrega el correo verificado. Es la opción más sencilla para los administradores.
- **Correo y contraseña:** hay que verificar el correo (con el enlace que envía Firebase) antes de tener permisos de administrador. Después de verificarlo, cierra sesión y vuelve a entrar para que el panel reciba un token nuevo con el correo verificado.
- **GitHub:** Firebase no siempre marca como verificado el correo que llega de GitHub. No lo uses para la cuenta de administrador.

Si inicias sesión con un correo de `ADMIN_EMAILS` que todavía no está verificado, la cuenta se trata como la de cualquier otro usuario (queda pendiente). En cuanto el correo esté verificado y vuelvas a iniciar sesión, tendrás acceso de administrador aunque la cuenta se haya creado como pendiente.

---

## Administrar usuarios desde el panel

Los administradores tienen una sección nueva en el panel: **`/admin/users`**.

Ahí puedes:

- **Ver la lista de usuarios**, filtrar por estado (pendientes, activos, suspendidos) y buscar por correo.
- **Aprobar** una cuenta pendiente: pasa a `active` y puede usar la API de inmediato.
- **Suspender** una cuenta: pasa a `suspended` y deja de funcionar, incluidas sus claves de API y sus teléfonos. Puedes reactivarla después.
- **Reactivar** una cuenta suspendida: vuelve a `active`.
- **Borrar** una cuenta con todos sus datos (teléfonos, mensajes, claves, webhooks…). **No se puede deshacer.**

Algunas reglas para evitar accidentes:

- Un administrador no puede cambiar su propio estado ni borrarse a sí mismo desde esta pantalla.
- Un administrador no puede suspender ni borrar a otro administrador. Para quitarle permisos a alguien, saca su correo de `ADMIN_EMAILS` y reinicia la API.

Estas acciones también están disponibles en la API (solo con la sesión de un administrador, nunca con claves de API):

| Método | Ruta | Uso |
| --- | --- | --- |
| `GET` | `/v1/admin/users?status=&query=&skip=&limit=` | Listar y buscar usuarios |
| `PUT` | `/v1/admin/users/{userID}/status` | Cambiar el estado: `{"status": "active"}` o `{"status": "suspended"}` |
| `DELETE` | `/v1/admin/users/{userID}` | Borrar el usuario y sus datos |

Respuestas: `401` sin autenticación, `403` si quien llama no es administrador (o usa una clave de API), `422` si intentas actuar sobre ti mismo o sobre otro administrador.

> **¿Suspender o borrar?** Si alguien no debería usar tu instancia, **suspéndelo**. Si lo borras y vuelve a iniciar sesión, se crea una cuenta nueva (pendiente) y te volverá a aparecer en la lista. Borra solo cuando quieras eliminar sus datos.

---

## Qué ven los usuarios pendientes o suspendidos

Cuando un usuario pendiente o suspendido inicia sesión, el panel no muestra las páginas normales. En su lugar ve un aviso:

- **Pendiente:** su cuenta espera la aprobación de un administrador. No hace falta que haga nada; cuando la aprueben podrá usar el panel con normalidad (puede tener que recargar la página).
- **Suspendida:** su cuenta fue suspendida por un administrador.

En ambos casos puede cerrar sesión o borrar su cuenta. Si alguien está esperando tu aprobación, conviene avisarle cuando la des.

---

## Efecto en claves de API y teléfonos

La aprobación aplica a **todas** las formas de autenticarse, no solo al panel:

- **Clave de API del usuario** (`x-api-key` o `Authorization: Bearer <clave>`): recibe `403` con `account_pending` o `account_suspended`. Solo funcionan `GET` y `DELETE /v1/users/me`.
- **Claves de API de teléfono** (las que usa la app de Android, empiezan con `pk_`): también reciben `403`. La app no puede registrar el token de notificaciones, enviar latidos (*heartbeats*), recoger mensajes pendientes ni reportar SMS recibidos.

Consecuencias prácticas:

- Un usuario pendiente o suspendido no puede enviar ni recibir SMS a través de httpSMS.
- Los SMS que lleguen al teléfono mientras la cuenta está suspendida **no se guardan** en httpSMS.
- Al suspender una cuenta que ya tenía teléfonos funcionando, esos teléfonos dejan de sincronizar. Al reactivarla vuelven a funcionar sin reconfigurar nada.
- Una clave de API **nunca** da permisos de administrador, aunque pertenezca a un usuario cuyo correo esté en `ADMIN_EMAILS`. Así, una clave filtrada no sirve para aprobar cuentas ni para borrar usuarios.

---

## Al actualizar desde una versión anterior

Al desplegar esta versión, la migración de la base de datos marca **todas las cuentas existentes como `active`**. Nadie pierde acceso por actualizar, y el usuario interno del sistema sigue funcionando.

Eso también significa que cualquier cuenta que alguien haya creado antes de activar la aprobación sigue activa. Después de configurar `ADMIN_EMAILS`:

1. Entra a `/admin/users` con tu cuenta de administrador.
2. Revisa la lista de usuarios activos.
3. **Suspende las cuentas que no reconozcas.** Si estás seguro de que no las necesitas, puedes borrarlas.

Mientras `ADMIN_EMAILS` siga vacía, el estado de las cuentas no tiene ningún efecto.

---

## Configuración recomendada de Firebase

En la consola de Firebase, en **Authentication**:

- **Sign-in method:** deja activados solo los métodos que realmente usas (Google, correo y contraseña, GitHub). Cada método extra es otra puerta de entrada.
- **Correo y contraseña:** si lo usas, pide a tus usuarios que verifiquen su correo (el panel tiene un botón en **Settings** para reenviar el enlace de verificación). Para la cuenta de administrador, la verificación es obligatoria.
- **Settings → User account linking:** mantén la opción por defecto, *Link accounts that use the same email* (una cuenta por correo). Evita que la misma dirección termine en varias cuentas con distintos proveedores.
- **Settings → Authorized domains:** deja solo `localhost` (si desarrollas en local) y los dominios reales de tu panel.
- **Cuenta de Google del administrador:** actívale la verificación en dos pasos. Quien controle esa cuenta controla tu instancia.

Firebase seguirá permitiendo que cualquiera se registre; la aprobación no lo impide, solo evita que esas cuentas usen la instancia. Si quieres cerrar el registro por completo, puedes desactivar la creación de cuentas desde la configuración de Firebase/Identity Platform (*User actions → Enable create (sign-up)*), pero entonces tendrás que crear tú cada cuenta nueva desde la consola de Firebase.

---

## Desactivar la aprobación

Deja `ADMIN_EMAILS` vacía (o quítala) y reinicia la API. La instancia vuelve al comportamiento anterior: todas las cuentas funcionan, incluidas las pendientes y suspendidas, y la sección `/admin/users` deja de estar disponible porque ya no hay administradores.

El estado de cada cuenta se conserva en la base de datos. Si vuelves a configurar `ADMIN_EMAILS`, las cuentas pendientes o suspendidas vuelven a quedar bloqueadas.

---

## Notas de seguridad

- **El administrador se decide en cada petición**, a partir del token de Firebase verificado. No es un campo guardado que alguien pueda modificar: depende solo de `ADMIN_EMAILS` y de que el correo esté verificado.
- **Protege dónde se configura `ADMIN_EMAILS`.** Quien pueda editar las variables de entorno de la API (tu `.env`, el panel de Coolify) puede hacerse administrador. Limita el acceso a Coolify y a tu servidor.
- **Mantén la lista corta.** Cada correo de `ADMIN_EMAILS` puede aprobar, suspender y borrar cuentas.
- **Las claves de API nunca son administradoras.** Si se filtra una clave de API (incluida la de un administrador), no sirve para usar `/v1/admin/...`. Aun así, rótala desde **Settings**.
- **Suspender bloquea el uso de la API desde las siguientes peticiones**, pero no cierra la sesión de Firebase del usuario: seguirá pudiendo iniciar sesión y verá el aviso de cuenta suspendida.
- **Borrar es irreversible.** Antes de borrar una cuenta, considera suspenderla.
- **No compartas tu cuenta de administrador.** Si otra persona necesita administrar, agrega su correo a `ADMIN_EMAILS` en lugar de prestarle tu sesión.
- Revisa de vez en cuando la lista de usuarios pendientes: las cuentas que nadie va a aprobar pueden borrarse.
