# CHANGELOG

_🌎 [Read in English](CHANGELOG.md)_

Este repositorio es un fork autoalojado y 100 % gratuito de
[NdoleStudio/httpsms](https://github.com/NdoleStudio/httpsms). La numeración
continúa después de la última versión del proyecto original (`v1.2.1`).

## [2.2.0] - 2026-11-24

### Agregado

- Los ejemplos de código de la documentación usan la misma ventana de código
  que la página de inicio, con una pestaña por lenguaje cuando hay varios
  ejemplos juntos.
- Ejemplos en Dart: enviar un SMS con Dio y recibir webhooks con Shelf.

### Cambiado

- Todas las vistas del panel tienen el mismo ancho de contenido (el de
  Contactos).
- Las tablas se convierten en una lista de tarjetas cuando sus columnas no
  caben en el ancho disponible, en cualquier tamaño de pantalla (antes solo en
  teléfonos), y vuelven a ser tabla cuando hay espacio. Los valores largos,
  como correos, ahora se ajustan.

### Corregido

- El pie de las tablas mostraba claves `$vuetify.dataFooter…` en lugar del
  texto traducido.
- Los avisos (toasts) salían en inglés en la vista en español ("Contact
  created successfully", mensajes de la API como "phone updated successfully"
  o "validation errors while …"); ahora se traducen.
- Las fechas relativas ("hace 2 minutos"), duraciones y fechas ahora siguen el
  idioma del panel en lugar de mostrarse siempre en inglés.
- El enlace de la tarjeta de cifrado en la página de inicio se salía de la
  tarjeta.
- Swagger UI ya no muestra la insignia "INVALID": se desactivó el validador en
  línea (validator.swagger.io), que no puede acceder a servidores
  autoalojados.

## [2.1.0] - 2026-11-24

### Agregado

- **Aprobación de cuentas y panel de administración.** Con la variable
  `ADMIN_EMAILS`, las cuentas nuevas quedan pendientes hasta que un
  administrador las aprueba en `/admin/users` (aprobar, suspender, reactivar o
  borrar). Las cuentas pendientes o suspendidas no pueden usar el panel, la API
  ni sus claves de teléfono. Los administradores deben iniciar sesión con un
  correo verificado y las claves API nunca dan permisos de administrador.
  Opcionalmente avisa por correo a los administradores y al usuario aprobado.
  Ver [docs/admin_ES.md](docs/admin_ES.md).
- Página `/account-pending` para cuentas pendientes o suspendidas.
- `GET /v1/users/me` incluye `status` e `is_admin`; nuevos endpoints
  `/v1/admin/users` documentados en Swagger (inglés y español).
- **Rediseño completo de la interfaz**, alineado con los demás productos
  Arquos: tipografía Plus Jakarta Sans, tarjetas, insignias, paleta azul y temas
  claro y oscuro. El panel tiene barra de navegación lateral en escritorio,
  lista y vista de conversaciones renovadas (se corrigieron los mensajes salientes
  que se cortaban a la derecha), estados de carga y vacíos, y Ajustes por
  secciones. El sitio público tiene nueva cabecera con menú móvil, portada, login,
  docs y páginas de error. Las tablas se convierten en tarjetas en el teléfono.
- Selector de apariencia en el panel (modo claro u oscuro, y paleta).
- El tema de Swagger usa la misma paleta que el panel.

### Cambiado

- El inicio de sesión ofrece solo **Google** y **correo electrónico**. Se quitó
  GitHub, que nunca estuvo habilitado en Firebase en esta instalación, por lo que
  no afecta a ninguna cuenta existente.
- Si el navegador bloquea la ventana emergente de Google, el inicio de sesión
  continúa por redirección.

### Corregido

- Los administradores perdían sus permisos porque la identidad de la clave API
  reemplazaba a la sesión de Firebase cuando el panel enviaba ambas.
- Los errores del inicio de sesión no se mostraban: ahora hay mensajes claros para
  ventanas bloqueadas, dominios no autorizados y métodos desactivados en Firebase.
- El workflow de releases fallaba al buscar pruebas unitarias de `release`.

### Seguridad

- La API ya no escribe en sus registros los tokens ni las claves API completas
  (solo `****` y los últimos 4 caracteres).

## [2.0.0] - 2026-11-24

Primera versión independiente del fork, basada en el original `v1.2.1` con sus
parches de seguridad y actualizaciones de dependencias.

### Agregado

- Documentación propia en `/docs` (español e inglés) y Swagger en `/index.html`
  y `/es/index.html` con tema oscuro y traducción completa.
- Filtro de conversaciones **Todos / No leídos** y mensajes de hasta 2048
  caracteres.
- Tiempo real con Soketi funcionando (una sola conexión WebSocket compartida).
- App Android: URL del servidor configurable al compilar y desde la app, QR que
  incluye servidor y clave, y soporte de `http://` con advertencia.
- Política de privacidad y Términos propios del fork.
- CI propio en GitHub Actions y publicación del APK firmado como release.
- Escáner de secretos en el pre-commit y en el CI, y guía de manejo de secretos.
- Variables `JWT_ISSUER` y `OTEL_SERVICE_NAME`.
- Incluye el trabajo previo del fork: traducción al español, guía de
  autoalojamiento con Docker y Coolify, y configuración del servidor en la app.

### Cambiado

- **Incompatible:** se eliminaron los planes de pago, los límites y los
  endpoints de suscripción y Lemon Squeezy; todo es gratuito y sin límites.
- **Incompatible:** el emisor (`iss`) de los JWT de webhooks y notificaciones es
  `JWT_ISSUER` (por defecto `APP_NAME` normalizado) en lugar de
  `api.httpsms.com`.
- **Incompatible:** el QR de inicio de sesión usa el formato
  `httpsms://login?server=…&key=…` (la app sigue aceptando los QR anteriores).
- `SOKETI_KEY` debe ser alfanumérica.
- PostgreSQL y Redis ya no se publican en el host (solo red interna de Docker).
- Sin dependencias del servicio alojado: sin enlaces ni recursos de httpsms.com,
  sin Microsoft Clarity, Sentry ni Axiom del proyecto original.

### Corregido

- `PUT /v1/phones` ya no borra el token FCM, el horario de envío ni la respuesta
  automática cuando no se envían en la petición.
- El monitor de heartbeat se reinicia cuando expira el historial.
- Vulnerabilidades XSS en textos con nombres de contactos, horarios y números.
- Pruebas intermitentes de contactos y de archivado de conversaciones.

[Sin publicar]: https://github.com/Ciudad-Valles-Open-Source/httpsms/compare/v2.2.0...HEAD
[2.2.0]: https://github.com/Ciudad-Valles-Open-Source/httpsms/compare/v2.1.0...v2.2.0
[2.1.0]: https://github.com/Ciudad-Valles-Open-Source/httpsms/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/Ciudad-Valles-Open-Source/httpsms/releases/tag/v2.0.0
