<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiGithub,
  mdiHeart,
  mdiLightbulbOn50,
  mdiEyeOffOutline,
  mdiPost,
  mdiScaleBalance,
  mdiEmailOutline,
  mdiBookOpenVariant,
} from '@mdi/js'

const { lgAndUp, mdAndUp } = useDisplay()
const authStore = useAuthStore()
const appStore = useAppStore()
const { githubIssuesUrl } = usePublicUrls()
</script>

<template>
  <v-app>
    <v-app-bar color="#121212" elevation="0">
      <v-container>
        <v-row>
          <v-col class="w-full d-flex">
            <NuxtLink to="/" class="text-decoration-none d-flex align-baseline">
              <v-avatar
                color="#121212"
                class="pa-1"
                rounded="0"
                :image="'/img/logo.svg'"
                :size="38"
              />
              <h3
                v-if="lgAndUp"
                class="text-headline-large mb-0 ml-1 mt-6 text-white"
              >
                httpSMS
              </h3>
            </NuxtLink>
            <v-spacer />
            <LanguageSwitcher class="mt-6" />
            <v-btn
              v-show="lgAndUp"
              size="large"
              variant="text"
              color="primary"
              class="my-5 mr-2"
              to="/docs"
            >
              {{ $t('website.docs') }}
            </v-btn>
            <v-btn
              v-show="lgAndUp"
              size="large"
              variant="text"
              color="primary"
              class="my-5 mr-2"
              :to="{ name: 'blog' }"
            >
              {{ $t('website.blog') }}
            </v-btn>
            <v-btn
              v-show="
                lgAndUp &&
                authStore.authStateChanged &&
                authStore.authUser === null
              "
              size="large"
              variant="text"
              color="primary"
              class="my-5 mr-2"
              :to="{ name: 'login' }"
            >
              {{ $t('website.login') }}
            </v-btn>
            <v-btn
              v-show="authStore.authStateChanged && authStore.authUser === null"
              color="primary"
              variant="flat"
              :class="{ 'mt-5': mdAndUp, 'mt-1': !mdAndUp }"
              :size="lgAndUp ? 'large' : 'default'"
              :to="{ name: 'login' }"
            >
              {{ $t('website.getStarted') }}
              <span v-show="lgAndUp">&nbsp;{{ $t('website.forFree') }}</span>
            </v-btn>
            <div
              v-show="authStore.authStateChanged && authStore.authUser !== null"
              class="position-relative d-inline-block"
            >
              <v-btn
                color="primary"
                variant="flat"
                :class="{ 'mt-5': mdAndUp, 'mt-1': !mdAndUp }"
                :size="lgAndUp ? 'large' : 'default'"
                :to="{ name: 'threads' }"
              >
                {{ $t('website.dashboard') }}
              </v-btn>
              <RedirectPromptPopover />
            </div>
          </v-col>
        </v-row>
      </v-container>
    </v-app-bar>
    <v-main>
      <AppToast />
      <slot />
    </v-main>
    <v-footer>
      <v-container>
        <v-row>
          <v-col cols="12" md="3">
            <NuxtLink to="/" class="text-decoration-none d-flex mt-n6">
              <v-avatar
                color="#212121"
                class="mt-8 pa-1"
                rounded="0"
                :image="'/img/logo.svg'"
                :size="38"
              />
              <h3 class="text-headline-large ml-1 mb-0 text-white">httpSMS</h3>
            </NuxtLink>
            <div class="text-title-medium mb-4 text-medium-emphasis">
              {{ $t('website.madeWith') }}
              <v-icon color="#cf1112" :icon="mdiHeart" />
              {{ $t('website.inTallinn') }}
            </div>
            <p class="mt-n3">
              <v-btn
                :href="appStore.appData.githubUrl"
                color="#ffffff"
                class="ml-n3"
                variant="text"
                :icon="mdiGithub"
              />
            </p>
          </v-col>
          <v-col cols="12" md="3">
            <h2 class="text-headline-small mb-2">
              {{ $t('website.resources') }}
            </h2>
            <ul style="list-style: none" class="pa-0">
              <li class="mb-2">
                <NuxtLink
                  class="text-white text-decoration-none footer-link"
                  to="/blog"
                >
                  {{ $t('website.blog') }}
                  <v-icon size="small" :icon="mdiPost" />
                </NuxtLink>
              </li>
            </ul>
          </v-col>
          <v-col cols="12" md="3">
            <h2 class="text-headline-small mb-2">
              {{ $t('website.developers') }}
            </h2>
            <ul style="list-style: none" class="pa-0">
              <li class="mb-2">
                <NuxtLink
                  to="/docs"
                  class="text-white text-decoration-none footer-link"
                >
                  {{ $t('website.documentation') }}
                  <v-icon size="small" :icon="mdiBookOpenVariant" />
                </NuxtLink>
              </li>
              <li class="mb-2">
                <a
                  :href="appStore.appData.githubUrl"
                  class="text-white text-decoration-none footer-link"
                >
                  {{ $t('website.github') }}
                  <v-icon size="small" :icon="mdiGithub" />
                </a>
              </li>
              <li class="mb-2">
                <a
                  :href="`${githubIssuesUrl}/new`"
                  class="text-white text-decoration-none footer-link"
                >
                  {{ $t('website.requestFeature') }}
                  <v-icon
                    size="small"
                    color="yellow"
                    :icon="mdiLightbulbOn50"
                  />
                </a>
              </li>
            </ul>
          </v-col>
          <v-col cols="12" md="3">
            <h2 class="text-headline-small mb-2">{{ $t('website.legal') }}</h2>
            <ul style="list-style: none" class="pa-0">
              <li class="mb-2">
                <NuxtLink
                  class="text-white text-decoration-none footer-link"
                  to="/terms-and-conditions"
                >
                  {{ $t('website.termsAndConditions') }}
                  <v-icon size="small" :icon="mdiScaleBalance" />
                </NuxtLink>
              </li>
              <li class="mb-2">
                <NuxtLink
                  class="text-white text-decoration-none footer-link"
                  to="/privacy-policy"
                >
                  {{ $t('website.privacyPolicy') }}
                  <v-icon size="small" :icon="mdiEyeOffOutline" />
                </NuxtLink>
              </li>
              <li class="mt-2">
                <a
                  class="text-white text-decoration-none footer-link"
                  :href="githubIssuesUrl"
                >
                  {{ $t('website.contactSupport') }}
                  <v-icon size="small" :icon="mdiEmailOutline" />
                </a>
              </li>
            </ul>
          </v-col>
        </v-row>
      </v-container>
    </v-footer>
  </v-app>
</template>

<style scoped>
.footer-link:hover {
  text-decoration: underline !important;
}
</style>
