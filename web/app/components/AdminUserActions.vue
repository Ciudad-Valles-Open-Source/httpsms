<script setup lang="ts">
import {
  mdiAccountCheck,
  mdiAccountCancel,
  mdiAccountReactivate,
  mdiDelete,
} from '@mdi/js'
import {
  adminUserActions,
  type AdminUser,
  type AdminUserAction,
} from '~/utils/admin-users'

const props = withDefaults(
  defineProps<{
    user: AdminUser
    disabled?: boolean
    busyAction?: AdminUserAction | null
    compact?: boolean
  }>(),
  { disabled: false, busyAction: null, compact: false },
)

const emit = defineEmits<{
  action: [action: AdminUserAction]
}>()

const ACTION_STYLE: Record<AdminUserAction, { icon: string; color: string }> = {
  approve: { icon: mdiAccountCheck, color: 'success' },
  suspend: { icon: mdiAccountCancel, color: 'warning' },
  reactivate: { icon: mdiAccountReactivate, color: 'success' },
  delete: { icon: mdiDelete, color: 'error' },
}

const actions = computed(() => adminUserActions(props.user))
</script>

<template>
  <div class="d-flex flex-wrap ga-2" :class="{ 'justify-end': compact }">
    <template v-for="action in actions" :key="action">
      <VTooltip v-if="compact" location="top">
        <template #activator="{ props: tooltipProps }">
          <VBtn
            v-bind="tooltipProps"
            icon
            size="small"
            variant="tonal"
            :color="ACTION_STYLE[action].color"
            :disabled="disabled || (!!busyAction && busyAction !== action)"
            :loading="busyAction === action"
            :aria-label="$t(`pages.adminUsers.actions.${action}`)"
            @click="emit('action', action)"
          >
            <VIcon :icon="ACTION_STYLE[action].icon" />
          </VBtn>
        </template>
        {{ $t(`pages.adminUsers.actions.${action}`) }}
      </VTooltip>
      <VBtn
        v-else
        size="small"
        variant="tonal"
        :color="ACTION_STYLE[action].color"
        :disabled="disabled || (!!busyAction && busyAction !== action)"
        :loading="busyAction === action"
        @click="emit('action', action)"
      >
        <VIcon start :icon="ACTION_STYLE[action].icon" />
        {{ $t(`pages.adminUsers.actions.${action}`) }}
      </VBtn>
    </template>
  </div>
</template>
