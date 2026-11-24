<script setup lang="ts">
import {
  mdiSend,
  mdiDotsVertical,
  mdiArrowLeft,
  mdiCheckAll,
  mdiDelete,
  mdiCallMissed,
  mdiPaperclip,
  mdiCheck,
  mdiAlert,
  mdiPackageUp,
  mdiPackageDown,
  mdiAccount,
  mdiRefresh,
  mdiContentCopy,
  mdiAccountPlus,
  mdiSquareEditOutline,
} from '@mdi/js'
import { bindRealtimeEvents } from '~/composables/useRealtime'
import { isValidPhoneNumber } from 'libphonenumber-js'
import { storeToRefs } from 'pinia'
import { getApiErrorMessage, toApiError } from '~/utils/api-error'
import type {
  EntitiesContact,
  EntitiesMessage,
  EntitiesMessageThread,
} from '~~/shared/types/api'
import {
  isUrlNotificationToken,
  formatWhatsappText,
} from '~/utils/whatsapp-format'

definePageMeta({
  middleware: ['auth'],
})

useHead({
  title: 'Messages - httpSMS',
})

const route = useRoute()
const router = useRouter()
const { mdAndDown } = useVDisplay()
const { formatPhoneNumber, startsWithLetter } = useFilters()
const notificationsStore = useNotificationsStore()
const authStore = useAuthStore()
const phonesStore = usePhonesStore()
const threadsStore = useThreadsStore()
const messagesStore = useMessagesStore()
const { t } = useI18n()
const { currentThread } = storeToRefs(threadsStore)

const formMessage = ref('')
const submitting = ref(false)
const loadingMessages = ref(false)
const hideMessages = ref(true)
const messages = ref<EntitiesMessage[]>([])
const selectedMenuItem = ref(-1)
const contactDialog = ref(false)
const messageBody = ref<HTMLElement | null>(null)
const form = ref<{ validate: () => Promise<{ valid: boolean }> } | null>(null)
let messageRequestId = 0

threadsStore.setThreadId(route.params.id as string)

const formMessageRules = [
  (v: string) =>
    v === '' || (v && v.length <= 2048) || t('threads.id.message_length_error'),
]

let unbindRealtime: () => void = () => {}

const contactIsPhoneNumber = computed(() => {
  const thread = currentThread.value
  if (!thread) return false
  return isValidPhoneNumber(thread.contact) || !isNaN(Number(thread.contact))
})

const messageVisibility = computed(() =>
  hideMessages.value ? 'hidden' : 'visible',
)
const contact = computed(() => currentThread.value?.contact ?? '')

function isMT(message: EntitiesMessage): boolean {
  return message.type === 'mobile-terminated'
}

function isMo(message: EntitiesMessage): boolean {
  return message.type === 'mobile-originated'
}

function isMissedCall(message: EntitiesMessage): boolean {
  return message.type === 'call/missed'
}

function isWhatsappFormatted(message: EntitiesMessage): boolean {
  const phone = phonesStore.phones.find((x) => x.phone_number === message.owner)
  return isUrlNotificationToken(phone?.fcm_token)
}

function isPending(message: EntitiesMessage): boolean {
  return ['sending', 'pending', 'scheduled'].includes(message.status)
}

function statusColor(message: EntitiesMessage): string {
  if (message.status === 'sending') return 'warning'
  if (message.status === 'scheduled') return 'teal'
  return 'primary'
}

function canResend(message: EntitiesMessage): boolean {
  return (
    isMT(message) &&
    (message.status === 'expired' || message.status === 'failed')
  )
}

function formatAttachmentName(url: string): string {
  const parts = url.split('/')
  if (parts.length >= 2) return '/' + parts.slice(-2).join('/')
  return url
}

function currentThreadContactTitle(): string {
  const thread = currentThread.value
  if (!thread) return ''
  return (
    thread.contact_details?.name?.trim() || formatPhoneNumber(thread.contact)
  )
}

async function refreshThreadContact(contacts: EntitiesContact[]) {
  await threadsStore.loadThreads(contacts)
}

function scrollToElement() {
  const el = messageBody.value
  if (el) {
    el.scrollTop = el.scrollHeight + 120
  }
  hideMessages.value = false
}

async function resetCurrentThreadUnreadCount(
  force = false,
  threadId = route.params.id as string,
) {
  try {
    await threadsStore.resetThreadUnreadCount(threadId, force)
  } catch (error) {
    console.error(error)
  }
}

