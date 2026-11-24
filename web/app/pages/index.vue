<script setup lang="ts">
import {
  mdiCheckCircle,
  mdiLockOutline,
  mdiLanguagePython,
  mdiCellphoneKey,
  mdiLanguageJavascript,
  mdiLanguagePhp,
  mdiPlus,
  mdiMinus,
  mdiLanguageCsharp,
  mdiLanguageJava,
  mdiMicrosoftExcel,
  mdiWebhook,
  mdiClockOutline,
  mdiArrowRight,
  mdiPowershell,
  mdiLanguageGo,
  mdiApi,
  mdiSpeedometer,
  mdiBellAlertOutline,
  mdiGithub,
  mdiPuzzleOutline,
  mdiBookOpenVariant,
  mdiShieldLockOutline,
} from '@mdi/js'

definePageMeta({
  layout: 'website',
  middleware: ['redirect-to-threads'],
})

const { t } = useI18n()

const {
  apiBaseUrl,
  assetUrl,
  docsUrl,
  githubUrl,
  sendMessageUrl,
  settingsUrl,
} = usePublicUrls()

useSeoMeta({
  title: () => t('pages.index.metaTitle'),
  description: () => t('pages.index.metaDescription'),
  ogTitle: () => t('pages.index.metaOgTitle'),
  ogDescription: () => t('pages.index.metaOgDescription'),
  ogImage: assetUrl('/header.png'),
  twitterCard: 'summary_large_image',
})

const config = useRuntimeConfig()
const { mdAndUp } = useVDisplay()
const { darkThemeName } = useThemePreference()

const selectedTab = ref('javascript')
const faqPanel = ref<number | undefined>(undefined)

interface Feature {
  icon: string
  title: string
  desc: string
  /** The description contains trusted markup from the locale files. */
  html?: boolean
  noCode?: boolean
  to?: string
  href?: string
  cta?: string
}

const features = computed<Feature[]>(() => [
  {
    icon: mdiMicrosoftExcel,
    title: 'pages.index.bulkSmsTitle',
    desc: 'pages.index.bulkSmsDesc',
    html: true,
    noCode: true,
    to: '/blog/how-to-send-sms-messages-from-excel',
    cta: 'pages.index.integrationGuideBtn',
  },
  {
    icon: mdiPuzzleOutline,
    title: 'pages.index.integrationsTitle',
    desc: 'pages.index.integrationsDesc',
    noCode: true,
    to: '/blog/send-sms-when-new-row-is-added-to-google-sheets-using-zapier',
    cta: 'pages.index.zapierIntegrationGuideBtn',
  },
  {
    icon: mdiWebhook,
    title: 'pages.index.webhooksTitle',
    desc: 'pages.index.webhooksDesc',
    to: '/docs/webhooks-introduction',
    cta: 'pages.index.documentationBtn',
  },
  {
    icon: mdiSpeedometer,
    title: 'pages.index.controlSendingTitle',
    desc: 'pages.index.controlSendingDesc',
    to: '/docs/control-sms-send-rate',
    cta: 'pages.index.documentationBtn',
  },
  {
    icon: mdiBellAlertOutline,
    title: 'pages.index.monitoringTitle',
    desc: 'pages.index.monitoringDesc',
  },
  {
    icon: mdiShieldLockOutline,
    title: 'pages.index.encryptionTitle',
    desc: 'pages.index.encryptionDesc',
    html: true,
    to: '/blog/end-to-end-encryption-to-sms-messages',
    cta: 'pages.index.setupEncryptionBtn',
  },
  {
    icon: mdiCellphoneKey,
    title: 'pages.index.multiplePhonesTitle',
    desc: 'pages.index.multiplePhonesDesc',
    to: '/docs/phone-api-keys',
    cta: 'pages.index.documentationBtn',
  },
  {
    icon: mdiClockOutline,
    title: 'pages.index.scheduleTextMessagesTitle',
    desc: 'pages.index.scheduleTextMessagesDesc',
    to: '/docs/scheduling-sms-messages',
    cta: 'pages.index.documentationBtn',
  },
  {
    icon: mdiGithub,
    title: 'pages.index.openSourceTitle',
    desc: 'pages.index.openSourceDesc',
    href: githubUrl,
    cta: 'pages.index.viewOnGithubBtn',
  },
])

