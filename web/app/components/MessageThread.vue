<script setup lang="ts">
import {
  mdiPlus,
  mdiDownload,
  mdiCheckAll,
  mdiCheck,
  mdiAlert,
  mdiAccount,
  mdiCellphoneArrowDown,
  mdiMessageTextOutline,
  mdiEmailOpenOutline,
  mdiPackageVariantClosed,
} from '@mdi/js'
import type { EntitiesMessageThread } from '~~/shared/types/api'
import {
  isUrlNotificationToken,
  formatWhatsappText,
} from '~/utils/whatsapp-format'

const { t } = useI18n()

const threadsStore = useThreadsStore()
const phonesStore = usePhonesStore()
const appStore = useAppStore()
const { githubIssuesUrl } = usePublicUrls()
const notificationsStore = useNotificationsStore()
const { formatPhoneNumber, startsWithLetter } = useFilters()

function threadDate(date: string): string {
  return new Date(date).toLocaleString(undefined, {
    month: 'short',
    day: 'numeric',
  })
}

function hasUnreadMessages(unreadCount: number): boolean {
  return unreadCount > 0
}

function unreadBadge(
  unreadCount: number,
): false | { color: string; content: string; dot: false } {
  if (unreadCount === 0) return false

  return {
    color: 'primary',
    content: unreadCount > 99 ? '99+' : String(unreadCount),
    dot: false,
  }
}

function onInstallApp() {
  notificationsStore.addNotification({
    type: 'info',
    message: t('components.messageThread.downloadingApp'),
  })
}

async function loadMoreThreads() {
  try {
    await threadsStore.loadMoreThreads()
  } catch (error) {
    notificationsStore.addNotification({
      type: 'error',
      message: t('components.messageThread.loadMoreError'),
    })
    console.error(error)
  }
}

function isWhatsappFormatted(thread: EntitiesMessageThread): boolean {
  const phone = phonesStore.phones.find((x) => x.phone_number === thread.owner)
  return isUrlNotificationToken(phone?.fcm_token)
}

function threadPreviewHtml(thread: EntitiesMessageThread): string {
  // Links are unwrapped because the preview sits inside a router link.
  return formatWhatsappText(thread.last_message_content ?? '').replace(
    /<a [^>]*>(.*?)<\/a>/g,
    '$1',
  )
}

function threadContactName(thread: EntitiesMessageThread): string {
  return thread.contact_details?.name?.trim() ?? ''
}

function threadContactTitle(thread: EntitiesMessageThread): string {
  return threadContactName(thread) || formatPhoneNumber(thread.contact)
}

function threadAvatarInitial(thread: EntitiesMessageThread): string {
  const contactName = threadContactName(thread)
  if (contactName) return contactName.substring(0, 1)

  return startsWithLetter(thread.contact) ? thread.contact.substring(0, 1) : ''
}
</script>

