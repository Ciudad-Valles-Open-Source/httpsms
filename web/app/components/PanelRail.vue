<script setup lang="ts">
import { mdiBookOpenVariant, mdiDownload, mdiLogout } from '@mdi/js'
import type { PanelNavItem } from '~/composables/usePanelNav'

const appStore = useAppStore()
const phonesStore = usePhonesStore()
const { items, isActive, badgeFor } = usePanelNav()
const logout = useLogout()
const { t } = useI18n()

function itemLabel(item: PanelNavItem): string {
  const badge = badgeFor(item)
  const label = t(item.label)
  return badge
    ? `${label} · ${t('components.messageThreadHeader.pendingUsers', { count: badge })}`
    : label
}
</script>

<template>
  <nav class="panel-rail" :aria-label="$t('components.panelNav.label')">
    <AppLogo to="/threads" :show-name="false" :size="40" class="mb-4" />

    <div class="panel-rail-items">
      <VTooltip
        v-for="item in items"
        :key="item.name"
        location="end"
        :text="itemLabel(item)"
      >
        <template #activator="{ props: tooltipProps }">
          <VBtn
            v-bind="tooltipProps"
            :to="{ name: item.name }"
            :aria-label="itemLabel(item)"
            :active="isActive(item)"
            :color="isActive(item) ? 'primary' : undefined"
            :variant="isActive(item) ? 'tonal' : 'text'"
            icon
            rounded="lg"
            class="panel-rail-btn"
          >
            <VBadge
              v-if="badgeFor(item)"
              color="warning"
              :content="badgeFor(item)"
              offset-x="-4"
              offset-y="-4"
            >
              <VIcon :icon="item.icon" />
            </VBadge>
            <VIcon v-else :icon="item.icon" />
          </VBtn>
        </template>
      </VTooltip>
    </div>

    <VSpacer />

    <div class="panel-rail-items">
      <VTooltip
        v-if="phonesStore.owner"
        location="end"
        :text="$t('components.messageThreadHeader.downloadApp')"
      >
        <template #activator="{ props: tooltipProps }">
          <VBtn
            v-bind="tooltipProps"
            :href="appStore.appData.appDownloadUrl"
            :aria-label="$t('components.messageThreadHeader.downloadApp')"
            icon
            variant="text"
            rounded="lg"
            class="panel-rail-btn"
          >
            <VIcon :icon="mdiDownload" />
          </VBtn>
        </template>
      </VTooltip>
      <VTooltip
        location="end"
        :text="$t('components.messageThreadHeader.documentation')"
      >
        <template #activator="{ props: tooltipProps }">
          <VBtn
            v-bind="tooltipProps"
            to="/docs"
            :aria-label="$t('components.messageThreadHeader.documentation')"
            icon
            variant="text"
            rounded="lg"
            class="panel-rail-btn"
          >
            <VIcon :icon="mdiBookOpenVariant" />
          </VBtn>
        </template>
      </VTooltip>
      <ThemeToggle size="default" location="end" with-palette />
      <LanguageSwitcher />
      <VTooltip
        location="end"
        :text="$t('components.messageThreadHeader.logout')"
      >
        <template #activator="{ props: tooltipProps }">
          <VBtn
            v-bind="tooltipProps"
            :aria-label="$t('components.messageThreadHeader.logout')"
            icon
            variant="text"
            rounded="lg"
            class="panel-rail-btn"
            @click="logout"
          >
            <VIcon :icon="mdiLogout" />
          </VBtn>
        </template>
      </VTooltip>
    </div>
  </nav>
</template>

<style scoped>
.panel-rail {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  padding: 1rem 0;
}

.panel-rail-items {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.35rem;
}

.panel-rail-btn {
  width: 48px;
  height: 48px;
}
</style>
