package docs

import (
	"encoding/json"
	"sync"
)

// Este archivo NO es generado por `swag init`. Contiene una superposición de
// traducción al español que se aplica en tiempo de ejecución sobre la
// especificación generada (docs.go), por lo que sobrevive a la regeneración.
//
// Las llaves del mapa tienen el formato "MÉTODO /ruta" (igual que en el spec).
// Si una operación no tiene traducción, se conserva el texto original en inglés.

type translation struct {
	Summary     string
	Description string
}

var spanishInfoTitle = "Referencia de la API de httpSMS"

var spanishInfoDescription = "Usa tu teléfono Android para enviar y recibir mensajes SMS mediante una API programable sencilla, con cifrado de extremo a extremo. Esta es una instalación gratuita y autoalojada: no hay planes de pago, límites de mensajes ni cobros. Para guías y ejemplos, consulta la sección /docs de tu panel web."

var spanishTags = map[string]string{
	"3CXIntegration":     "Integración 3CX",
	"Attachments":        "Adjuntos",
	"Billing":            "Uso de mensajes",
	"BulkSMS":            "SMS masivos",
	"Contacts":           "Contactos",
	"Discord":            "Discord",
	"DiscordIntegration": "Integración de Discord",
	"Heartbeats":         "Latidos (heartbeats)",
	"MessageThreads":     "Conversaciones",
	"Messages":           "Mensajes",
	"PhoneAPIKeys":       "Claves API de teléfono",
	"Phones":             "Teléfonos",
	"SendSchedules":      "Horarios de envío",
	"Users":              "Usuarios",
	"Webhooks":           "Webhooks",
}