const faqs = [
  { q: 'pages.index.faq1Question', a: 'pages.index.faq1Answer' },
  { q: 'pages.index.faq2Question', a: 'pages.index.faq2Answer' },
  { q: 'pages.index.faq3Question', a: 'pages.index.faq3Answer', html: true },
  { q: 'pages.index.faq4Question', a: 'pages.index.faq4Answer' },
]
</script>

<template>
  <div>
    <!-- Hero -->
    <section class="hero">
      <VContainer class="site-container">
        <VRow align="center" class="py-md-16 py-8">
          <VCol cols="12" md="6" class="text-center text-md-left">
            <span class="hero-badge mb-6">
              <span class="dot" />
              {{ $t('pages.index.heroBadge') }}
            </span>
            <h1 class="hero-title">
              {{ $t('pages.index.heroTitleLead') }}
              <span class="text-gradient">{{
                $t('pages.index.heroTitleAccent')
              }}</span>
            </h1>
            <p
              class="hero-subtitle text-medium-emphasis"
              v-html="$t('pages.index.heroSubtitle')"
            ></p>
            <div
              class="d-flex flex-wrap ga-3 justify-center justify-md-start mt-8"
            >
              <VBtn size="x-large" class="btn-gradient" to="/login">
                {{ $t('pages.index.getStartedBtn') }}
                <VIcon end :icon="mdiArrowRight" />
              </VBtn>
              <VBtn
                size="x-large"
                variant="outlined"
                color="primary"
                to="/docs"
                :prepend-icon="mdiBookOpenVariant"
              >
                {{ $t('pages.index.documentationBtn') }}
              </VBtn>
            </div>
            <ul class="hero-checks justify-center justify-md-start">
              <li>
                <VIcon size="18" color="success" :icon="mdiCheckCircle" />
                {{ $t('pages.index.freeToUse') }}
              </li>
              <li>
                <VIcon size="18" color="success" :icon="mdiCheckCircle" />
                {{ $t('pages.index.openSourceBadge') }}
              </li>
              <li>
                <VIcon size="18" color="success" :icon="mdiCheckCircle" />
                {{ $t('pages.index.selfHostedBadge') }}
              </li>
            </ul>
          </VCol>
          <VCol cols="12" md="6">
            <div class="hero-visual">
              <div class="hero-frame">
                <img
                  src="/header.png"
                  alt="httpSMS"
                  width="1280"
                  height="720"
                  loading="eager"
                />
              </div>
              <span v-if="mdAndUp" class="float-chip float-chip-1">
                <VIcon size="16" :icon="mdiApi" />
                {{ $t('pages.index.chipApi') }}
              </span>
              <span v-if="mdAndUp" class="float-chip float-chip-2">
                <VIcon size="16" :icon="mdiWebhook" />
                {{ $t('pages.index.webhooksTitle') }}
              </span>
              <span v-if="mdAndUp" class="float-chip float-chip-3">
                <VIcon size="16" :icon="mdiLockOutline" />
                {{ $t('pages.index.chipEncryption') }}
              </span>
            </div>
          </VCol>
        </VRow>
      </VContainer>
    </section>

    <!-- Features -->
    <section class="section section-alt">
      <VContainer class="site-container">
        <div class="section-head">
          <span class="kicker">{{ $t('pages.index.featuresKicker') }}</span>
          <h2 class="section-title">{{ $t('pages.index.featuresTitle') }}</h2>
        </div>
        <VRow>
          <VCol
            v-for="feature in features"
            :key="feature.title"
            cols="12"
            sm="6"
            lg="4"
          >
            <article class="card-soft card-hover feature-card">
              <div class="d-flex align-center ga-3 mb-4">
                <span class="icon-tile">
                  <VIcon :icon="feature.icon" />
                </span>
                <VChip
                  v-if="feature.noCode"
                  size="small"
                  color="secondary"
                  variant="tonal"
                  class="font-weight-bold"
                >
                  {{ $t('pages.index.noCodeBadge') }}
                </VChip>
              </div>
              <h3 class="feature-card-title">{{ $t(feature.title) }}</h3>
              <p
                v-if="feature.html"
                class="feature-card-desc"
                v-html="$t(feature.desc)"
              ></p>
              <p v-else class="feature-card-desc">{{ $t(feature.desc) }}</p>
              <VBtn
                v-if="feature.cta"
                variant="text"
                color="primary"
                class="feature-card-cta px-0"
                :to="feature.to"
                :href="feature.href"
              >
                {{ $t(feature.cta) }}
                <VIcon end :icon="mdiArrowRight" />
              </VBtn>
            </article>
          </VCol>
        </VRow>
      </VContainer>
    </section>

    <!-- Get Started Section -->
    <section class="section">
      <VContainer class="site-container">
        <div class="section-head">
          <span class="kicker">{{ $t('pages.index.howItWorksKicker') }}</span>
          <h2 class="section-title">
            {{ $t('pages.index.getStartedSectionTitle') }}
          </h2>
        </div>
        <VRow>
          <VCol cols="12" md="4">
            <article class="card-soft step-card">
              <span class="step-card-num">1</span>
              <h3 class="feature-card-title">
                {{ $t('pages.index.step1Title') }}
              </h3>
              <p
                class="feature-card-desc"
                v-html="$t('pages.index.step1Desc')"
              ></p>
            </article>
          </VCol>
          <VCol cols="12" md="4">
            <article class="card-soft step-card">
              <span class="step-card-num">2</span>
              <h3 class="feature-card-title">
                {{ $t('pages.index.step2Title') }}
              </h3>
              <p
                class="feature-card-desc"
                v-html="
                  $t('pages.index.step2Desc').replace(
                    '{url}',
                    config.public.appDownloadUrl,
                  )
                "
              ></p>
            </article>
          </VCol>
          <VCol cols="12" md="4">
            <article class="card-soft step-card">
              <span class="step-card-num">3</span>
              <h3 class="feature-card-title">
                {{ $t('pages.index.step3Title') }}
              </h3>
              <p class="feature-card-desc">
                {{ $t('pages.index.step3Desc') }}
                <NuxtLink class="text-decoration-none" to="/docs">
                  {{ docsUrl }}
                </NuxtLink>
              </p>
            </article>
          </VCol>
        </VRow>
        <VThemeProvider
          :theme="darkThemeName"
          with-background
          class="code-window mt-8"
        >
          <div class="code-window-bar"><span /><span /><span /></div>
          <VTabs
            v-model="selectedTab"
            color="primary"
            density="compact"
            show-arrows
          >
            <VTab value="javascript">
              <VIcon
                color="#efd81d"
                class="mr-1"
                :icon="mdiLanguageJavascript"
              />
              Javascript
            </VTab>
            <VTab value="php">
              <VIcon color="#777bb3" class="mr-2" :icon="mdiLanguagePhp" />
              PHP
            </VTab>
            <VTab value="python">
              <VIcon color="#ffffff" class="mr-2" :icon="mdiLanguagePython" />
              Python
            </VTab>
            <VTab value="go">
              <VIcon color="#00aed8" class="mr-2" :icon="mdiLanguageGo" />
              Go
            </VTab>
            <VTab value="java">
              <VIcon color="#0c89c7" class="mr-2" :icon="mdiLanguageJava" />
              Java
            </VTab>
            <VTab value="curl">
              <VIcon color="primary" class="mr-2" :icon="mdiPowershell" />
              cURL
            </VTab>
            <VTab value="c-sharp">
              <VIcon color="#68217a" class="mr-2" :icon="mdiLanguageCsharp" />
              C#
            </VTab>
          </VTabs>
          <VTabsWindow v-model="selectedTab" v-highlight>
            <VTabsWindowItem value="javascript">
              <pre
                class="code-window-pre"
              ><code class="language-javascript">import HttpSms from 'httpsms'

