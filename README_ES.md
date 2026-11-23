[![Discord](https://img.shields.io/discord/1095778291488653372?label=Discord)](https://discord.gg/kGk8HVqeEZ)

*🌎 [Read in English](README.md)*

httpSMS es un servicio gratuito, de código abierto y autoalojado que te permite usar tu teléfono Android como un Gateway de SMS para enviar y recibir mensajes SMS. Este repositorio es un fork 100 % gratuito de [NdoleStudio/httpsms](https://github.com/NdoleStudio/httpsms): sin planes de pago, sin límites y sin depender del servicio alojado.
Realizas una solicitud a una API HTTP simple y esto activa tu teléfono Android para enviar un SMS. Los mensajes SMS recibidos en tu teléfono android también pueden ser reenviados a tu endpoint de webhook.

Guía de Inicio Rápido: abre la sección **/docs** de tu panel web (por ejemplo `https://tu-dominio/docs`). El contenido está en [`web/docs-content/`](./web/docs-content), y la referencia de la API (Swagger) la sirve la propia API en `/index.html` (en español: `/es/index.html`).

<img width="1115" alt="header" src="https://user-images.githubusercontent.com/4196457/194767449-f12d84a0-22f1-4787-afb2-17398fb459f6.png">

## Tabla de Contenidos

- [¿Por qué?](#por-qué)
- [Interfaz Web](#interfaz-web)
- [API](#api)
- [Aplicación Android](#aplicación-android)
- [Chat/foro](#chatforo)
- [Características](#características)
  - [Cifrado de extremo a extremo](#cifrado-de-extremo-a-extremo)
  - [Webhook](#webhook)
  - [Presión de retroceso (Back Pressure)](#presión-de-retroceso-back-pressure)
  - [Expiración de Mensajes](#expiración-de-mensajes)
- [Clientes API](#clientes-api)
- [Flujos](#flujos)
  - [Enviando un Mensaje SMS](#enviando-un-mensaje-sms)
- [Configuración de Auto-Alojamiento - Docker](#configuración-de-auto-alojamiento---docker)
- [Pruebas de Integración](#pruebas-de-integración)
- [Licencia](#licencia)

## ¿Por qué?

Soy originario de Camerún y quería una forma automatizada de enviar y recibir mensajes SMS utilizando una API.
Desafortunadamente, muchos países no soportan la capacidad de comprar números de teléfono virtuales, y no pude encontrar una buena solución lista para usar que me ayudara a enviar/recibir mensajes SMS usando un teléfono móvil utilizando una API http intuitiva.

## Interfaz Web

La interfaz web (servida en tu `APP_URL`) está construida utilizando [Nuxt](https://nuxtjs.org/) y [Vuetify](https://vuetifyjs.com/en/).
Está alojada como una aplicación de una sola página (SPA) en firebase. El código fuente está en el directorio [web](./web).

## API

La API (servida en tu `API_BASE_URL`) está construida utilizando [Fiber](https://gofiber.io/), Go y [CockroachDB](https://www.cockroachlabs.com/) para la base de datos.
Se ejecuta en Docker junto a la interfaz web. La referencia interactiva de la API (Swagger) la sirve la propia API en `/index.html` (en español: `/es/index.html`).

```go
// Enviando un Mensaje SMS usando Go
client := httpsms.New(
    httpsms.WithAPIKey(/* Clave API de la página de Configuración de tu panel */),
    httpsms.WithBaseURL("https://api.tu-dominio.com"),
)

client.Messages.Send(context.Background(), &httpsms.MessageSendParams{
    Content: "Este es un mensaje de texto de prueba",
    From:    "+18005550199",
    To:      "+18005550100",
})
```

## Aplicación Android

[La Aplicación Android](https://github.com/Ciudad-Valles-Open-Source/httpsms/releases/latest/download/HttpSms.apk) es una aplicación nativa construida utilizando Kotlin con principios de material design.
Esta aplicación debe estar instalada en un teléfono Android antes de que puedas comenzar a enviar y recibir mensajes SMS.

[<img src=".github/ghbadge.png" alt="Consíguelo en GitHub" height="80">](https://github.com/Ciudad-Valles-Open-Source/httpsms/releases/)

## Chat/foro

Hay varias formas de ponerte en contacto conmigo y/o con el resto de la comunidad. Siéntete libre de usar cualquiera de estos métodos. El que mejor funcione para ti:

- [Servidor de Discord](https://discord.gg/kGk8HVqeEZ) - chat directo con la comunidad
- [Problemas en GitHub (Issues)](https://github.com/Ciudad-Valles-Open-Source/httpsms/issues) - preguntas, características, errores

## Características

### Cifrado de extremo a extremo

Puedes cifrar tus mensajes de extremo a extremo utilizando el algoritmo [Cifrado AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) de grado militar. Tu clave de cifrado se almacena únicamente en tu teléfono móvil, por lo que incluso el servidor no tendrá forma de ver el contenido de tus mensajes SMS que se envían y reciben en tu teléfono Android.

### Webhook

Si quieres construir integraciones avanzadas, soportamos webhooks. La plataforma httpSMS puede reenviar mensajes SMS recibidos en el teléfono android a tu servidor utilizando una URL de callback que tú proporcionas.

### Presión de retroceso (Back Pressure)

Para no abusar de la API de SMS en android, puedes establecer un límite de tasa, por ejemplo 3 mensajes por minuto. De tal forma que incluso si llamas a la API para enviar mensajes a 100 personas, sólo enviará los mensajes a un ritmo de 3 mensajes por minuto.

### Expiración de Mensajes

A veces ocurre que el teléfono no recibe la notificación push a tiempo y no puede enviar el mensaje SMS. Es posible establecer un tiempo de espera durante el cual un mensaje es válido y si un mensaje expira después de que transcurra el tiempo de espera, serás notificado.

## Clientes API

- [x] Go: https://github.com/NdoleStudio/httpsms-go
- [x] JavaScript/TypeScript: https://github.com/NdoleStudio/httpsms-node

## Flujos

### Enviando un Mensaje SMS

```mermaid
sequenceDiagram
Usuario->>+httpSMS API: Llama a la API /v1/messages/send
httpSMS API-->>+Cola Push: Programa la notificación sobre un nuevo mensaje
httpSMS API-->>-Usuario: Responde con 202 (Aceptado)
Cola Push-->>+httpSMS API: [Asíncrono] Envía solicitud de notificación
httpSMS API-->>-Aplicación Android: Envía notificación push sobre el nuevo mensaje
Aplicación Android-->>httpSMS API: [Asíncrono] Obtiene mensaje
Aplicación Android-->>Aplicación Android: Envía Mensaje usando la API de SMS de Android
Aplicación Android-->>httpSMS API: [Asíncrono] Envía resultado del envío de SMS
Aplicación Android-->>httpSMS API: [Asíncrono] Envía Reporte de Entrega
```

## Configuración de Auto-Alojamiento - Docker

La aplicación se distribuye como imágenes de Docker y se puede ejecutar en cualquier host con Docker Engine o desplegarse en plataformas como Coolify, Portainer, CapRover, y Docker Swarm.

La documentación completa de auto-alojamiento está organizada en el directorio [`docs/`](./docs):

| Documento | Descripción |
|---|---|
| [docs/self-hosting_ES.md](./docs/self-hosting_ES.md) | Guía completa de configuración paso a paso |
| [docs/environment-variables_ES.md](./docs/environment-variables_ES.md) | Referencia para todas las variables de entorno |
| [docs/deployment-platforms_ES.md](./docs/deployment-platforms_ES.md) | Instrucciones específicas de plataforma (Coolify, Portainer, CapRover, Swarm) |
| [docs/docker-configuration_ES.md](./docs/docker-configuration_ES.md) | Referencia técnica de Dockerfile y docker-compose.yml |
| [docs/secrets_ES.md](./docs/secrets_ES.md) | Manejo de secretos: qué nunca subir, Coolify, GitHub Secrets y el escáner de pre-commit |
| [docs/realtime_ES.md](./docs/realtime_ES.md) | Tiempo real con Soketi: configuración, `SOKETI_KEY` alfanumérica y diagnóstico |

### Inicio Rápido (desarrollo local)

```bash
git clone https://github.com/Ciudad-Valles-Open-Source/httpsms.git
cd httpsms
cp .env.example .env
# Edita .env con tus credenciales de Firebase, SMTP y Cloudflare Turnstile
docker compose up --build
```

La API ejecuta automáticamente migraciones de base de datos en el primer arranque. El usuario de sistema interno requerido por la cola de eventos también se crea automáticamente en el primer arranque. Revisa los registros del servicio `api` para ver las credenciales generadas si `EVENTS_QUEUE_USER_ID` y `EVENTS_QUEUE_USER_API_KEY` no se establecieron antes de iniciar.

Cuando el stack se esté ejecutando:

- Interfaz Web: http://localhost:3000
- API: http://localhost:8000
- Swagger UI: http://localhost:8000/index.html

Para instrucciones detalladas de configuración incluyendo la configuración de Firebase, configuración SMTP y pasos de construcción de la aplicación Android, consulta [docs/self-hosting_ES.md](./docs/self-hosting_ES.md).

## Pruebas de Integración

El proyecto incluye pruebas de integración de extremo a extremo que validan el ciclo de vida completo de envío/recepción de SMS. Las pruebas ejecutan todo el stack (API, PostgreSQL, Redis) en Docker junto con un emulador de teléfono que simula un dispositivo Android.

Documentación completa: [`tests/README.md`](tests/README.md)

**Ejecución rápida:**

```bash
cd tests
bash generate-firebase-credentials.sh
export FIREBASE_CREDENTIALS=$(jq -c . firebase-credentials.json)
docker compose up -d --build --wait
docker compose wait seed && sleep 2
go test -v -timeout 120s ./...
docker compose down -v
```

Las pruebas de integración también se ejecutan automáticamente en CI en cada push/PR a `main`.

## Licencia

Este proyecto está licenciado bajo la licencia GNU AFFERO GENERAL PUBLIC LICENSE Versión 3 - consulta el archivo [LICENSE](LICENSE) para más detalles