var spanishOperations = map[string]translation{
	"GET /attachments/{userID}/{messageID}/{attachmentIndex}/{filename}": {"Descargar un adjunto de mensaje", "Descarga un adjunto MMS a partir de los componentes de su ruta."},

	"GET /billing/usage":         {"Obtener el uso de mensajes del mes actual", "Obtiene el resumen de mensajes enviados y recibidos por el usuario en el mes actual. El uso se registra solo con fines estadísticos; no hay límites."},
	"GET /billing/usage-history": {"Obtener el historial de uso de mensajes", "Obtiene los registros mensuales de mensajes enviados y recibidos por el usuario en meses anteriores, ordenados por fecha de forma descendente. El uso se registra solo con fines estadísticos; no hay límites."},

	"GET /bulk-messages":  {"Listar envíos masivos", "Obtiene el resumen de los últimos 10 envíos masivos del usuario autenticado, con el número de mensajes en cada estado."},
	"POST /bulk-messages": {"Subir archivo de SMS masivos", "Envía mensajes SMS a varios destinatarios a partir de la plantilla CSV (`/templates/httpsms-bulk.csv`) o la plantilla de Excel (`/templates/httpsms-bulk.xlsx`) que ofrece tu panel web."},

	"GET /contacts":                {"Listar contactos", "Devuelve la lista paginada de contactos del usuario autenticado. El campo \"total\" de nivel superior indica cuántos contactos coinciden con el filtro de búsqueda, sin importar skip/limit, para que los clientes puedan paginar desde el servidor."},
	"POST /contacts":               {"Crear uno o varios contactos", "Crea un solo contacto o un lote de contactos. Acepta un arreglo JSON o un objeto con un arreglo \"contacts\"."},
	"POST /contacts/upload":        {"Importar contactos desde CSV", "Sube un archivo CSV con contactos en el campo multipart \"document\". Columnas: Name, Emails, PhoneNumbers (los valores múltiples se separan con \";\")."},
	"PUT /contacts/{contactID}":    {"Actualizar un contacto", "Actualiza los datos de un solo contacto."},
	"DELETE /contacts/{contactID}": {"Eliminar un contacto", "Elimina un solo contacto de la base de datos."},

	"GET /discord-integrations":                {"Listar integraciones de Discord", "Obtiene las integraciones de Discord del usuario autenticado."},
	"POST /discord-integrations":               {"Crear una integración de Discord", "Guarda una nueva integración de Discord para el usuario autenticado."},
	"PUT /discord-integrations/{discordID}":    {"Actualizar una integración de Discord", "Actualiza una integración de Discord del usuario autenticado."},
	"DELETE /discord-integrations/{discordID}": {"Eliminar una integración de Discord", "Elimina una integración de Discord del usuario autenticado."},
	"POST /discord/event":                      {"Recibir un evento de Discord", "Publica un evento de Discord para los oyentes registrados. Es el endpoint de interacciones de Discord: se sirve en `/discord/event` del host de la API, fuera de la ruta base `/v1`."},

	"GET /heartbeats":  {"Listar latidos de un número de teléfono", "Obtiene las últimas veces que un número de teléfono consultó mensajes pendientes, ordenadas por fecha de forma descendente."},
	"POST /heartbeats": {"Registrar un latido de un número de teléfono", "Guarda un latido para indicar que el número de teléfono sigue activo."},

	"POST /integration/3cx/messages": {"Enviar un SMS desde 3CX", "Envía un mensaje SMS desde la plataforma 3CX. Este endpoint se sirve en `/integration/3cx/messages` del host de la API, fuera de la ruta base `/v1`, y se autentica con tu clave API como token bearer."},

	"GET /message-threads":                      {"Listar conversaciones de un número de teléfono", "Obtiene la lista de contactos con los que se ha comunicado un número de teléfono (conversaciones), ordenada por fecha de forma descendente."},
	"GET /message-threads/{messageThreadID}":    {"Obtener una conversación", "Devuelve una conversación del usuario autenticado, incluidos los datos del contacto correspondiente cuando existen."},
	"PUT /message-threads/{messageThreadID}":    {"Actualizar una conversación", "Actualiza los datos de una conversación."},
	"DELETE /message-threads/{messageThreadID}": {"Eliminar una conversación", "Elimina una conversación de la base de datos junto con todos sus mensajes."},

	"GET /messages":                     {"Listar mensajes entre 2 números de teléfono", "Obtiene la lista de mensajes intercambiados entre 2 números de teléfono, ordenada por fecha de forma descendente."},
	"POST /messages/bulk-send":          {"Enviar varios SMS a la vez", "Agrega varios mensajes SMS a la cola para que los envíe el teléfono Android."},
	"POST /messages/calls/missed":       {"Registrar una llamada perdida", "La app Android de httpSMS llama a este endpoint para registrar una llamada perdida en el teléfono."},
	"GET /messages/outstanding":         {"Obtener un mensaje pendiente", "Obtiene un mensaje pendiente que debe enviar un teléfono Android."},
	"POST /messages/receive":            {"Recibir un SMS desde un teléfono", "Registra un nuevo mensaje recibido en un teléfono móvil."},
	"GET /messages/search":              {"Buscar mensajes del usuario", "Devuelve todos los mensajes que cumplen los criterios de búsqueda, incluidas las llamadas perdidas."},
	"POST /messages/send":               {"Enviar un SMS", "Agrega un nuevo mensaje SMS a la cola para que lo envíe tu teléfono Android."},
	"GET /messages/{messageID}":         {"Obtener un mensaje", "Obtiene un mensaje de la base de datos a partir de su ID."},
	"DELETE /messages/{messageID}":      {"Eliminar un mensaje", "Elimina un mensaje de la base de datos y quita su contenido de la lista de conversaciones."},
	"POST /messages/{messageID}/events": {"Registrar un evento de un mensaje", "Usa este endpoint para informar los eventos de un mensaje cuando el teléfono lo marca como fallido, enviado o entregado."},

	"GET /phone-api-keys":                                     {"Listar claves API de teléfono", "Obtiene la lista de claves API de teléfono que el usuario ha creado en httpSMS."},
	"POST /phone-api-keys":                                    {"Crear una clave API de teléfono", "Crea una nueva clave API de teléfono para iniciar sesión en la app httpSMS de tu teléfono Android."},
	"DELETE /phone-api-keys/{phoneAPIKeyID}":                  {"Eliminar una clave API de teléfono", "Elimina una clave API de teléfono de la base de datos; ya no podrá usarse para autenticarse."},
	"DELETE /phone-api-keys/{phoneAPIKeyID}/phones/{phoneID}": {"Desvincular un teléfono de una clave API de teléfono", "Quita la asociación entre el teléfono y la clave API. Tendrás que volver a iniciar sesión en la app httpSMS de tu teléfono Android con una nueva clave API de teléfono."},

	"GET /phones":              {"Listar teléfonos", "Obtiene la lista de teléfonos que el usuario ha registrado en httpSMS."},
	"PUT /phones":              {"Crear o actualizar un teléfono", "Actualiza las propiedades del teléfono de un usuario. Si no existe un teléfono con ese número, se crea uno nuevo (funciona como un \"upsert\"). Los campos que se omiten en el cuerpo conservan su valor actual. Las pasarelas de teléfono basadas en URL reciben activaciones HTTP compatibles con FCM."},
	"PUT /phones/fcm-token":    {"Crear o actualizar el token FCM de un teléfono", "Actualiza el token FCM o la URL de callback del adaptador de un teléfono. Si no existe un teléfono con ese número, se crea uno nuevo (funciona como un \"upsert\"). Las pasarelas de teléfono basadas en URL reciben activaciones HTTP compatibles con FCM."},
	"DELETE /phones/{phoneID}": {"Eliminar un teléfono", "Elimina un teléfono guardado en la base de datos."},

	"GET /send-schedules":                 {"Listar horarios de envío", "Lista todos los horarios de envío del usuario autenticado."},
	"POST /send-schedules":                {"Crear un horario de envío", "Crea un nuevo horario de envío para el usuario autenticado."},
	"PUT /send-schedules/{scheduleID}":    {"Actualizar un horario de envío", "Actualiza un horario de envío del usuario autenticado."},
	"DELETE /send-schedules/{scheduleID}": {"Eliminar un horario de envío", "Elimina un horario de envío del usuario autenticado."},

	"GET /users/me":                     {"Obtener el usuario actual", "Obtiene los datos del usuario autenticado."},
	"PUT /users/me":                     {"Actualizar el usuario actual", "Actualiza los datos del usuario autenticado."},
	"DELETE /users/me":                  {"Eliminar el usuario actual", "Elimina al usuario autenticado junto con todos sus datos."},
	"DELETE /users/{userID}/api-keys":   {"Rotar la clave API del usuario", "Genera una nueva clave API para el usuario e invalida la actual, por ejemplo si se ha visto comprometida."},
	"PUT /users/{userID}/notifications": {"Actualizar la configuración de notificaciones", "Actualiza la configuración de notificaciones por correo electrónico del usuario."},

	"GET /webhooks":                {"Listar webhooks", "Obtiene los webhooks del usuario autenticado."},
	"POST /webhooks":               {"Crear un webhook", "Guarda un nuevo webhook para el usuario autenticado."},
	"PUT /webhooks/{webhookID}":    {"Actualizar un webhook", "Actualiza un webhook del usuario autenticado."},
	"DELETE /webhooks/{webhookID}": {"Eliminar un webhook", "Elimina un webhook del usuario autenticado."},
}

