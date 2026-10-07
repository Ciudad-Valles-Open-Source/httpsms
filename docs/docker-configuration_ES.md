# Configuración de Docker

*🌎 [Read in English](docker-configuration.md)*

Este documento detalla el diseño de las imágenes de Docker del proyecto y el archivo `docker-compose.yml`. Está destinado a ingenieros que necesiten modificar el proceso de construcción, integrar httpSMS en una canalización (pipeline) de CI/CD existente, o depurar el despliegue.

## Vista General de Archivos

| Archivo | Descripción |
|---|---|
| [`api/Dockerfile`](../api/Dockerfile) | Construye el binario Go de la API y lo empaqueta en una imagen mínima Alpine |
| [`web/Dockerfile`](../web/Dockerfile) | Construye el frontend web estático Nuxt y lo sirve utilizando nginx |
| `docker-compose.yml` | Orquesta los cuatro servicios para despliegues locales y en producción |
| `.env.example` | Plantilla que enumera cada variable configurable con su documentación |

---

## Dockerfile de la API

**Ubicación:** [`api/Dockerfile`](../api/Dockerfile)

### Etapas de compilación (Build stages)

La imagen de la API utiliza una compilación de dos etapas:

**Etapa 1 — Constructor (Builder) (`golang:1.26-alpine`)**

- Descarga las dependencias del módulo de Go antes de copiar el código fuente, aprovechando la caché de capas de Docker. Las compilaciones subsecuentes reutilizan la capa de dependencias si `go.mod` y `go.sum` no han cambiado.
- Instala y ejecuta `swag` para generar la especificación de Swagger a partir de las anotaciones en el código fuente de Go. El resultado generado se incrusta en el binario.
- Compila la aplicación con `CGO_ENABLED=0` para producir un binario vinculado estáticamente. Las banderas del enlazador `-s -w` eliminan los símbolos de depuración e información DWARF, reduciendo el tamaño del binario.
- El argumento de compilación `GIT_COMMIT` se inyecta en el binario en tiempo de compilación a través de `-X main.Version`, haciendo que la versión desplegada sea identificable en la UI de Swagger y en los registros.

**Etapa 2 — Tiempo de ejecución (Runtime) (`alpine:3.20`)**

- Crea un usuario y grupo no root (`httpsms`) para el proceso.
- Copia solo el binario compilado y el archivo de certificado `root.crt` de la etapa del constructor.
- Instala `curl` y `ca-certificates` para controles de salud (health checks) y verificación TLS respectivamente, y `tzdata` para el soporte de zona horaria.
- Ejecuta el binario con `--dotenv=false`, lo cual deshabilita la carga de archivos `.env`. Toda la configuración debe proporcionarse a través de variables de entorno del contenedor.

### Argumentos de compilación

| Argumento | Por defecto | Descripción |
|---|---|---|
| `GIT_COMMIT` | `dev` | SHA del commit de Git inyectado en la cadena de versión del binario. |

---

## Dockerfile Web

**Ubicación:** [`web/Dockerfile`](../web/Dockerfile)

### Etapas de compilación

**Etapa 1 — Constructor (Builder) (`node:22-alpine`)**

- Usa `corepack` (integrado en Node 22) para activar `pnpm` sin un paso de instalación global separado.
- Instala dependencias con `--frozen-lockfile` para garantizar compilaciones reproducibles. La compilación falla si `pnpm-lock.yaml` no está sincronizado con `package.json`.
- Acepta toda la configuración de la aplicación como argumentos de compilación de Docker (`ARG`), luego los promueve a variables de entorno (`ENV`) para que el proceso de compilación de Nuxt pueda leerlos.
- Ejecuta `pnpm run generate`, el cual ejecuta `nuxi generate`. Esto produce una salida estática completa en HTML/CSS/JS en `.output/public`.

**Etapa 2 — Tiempo de ejecución (Runtime) (`nginx:stable-alpine`)**

- Copia la salida estática de la etapa del constructor a la raíz de documentos de nginx.
- Utiliza el archivo `nginx.conf` proporcionado en el proyecto, el cual configura un único bloque de servidor en el puerto 3000 con enrutamiento `try_files` para soportar la navegación del lado del cliente en la SPA.

### Por qué se requieren argumentos de compilación

La configuración de Nuxt lee las variables de entorno en tiempo de compilación usando `process.env`. Dado que la salida es completamente estática, no hay proceso de servidor en tiempo de ejecución que pueda inyectar la configuración. Todo valor de configuración que deba estar disponible en el navegador debe estar presente cuando se ejecuta `pnpm run generate`.

Una consecuencia de este diseño es que **cambiar cualquier variable de entorno web requiere recompilar la imagen de Docker**. No hay forma de aplicar cambios de configuración web reiniciando el contenedor.

### Argumentos de compilación