const client = new HttpSms('' /* Get the API Key from {{ settingsUrl }} */, '{{ apiBaseUrl }}');

client.messages.postSend({
    content:   'This is a sample text message',
    from:      '+18005550199', // Put the correct phone number here
    to:        '+18005550100', // Put the correct phone number here
})
.then((message) => {
    console.log(message.id); // log the ID of the sent message
})</code></pre>
            </VTabsWindowItem>
            <VTabsWindowItem value="php">
              <pre class="code-window-pre"><code class="language-php">&lt;?php
$apiKey = "Get API Key from {{ settingsUrl }}";

$options = array(
  'http' => array(
    'method'  => 'POST',
    'content' => json_encode( [
        'content' => 'This is a sample text message',
        'from'    => "+18005550199",
        'to'      => "+18005550100"
    ]),
    'header'=>  "Content-Type: application/json\r\n" .
                "Accept: application/json\r\n" .
                "x-api-key: $apiKey\r\n"
    )
);

$context  = stream_context_create( $options );
$result = file_get_contents( "{{ sendMessageUrl }}", false, $context );

echo $result;</code></pre>
            </VTabsWindowItem>
            <VTabsWindowItem value="python">
              <pre
                class="code-window-pre"
              ><code class="language-python">import requests
import json

api_key = "Get API Key from {{ settingsUrl }}"

