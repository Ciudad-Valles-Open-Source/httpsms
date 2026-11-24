# Integración continua (GitHub Actions)

El repositorio tiene dos workflows en `.github/workflows/`. Ninguno despliega tu servidor: **Coolify despliega la API y la web directamente desde el repositorio**. GitHub Actions solo verifica el código y publica el APK de Android.

| Workflow | Cuándo corre | Qué hace |
| --- | --- | --- |
| `ci.yaml` | Cada push a `main` y cada pull request | Escanea secretos, prueba la API, corre las pruebas de integración, compila la web y la app Android, y construye las imágenes Docker |
| `release.yaml` | Al subir una etiqueta `v*` (o a mano) | Compila el APK firmado y lo publica como `HttpSms.apk` en una Release de GitHub |

---

## `ci.yaml`: qué revisa cada job

| Job | Qué verifica | ¿Necesita secretos? |
| --- | --- | --- |
| **Secret Scan** | Que ningún commit del push/PR contenga claves, contraseñas, keystores o credenciales (`scripts/check-secrets.sh --range`). Es la red de seguridad por si alguien se salta el hook local. | No |
| **API (Go)** | `gofmt`, `go vet`, compilación y pruebas unitarias (incluye la cobertura de la traducción de Swagger) | No |
| **Integration Tests** | Levanta la API con emuladores (Firebase, Pub/Sub, etc.) en Docker y corre las pruebas de punta a punta. Las credenciales de Firebase se generan al vuelo y se descartan. | No |
| **Web (Nuxt)** | Instalación con lockfile congelado, lint (eslint, stylelint, prettier), pruebas y `pnpm run generate` | No (variables opcionales) |
| **Android App** | Pruebas unitarias y compilación del APK de depuración | Opcional: `GOOGLE_SERVICES_JSON_B64` |
| **Docker Images** | Que los `Dockerfile` de la API y la web construyan (no publica nada) | No |

> **Pull requests de forks:** GitHub no les pasa secretos. Por eso el job de Android usa un `google-services.json` provisional que solo sirve para compilar, y el CI funciona igual.

---

## `release.yaml`: publicar una versión del APK

1. Configura los secretos de la tabla de abajo (una sola vez).
2. Crea y sube una etiqueta:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. El workflow compila el APK firmado, corre las pruebas y crea la Release `v1.0.0` con `HttpSms.apk` y su suma `HttpSms.apk.sha256`.

Como siempre se llama `HttpSms.apk`, el enlace
`https://github.com/<tu-org>/<tu-repo>/releases/latest/download/HttpSms.apk`
apunta siempre a la última versión. Es el enlace que usa por defecto el botón de descarga del panel (`APP_DOWNLOAD_URL`).

También puedes lanzarlo desde **Actions → Release Android App → Run workflow**, indicando una etiqueta que ya exista.

---

## Secretos y variables de GitHub

Créalos en **Settings → Secrets and variables → Actions**. Los *secrets* nunca se muestran en los logs; las *variables* son valores públicos.

### Secrets

| Nombre | Obligatorio para | Cómo obtenerlo |
| --- | --- | --- |
| `GOOGLE_SERVICES_JSON_B64` | Release (opcional en CI) | `base64 -i android/app/google-services.json` |
| `ANDROID_KEYSTORE_B64` | Release | `base64 -i ruta/a/tu-release.jks` |
| `ANDROID_KEYSTORE_PASSWORD` | Release | La contraseña del keystore |
| `ANDROID_KEY_ALIAS` | Release | El alias de la llave dentro del keystore |
| `ANDROID_KEY_PASSWORD` | Release | La contraseña de esa llave |
| `HTTPSMS_SENTRY_DSN` | Opcional | Tu DSN de Sentry, si quieres reportes de fallos |
| `HTTPSMS_LOGS_URL` / `HTTPSMS_LOGS_TOKEN` | Opcional | Tu endpoint y token para registros remotos |

En Linux usa `base64 -w0 archivo` para obtener una sola línea. En macOS, `base64 -i archivo | pbcopy` lo copia directo al portapapeles.

### Variables (públicas)

| Nombre | Para qué |
| --- | --- |
| `HTTPSMS_DEFAULT_SERVER_URL` | URL de tu API que trae precargada la app (por ejemplo `https://apisms.tu-dominio.com`) |
| `API_BASE_URL`, `APP_URL` | Opcionales: con qué URLs compila la web en el CI (si no, usa `example.com`) |

### Crear un keystore (si aún no tienes uno)

```bash
keytool -genkeypair -v -keystore httpsms-release.jks -alias httpsms \
  -keyalg RSA -keysize 4096 -validity 10000
```

> **Guárdalo fuera del repositorio y haz un respaldo seguro.** Si lo pierdes, no podrás publicar actualizaciones que se instalen sobre la app existente: los usuarios tendrían que desinstalarla primero. El `.gitignore` y el escáner de secretos impiden subirlo por accidente.

---

## Protecciones recomendadas en GitHub

- **Settings → Code security:** activa *Secret scanning* y *Push protection*.
- **Settings → Branches → `main`:** exige que pasen los checks `Secret Scan`, `API (Go)`, `Web (Nuxt)` y `Android App` antes de fusionar.
- **`.github/CODEOWNERS`:** agrega tu usuario para que los cambios a los workflows requieran tu revisión.

Ver también [secrets_ES.md](./secrets_ES.md).