Todos los argumentos de compilación corresponden directamente a las variables de entorno descritas en [environment-variables_ES.md](./environment-variables_ES.md#variables-de-entorno-del-frontend-web). Por defecto son cadenas vacías cuando no se proporcionan, lo que puede resultar en que características como la autenticación de Firebase no funcionen.

---

## docker-compose.yml

**Ubicación:** [`docker-compose.yml`](../docker-compose.yml)

### Principios de diseño

- **Toda la configuración se proporciona a través de variables de entorno.** Ninguna directiva `env_file` hace referencia a archivos que estarían ausentes en un clon del repositorio o un runner de CI/CD. Esto hace que el archivo sea compatible con plataformas de orquestación como Coolify, Portainer y Docker Swarm, que inyectan configuración desde sus propios almacenes de variables.
- **Se proporcionan valores por defecto para variables no sensibles.** Las variables usan la sintaxis `${VAR:-default}`. Esto permite que el stack se inicie para desarrollo local sin requerir que cada variable esté establecida, a la vez que permite la sobrescritura completa de la configuración en producción.
- **Se definen controles de salud (health checks) para todos los servicios.** Los servicios dependientes esperan a que sus dependencias estén sanas antes de iniciar, reduciendo condiciones de carrera en el inicio. El servicio `api` espera a que PostgreSQL acepte conexiones; el servicio `web` espera a que el servicio `api` esté sano.

### Servicios

#### postgres

Ejecuta PostgreSQL 16 con un volumen nombrado para persistencia de datos. El control de salud usa `pg_isready` para verificar que la base de datos está aceptando conexiones antes de que inicien los servicios dependientes.

#### redis

Ejecuta Redis 7 con persistencia append-only habilitada (`--appendonly yes`), que sobrevive a los reinicios del contenedor. El control de salud usa `redis-cli ping`.

#### api

Construye desde `api/Dockerfile`. Recibe toda la configuración como variables de entorno. Los valores de `DATABASE_URL` y `REDIS_URL` se construyen a partir de los conjuntos de variables de PostgreSQL y Redis para garantizar que se mantengan sincronizados con la configuración de esos servicios.

La variable `EVENTS_QUEUE_ENDPOINT` tiene el valor por defecto `http://api:8000/v1/events`, utilizando el nombre del servicio de red de Docker `api` para que la cola de eventos pueda alcanzar la API desde dentro de la red de compose.

#### web

Construye desde `web/Dockerfile`. Toda la configuración se pasa como argumentos de compilación que el Dockerfile acepta mediante declaraciones `ARG`. El servicio `web` depende de que el servicio `api` esté sano para evitar iniciar antes de que la API esté lista para atender solicitudes.

### Volúmenes

Se definen dos volúmenes nombrados:

- `postgres_data` — persiste el directorio de datos de PostgreSQL a través de reinicios y actualizaciones de contenedores.
- `redis_data` — persiste el registro append-only de Redis.

---

## Manejo de Variables de Entorno

### API (runtime)

El servidor de la API Go lee toda la configuración de variables de entorno de proceso en el momento del inicio. La bandera `--dotenv=false` pasada al binario garantiza que no intente cargar un archivo `.env` desde el disco. Este es el comportamiento esperado en despliegues contenedorizados donde el orquestador es responsable de inyectar las variables de entorno.

### Web (build time)

El bloque `runtimeConfig.public` de Nuxt en `nuxt.config.ts` lee las variables de entorno a través de `process.env` en el momento en que se ejecuta `nuxi generate`. Los valores resultantes se incrustan en el paquete JavaScript generado. Esto significa que:

- Los valores son visibles en el navegador (trátalos como públicos).
- Los valores secretos nunca deben pasarse como argumentos de compilación web.
- Los cambios de configuración requieren una nueva construcción de imagen y redistribución.

---

## Controles de Salud (Health Checks)

Cada servicio en `docker-compose.yml` define un control de salud que Docker utiliza para determinar la disponibilidad:

| Servicio | Método del Control de Salud | Condición Sano (Healthy) |
|---|---|---|
| `postgres` | `pg_isready` | Base de datos aceptando conexiones |
| `redis` | `redis-cli ping` | Devuelve `PONG` |
| `api` | `curl -f http://localhost:8000/` | Respuesta HTTP 200 |
| `web` | `wget --spider http://localhost:3000/` | Respuesta HTTP 200 |

La cadena de dependencias es: `postgres` y `redis` deben estar sanos antes de que inicie `api`; `api` debe estar sana antes de que inicie `web`.

---

## Argumentos de Compilación vs Variables en Tiempo de Ejecución

La siguiente tabla resume qué variables se consumen en tiempo de compilación frente a tiempo de ejecución para cada servicio:

| Tipo de variable | API | Web |
|---|---|---|
| Credenciales de Base de datos | Tiempo de Ejec. | No aplicable |
| Cuenta de servicio de Firebase | Tiempo de Ejec. | No aplicable |
| Config SDK Web Firebase | No aplicable | **Tiempo compilación** |
| Credenciales SMTP | Tiempo de Ejec. | No aplicable |
| URL base de la API | Tiempo Ejec. (Swagger) | **Tiempo compilación** |
| URL de la aplicación | Tiempo de Ejec. | **Tiempo compilación** |
| Configuración de Pusher | Tiempo de Ejec. | **Tiempo compilación** |

Al reconfigurar un despliegue, identifica si la variable cambiada se consume en tiempo de compilación o ejecución para determinar si es necesaria una reconstrucción completa o si basta con un simple reinicio.
