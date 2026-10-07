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

var spanishInfoDescription = "Usa tu teléfono Android para enviar y recibir mensajes SMS mediante una API programable simple, con cifrado de extremo a extremo."

var spanishTags = map[string]string{
	"Contacts":           "Contactos",
	"DiscordIntegration": "Integración de Discord",
	"Discord":            "Discord",
	"Heartbeats":         "Latidos (Heartbeats)",
	"3CXIntegration":     "Integración 3CX",
	"MessageThreads":     "Hilos de mensajes",
	"Messages":           "Mensajes",
	"PhoneAPIKeys":       "Claves API de teléfono",
	"Phones":             "Teléfonos",
	"SendSchedules":      "Horarios de envío",
	"Users":              "Usuarios",
	"Attachments":        "Adjuntos",
	"Webhooks":           "Webhooks",
	"BulkMessages":       "Mensajes masivos",
}

var spanishOperations = map[string]translation{
	"GET /bulk-messages":  {"Listar órdenes de mensajes masivos", "Obtiene los resúmenes de las últimas 10 órdenes de mensajes masivos del usuario autenticado, con el conteo por estado."},
	"POST /bulk-messages": {"Guardar archivo de SMS masivos", "Envía mensajes SMS masivos a varios destinatarios con base en nuestra [plantilla CSV](https://httpsms.com/templates/httpsms-bulk.csv) o nuestra [plantilla Excel](https://httpsms.com/templates/httpsms-bulk.xlsx)."},

	"GET /contacts":                {"Listar contactos", "Devuelve la lista paginada de contactos del usuario autenticado. El campo \"total\" de nivel superior es el número de contactos que coinciden con el filtro de búsqueda, independientemente de skip/limit, para que los clientes puedan paginar desde el servidor."},
	"POST /contacts":               {"Crear uno o varios contactos", "Crea un solo contacto o un lote de contactos. Acepta un arreglo JSON o un objeto con un arreglo \"contacts\"."},
	"POST /contacts/upload":        {"Importar contactos desde CSV", "Sube un archivo CSV (campo multipart \"document\") con contactos. Columnas: Name, Emails, PhoneNumbers (los valores múltiples se separan con \";\")."},
	"PUT /contacts/{contactID}":    {"Actualizar un contacto", "Actualiza los detalles de un solo contacto."},
	"DELETE /contacts/{contactID}": {"Eliminar un contacto", "Elimina un solo contacto de la base de datos."},

	"GET /discord-integrations":                {"Obtener las integraciones de Discord de un usuario", "Obtiene las integraciones de Discord de un usuario."},
	"POST /discord-integrations":               {"Guardar integración de Discord", "Guarda una integración de Discord para el usuario autenticado."},
	"PUT /discord-integrations/{discordID}":    {"Actualizar una integración de Discord", "Actualiza una integración de Discord para el usuario actualmente autenticado."},
	"DELETE /discord-integrations/{discordID}": {"Eliminar integración de Discord", "Elimina una integración de Discord de un usuario."},
	"POST /discord/event":                      {"Consumir un evento de Discord", "Publica un evento de Discord a los oyentes registrados."},

	"GET /heartbeats":  {"Obtener los latidos de un número de teléfono propietario", "Obtiene la última vez que un número de teléfono solicitó mensajes pendientes. Se ordena por marca de tiempo en orden descendente."},
	"POST /heartbeats": {"Registrar el latido de un número de teléfono propietario", "Guarda el latido para notificar que un número de teléfono sigue activo."},

	"POST /integration/3cx/messages": {"Enviar un mensaje SMS de 3CX", "Envía un mensaje SMS desde la plataforma 3CX."},

	"GET /message-threads":                      {"Obtener los hilos de mensajes de un número de teléfono", "Obtiene la lista de contactos con los que se ha comunicado un número de teléfono (hilos). Se ordena por marca de tiempo en orden descendente."},
	"PUT /message-threads/{messageThreadID}":    {"Actualizar un hilo de mensajes", "Actualiza los detalles de un hilo de mensajes."},
	"DELETE /message-threads/{messageThreadID}": {"Eliminar un hilo de mensajes de la base de datos.", "Elimina un hilo de mensajes de la base de datos y también todos los mensajes del hilo."},

	"GET /messages":                     {"Obtener los mensajes enviados entre 2 números de teléfono", "Obtiene la lista de mensajes enviados entre 2 números de teléfono. Se ordena por marca de tiempo en orden descendente."},
	"POST /messages/bulk-send":          {"Enviar mensajes SMS masivos", "Agrega mensajes SMS masivos para que sean enviados por el teléfono Android."},
	"POST /messages/calls/missed":       {"Registrar una llamada perdida en el teléfono móvil", "Este endpoint es llamado por la app Android de httpSMS para registrar una llamada perdida en el teléfono móvil."},
	"GET /messages/outstanding":         {"Obtener un mensaje pendiente", "Obtiene un mensaje pendiente que debe ser enviado por un teléfono Android."},
	"POST /messages/receive":            {"Recibir un nuevo mensaje SMS desde un teléfono móvil", "Agrega un nuevo mensaje recibido desde un teléfono móvil."},
	"GET /messages/search":              {"Buscar en todos los mensajes de un usuario", "Devuelve la lista de todos los mensajes según los criterios de filtro, incluyendo llamadas perdidas."},
	"POST /messages/send":               {"Enviar un mensaje SMS", "Agrega un nuevo mensaje SMS para que sea enviado por tu teléfono Android."},
	"GET /messages/{messageID}":         {"Obtener un mensaje de la base de datos.", "Obtiene un mensaje de la base de datos por su ID."},
	"DELETE /messages/{messageID}":      {"Eliminar un mensaje de la base de datos.", "Elimina un mensaje de la base de datos y quita su contenido de la lista de hilos."},
	"POST /messages/{messageID}/events": {"Registrar (upsert) un evento de un mensaje en el teléfono móvil", "Usa este endpoint para enviar eventos de un mensaje cuando el teléfono móvil lo marca como fallido, enviado o entregado."},

	"GET /phone-api-keys":                                     {"Obtener las claves API de teléfono de un usuario", "Obtiene la lista de claves API de teléfono que un usuario ha registrado en la aplicación httpSMS."},
	"POST /phone-api-keys":                                    {"Guardar clave API de teléfono", "Crea una nueva clave API de teléfono que puede usarse para iniciar sesión en la app httpSMS de tu teléfono Android."},
	"DELETE /phone-api-keys/{phoneAPIKeyID}":                  {"Eliminar una clave API de teléfono de la base de datos.", "Elimina una clave API de teléfono de la base de datos; ya no podrá usarse para autenticación."},
	"DELETE /phone-api-keys/{phoneAPIKeyID}/phones/{phoneID}": {"Quitar la asociación de un teléfono con la clave API de teléfono.", "Deberás iniciar sesión de nuevo en la app httpSMS de tu teléfono Android con una nueva clave API de teléfono."},

	"GET /phones":              {"Obtener los teléfonos de un usuario", "Obtiene la lista de teléfonos que un usuario ha registrado en la aplicación httpSMS."},
	"PUT /phones":              {"Crear o actualizar teléfono (upsert)", "Actualiza las propiedades del teléfono de un usuario. Si no existe un teléfono con este número, se creará uno nuevo. Piensa en este método como un 'upsert'. Las pasarelas de teléfono basadas en URL reciben activaciones HTTP compatibles con FCM."},
	"PUT /phones/fcm-token":    {"Crear o actualizar el token FCM de un teléfono", "Actualiza el token FCM o la URL de callback del adaptador de un teléfono. Si no existe un teléfono con este número, se creará uno nuevo. Piensa en este método como un 'upsert'. Las pasarelas de teléfono basadas en URL reciben activaciones HTTP compatibles con FCM."},
	"DELETE /phones/{phoneID}": {"Eliminar teléfono", "Elimina un teléfono almacenado en la base de datos."},

	"GET /send-schedules":                 {"Listar horarios de envío", "Lista todos los horarios de envío del usuario autenticado."},
	"POST /send-schedules":                {"Crear horario de envío", "Crea un nuevo horario de envío para el usuario autenticado."},
	"PUT /send-schedules/{scheduleID}":    {"Actualizar horario de envío", "Actualiza un horario de envío del usuario autenticado."},
	"DELETE /send-schedules/{scheduleID}": {"Eliminar horario de envío", "Elimina un horario de envío del usuario autenticado."},

	"GET /users/me":                      {"Obtener el usuario actual", "Obtiene los detalles del usuario actualmente autenticado."},
	"PUT /users/me":                      {"Actualizar un usuario", "Actualiza los detalles del usuario actualmente autenticado."},
	"DELETE /users/me":                   {"Eliminar un usuario", "Elimina al usuario actualmente autenticado junto con todos sus datos."},
	"DELETE /users/subscription":         {"Cancelar la suscripción del usuario", "Cancela la suscripción del usuario autenticado."},
	"GET /users/subscription-update-url": {"URL de actualización de suscripción del usuario autenticado", "Obtiene la URL de suscripción del usuario autenticado."},
	"POST /users/subscription/invoices/{subscriptionInvoiceID}": {"Generar una factura de pago de suscripción", "Genera un nuevo archivo PDF de factura para el pago de suscripción indicado, con los parámetros dados."},
	"GET /users/subscription/payments":                          {"Obtener los últimos 10 pagos de suscripción.", "Los pagos de suscripción se generan durante todo el ciclo de vida de una suscripción; normalmente hay uno al momento de la compra y uno por cada renovación."},
	"DELETE /users/{userID}/api-keys":                           {"Rotar la clave API del usuario", "Rota la clave API del usuario en caso de que la clave actual esté comprometida."},
	"PUT /users/{userID}/notifications":                         {"Actualizar la configuración de notificaciones", "Actualiza la configuración de notificaciones por correo electrónico de un usuario."},

	"GET /v1/attachments/{userID}/{messageID}/{attachmentIndex}/{filename}": {"Descargar un adjunto de mensaje", "Descarga un adjunto MMS usando los componentes de su ruta."},

	"GET /webhooks":                {"Obtener los webhooks de un usuario", "Obtiene los webhooks de un usuario."},
	"POST /webhooks":               {"Guardar un webhook", "Guarda un webhook para el usuario autenticado."},
	"PUT /webhooks/{webhookID}":    {"Actualizar un webhook", "Actualiza un webhook para el usuario actualmente autenticado."},
	"DELETE /webhooks/{webhookID}": {"Eliminar webhook", "Elimina un webhook de un usuario."},
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