async function handleInboundMessage() {
  if (loadingMessages.value) return

  try {
    const thread = currentThread.value
    await threadsStore.loadThreads()
    if (!thread || loadingMessages.value) {
      return
    }

    await resetCurrentThreadUnreadCount(true)
    loadMessages(false, false)
  } catch (error) {
    console.error(error)
  }
}

async function loadMessages(hide = true, resetUnreadCount = true) {
  const requestId = ++messageRequestId
  loadingMessages.value = true
  const threadId = route.params.id as string
  hideMessages.value = hide

  if (resetUnreadCount) {
    void resetCurrentThreadUnreadCount(false, threadId)
  }

  try {
    const response = await threadsStore.loadThreadMessages(threadId)
    if (requestId !== messageRequestId) return

    messages.value = [...response].reverse()
    await nextTick()
    if (requestId !== messageRequestId) return

    scrollToElement()
  } finally {
    if (requestId === messageRequestId) {
      loadingMessages.value = false
    }
  }
}

async function loadData() {
  const threadId = route.params.id as string

  if (!threadsStore.hasThreadId(threadId)) {
    if (!authStore.user) await authStore.loadUser()
    await phonesStore.loadPhones()
    await threadsStore.loadThreads()
  }

  if (route.params.id !== threadId) return
  if (!threadsStore.hasThreadId(threadId)) {
    let requestedThread: EntitiesMessageThread
    try {
      requestedThread = await threadsStore.fetchThread(threadId)
    } catch (error) {
      if (route.params.id !== threadId) return
      const status = toApiError(error).status
      if (status === 404 || status === 422) {
        notificationsStore.addNotification({
          type: 'error',
          message: 'The conversation could not be found',
        })
        await router.push('/threads')
        return
      }
      notificationsStore.addNotification({
        type: 'error',
        message: getApiErrorMessage(
          error,
          'The conversation could not be loaded',
        ),
      })
      throw error
    }
    if (route.params.id !== threadId) return
    const contextChanged =
      phonesStore.owner !== requestedThread.owner ||
      threadsStore.archivedThreads !== requestedThread.is_archived
    if (phonesStore.owner !== requestedThread.owner) {
      await authStore.updateUser({ owner: requestedThread.owner })
    }
    if (route.params.id !== threadId) return
    threadsStore.archivedThreads = requestedThread.is_archived
    if (contextChanged) await threadsStore.loadThreads()
    if (route.params.id !== threadId) return
    threadsStore.setThreadId(threadId)
  }

  if (!threadsStore.hasThreadId(threadId)) {
    await router.push('/threads')
    return
  }
  await loadMessages()
}

async function archiveThread() {
  await threadsStore.updateThread({
    threadId: currentThread.value!.id,
    isArchived: true,
  })
  await router.push('/threads')
}

async function unArchiveThread() {
  await threadsStore.updateThread({
    threadId: currentThread.value!.id,
    isArchived: false,
  })
  await router.push('/threads')
}

async function resendMessage(message: EntitiesMessage) {
  await messagesStore.sendMessage({
    from: message.owner,
    to: message.contact,
    content: message.content,
    sim: 'DEFAULT',
  })
  setTimeout(() => {
    selectedMenuItem.value = -1
  }, 1000)
  loadMessages(false)
}

async function deleteMessage(message: EntitiesMessage) {
  await messagesStore.deleteMessage(message.id)
  setTimeout(() => {
    selectedMenuItem.value = -1
  }, 1000)
  loadMessages(false)
}

async function copyMessageId(message: EntitiesMessage) {
  await navigator.clipboard.writeText(message.id)
  notificationsStore.addNotification({
    message: t('threads.id.message_id_copied'),
    type: 'success',
  })
  setTimeout(() => {
    selectedMenuItem.value = -1
  }, 1000)
}

async function deleteThread(threadID: string) {
  await threadsStore.deleteThread(threadID)
  await router.push('/threads')
}

async function sendMessage(event: KeyboardEvent | Event) {
  if (event instanceof KeyboardEvent && event.shiftKey) return
  if (!formMessage.value.trim()) return

  if (form.value) {
    const { valid } = await form.value.validate()
    if (!valid) return
  }

  submitting.value = true
  await messagesStore.sendMessage({
    from: phonesStore.owner!,
    to: currentThread.value!.contact,
    content: formMessage.value,
    sim: 'DEFAULT',
  })
  loadMessages(false)
  formMessage.value = ''
  submitting.value = false
}