<template>
  <div class="thread-list">
    <div v-if="threadsStore.loadingThreads" class="pa-2" aria-busy="true">
      <v-skeleton-loader
        v-for="n in 6"
        :key="n"
        type="list-item-avatar-two-line"
        color="transparent"
      />
    </div>
    <div
      v-if="!threadsStore.loadingThreads && threadsStore.archivedThreads"
      class="archived-banner"
    >
      <v-icon size="small" :icon="mdiPackageVariantClosed" class="mr-2" />
      {{ $t('components.messageThread.archivedMessages') }}
    </div>

    <!-- No phone connected yet -->
    <div
      v-if="phonesStore.phones.length === 0 && !threadsStore.loadingThreads"
      class="empty-state"
    >
      <span class="icon-tile icon-tile-lg">
        <v-icon size="32" :icon="mdiCellphoneArrowDown" />
      </span>
      <p class="empty-state-title">
        {{ $t('components.messageThread.noPhoneTitle') }}
      </p>
      <p class="empty-state-text">
        {{ $t('components.messageThread.installAppPrompt') }}
        <a
          :href="githubIssuesUrl"
          target="_blank"
          class="text-decoration-none hover:text-decoration-underline"
          >{{ $t('components.messageThread.messageDiscord') }}</a
        >
        {{ $t('components.messageThread.setupHelp') }}
      </p>
      <v-btn
        class="btn-gradient"
        :href="appStore.appData.appDownloadUrl"
        @click="onInstallApp"
      >
        <v-icon :icon="mdiDownload" start />
        {{ $t('components.messageThread.downloadApp') }}
      </v-btn>
    </div>

    <!-- Phone connected but no conversations -->
    <div
      v-else-if="
        !threadsStore.loadingThreads &&
        threadsStore.visibleThreads.length === 0 &&
        !threadsStore.unreadOnly &&
        !threadsStore.archivedThreads
      "
      class="empty-state"
    >
      <span class="icon-tile icon-tile-lg">
        <v-icon size="32" :icon="mdiMessageTextOutline" />
      </span>
      <p class="empty-state-title">
        {{ $t('components.messageThread.emptyTitle') }}
      </p>
      <p v-if="phonesStore.owner" class="empty-state-text">
        {{ $t('components.messageThread.startSendingMessages') }}
      </p>
      <v-btn
        v-if="phonesStore.owner"
        class="btn-gradient"
        :to="{ name: 'messages' }"
      >
        <v-icon :icon="mdiPlus" start />
        {{ $t('components.messageThread.newMessage') }}
      </v-btn>
    </div>

    <!-- Archive is empty -->
    <div
      v-else-if="
        !threadsStore.loadingThreads &&
        threadsStore.visibleThreads.length === 0 &&
        threadsStore.archivedThreads &&
        !threadsStore.hasMoreThreads
      "
      class="empty-state"
    >
      <span class="icon-tile icon-tile-lg">
        <v-icon size="32" :icon="mdiPackageVariantClosed" />
      </span>
      <p class="empty-state-title">
        {{ $t('components.messageThread.emptyArchiveTitle') }}
      </p>
    </div>

    <!-- Unread filter is empty -->
    <div
      v-if="
        !threadsStore.loadingThreads &&
        threadsStore.unreadOnly &&
        !threadsStore.archivedThreads &&
        threadsStore.visibleThreads.length === 0 &&
        !threadsStore.hasMoreThreads &&
        phonesStore.owner
      "
      class="empty-state"
    >
      <span class="icon-tile icon-tile-lg">
        <v-icon size="32" :icon="mdiEmailOpenOutline" />
      </span>
      <p class="empty-state-title">
        {{ $t('components.messageThread.noUnread') }}
      </p>
    </div>

    <v-list
      v-if="
        !threadsStore.loadingThreads && threadsStore.visibleThreads.length > 0
      "
      class="py-2 px-2"
      lines="two"
      bg-color="transparent"
    >
      <v-list-item
        v-for="thread in threadsStore.visibleThreads"
        :key="thread.id"
        color="primary"
        rounded="lg"
        class="thread-item mb-1"
        :to="{ name: 'threads-id', params: { id: thread.id } }"
        :active="threadsStore.threadId === thread.id"
      >
        <template #prepend>
          <v-avatar
            :color="thread.color"
            size="44"
            :badge="unreadBadge(thread.unread_count)"
          >
            <v-icon v-if="!threadAvatarInitial(thread)" color="white">{{
              mdiAccount
            }}</v-icon>
            <span v-else class="text-white text-title-large">{{
              threadAvatarInitial(thread)
            }}</span>
          </v-avatar>
        </template>
        <v-list-item-title
          class="thread-item-title"
          :class="{
            'font-weight-bold': hasUnreadMessages(thread.unread_count),
          }"
        >
          {{ threadContactTitle(thread) }}
        </v-list-item-title>
        <v-list-item-subtitle
          class="text-truncate mt-1 thread-item-preview"
          :class="{
            'font-weight-bold opacity-100 text-high-emphasis':
              hasUnreadMessages(thread.unread_count),
          }"
        >
          <!-- eslint-disable vue/no-v-html -->
          <span
            v-if="isWhatsappFormatted(thread)"
            v-html="threadPreviewHtml(thread)"
          />
          <!-- eslint-enable vue/no-v-html -->
          <template v-else>{{ thread.last_message_content }}</template>
        </v-list-item-subtitle>
        <template #append>
          <div class="d-flex flex-column align-end align-self-start pt-1">
            <span
              class="text-body-small"
              :class="
                hasUnreadMessages(thread.unread_count)
                  ? 'text-primary font-weight-bold'
                  : 'text-medium-emphasis'
              "
            >
              {{ threadDate(thread.order_timestamp) }}
            </span>
            <VTooltip location="bottom">
              <template #activator="{ props }">
                <div v-bind="props" class="mt-1">
                  <v-icon
                    v-if="thread.status === 'expired'"
                    color="warning"
                    size="x-small"
                    :icon="mdiAlert"
                  />
                  <v-icon
                    v-else-if="thread.status === 'delivered'"
                    color="primary"
                    size="x-small"
                    :icon="mdiCheckAll"
                  />
                  <v-icon
                    v-else-if="thread.status === 'received'"
                    color="success"
                    size="x-small"
                    :icon="mdiCheckAll"
                  />
                  <v-icon
                    v-else-if="thread.status === 'sent'"
                    size="x-small"
                    :icon="mdiCheck"
                  />
                  <v-icon
                    v-else-if="thread.status === 'failed'"
                    color="error"
                    size="x-small"
                    :icon="mdiAlert"
                  />
                </div>
              </template>
              <span>{{
                $t(`components.messageThread.status.${thread.status}`)
              }}</span>
            </VTooltip>
          </div>
        </template>
      </v-list-item>
    </v-list>
    <div
      v-if="!threadsStore.loadingThreads && threadsStore.hasMoreThreads"
      class="d-flex justify-center pa-4"
    >
      <v-btn
        color="primary"
        variant="tonal"
        :loading="threadsStore.loadingMoreThreads"
        :disabled="threadsStore.loadingMoreThreads"
        @click="loadMoreThreads"
      >
        {{ $t('components.messageThread.loadMore') }}
      </v-btn>
    </div>
  </div>
</template>

<style scoped>
.archived-banner {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0.75rem 0.75rem 0;
  padding: 0.45rem 0.75rem;
  border-radius: 0.6rem;
  background: rgb(var(--v-theme-warning), 0.14);
  color: rgb(var(--v-theme-warning));
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.thread-item {
  padding-block: 0.6rem !important;
}

.thread-item-title {
  font-weight: 600;
}

.thread-item-preview {
  max-width: 230px;
}
</style>