url = '{{ sendMessageUrl }}'

headers = {
    'x-api-key': api_key,
    'Accept': 'application/json',
    'Content-Type': 'application/json'
}

payload = {
    "content": "This is a sample text message",
    "from": "+18005550199",
    "to": "+18005550100"
}

response = requests.post(url, headers=headers, data=json.dumps(payload))

print(json.dumps(response.json(), indent=4))</code></pre>
            </VTabsWindowItem>
            <VTabsWindowItem value="go">
              <pre
                class="code-window-pre"
              ><code class="language-go">import "github.com/NdoleStudio/httpsms-go"

client := httpsms.New(
    httpsms.WithAPIKey(""/* API Key from {{ settingsUrl }} */),
    httpsms.WithBaseURL("{{ apiBaseUrl }}"),
)

client.Messages.Send(context.Background(), &amp;httpsms.MessageSendParams{
    Content: "This is a sample text message",
    From:    "+18005550199",
    To:      "+18005550100",
})</code></pre>
            </VTabsWindowItem>
            <VTabsWindowItem value="java">
              <pre
                class="code-window-pre"
              ><code class="language-java">var client = HttpClient.newHttpClient();
var apiKey = "Get API Key from {{ settingsUrl }}";

var payload = """
        {
            "content": "This is a sample text message",
            "from": "+18005550199",
            "to": "+18005550100"
        }
        """;

var request = HttpRequest.newBuilder()
        .uri(URI.create("{{ sendMessageUrl }}"))
        .header("accept", "application/json")
        .header("Content-Type", "application/json")
        .header("x-api-key", apiKey)
        .POST(HttpRequest.BodyPublishers.ofString(payload))
        .build();

var response = client.send(request, HttpResponse.BodyHandlers.ofString());
System.out.println(response.body());</code></pre>
            </VTabsWindowItem>
            <VTabsWindowItem value="curl">
              <pre
                class="code-window-pre"
              ><code class="language-bash">curl --location --request POST '{{ sendMessageUrl }}' \
--header 'x-api-key: Get API Key from {{ settingsUrl }}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from": "+18005550199",
    "to": "+18005550100",
    "content": "This is a sample text message"
}'</code></pre>
            </VTabsWindowItem>
            <VTabsWindowItem value="c-sharp">
              <pre
                class="code-window-pre"
              ><code class="language-csharp">var client = new HttpClient();
client.DefaultRequestHeaders.Add("x-api-key", ""/* Get API Key from {{ settingsUrl }} */);

var response = await client.PostAsync(
    "{{ sendMessageUrl }}",
    new StringContent(
        JsonSerializer.Serialize(new {
            from = "+18005550199",
            To = "+18005550100",
            Content = "This is a sample text message",
        }),
        Encoding.UTF8,
        "application/json"
    )
);

