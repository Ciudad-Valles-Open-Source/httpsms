<script setup lang="ts">
import { mdiCheck, mdiWeatherNight, mdiWhiteBalanceSunny } from '@mdi/js'
import { THEME_PALETTES } from '~/composables/useThemePreference'

const { mode, palette, setMode, setPalette } = useThemePreference()
</script>

<template>
  <div class="appearance-picker">
    <p class="appearance-label">{{ $t('components.theme.mode') }}</p>
    <VBtnToggle
      :model-value="mode"
      mandatory
      divided
      variant="outlined"
      density="comfortable"
      rounded="lg"
      color="primary"
      class="appearance-modes"
      @update:model-value="setMode"
    >
      <VBtn value="dark" :prepend-icon="mdiWeatherNight">
        {{ $t('components.theme.dark') }}
      </VBtn>
      <VBtn value="light" :prepend-icon="mdiWhiteBalanceSunny">
        {{ $t('components.theme.light') }}
      </VBtn>
    </VBtnToggle>

    <p class="appearance-label mt-4">{{ $t('components.theme.palette') }}</p>
    <div
      class="appearance-palettes"
      role="radiogroup"
      :aria-label="$t('components.theme.palette')"
    >
      <button
        v-for="option in THEME_PALETTES"
        :key="option.value"
        type="button"
        role="radio"
        class="palette-option"
        :class="{ 'palette-option-active': palette === option.value }"
        :aria-checked="palette === option.value"
        @click="setPalette(option.value)"
      >
        <span
          class="palette-swatch"
          :style="{
            background: `linear-gradient(135deg, ${option.swatch[0]}, ${option.swatch[1]})`,
          }"
        >
          <VIcon
            v-if="palette === option.value"
            size="16"
            color="white"
            :icon="mdiCheck"
          />
        </span>
        {{ $t(option.label) }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.appearance-label {
  margin: 0 0 0.5rem;
  color: rgb(var(--v-theme-on-surface), var(--v-medium-emphasis-opacity));
  font-size: 0.72rem;
  font-weight: 800;
  letter-spacing: 0.1em;
  text-transform: uppercase;
}

.appearance-modes {
  width: 100%;
}

.appearance-modes .v-btn {
  flex: 1;
}

.appearance-palettes {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.5rem;
}

.palette-option {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.5rem 0.6rem;
  border: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  border-radius: 0.75rem;
  background: transparent;
  color: rgb(var(--v-theme-on-surface));
  font: inherit;
  font-size: 0.875rem;
  font-weight: 700;
  cursor: pointer;
  transition: border-color 0.15s;
}

.palette-option:hover,
.palette-option-active {
  border-color: rgb(var(--v-theme-primary));
}

.palette-swatch {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  width: 1.75rem;
  height: 1.75rem;
  border-radius: 50%;
}
</style>
