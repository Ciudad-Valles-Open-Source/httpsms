<script setup lang="ts">
import { mdiApi, mdiOpenInNew } from '@mdi/js'

defineProps<{
  currentSlug: string
}>()

const placeholders = useDocPlaceholders()
const groups = DOCS_GROUPS.map((group) => ({
  ...group,
  pages: group.pages.map((page) => ({ ...page, to: docPath(page.slug) })),
}))
</script>

<template>
  <VList nav density="compact" class="docs-sidebar py-4" color="primary">
    <template v-for="group in groups" :key="group.titleKey">
      <VListSubheader class="text-uppercase font-weight-bold">
        {{ $t(group.titleKey) }}
      </VListSubheader>
      <VListItem
        v-for="page in group.pages"
        :key="page.slug"
        :to="page.to"
        :active="page.slug === currentSlug"
        :title="$t(page.titleKey)"
        rounded="lg"
      />
    </template>
    <VDivider class="my-3" />
    <VListItem
      :href="placeholders.API_DOCS_URL"
      target="_blank"
      rel="noopener noreferrer"
      rounded="lg"
      :prepend-icon="mdiApi"
      :title="$t('docs.apiReference')"
    >
      <template #append>
        <VIcon size="small" :icon="mdiOpenInNew" />
      </template>
    </VListItem>
  </VList>
</template>

<style scoped>
.docs-sidebar :deep(.v-list-subheader) {
  min-height: 36px;
  margin-top: 8px;
  font-size: 0.75rem;
  letter-spacing: 0.06em;
}
</style>