var (
	spanishOnce sync.Once
	spanishDoc  string
)

// SpanishDoc devuelve la especificación OpenAPI/Swagger con textos en español.
// Se construye una sola vez a partir de la especificación generada y se guarda
// en caché. Si algo falla, devuelve la especificación original en inglés.
func SpanishDoc() string {
	spanishOnce.Do(func() {
		original := SwaggerInfo.ReadDoc()

		var spec map[string]any
		if err := json.Unmarshal([]byte(original), &spec); err != nil {
			spanishDoc = original
			return
		}

		if info, ok := spec["info"].(map[string]any); ok {
			info["title"] = spanishInfoTitle
			info["description"] = spanishInfoDescription
		}

		if paths, ok := spec["paths"].(map[string]any); ok {
			for path, methods := range paths {
				methodMap, ok := methods.(map[string]any)
				if !ok {
					continue
				}
				for method, op := range methodMap {
					operation, ok := op.(map[string]any)
					if !ok {
						continue
					}
					if tr, found := spanishOperations[upper(method)+" "+path]; found {
						operation["summary"] = tr.Summary
						operation["description"] = tr.Description
					}
					if tags, ok := operation["tags"].([]any); ok {
						for i, tag := range tags {
							if name, ok := tag.(string); ok {
								if translated, found := spanishTags[name]; found {
									tags[i] = translated
								}
							}
						}
					}
				}
			}
		}

		translated, err := json.Marshal(spec)
		if err != nil {
			spanishDoc = original
			return
		}
		spanishDoc = string(translated)
	})
	return spanishDoc
}

func upper(method string) string {
	out := []byte(method)
	for i, c := range out {
		if c >= 'a' && c <= 'z' {
			out[i] = c - 32
		}
	}
	return string(out)
}
