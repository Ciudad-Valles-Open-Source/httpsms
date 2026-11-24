// Toast messages come from two places that only speak English: the API
// (`response.message`, error `data.message`) and a few fixed strings in the
// stores. The toast looks them up here and shows the `notices.*` translation;
// anything unknown is shown as-is.

export interface NoticeTranslation {
  key: string
  params?: Record<string, string>
}

// Exact messages, compared lowercased and without a trailing period.
const EXACT: Record<string, string> = {
  // API success messages
  'api key rotated successfully': 'apiKeyRotated',
  'contact deleted successfully': 'contactDeleted',
  'contact updated successfully': 'contactUpdated',
  'discord integration created successfully': 'discordCreated',
  'discord integration deleted successfully': 'discordDeleted',
  'discord integration updated successfully': 'discordUpdated',
  'message added to queue': 'messageQueued',
  'message deleted successfully': 'messageDeleted',
  'message thread updated successfully': 'threadUpdated',
  'phone api key created successfully': 'phoneApiKeyCreated',
  'phone api key deleted successfully': 'phoneApiKeyDeleted',
  'phone deleted successfully': 'phoneDeleted',
  'phone has been dissociated from phone api key successfully':
    'phoneRemovedFromKey',
  'phone updated successfully': 'phoneUpdated',
  'send schedule created successfully': 'scheduleCreated',
  'send schedule deleted successfully': 'scheduleDeleted',
  'send schedule updated successfully': 'scheduleUpdated',
  'thread thread deleted successfully': 'threadDeleted',
  'user deleted successfully': 'accountDeleted',
  'user notification settings updated successfully': 'notificationsUpdated',
  'user updated successfully': 'userUpdated',
  'webhook created successfully': 'webhookCreated',
  'webhook deleted successfully': 'webhookDeleted',
  'webhook updated successfully': 'webhookUpdated',
  'attachment not found': 'attachmentNotFound',
  // Fixed strings from the stores
  'the message has been deleted successfully': 'messageDeleted',
  'bulk messages sent successfully': 'bulkSent',
  archived: 'threadArchived',
  unarchived: 'threadUnarchived',
  'the message thread could not be marked as read': 'threadMarkReadError',
  'the message thread has been deleted successfully': 'threadDeleted',
  'the conversation could not be found': 'threadNotFound',
  'the conversation could not be loaded': 'threadLoadError',
  'contacts created successfully': 'contactsCreated',
  'contact created successfully': 'contactCreated',
  'contact updated': 'contactUpdated',
  'contacts imported successfully': 'contactsImported',
  'your account has been deleted successfully': 'accountDeleted',
  'error while loading contacts': 'contactsLoadError',
  'error while saving contacts': 'contactsSaveError',
  'error while updating contact': 'contactUpdateError',
  'error while deleting contact': 'contactDeleteError',
  'error while importing contacts': 'contactsImportError',
  'error while sending message': 'messageSendError',
  'error while updating phone': 'phoneUpdateError',
}

// Messages with values in them.
const PATTERNS: Array<{
  pattern: RegExp
  key: string
  params?: (match: RegExpMatchArray) => Record<string, string>
}> = [
  {
    pattern: /^added (\d+) out of (\d+) messages to the queue$/i,
    key: 'bulkQueued',
    params: (m) => ({ added: m[1]!, total: m[2]! }),
  },
  { pattern: /^validation errors while /i, key: 'validationErrors' },
  { pattern: /^cannot find .+ with id \[/i, key: 'notFound' },
]

/** The `notices.*` translation for a known English message, if any. */
export function findNoticeTranslation(
  message: string,
): NoticeTranslation | null {
  const normalized = message.trim().replace(/\.$/, '').toLowerCase()
  const exact = EXACT[normalized]
  if (exact) return { key: `notices.${exact}` }
  for (const rule of PATTERNS) {
    const match = message.trim().match(rule.pattern)
    if (match) {
      return { key: `notices.${rule.key}`, params: rule.params?.(match) }
    }
  }
  return null
}

/** Translates a toast message, or returns it unchanged when it isn't known. */
export function translateNotice(
  message: string,
  t: (key: string, params?: Record<string, string>) => string,
): string {
  const translation = findNoticeTranslation(message)
  return translation ? t(translation.key, translation.params ?? {}) : message
}
