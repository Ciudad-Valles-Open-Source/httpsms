<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiGithub,
  mdiHeart,
  mdiMenu,
  mdiArrowRight,
  mdiBookOpenVariant,
  mdiPost,
  mdiLogin,
  mdiViewDashboardOutline,
} from '@mdi/js'

useThemeScope('public')

const { mdAndUp } = useDisplay()
const authStore = useAuthStore()
const appStore = useAppStore()
const route = useRoute()
const { githubIssuesUrl } = usePublicUrls()

const mobileMenu = ref(false)

const isGuest = computed(
  () => authStore.authStateChanged && authStore.authUser === null,
)
const isSignedIn = computed(
  () => authStore.authStateChanged && authStore.authUser !== null,
)

watch(
  () => route.fullPath,
  () => {
    mobileMenu.value = false
  },
)
</script>

<template>
  <v-app>
    <v-app-bar class="app-glass" height="68">
      <div class="site-bar">
        <AppLogo />
        <v-spacer />
        <nav
          v-if="mdAndUp"
          class="site-bar-links"
          :aria-label="$t('website.mainNavigation')"
        >
          <v-btn variant="text" to="/docs">{{ $t('website.docs') }}</v-btn>
          <v-btn variant="text" :to="{ name: 'blog' }">
            {{ $t('website.blog') }}
          </v-btn>
        </nav>
        <ThemeToggle />
        <LanguageSwitcher />
        <template v-if="mdAndUp">
          <v-btn
            v-show="isGuest"
            variant="text"
            class="ml-1"
            :to="{ name: 'login' }"
          >
            {{ $t('website.login') }}
          </v-btn>
          <v-btn
            v-show="isGuest"
            class="btn-gradient ml-2"
            :to="{ name: 'login' }"
          >
            {{ $t('website.getStarted') }}
            <v-icon end :icon="mdiArrowRight" />
          </v-btn>
          <div
            v-show="isSignedIn"
            class="position-relative d-inline-block ml-2"
          >
            <v-btn class="btn-gradient" :to="{ name: 'threads' }">
              {{ $t('website.dashboard') }}
              <v-icon end :icon="mdiArrowRight" />
            </v-btn>
            <RedirectPromptPopover />
          </div>
        </template>
        <v-btn
          v-else
          icon
          variant="text"
          :aria-label="$t('website.openMenu')"
          :aria-expanded="mobileMenu"
          @click="mobileMenu = !mobileMenu"
        >
          <v-icon :icon="mdiMenu" />
        </v-btn>
      </div>
    </v-app-bar>

    <v-navigation-drawer
      v-if="!mdAndUp"
      v-model="mobileMenu"
      location="top"
      temporary
      color="background"
      class="site-mobile-menu"
    >
      <v-list nav class="pa-3">
        <v-list-item
          to="/docs"
          :prepend-icon="mdiBookOpenVariant"
          :title="$t('website.docs')"
        />
        <v-list-item
          :to="{ name: 'blog' }"
          :prepend-icon="mdiPost"
          :title="$t('website.blog')"
        />
        <v-list-item
          :href="appStore.appData.githubUrl"
          :prepend-icon="mdiGithub"
          :title="$t('website.github')"
        />
      </v-list>
      <div class="px-4 pb-4 d-flex flex-column ga-2">
        <template v-if="isGuest">
          <v-btn
            block
            size="large"
            class="btn-gradient"
            :to="{ name: 'login' }"
          >
            {{ $t('website.getStarted') }}
            <v-icon end :icon="mdiArrowRight" />
          </v-btn>
          <v-btn
            block
            size="large"
            variant="outlined"
            :prepend-icon="mdiLogin"
            :to="{ name: 'login' }"
          >
            {{ $t('website.login') }}
          </v-btn>
        </template>
        <v-btn
          v-if="isSignedIn"
          block
          size="large"
          class="btn-gradient"
          :prepend-icon="mdiViewDashboardOutline"
          :to="{ name: 'threads' }"
        >
          {{ $t('website.dashboard') }}
        </v-btn>
      </div>
    </v-navigation-drawer>

    <v-main>
      <AppToast />
      <slot />
    </v-main>

    <footer class="site-footer">
      <v-container class="site-container">
        <v-row>
          <v-col cols="12" md="4">
            <AppLogo />
            <p class="text-body-medium text-medium-emphasis mt-4 mb-3 pr-md-8">
              {{ $t('website.footerTagline') }}
            </p>
            <div class="text-body-small text-medium-emphasis">
              {{ $t('website.madeWith') }}
              <v-icon size="14" color="error" :icon="mdiHeart" />
              {{ $t('website.inTallinn') }}
            </div>
          </v-col>
          <v-col cols="6" md="2" offset-md="2">
            <h2 class="footer-heading">{{ $t('website.resources') }}</h2>
            <ul class="footer-links">
              <li>
                <NuxtLink to="/blog">{{ $t('website.blog') }}</NuxtLink>
              </li>
              <li>
                <NuxtLink to="/docs">{{
                  $t('website.documentation')
                }}</NuxtLink>
              </li>
            </ul>
          </v-col>
          <v-col cols="6" md="2">
            <h2 class="footer-heading">{{ $t('website.developers') }}</h2>
            <ul class="footer-links">
              <li>
                <a :href="appStore.appData.githubUrl">{{
                  $t('website.github')
                }}</a>
              </li>
              <li>
                <a :href="`${githubIssuesUrl}/new`">{{
                  $t('website.requestFeature')
                }}</a>
              </li>
            </ul>
          </v-col>
          <v-col cols="12" md="2">
            <h2 class="footer-heading">{{ $t('website.legal') }}</h2>
            <ul class="footer-links">
              <li>
                <NuxtLink to="/terms-and-conditions">{{
                  $t('website.termsAndConditions')
                }}</NuxtLink>
              </li>
              <li>
                <NuxtLink to="/privacy-policy">{{
                  $t('website.privacyPolicy')
                }}</NuxtLink>
              </li>
              <li>
                <a :href="githubIssuesUrl">{{
                  $t('website.contactSupport')
                }}</a>
              </li>
            </ul>
          </v-col>
        </v-row>
      </v-container>
    </footer>
  </v-app>
</template>

<style scoped>
.site-bar {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.site-bar-links {
  display: flex;
  gap: 0.25rem;
  margin-right: 0.5rem;
}

.site-container {
  max-width: 1200px;
}

.site-footer {
  padding: 3rem 0 2rem;
  border-top: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  background: rgb(var(--v-theme-background));
}

.footer-heading {
  margin-bottom: 0.75rem;
  font-size: 0.78rem;
  font-weight: 800;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: rgb(var(--v-theme-on-background), var(--v-medium-emphasis-opacity));
}

.footer-links {
  padding: 0;
  list-style: none;
}

.footer-links li + li {
  margin-top: 0.5rem;
}

.footer-links a {
  color: rgb(var(--v-theme-on-background)) !important;
  text-decoration: none;
  transition: color 0.15s;
}

.footer-links a:hover {
  color: rgb(var(--v-theme-primary)) !important;
}

.site-mobile-menu {
  height: auto !important;
}
</style>