onMounted(async () => {
  await loadData()

  unbindRealtime = bindRealtimeEvents(authStore.user?.id, {
    'message.phone.sent': () => {
      if (!loadingMessages.value) loadMessages(false)
    },
    'message.send.failed': () => {
      if (!loadingMessages.value) loadMessages(false)
    },
    'message.phone.received': () => void handleInboundMessage(),
    'message.call.missed': () => void handleInboundMessage(),
  })
})

watch(
  () => route.params.id as string,
  (threadId, previousThreadId) => {
    if (!previousThreadId || threadId === previousThreadId) return

    threadsStore.setThreadId(threadId)
    void loadData()
  },
)

onBeforeUnmount(() => unbindRealtime())
</script>

<template>
  <div class="thread-view">
    <header class="thread-view-header">
      <VBtn
        v-if="mdAndDown"
        icon
        variant="text"
        to="/threads"
        :aria-label="$t('threads.id.back')"
      >
        <VIcon :icon="mdiArrowLeft" />
      </VBtn>
      <template v-if="currentThread">
        <VAvatar :color="currentThread.color" size="40" class="flex-shrink-0">
          <VIcon
            v-if="!startsWithLetter(currentThreadContactTitle())"
            color="white"
            :icon="mdiAccount"
          />
          <span v-else class="text-white text-title-large">{{
            currentThreadContactTitle().substring(0, 1)
          }}</span>
        </VAvatar>
        <div class="thread-view-title">
          <div class="text-title-medium font-weight-bold text-truncate">
            {{ currentThreadContactTitle() }}
          </div>
          <div class="text-body-small text-medium-emphasis text-truncate">
            <template v-if="currentThread.contact_details?.name?.trim()">
              {{ formatPhoneNumber(currentThread.contact) }}
            </template>
            <template v-else>
              {{
                $t('threads.id.via_phone', {
                  phone: formatPhoneNumber(currentThread.owner),
                })
              }}
            </template>
          </div>
        </div>
      </template>
      <VSpacer />
      <VMenu>
        <template #activator="{ props }">
          <VBtn
            icon
            variant="text"
            v-bind="props"
            :aria-label="$t('threads.id.thread_actions')"
          >
            <VIcon :icon="mdiDotsVertical" />
          </VBtn>
        </template>
        <VList prepend-gap="16" density="compact" min-width="220">
          <VListItem
            v-if="currentThread?.contact_details"
            @click="contactDialog = true"
          >
            <template #prepend>
              <VIcon :icon="mdiSquareEditOutline" />
            </template>
            <VListItemTitle>{{ $t('threads.id.edit_contact') }}</VListItemTitle>
          </VListItem>
          <VListItem
            v-else-if="contactIsPhoneNumber"
            @click="contactDialog = true"
          >
            <template #prepend>
              <VIcon :icon="mdiAccountPlus" />
            </template>
            <VListItemTitle>{{ $t('threads.id.add_contact') }}</VListItemTitle>
          </VListItem>
          <VListItem
            v-if="currentThread && !currentThread.is_archived"
            @click.prevent="archiveThread"
          >
            <template #prepend>
              <VIcon :icon="mdiPackageDown" />
            </template>
            <VListItemTitle>{{ $t('threads.id.archive') }}</VListItemTitle>
          </VListItem>
          <VListItem
            v-if="currentThread && currentThread.is_archived"
            @click.prevent="unArchiveThread"
          >
            <template #prepend>
              <VIcon :icon="mdiPackageUp" />
            </template>
            <VListItemTitle>{{ $t('threads.id.unarchive') }}</VListItemTitle>
          </VListItem>
          <VListItem
            v-if="currentThread"
            base-color="error"
            @click.prevent="deleteThread(currentThread.id)"
          >
            <template #prepend>
              <VIcon :icon="mdiDelete" />
            </template>
            <VListItemTitle>{{
              $t('threads.id.delete_thread')
            }}</VListItemTitle>
          </VListItem>
        </VList>
      </VMenu>
    </header>
    <VProgressLinear
      :active="loadingMessages"
      color="primary"
      height="2"
      indeterminate
    />

    <div ref="messageBody" class="thread-view-body">
      <div v-if="currentThread" class="thread-view-messages">
        <div
          v-for="message in messages"
          :key="message.id"
          class="msg-row"
          :class="isMT(message) ? 'msg-row-out' : 'msg-row-in'"
          :style="{ visibility: messageVisibility }"
        >
          <VAvatar
            v-if="isMo(message)"
            :color="currentThread.color"
            size="32"
            class="msg-row-avatar"
          >
            <VIcon
              v-if="!startsWithLetter(message.contact)"
              size="18"
              color="white"
              :icon="mdiAccount"
            />
            <span v-else class="text-white text-body-medium font-weight-bold">{{
              message.contact.substring(0, 1)
            }}</span>
          </VAvatar>
          <VAvatar
            v-if="isMissedCall(message)"
            size="32"
            color="error"
            variant="tonal"
            class="msg-row-avatar"
          >
            <VIcon size="18" :icon="mdiCallMissed" />
          </VAvatar>

          <div class="msg-stack">
            <div
              v-if="message.content"
              class="msg-bubble text-break"
              :class="{
                'msg-bubble-out': isMT(message),
                'msg-bubble-failed':
                  isMT(message) &&
                  (message.status === 'failed' || message.status === 'expired'),
              }"
            >
              <span
                v-if="!isMissedCall(message) && isWhatsappFormatted(message)"
                v-html="formatWhatsappText(message.content)"
              />
              <span v-else-if="!isMissedCall(message)">{{
                message.content
              }}</span>
              <span v-else class="text-medium-emphasis">{{
                $t('threads.id.missed_phone_call')
              }}</span>
            </div>
            <div v-if="message.attachments?.length" class="msg-attachments">
              <a
                v-for="(attachment, index) in message.attachments"
                :key="index"
                target="_blank"
                rel="noopener noreferrer"
                :href="attachment"
                class="text-decoration-none hover:text-decoration-underline text-body-medium d-flex align-center"
              >
                <VIcon size="x-small" class="mr-1" :icon="mdiPaperclip" />
                {{ formatAttachmentName(attachment) }}
              </a>
            </div>
            <div class="msg-meta">
              <span>{{
                new Date(message.order_timestamp).toLocaleString()
              }}</span>
              <VTooltip location="bottom">
                <template #activator="{ props }">
                  <span v-bind="props" class="d-inline-flex">
                    <VIcon
                      v-if="message.status === 'expired'"
                      color="warning"
                      size="16"
                      :icon="mdiAlert"
                    />
                    <VProgressCircular
                      v-else-if="isPending(message)"
                      indeterminate
                      :size="12"
                      :width="1.5"
                      :color="statusColor(message)"
                    />
                    <VIcon
                      v-else-if="message.status === 'delivered'"
                      color="primary"
                      size="16"
                      :icon="mdiCheckAll"
                    />
                    <VIcon
                      v-else-if="message.status === 'sent'"
                      size="16"
                      :icon="mdiCheck"
                    />
                    <VIcon
                      v-else-if="message.status === 'failed'"
                      color="error"
                      size="16"
                      :icon="mdiAlert"
                    />
                  </span>
                </template>
                <span>{{ message.failure_reason || message.status }}</span>
              </VTooltip>
            </div>
          </div>

          <VMenu>
            <template #activator="{ props }">
              <VBtn
                icon
                variant="text"
                size="small"
                class="msg-row-actions"
                :aria-label="$t('threads.id.message_actions')"
                v-bind="props"
              >
                <VIcon :icon="mdiDotsVertical" />
              </VBtn>
            </template>
            <VList prepend-gap="16" density="compact">
              <VListItem
                v-if="canResend(message)"
                @click.prevent="resendMessage(message)"
              >
                <template #prepend
                  ><VIcon size="small" :icon="mdiRefresh"
                /></template>
                <VListItemTitle class="text-body-medium">{{
                  $t('threads.id.resend_message')
                }}</VListItemTitle>
              </VListItem>
              <VListItem @click.prevent="copyMessageId(message)">
                <template #prepend
                  ><VIcon size="small" :icon="mdiContentCopy"
                /></template>
                <VListItemTitle class="text-body-medium">{{
                  $t('threads.id.copy_message_id')
                }}</VListItemTitle>
              </VListItem>
              <VListItem
                base-color="error"
                @click.prevent="deleteMessage(message)"
              >
                <template #prepend
                  ><VIcon size="small" :icon="mdiDelete"
                /></template>
                <VListItemTitle class="text-body-medium">{{
                  $t('threads.id.delete_message')
                }}</VListItemTitle>
              </VListItem>
            </VList>
          </VMenu>
        </div>
      </div>
    </div>

    <footer v-if="currentThread" class="thread-view-composer">
      <VForm ref="form" class="thread-view-form" @submit.prevent>
        <VTextField
          v-model="formMessage"
          :disabled="submitting || !contactIsPhoneNumber"
          variant="solo-filled"
          flat
          rounded="pill"
          hide-details="auto"
          :rules="formMessageRules"
          :placeholder="
            contactIsPhoneNumber
              ? $t('threads.id.type_message_here')
              : $t('threads.id.cannot_send_messages', { contact })
          "
          :aria-label="$t('threads.id.type_message_here')"
          @keydown.enter="sendMessage"
        />
        <VBtn
          :disabled="submitting || !contactIsPhoneNumber"
          type="submit"
          class="btn-gradient ml-2 flex-shrink-0"
          icon
          size="large"
          rounded="circle"
          :loading="submitting"
          :aria-label="$t('threads.id.send')"
          @click="sendMessage"
        >
          <VIcon :icon="mdiSend" />
        </VBtn>
      </VForm>
    </footer>
    <ContactDialog
      v-model="contactDialog"
      :contact="currentThread?.contact_details"
      :initial-phone-number="currentThread?.contact ?? ''"
      :refresh-contacts="false"
      @saved="refreshThreadContact"
    />
  </div>
