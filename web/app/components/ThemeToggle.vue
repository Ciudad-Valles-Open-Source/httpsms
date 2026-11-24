<script setup lang="ts">
import {
  mdiPaletteOutline,
  mdiWeatherNight,
  mdiWhiteBalanceSunny,
} from '@mdi/js'

withDefaults(
  defineProps<{
    size?: string
    location?: 'bottom' | 'end'
    /** Open the full appearance menu (mode + palette) instead of toggling. */
    withPalette?: boolean
  }>(),
  { size: 'small', location: 'bottom', withPalette: false },
)

const { isDark, toggleTheme } = useThemePreference()
</script>

<template>
  <VMenu
    v-if="withPalette"
    :location="location"
    :close-on-content-click="false"
    offset="8"
  >
    <template #activator="{ props: menuProps }">
      <VTooltip :location="location" :text="$t('components.theme.title')">
        <template #activator="{ props: tooltipProps }">
          <VBtn
            v-bind="{ ...menuProps, ...tooltipProps }"
            icon
            variant="text"
            :size="size"
            :aria-label="$t('components.theme.title')"
          >
            <VIcon :icon="mdiPaletteOutline" />
          </VBtn>
        </template>
      </VTooltip>
    </template>
    <VCard min-width="260" class="pa-4">
      <p class="text-title-small font-weight-bold mb-3">
        {{ $t('components.theme.title') }}
      </p>
      <AppearancePicker />
    </VCard>
  </VMenu>
  <VTooltip
    v-else
    :location="location"
    :text="
      isDark ? $t('components.theme.toLight') : $t('components.theme.toDark')
    "
  >
    <template #activator="{ props: tooltipProps }">
      <VBtn
        v-bind="tooltipProps"
        icon
        variant="text"
        :size="size"
        :aria-label="
          isDark
            ? $t('components.theme.toLight')
            : $t('components.theme.toDark')
        "
        @click="toggleTheme"
      >
        <VIcon :icon="isDark ? mdiWhiteBalanceSunny : mdiWeatherNight" />
      </VBtn>
    </template>
  </VTooltip>
</template>
