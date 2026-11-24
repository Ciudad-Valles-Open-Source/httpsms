<script setup lang="ts">
import { useDisplay } from 'vuetify'

useThemeScope('public')

const route = useRoute()
const { lgAndUp, mdAndUp } = useDisplay()
const authStore = useAuthStore()
const docsHome = DOCS_BASE_PATH

// null lets Vuetify open the drawer on desktop and keep it closed on mobile.
const drawer = ref<boolean | null>(null)

const currentSlug = computed(
  () => (route.params.slug as string | undefined) || DOCS_DEFAULT_SLUG,
)

// Close the mobile drawer after picking a page.
watch(
  () => route.path,
  () => {
    if (!lgAndUp.value) drawer.value = false
  },
)
</script>

<template>
  <v-app>
    <v-app-bar class="app-glass" height="64">
      <v-app-bar-nav-icon
        v-if="!lgAndUp"
        :aria-label="$t('docs.toggleNavigation')"
        @click="drawer = !drawer"
      />
      <AppLogo :show-name="mdAndUp" :size="34" :class="{ 'ml-4': lgAndUp }" />
      <NuxtLink :to="docsHome" class="text-decoration-none ml-3">
        <v-chip
          color="primary"
          variant="tonal"
          size="small"
          class="font-weight-bold"
        >
          {{ $t('docs.title') }}
        </v-chip>
      </NuxtLink>
      <v-spacer />
      <ThemeToggle />
      <LanguageSwitcher />
      <v-btn v-if="mdAndUp" variant="text" class="mr-1" :to="{ name: 'blog' }">
        {{ $t('website.blog') }}
      </v-btn>
      <v-btn
        v-show="authStore.authStateChanged && authStore.authUser === null"
        class="btn-gradient mr-4"
        :to="{ name: 'login' }"
      >
        {{ $t('website.login') }}
      </v-btn>
      <v-btn
        v-show="authStore.authStateChanged && authStore.authUser !== null"
        class="btn-gradient mr-4"
        :to="{ name: 'threads' }"
      >
        {{ $t('website.dashboard') }}
      </v-btn>
    </v-app-bar>
    <v-navigation-drawer
      v-model="drawer"
      color="background"
      mobile-breakpoint="lg"
      :width="280"
    >
      <DocsSidebar :current-slug="currentSlug" />
    </v-navigation-drawer>
    <v-main>
      <AppToast />
      <slot />
    </v-main>
  </v-app>
</template>