</template>

<style lang="scss">
.thread-view {
  display: flex;
  flex-direction: column;
  height: 100dvh;
  background:
    radial-gradient(
      circle at 85% 0%,
      rgb(var(--v-theme-primary), 0.07),
      transparent 40%
    ),
    rgb(var(--v-theme-background));
}

.thread-view-header {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 0.75rem;
  min-height: 68px;
  padding: 0.5rem 0.75rem 0.5rem 1rem;
  border-bottom: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  background: rgb(var(--v-theme-surface));
}

.thread-view-title {
  min-width: 0;
  line-height: 1.3;
}

.thread-view-body {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
}

.thread-view-messages {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  gap: 0.9rem;
  max-width: 920px;
  min-height: 100%;
  margin: 0 auto;
  padding: 1.5rem 1rem;
}

.msg-row {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
}

.msg-row-out {
  flex-direction: row-reverse;
}

.msg-row-avatar {
  margin-top: 0.15rem;
}

.msg-row-actions {
  opacity: 0.55;
  transition: opacity 0.15s;
}

.msg-row:hover .msg-row-actions,
.msg-row-actions[aria-expanded='true'] {
  opacity: 1;
}

.msg-stack {
  display: flex;
  flex-direction: column;
  max-width: min(72%, 560px);
}