Console.WriteLine(await response.Content.ReadAsStringAsync());</code></pre>
            </VTabsWindowItem>
          </VTabsWindow>
        </VThemeProvider>
      </VContainer>
    </section>

    <!-- FAQ Section -->
    <section class="section section-alt">
      <VContainer class="site-container">
        <VRow justify="center">
          <VCol cols="12" md="9" lg="8">
            <div class="section-head">
              <span class="kicker">{{ $t('pages.index.faqKicker') }}</span>
              <h2 class="section-title">{{ $t('pages.index.faqTitle') }}</h2>
              <p
                class="text-body-large text-medium-emphasis mt-3"
                v-html="$t('pages.index.faqDesc')"
              ></p>
            </div>
            <VExpansionPanels v-model="faqPanel">
              <VExpansionPanel v-for="(faq, index) in faqs" :key="faq.q">
                <VExpansionPanelTitle class="text-title-medium py-5">
                  {{ $t(faq.q) }}
                  <template #actions>
                    <VIcon
                      color="primary"
                      :icon="faqPanel === index ? mdiMinus : mdiPlus"
                    />
                  </template>
                </VExpansionPanelTitle>
                <VExpansionPanelText class="text-medium-emphasis">
                  <p v-if="faq.html" v-html="$t(faq.a)"></p>
                  <p v-else>{{ $t(faq.a) }}</p>
                </VExpansionPanelText>
              </VExpansionPanel>
            </VExpansionPanels>
          </VCol>
        </VRow>
      </VContainer>
    </section>

    <!-- Call to action -->
    <section class="section pt-0">
      <VContainer class="site-container">
        <div class="cta-panel">
          <h2 class="cta-panel-title">{{ $t('pages.index.ctaTitle') }}</h2>
          <p class="cta-panel-text">{{ $t('pages.index.ctaDesc') }}</p>
          <div class="d-flex flex-wrap justify-center ga-3">
            <VBtn size="x-large" color="white" variant="flat" to="/login">
              <span class="text-primary">{{
                $t('pages.index.getStartedBtn')
              }}</span>
              <VIcon end color="primary" :icon="mdiArrowRight" />
            </VBtn>
            <VBtn
              size="x-large"
              variant="outlined"
              color="white"
              :href="githubUrl"
              :prepend-icon="mdiGithub"
            >
              {{ $t('pages.index.viewOnGithubBtn') }}
            </VBtn>
          </div>
        </div>
      </VContainer>
    </section>
  </div>
</template>

<style scoped>
.site-container {
  max-width: 1200px;
}

.hero {
  position: relative;
  overflow: hidden;
}

.hero::before,
.hero::after {
  content: '';
  position: absolute;
  z-index: 0;
  border-radius: 50%;
  filter: blur(90px);
  pointer-events: none;
}

.hero::before {
  top: -10rem;
  right: -8rem;
  width: 28rem;
  height: 28rem;
  background: rgb(var(--v-theme-primary));
  opacity: 0.25;
}

.hero::after {
  bottom: -8rem;
  left: -8rem;
  width: 22rem;
  height: 22rem;
  background: rgb(var(--v-theme-secondary));
  opacity: 0.18;
}

.hero > * {
  position: relative;
  z-index: 1;
}

.hero-title {
  margin: 0;
  font-size: clamp(2.3rem, 5.2vw, 3.6rem);
  font-weight: 800;
  line-height: 1.08;
  letter-spacing: -0.03em;
}

.hero-subtitle {
  max-width: 34rem;
  margin: 1.25rem 0 0;
  font-size: 1.15rem;
  line-height: 1.65;
}

.text-md-left .hero-subtitle {
  margin-inline: auto;
}

@media (width >= 960px) {
  .text-md-left .hero-subtitle {
    margin-inline: 0;
  }
}

.hero-checks {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem 1.25rem;
  margin: 1.75rem 0 0;
  padding: 0;
  font-size: 0.9rem;
  font-weight: 600;
  list-style: none;
}

.hero-checks li {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
}

.hero-visual {
  position: relative;
  max-width: 600px;
  margin: 0 auto;
}

