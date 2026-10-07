# Guía de Configuración de Auto-Alojamiento

*🌎 [Read in English](self-hosting.md)*

Esta guía detalla los pasos requeridos para ejecutar httpSMS utilizando tus propios recursos de infraestructura y servicios de terceros. Desplegar toda la pila (stack) utilizando Docker Compose garantiza entornos reproducibles con una mínima configuración a nivel de host.

## Requisitos Previos

- Docker Engine y Docker Compose instalados.
- Un proyecto de Firebase (para notificaciones de la aplicación Android y autenticación web).
- Un servidor SMTP (para enviar alertas y notificaciones del sistema).
- Cloudflare Turnstile (para protección de rutas API).

---

## Paso 1 — Configurar Firebase

httpSMS usa Firebase para dos propósitos:

- **Firebase Cloud Messaging (FCM):** Envía notificaciones push a la aplicación Android para activar el envío de SMS salientes.
- **Autenticación de Firebase:** Autentica usuarios en la interfaz web.

### Crear un proyecto de Firebase

1. Abre la [Consola de Firebase](https://console.firebase.google.com/) y crea un nuevo proyecto.
2. Registra una **Web App** dentro del proyecto. Sigue la [Guía de configuración de la Web App](https://firebase.google.com/docs/web/setup#register-app) y guarda los valores de configuración del SDK:

```js
const firebaseConfig = {
  apiKey: "AIzaSy...",
  authDomain: "your-project.firebaseapp.com",
  projectId: "your-project",
  storageBucket: "your-project.appspot.com",
  messagingSenderId: "123456789",
  appId: "1:123456789:web:abc123",
  measurementId: "G-XXXXX",
};
```

Estos valores son utilizados por el frontend web en tiempo de construcción (build time).

### Habilitar inicio de sesión por Correo/Contraseña

1. En la Consola de Firebase, ve a **Autenticación > Método de inicio de sesión** (Sign-in method).
2. Habilita **Correo electrónico/Contraseña** y haz clic en **Guardar**.

> **Problema Conocido:** El proveedor de inicio de sesión por correo/contraseña de Firebase tiene un error con la protección contra enumeración de correos que puede impedir el inicio de sesión. Si los usuarios no pueden iniciar sesión, [deshabilita la protección de enumeración de correo electrónico](https://cloud.google.com/identity-platform/docs/admin/email-enumeration-protection#disable) en la configuración de Identity Platform.

### Generar una clave de cuenta de servicio

1. Sigue la [Guía de configuración del SDK de Administración de Firebase](https://firebase.google.com/docs/admin/setup#initialize_the_sdk_in_non-google_environments).
2. Descarga el archivo JSON de la cuenta de servicio.
3. Conviértelo a una cadena de una sola línea para usarlo en la variable de entorno `FIREBASE_CREDENTIALS`:

```bash
cat firebase-credentials.json | tr -d '\n'
```

### Generar google-services.json para la aplicación Android

Sigue las [Instrucciones de configuración de Android](https://support.google.com/firebase/answer/7015592?hl=es#android) para descargar `google-services.json`. Este archivo es necesario al compilar la aplicación de Android en el Paso 8.

---

## Paso 2 — Configurar SMTP

httpSMS envía correos electrónicos transaccionales para eventos como cuando un teléfono Android ha estado fuera de línea durante un período prolongado.

Cualquier servidor SMTP es compatible. Para desarrollo local, [Mailtrap](https://mailtrap.io/) proporciona un entorno de pruebas gratuito. Para producción, servicios como Mailgun, Postmark o SendGrid son adecuados.

Registra los siguientes valores para su uso en variables de entorno:

- `SMTP_HOST`
- `SMTP_PORT` (típicamente `587` para TLS o `465` para SSL)
- `SMTP_USERNAME`
- `SMTP_PASSWORD`
- `SMTP_FROM_EMAIL` (la dirección del remitente)

---

## Paso 3 — Configurar Cloudflare Turnstile

La ruta `/v1/messages/search` está protegida por un CAPTCHA de [Cloudflare Turnstile](https://developers.cloudflare.com/turnstile/get-started/) para prevenir abusos.

1. Abre el [Panel de Cloudflare](https://dash.cloudflare.com/) y navega a **Turnstile**.
2. Crea un nuevo widget para el sitio. Para desarrollo local, usa `localhost` como el dominio.
3. Registra la **Clave del Sitio** (`CLOUDFLARE_TURNSTILE_SITE_KEY`, usada por el frontend web) y la **Clave Secreta** (`CLOUDFLARE_TURNSTILE_SECRET_KEY`, usada por la API).

---

## Paso 4 — Clonar el Repositorio

```bash
git clone https://github.com/NdoleStudio/httpsms.git
cd httpsms
```

---

## Paso 5 — Establecer Variables de Entorno

Copia la plantilla proporcionada y rellena todos los valores requeridos:

```bash
cp .env.example .env
```

Abre `.env` y configura las siguientes variables requeridas:

```dotenv
# Base de datos
POSTGRES_PASSWORD=<a-strong-random-password>

# Firebase (API — lado del servidor)
GCP_PROJECT_ID=<your-firebase-project-id>
FIREBASE_CREDENTIALS=<single-line-json-from-step-1>

# Firebase (Web — inyectado en tiempo de build)
FIREBASE_API_KEY=
FIREBASE_AUTH_DOMAIN=
FIREBASE_PROJECT_ID=
FIREBASE_STORAGE_BUCKET=
FIREBASE_MESSAGING_SENDER_ID=
FIREBASE_APP_ID=
FIREBASE_MEASUREMENT_ID=

# SMTP
SMTP_HOST=
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM_EMAIL=

# Cloudflare Turnstile
CLOUDFLARE_TURNSTILE_SECRET_KEY=
CLOUDFLARE_TURNSTILE_SITE_KEY=

# URLs de Aplicación (actualiza a tu dominio para producción)
APP_URL=http://localhost:3000
API_BASE_URL=http://localhost:8000
```

Para una referencia completa de todas las variables disponibles, consulta [environment-variables_ES.md](./environment-variables_ES.md).

---

## Paso 6 — Compilar e Iniciar

```bash
docker compose up --build
```

La primera compilación descarga las imágenes base y compila el binario Go; esto puede llevar varios minutos. Cuando se complete:

- Interfaz Web: [http://localhost:3000](http://localhost:3000)
- API: [http://localhost:8000](http://localhost:8000)
- Swagger UI: [http://localhost:8000/index.html](http://localhost:8000/index.html)

---

## Paso 7 — Inicialización del Usuario del Sistema

httpSMS usa un usuario del sistema interno dedicado para procesar eventos asincrónicos a través de la cola de eventos. Este usuario no corresponde a ningún operador humano.

**Este paso es completamente automático.** Cuando la API se inicia contra una base de datos vacía, detecta la ausencia de usuarios y crea el usuario del sistema automáticamente a partir de los valores de `EVENTS_QUEUE_USER_ID` y `EVENTS_QUEUE_USER_API_KEY`. Si esas variables no están configuradas, se generan valores seguros y se imprimen en los registros del contenedor de la API:

```
╔══════════════════════════════════════════════════════════════╗
║              USUARIO DEL SISTEMA CREADO EXITOSAMENTE         ║
╠══════════════════════════════════════════════════════════════╣
║  ID Usuario: a1b2c3d4-e5f6-7890-abcd-ef1234567890            ║
║  Clave API : uk_Abc123...                                     ║
║  Correo    : system@httpsms.local                             ║
╠══════════════════════════════════════════════════════════════╣
║  Guarda estos valores en tus variables de entorno:           ║
║    EVENTS_QUEUE_USER_ID                                    ║
║    EVENTS_QUEUE_USER_API_KEY                               ║
║                                                            ║
║  Si se autogeneró, configúralos antes del próximo reinicio   ║
║  para evitar crear un nuevo usuario del sistema.           ║
╚══════════════════════════════════════════════════════════════╝
```

**Después del primer arranque**, recupera estos valores de los registros y persiste en tu configuración de entorno:

```dotenv
EVENTS_QUEUE_USER_ID=<valor-de-los-logs>
EVENTS_QUEUE_USER_API_KEY=<valor-de-los-logs>
```

Reinicia el servicio de API para que recoja los cambios:

```bash
docker compose restart api
```

Si prefieres predefinir estos valores antes del primer arranque, establécelos en tu archivo `.env` de antemano. El proceso de inicialización (seed) usará los valores proporcionados en lugar de autogenerarlos.

---

## Paso 8 — Compilar la Aplicación Android

Antes de compilar la aplicación Android, reemplaza `android/app/google-services.json` con el archivo descargado en el Paso 1. Esto habilita las notificaciones push de FCM desde tu propio proyecto de Firebase.

Abre el proyecto en [Android Studio](https://developer.android.com/studio) y compila un APK de release. Instálalo en el teléfono Android que funcionará como tu Gateway SMS.

---

## Verificación Post-Despliegue

1. Abre la interfaz web y crea una cuenta de usuario usando el formulario de registro.
2. Navega a **Configuración** y verifica que se muestre tu Clave API.
3. Abre la aplicación de Android, inicia sesión con las mismas credenciales, y confirma que el teléfono aparezca en el panel de control.
4. Envía un mensaje SMS de prueba a través de la API o la interfaz web.

Para instrucciones de despliegue específicas de la plataforma (Coolify, Portainer, CapRover, Docker Swarm), consulta [deployment-platforms_ES.md](./deployment-platforms_ES.md).