.msg-row-out .msg-stack {
  align-items: flex-end;
}

.msg-bubble {
  padding: 0.65rem 0.95rem;
  border: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  border-radius: 1.1rem 1.1rem 1.1rem 0.35rem;
  background: rgb(var(--v-theme-surface));
  line-height: 1.5;
  white-space: pre-line;
}

.msg-bubble-out {
  border: 0;
  border-radius: 1.1rem 1.1rem 0.35rem;
  background: var(--brand-grad);
  color: rgb(var(--v-theme-on-primary));
}

.msg-bubble-out a {
  color: inherit;
  text-decoration: underline;
}

.msg-bubble-failed {
  outline: 2px solid rgb(var(--v-theme-error), 0.55);
  outline-offset: 2px;
}

.msg-attachments {
  margin-top: 0.35rem;
  padding: 0.5rem 0.75rem;
  border: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  border-radius: 0.8rem;
  background: rgb(var(--v-theme-surface));
}

.msg-meta {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  margin-top: 0.3rem;
  padding-inline: 0.25rem;
  color: rgb(var(--v-theme-on-surface), var(--v-medium-emphasis-opacity));
  font-size: 0.75rem;
}

.thread-view-composer {
  flex-shrink: 0;
  padding: 0.75rem 1rem calc(0.75rem + env(safe-area-inset-bottom));
  border-top: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  background: rgb(var(--v-theme-surface));
}

.thread-view-form {
  display: flex;
  align-items: flex-start;
  max-width: 920px;
  margin: 0 auto;
}

@media (width <= 600px) {
  .thread-view-messages {
    padding: 1rem 0.5rem;
  }

  .msg-stack {
    max-width: 82%;
  }
}
</style>