.hero-frame {
  padding: 0.5rem;
  border: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  border-radius: 1.5rem;
  background: rgb(var(--v-theme-surface));
  box-shadow: 0 25px 60px rgb(var(--v-theme-primary), 0.25);
}

.hero-frame img {
  display: block;
  width: 100%;
  height: auto;
  border-radius: 1.1rem;
}

.float-chip {
  position: absolute;
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.45rem 0.8rem;
  border: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  border-radius: 0.8rem;
  background: rgb(var(--v-theme-surface));
  box-shadow: var(--card-shadow);
  color: rgb(var(--v-theme-on-surface));
  font-size: 0.8rem;
  font-weight: 700;
  animation: float 6s ease-in-out infinite;
}

.float-chip :deep(.v-icon) {
  color: rgb(var(--v-theme-primary));
}

.float-chip-1 {
  top: -1rem;
  left: -1.5rem;
}

.float-chip-2 {
  top: 45%;
  right: -1.75rem;
  animation-delay: -2s;
}

.float-chip-3 {
  bottom: -1rem;
  left: 10%;
  animation-delay: -4s;
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0);
  }

  50% {
    transform: translateY(-8px);
  }
}

.section {
  padding: 5rem 0;
}

.section-alt {
  background: rgb(var(--v-theme-surface-light), 0.45);
}

.section-head {
  max-width: 44rem;
  margin: 0 auto 2.5rem;
  text-align: center;
}

.section-title {
  margin: 0;
  font-size: clamp(1.8rem, 3.6vw, 2.4rem);
  font-weight: 800;
  line-height: 1.15;
  letter-spacing: -0.02em;
}

.feature-card,
.step-card {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 1.5rem;
}

.feature-card-title {
  margin: 0 0 0.5rem;
  font-size: 1.15rem;
  font-weight: 800;
  letter-spacing: -0.01em;
}

.feature-card-desc {
  flex: 1;
  margin: 0;
  color: rgb(var(--v-theme-on-surface), var(--v-medium-emphasis-opacity));
  font-size: 0.95rem;
  line-height: 1.6;
}

.feature-card-cta {
  align-self: flex-start;
  margin-top: 1rem;
}

.step-card-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  margin-bottom: 1rem;
  border-radius: 0.8rem;
  background: var(--brand-grad);
  color: rgb(var(--v-theme-on-primary));
  font-weight: 800;
}

.code-window {
  overflow: hidden;
  border: 1px solid rgb(255 255 255 / 10%);
  border-radius: var(--card-radius);
  background: #0d1117 !important;
  box-shadow: var(--card-shadow);
}

.code-window-bar {
  display: flex;
  gap: 0.4rem;
  padding: 0.8rem 1rem 0.4rem;
}

.code-window-bar span {
  width: 0.7rem;
  height: 0.7rem;
  border-radius: 50%;
  background: rgb(255 255 255 / 18%);
}

.code-window :deep(.v-tabs) {
  background: transparent;
  border-bottom: 1px solid rgb(255 255 255 / 10%);
}

.code-window :deep(.code-window-pre) {
  margin: 0;
  padding: 1.25rem 1.5rem;
  overflow-x: auto;
  font-size: 0.875rem;
  line-height: 1.6;
}

.code-window :deep(code.hljs) {
  padding: 0;
  background: transparent;
  font-size: 0.875rem;
}

.cta-panel {
  position: relative;
  overflow: hidden;
  padding: 3.5rem 1.5rem;
  border-radius: 1.5rem;

  /* Deep brand shades keep white text above 4.5:1 in every theme. */
  background: var(--brand-deep-grad);
  color: #fff;
  text-align: center;
}

.cta-panel-title {
  margin: 0 0 0.75rem;
  font-size: clamp(1.7rem, 3.4vw, 2.3rem);
  font-weight: 800;
  letter-spacing: -0.02em;
}

.cta-panel-text {
  max-width: 36rem;
  margin: 0 auto 2rem;
  opacity: 0.92;
  font-size: 1.05rem;
}

@media (width <= 600px) {
  .section {
    padding: 3.5rem 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .float-chip {
    animation: none;
  }
}
</style>
