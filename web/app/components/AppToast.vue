<script setup lang="ts">
import { useDisplay } from 'vuetify'
import { mdiCheck, mdiInformation } from '@mdi/js'

const { lgAndUp } = useDisplay()
const notificationsStore = useNotificationsStore()

const { t } = useI18n()
// API responses and some stores send English text; show it translated.
const message = computed(() =>
  translateNotice(notificationsStore.notification.message, t),
)

const notificationActive = computed({
  get: () => notificationsStore.notification.active,
  set: () => notificationsStore.disableNotification(),
})
</script>

<template>
  <v-snackbar
    v-model="notificationActive"
    variant="tonal"
    :color="notificationsStore.notification.type"
    :timeout="notificationsStore.notification.timeout"
    :z-index="10000"
  >
    <v-icon
      v-if="notificationsStore.notification.type === 'success'"
      :color="notificationsStore.notification.type"
      :icon="mdiCheck"
    />
    <v-icon
      v-if="notificationsStore.notification.type === 'info'"
      :color="notificationsStore.notification.type"
      :icon="mdiInformation"
    />
    {{ message }}
    <template #actions>
      <v-btn
        v-if="lgAndUp"
        :color="notificationsStore.notification.type"
        variant="text"
        @click="notificationsStore.disableNotification()"
      >
        <span class="font-weight-bold">{{
          $t('components.appToast.close')
        }}</span>
      </v-btn>
    </template>
  </v-snackbar>
</template>
