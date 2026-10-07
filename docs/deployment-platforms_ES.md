# Plataformas de Despliegue

*🌎 [Read in English](deployment-platforms.md)*

Este documento proporciona instrucciones específicas de la plataforma para desplegar httpSMS en plataformas comunes de orquestación de contenedores. Todas las plataformas utilizan las mismas imágenes de Docker y el mismo conjunto de variables de entorno descritas en [environment-variables_ES.md](./environment-variables_ES.md).

## Tabla de Contenidos

- [Coolify](#coolify)
- [Portainer](#portainer)
- [CapRover](#caprover)
- [Docker Swarm](#docker-swarm)
- [Docker Compose Genérico](#docker-compose-genérico)
- [Notas Importantes sobre el Frontend Web](#notas-importantes-sobre-el-frontend-web)

---

## Coolify

[Coolify](https://coolify.io/) es una alternativa auto-alojada a Heroku/Netlify. Soporta despliegues basados en Git de stacks de Docker Compose.

### Despliegue con Docker Compose (recomendado)

Este enfoque despliega los cuatro servicios (API, web, PostgreSQL, Redis) como un solo stack.

1. En Coolify, navega a **New Resource > Docker Compose**.
2. Conecta tu repositorio de Git y selecciona la rama a desplegar.
3. Coolify detectará `docker-compose.yml` automáticamente.
4. Navega al panel de **Environment Variables** y establece todas las variables requeridas. Consulta [environment-variables_ES.md](./environment-variables_ES.md) para ver la lista completa.

**Variables mínimas requeridas para un despliegue funcional:**

```dotenv
POSTGRES_PASSWORD=<strong-password>

GCP_PROJECT_ID=<firebase-project-id>
FIREBASE_CREDENTIALS=<single-line-service-account-json>

FIREBASE_API_KEY=<web-sdk-value>
FIREBASE_AUTH_DOMAIN=<web-sdk-value>
FIREBASE_PROJECT_ID=<web-sdk-value>
FIREBASE_STORAGE_BUCKET=<web-sdk-value>
FIREBASE_MESSAGING_SENDER_ID=<web-sdk-value>
FIREBASE_APP_ID=<web-sdk-value>

APP_URL=https://your-domain.com
API_BASE_URL=https://api.your-domain.com
SWAGGER_HOST=api.your-domain.com

SMTP_HOST=<smtp-host>
SMTP_USERNAME=<smtp-user>
SMTP_PASSWORD=<smtp-password>
SMTP_FROM_EMAIL=<sender-address>

CLOUDFLARE_TURNSTILE_SECRET_KEY=<secret>
CLOUDFLARE_TURNSTILE_SITE_KEY=<site-key>
```

5. En la configuración del servicio de Coolify, asigna dominios públicos a los servicios `web` (puerto 3000) y `api` (puerto 8000). Coolify provee TLS automáticamente a través de Traefik.
6. Haz clic en **Deploy**. El usuario del sistema se crea automáticamente en el primer arranque. Revisa los registros del servicio `api` para recuperar las credenciales generadas si `EVENTS_QUEUE_USER_ID` y `EVENTS_QUEUE_USER_API_KEY` no se establecieron de antemano.

### Despliegue de Dockerfile Individual

Usa este enfoque para desplegar la API y el frontend web como aplicaciones separadas en Coolify, con PostgreSQL y Redis administrados como recursos de base de datos de Coolify.

**Aplicación API:**

- Tipo de recurso: Application (Dockerfile)
- Ruta del Dockerfile: `api/Dockerfile`
- Contexto de build: `api/`
- Puerto: `8000`

**Aplicación Web:**

- Tipo de recurso: Application (Dockerfile)
- Ruta del Dockerfile: `web/Dockerfile`
- Contexto de build: `web/`
- Puerto: `3000`

Establece las variables de entorno para cada aplicación por separado. Ten en cuenta que las variables del frontend web son variables de tiempo de compilación (build-time) y deben configurarse como **Build Variables** (no como variables de entorno de tiempo de ejecución) en la configuración de la aplicación de Coolify.

---

## Portainer

[Portainer](https://www.portainer.io/) es una interfaz de usuario para gestión de contenedores para Docker y Kubernetes.

### Stack de Docker Compose

1. En Portainer, navega a **Stacks > Add Stack**.
2. Selecciona **Repository** y conecta tu repositorio de Git.
3. Establece la ruta del archivo Compose en `docker-compose.yml`.
4. En **Environment variables**, agrega todas las variables requeridas.
5. Haz clic en **Deploy the stack**.

### Gestionando stacks

Para actualizar el despliegue después de hacer push de nuevos commits:

1. Navega al stack en Portainer.
2. Haz clic en **Pull and redeploy**.

---

## CapRover

[CapRover](https://caprover.com/) proporciona una capa PaaS simple encima de Docker Swarm.

### Configuración

Crea las siguientes aplicaciones en CapRover antes de desplegar:

1. **Servicios de Base de datos** — usa las Apps de un Clic de CapRover (One-Click Apps) para crear instancias de PostgreSQL y Redis.
2. **Aplicación API** — crea una nueva aplicación llamada `httpsms-api`.
3. **Aplicación Web** — crea una nueva aplicación llamada `httpsms-web`.

### Desplegar la API

1. En la configuración de la aplicación `httpsms-api`, ve a **Deployment** y conecta tu repositorio de Git.
2. Establece la **ruta del Dockerfile** en `api/Dockerfile` y el **contexto de build** en `api/`.
3. En **App Configs**, establece todas las variables de entorno de la API. Usa los nombres de los servicios internos proporcionados por CapRover para `DATABASE_URL` y `REDIS_URL`.
4. Habilita HTTPS en la pestaña **HTTP Settings**.
5. Despliega.

### Desplegar el frontend web

1. En la configuración de la aplicación `httpsms-web`, conecta el mismo repositorio.
2. Establece la **ruta del Dockerfile** en `web/Dockerfile` y el **contexto de build** en `web/`.
3. En **App Configs**, establece todos los argumentos de build web como variables de entorno. CapRover los pasa al proceso de compilación de Docker como argumentos de compilación (build arguments) cuando los nombres de las variables coinciden con las declaraciones `ARG` en el Dockerfile.
4. Habilita HTTPS y despliega.

---

## Docker Swarm

Docker Swarm es adecuado para despliegues de producción en múltiples nodos.

### Desplegar como un stack

```bash
# Inicializar Swarm (si no está activo aún)
docker swarm init

# Desplegar el stack desde la raíz del repositorio
docker stack deploy --compose-file docker-compose.yml httpsms
```

Las variables de entorno pueden ser provistas exportándolas en el shell antes de correr el comando de despliegue, o cargando un archivo `.env`:

```bash
set -a && source .env && set +a
docker stack deploy --compose-file docker-compose.yml httpsms
```

### Gestión de secretos

Para despliegues de Swarm en producción, usa Docker Secrets en lugar de variables de entorno simples para valores sensibles. Crea secretos para `POSTGRES_PASSWORD` y `FIREBASE_CREDENTIALS`:

```bash
echo "strong-password" | docker secret create postgres_password -
cat firebase-credentials.json | tr -d '\n' | docker secret create firebase_credentials -
```

Luego referéncialos en el archivo Compose usando la clave `secrets:`. Esto requiere modificar `docker-compose.yml` para montar secretos como archivos y leerlos desde el código de la API, lo cual no está configurado por defecto.

### Actualizar un stack en ejecución

```bash
docker stack deploy --compose-file docker-compose.yml httpsms
```

Docker Swarm realiza una actualización progresiva (rolling update) automáticamente.

---

## Docker Compose Genérico

Para cualquier host Linux con Docker Engine instalado:

```bash
# Clonar el repositorio
git clone https://github.com/NdoleStudio/httpsms.git
cd httpsms

# Crear y configurar el archivo de entorno
cp .env.example .env
# Edita .env con tus valores

# Compilar e iniciar todos los servicios en segundo plano
docker compose up -d --build

# Ver los registros (logs)
docker compose logs -f api

# Detener todos los servicios
docker compose down
```

Para actualizar después de descargar el código nuevo (pull):

```bash
git pull
docker compose up -d --build
```

---

## Notas Importantes sobre el Frontend Web

La aplicación web Nuxt es una aplicación de una sola página (SPA) generada estáticamente. Toda la configuración, incluyendo la URL base de la API y las credenciales de Firebase, se incrusta en el paquete JavaScript en el momento de construir la imagen.

Esto tiene dos consecuencias prácticas:

1. **Cambiar una variable de entorno web requiere recompilar la imagen.** Reiniciar el contenedor sin recompilar no aplica los cambios de configuración al frontend web.

2. **`API_BASE_URL` debe ser la URL pública de la API**, accesible desde el navegador del usuario final. Establecerla a una dirección de red interna de Docker (como `http://api:8000`) hará que todas las peticiones del navegador fallen.

Cuando uses Docker Compose localmente, establece `API_BASE_URL=http://localhost:8000`. En producción, establécela a la URL HTTPS de tu API (por ejemplo, `https://api.tu-dominio.com`).
