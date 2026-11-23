<script setup lang="ts">
import { useDisplay } from 'vuetify'

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
    <v-app-bar color="#121212" elevation="0" border="b">
      <v-app-bar-nav-icon
        v-if="!lgAndUp"
        :aria-label="$t('docs.toggleNavigation')"
        @click="drawer = !drawer"
      />
      <NuxtLink
        to="/"
        class="text-decoration-none d-flex align-center"
        :class="{ 'ml-4': lgAndUp }"
      >
        <v-avatar
          color="#121212"
          class="pa-1"
          rounded="0"
          :image="'/img/logo.svg'"
          :size="34"
        />
        <span v-if="mdAndUp" class="text-headline-small ml-1 text-white">
          httpSMS
        </span>
      </NuxtLink>
      <NuxtLink :to="docsHome" class="text-decoration-none ml-3">
        <v-chip color="primary" variant="tonal" size="small" label>
          {{ $t('docs.title') }}
        </v-chip>
      </NuxtLink>
      <v-spacer />
      <LanguageSwitcher />
      <v-btn
        v-if="mdAndUp"
        variant="text"
        color="primary"
        class="mr-1"
        :to="{ name: 'blog' }"
      >
        {{ $t('website.blog') }}
      </v-btn>
      <v-btn
        v-show="authStore.authStateChanged && authStore.authUser === null"
        color="primary"
        variant="flat"
        class="mr-4"
        :to="{ name: 'login' }"
      >
        {{ $t('website.login') }}
      </v-btn>
      <v-btn
        v-show="authStore.authStateChanged && authStore.authUser !== null"
        color="primary"
        variant="flat"
        class="mr-4"
        :to="{ name: 'threads' }"
      >
        {{ $t('website.dashboard') }}
      </v-btn>
    </v-app-bar>
    <v-navigation-drawer
      v-model="drawer"
      color="#121212"
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
