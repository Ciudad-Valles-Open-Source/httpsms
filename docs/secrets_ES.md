# Manejo de secretos

Un **secreto** es cualquier valor que da acceso a algo: contraseñas, claves de API, tokens, claves privadas, credenciales de Firebase/Google, keystores de Android. **Nunca deben estar en el repositorio**, aunque sea privado: el historial de git es permanente y basta un `git push` o un fork para que queden públicos.

Este proyecto sigue una regla simple:

> **En git solo hay plantillas con valores vacíos o de ejemplo. Los valores reales viven en `.env` (ignorado por git), en las variables de entorno de Coolify o en los secretos de GitHub Actions.**

---

## Dónde va cada cosa

| Entorno | Dónde guardas los secretos | ¿Está en git? |
| --- | --- | --- |
| Tu computadora (desarrollo, `docker compose`) | `.env` en la raíz del repo (copia de `.env.example`) | No, ignorado |
| Producción (Coolify) | *Environment Variables* del recurso, marcadas como secretas | No |
| Compilación de la app Android en tu computadora | Variables de entorno o `~/.gradle/gradle.properties` | No |
| Compilación de la app Android en GitHub Actions | *Settings → Secrets and variables → Actions* | No |
| Plantillas (`.env.example`, `web/.env.docker`, `web/.env.production`, `api/.env.*`, `tests/.env.test`) | Solo nombres de variables con valores vacíos o de ejemplo | Sí |

### Variables secretas del servidor

Van en `.env` (local) y en Coolify (producción). Nunca en `.env.example` ni en los archivos `.env.docker` / `.env.production`.

| Variable | Qué es |
| --- | --- |
| `POSTGRES_PASSWORD` | Contraseña de la base de datos |
| `FIREBASE_CREDENTIALS` | JSON de la cuenta de servicio de Firebase (contiene una clave privada) |
| `EVENTS_QUEUE_USER_API_KEY` | Clave del usuario del sistema que procesa la cola de eventos |
| `SOKETI_SECRET` | Firma los eventos en tiempo real (`SOKETI_KEY` es pública, pero mantenla junto) |
| `SMTP_PASSWORD` | Contraseña del servidor de correo |
| `CLOUDFLARE_TURNSTILE_SECRET_KEY` | Clave secreta de Turnstile (la *site key* sí es pública) |
| `AXIOM_TOKEN` | Token de observabilidad, si lo usas |

> **Sobre `FIREBASE_API_KEY` (web):** la clave web de Firebase no es secreta en sentido estricto (viaja dentro del JavaScript que descarga el navegador), pero **igual no la subas al repo**: así nadie puede usar tu proyecto de Firebase desde otro sitio. Restríngela en Google Cloud Console → *APIs y servicios → Credenciales* a tus dominios.

### Secretos de la app Android

| Qué | Cómo pasarlo a la compilación |
| --- | --- |
| `google-services.json` (tu proyecto de Firebase) | Archivo en `android/app/google-services.json`, ignorado por git |
| Keystore de firma (`*.jks`) y sus contraseñas | Archivo fuera del repo + `android/key.properties` (ignorado) |
| `HTTPSMS_SENTRY_DSN`, `HTTPSMS_LOGS_URL`, `HTTPSMS_LOGS_TOKEN` (opcionales) | Variables de entorno o `-P` de Gradle |
| `HTTPSMS_DEFAULT_SERVER_URL` | No es secreto, pero se pasa igual |

---

## Guardarlos como secretos

### Coolify

1. Abre tu recurso → **Environment Variables**.
2. Agrega cada variable con su valor real y marca **Is Secret** (o *Locked*) para que no se muestre en los logs ni en la interfaz.
3. Las que usa la web al compilar (`FIREBASE_*`, `SOKETI_KEY`, `API_BASE_URL`…) deben estar marcadas como **Build Variable** también, porque la web las integra al compilar.
4. Vuelve a desplegar.

### GitHub Actions

1. En GitHub: **Settings → Secrets and variables → Actions → New repository secret**.
2. Para archivos (como `google-services.json` o el keystore), guárdalos en base64:

   ```bash
   base64 -i android/app/google-services.json | pbcopy   # macOS
   base64 -w0 android/app/google-services.json            # Linux
   ```

   Pega el resultado como secreto (por ejemplo `GOOGLE_SERVICES_JSON_B64`).
3. En el workflow, reconstrúyelo antes de compilar:

   ```yaml
   - name: Restore google-services.json
     run: echo "${{ secrets.GOOGLE_SERVICES_JSON_B64 }}" | base64 -d > android/app/google-services.json
   - name: Build APK
     env:
       HTTPSMS_DEFAULT_SERVER_URL: ${{ vars.HTTPSMS_DEFAULT_SERVER_URL }}
       HTTPSMS_SENTRY_DSN: ${{ secrets.HTTPSMS_SENTRY_DSN }}
     run: cd android && ./gradlew assembleRelease
   ```

GitHub oculta automáticamente el valor de los secretos en los logs.

---

## Protecciones incluidas en el repo

1. **`.gitignore`** ignora `.env`, `google-services.json`, keystores (`*.jks`, `*.keystore`, `*.p12`), claves privadas (`*.pem`, `*.key`), `key.properties` y credenciales de Firebase/GCP (`*firebase-adminsdk*.json`, `*service-account*.json`, `*credentials*.json`).
2. **Escáner de secretos en el pre-commit** (`scripts/check-secrets.sh`). Antes de cada commit revisa **solo lo que vas a commitear** y lo bloquea si encuentra:
   - archivos sensibles (los de arriba, aunque los fuerces con `git add -f`);
   - claves con formato conocido: Google/Firebase (`AIza…`), bloques `PRIVATE KEY`, tokens de GitHub, Slack, Stripe, AWS, DSN de Sentry, tokens de Axiom, JWT;
   - en archivos `.env*`, cualquier variable `…_SECRET`, `…_PASSWORD`, `…_TOKEN`, `…_API_KEY`, `…_CREDENTIALS` o `…_DSN` con un valor que no sea un ejemplo (vacío, `changeme`, `your-…`, `<valor>`, `${VAR:-…}`).

   Puedes ejecutarlo a mano: `bash scripts/check-secrets.sh`. Si da un falso positivo, agrega el comentario `secret-scan: ignore` en esa línea.

> El hook se instala con `pnpm install` dentro de `web/` (script `prepare` de husky). Si clonas el repo en otra computadora, ejecuta `cd web && pnpm install` antes de hacer commits.

### Activa también la protección de GitHub

En tu repositorio: **Settings → Code security → Secret Protection** (o *Secret scanning*) → activa **Secret scanning** y **Push protection**. GitHub rechazará cualquier `git push` que contenga un secreto conocido, aunque alguien se salte el hook local.

---

## Si un secreto se filtra

1. **Primero, revócalo y genera uno nuevo** (en Firebase, Google Cloud, tu proveedor de correo, Soketi…). Borrarlo del repo **no** sirve: ya pudo ser copiado.
2. Actualiza el valor nuevo en `.env`, Coolify y GitHub Secrets, y vuelve a desplegar.
3. Después, si quieres, limpia el historial con [`git filter-repo`](https://github.com/newren/git-filter-repo) y fuerza el push. Avisa a quien tenga clones, porque el historial cambia.

> **Nota sobre el historial heredado:** el historial de este fork incluye commits del proyecto original con configuraciones de **sus** servicios (por ejemplo, su `google-services.json` y su configuración web de Firebase). No son tuyos ni dan acceso a tu instalación; ya se eliminaron de la versión actual. Ninguno de los valores reales de tu `.env` aparece en el historial.
