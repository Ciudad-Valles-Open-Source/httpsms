<script setup lang="ts">
import {
  mdiCheckCircle,
  mdiSend,
  mdiGift,
  mdiLightbulbOn60,
  mdiCreation,
  mdiLockOutline,
  mdiLanguagePython,
  mdiCellphoneKey,
  mdiTallyMark1,
  mdiTallyMark3,
  mdiTallyMark2,
  mdiLabel,
  mdiLanguageJavascript,
  mdiLanguagePhp,
  mdiPlus,
  mdiMinus,
  mdiLanguageCsharp,
  mdiLanguageJava,
  mdiMicrosoftExcel,
  mdiWebhook,
  mdiClockOutline,
  mdiArrowRightThin,
  mdiPowershell,
  mdiLanguageGo,
} from '@mdi/js'

definePageMeta({
  layout: 'website',
  middleware: ['redirect-to-threads'],
})

const { t } = useI18n()

useSeoMeta({
  title: () => t('pages.index.metaTitle'),
  description: () => t('pages.index.metaDescription'),
  ogTitle: () => t('pages.index.metaOgTitle'),
  ogDescription: () => t('pages.index.metaOgDescription'),
  ogImage: 'https://httpsms.com/header.png',
  twitterCard: 'summary_large_image',
})

const config = useRuntimeConfig()
const { lgAndUp, mdAndUp, mdAndDown, md, smAndDown, xl } = useVDisplay()

const selectedTab = ref('javascript')
const yearlyPricing = ref(false)
const faqPanel = ref<number | undefined>(undefined)
const pricing = ref(0)

const pricingLabels = ['10K', '20K', '50K', '100K', '200K']
const pricingLabelsFull = ['10,000', '20,000', '50,000', '100,000', '200,000']

const planMessages = computed(() =>
  pricingLabels[pricing.value].replace('K', ',000'),
)
const planMonthlyPrice = computed(() => [20, 35, 89, 175, 350][pricing.value])
const planYearlyPrice = computed(
  () => [200, 350, 1068, 2100, 4200][pricing.value],
)
const planYearlyMonthlyPrice = computed(
  () => [16.66, 29.16, 89, 175, 350][pricing.value],
)
</script>

<template>
  <div>
    <VContainer>
      <VRow :class="{ 'py-4': lgAndUp }">
        <VCol
          cols="12"
          md="6"
          class="pt-8 pb-16"
          :class="{
            'text-center': mdAndDown,
          }"
        >
          <h1
            class="text-display-large font-weight-bold pb-1 gradient-header"
            :class="{
              'mt-16 font-size-45': lgAndUp,
              'mt-10': md,
              'mt-n8': smAndDown,
            }"
          >
            {{ $t('pages.index.heroTitle') }}
          </h1>
          <h2 class="text-medium-emphasis text-headline-small mt-8 mb-8" v-html="$t('pages.index.heroSubtitle')">
          </h2>
          <div :class="{ 'text-center': mdAndDown }">
            <VBtn color="primary" size="large" class="mt-4 mb-4" to="/login">
              <VIcon v-if="lgAndUp" start :icon="mdiSend" />
              {{ $t('pages.index.getStartedBtn') }}
            </VBtn>
            <VBtn
              size="large"
              variant="tonal"
              class="mt-4 mb-4 ml-4"
              href="https://sandbox.httpsms.com"
            >
              <VIcon v-if="lgAndUp" start :icon="mdiCreation" color="#ffe500" />
              {{ $t('pages.index.liveDemoBtn') }}
            </VBtn>
          </div>
          <p class="text-body-medium mt-2" v-html="$t('pages.index.trustedBy')"></p>
          <div class="mt-4" :class="{ 'text-center': mdAndDown }">
            <VIcon color="success" :icon="mdiCheckCircle" />
            {{ $t('pages.index.freeToUse') }}
            <VIcon class="ml-4" color="success" :icon="mdiCheckCircle" />
            {{ $t('pages.index.openSourceBadge') }}
          </div>
          <div v-if="xl" class="mt-4">
            <a href="https://www.uneed.best/tool/httpsmscom">
              <img
                src="https://www.uneed.best/POTD1A.png"
                style="width: 250px"
                alt="Uneed POTD1 Badge"
              />
            </a>
          </div>
          <VDivider
            v-if="mdAndDown"
            class="mt-6 mr-16 bg-success"
            :class="{ 'ml-16': mdAndDown }"
          />
        </VCol>
        <VCol v-if="mdAndUp" cols="12" md="6" class="d-flex align-center">
          <div
            class="mx-auto"
            style="max-width: 98%; width: 100%; aspect-ratio: 16/9"
          >
            <iframe
              src="https://www.youtube-nocookie.com/embed/XTj17RA5txQ?rel=0&modestbranding=1"
              title="httpSMS demo video"
              width="100%"
              height="100%"
              loading="lazy"
              style="border: none; border-radius: 8px"
              allow="
                accelerometer;
                autoplay;
                clipboard-write;
                encrypted-media;
                gyroscope;
                picture-in-picture;
                web-share;
              "
              allowfullscreen
            />
          </div>
        </VCol>
      </VRow>
    </VContainer>

    <!-- Features Section -->
    <VSheet class="py-16">
      <VContainer>
        <!-- Bulk SMS -->
        <VRow class="mb-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="2">
            <div>
              <h3
                class="text-display-medium mb-1"
                :class="{ 'mt-n8': mdAndUp }"
              >
                {{ $t('pages.index.bulkSmsTitle') }}
                <VChip class="ma-2" color="pink" label>
                  <VIcon start :icon="mdiLabel" />
                  {{ $t('pages.index.noCodeBadge') }}
                </VChip>
              </h3>
              <h5 class="text-title-large font-weight-light my-2" v-html="$t('pages.index.bulkSmsDesc')">
              </h5>
              <VBtn
                to="/blog/how-to-send-sms-messages-from-excel"
                color="primary"
              >
                <VIcon start :icon="mdiMicrosoftExcel" />
                {{ $t('pages.index.integrationGuideBtn') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="1">
            <VImg
              class="mb-4"
              max-height="400"
              :src="'/img/bulk-sms-template.png'"
            />
          </VCol>
        </VRow>

        <!-- Integrations -->
        <VRow class="mb-16 mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="1">
            <div>
              <h3 class="text-display-medium mb-1">
                {{ $t('pages.index.integrationsTitle') }}
                <VChip class="ma-2" color="pink" label>
                  <VIcon start :icon="mdiLabel" />
                  {{ $t('pages.index.noCodeBadge') }}
                </VChip>
              </h3>
              <h5 class="text-title-large font-weight-light my-2">
                {{ $t('pages.index.integrationsDesc') }}
              </h5>
              <VBtn
                to="/blog/send-sms-when-new-row-is-added-to-google-sheets-using-zapier"
                color="primary"
              >
                {{ $t('pages.index.zapierIntegrationGuideBtn') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="2">
            <VImg
              class="mb-4"
              :class="{ 'mt-16': mdAndUp }"
              max-height="400"
              :src="'/img/zapier-logo.svg'"
            />
          </VCol>
        </VRow>

        <!-- Webhooks -->
        <VRow class="mb-16 mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="2">
            <div>
              <h3 class="text-display-medium mb-1">{{ $t('pages.index.webhooksTitle') }}</h3>
              <h5 class="text-title-large font-weight-light my-2">
                {{ $t('pages.index.webhooksDesc') }}
              </h5>
              <VBtn
                target="_blank"
                href="https://docs.httpsms.com/webhooks/introduction"
                color="primary"
              >
                <VIcon start :icon="mdiWebhook" />
                {{ $t('pages.index.documentationBtn') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="1">
            <VImg class="mb-4" max-height="300" :src="'/img/connection.svg'" />
          </VCol>
        </VRow>

        <!-- Control Sending -->
        <VRow class="mb-16 mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="1">
            <div>
              <h3 class="text-display-medium mb-1">{{ $t('pages.index.controlSendingTitle') }}</h3>
              <h5 class="text-title-large font-weight-light my-2">
                {{ $t('pages.index.controlSendingDesc') }}
              </h5>
              <VBtn
                href="https://docs.httpsms.com/features/control-sms-send-rate"
                color="primary"
              >
                <VIcon start :icon="mdiArrowRightThin" />
                {{ $t('pages.index.documentationBtn') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="2">
            <VImg class="mb-4" max-height="300" :src="'/img/queue.svg'" />
          </VCol>
        </VRow>

        <!-- Monitoring -->
        <VRow class="mb-16 mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="2">
            <div>
              <h3 class="text-display-medium mb-1">{{ $t('pages.index.monitoringTitle') }}</h3>
              <h5 class="text-title-large font-weight-light my-2">
                {{ $t('pages.index.monitoringDesc') }}
              </h5>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="1">
            <VImg class="mb-4" max-height="300" :src="'/img/alert.svg'" />
          </VCol>
        </VRow>

        <!-- Open Source -->
        <VRow class="mt-16 mb-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="1">
            <div>
              <h3 class="text-display-medium mb-1">{{ $t('pages.index.openSourceTitle') }}</h3>
              <h5 class="text-title-large mb-3 font-weight-light my-2">
                {{ $t('pages.index.openSourceDesc') }}
              </h5>
              <a
                class="text-decoration-none"
                :href="config.public.appGithubUrl"
              >
                <img
                  alt="GitHub Repo stars"
                  height="32"
                  src="https://img.shields.io/github/stars/NdoleStudio/httpsms?style=social"
                />
              </a>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="2">
            <VImg
              class="mb-4"
              max-height="400"
              :src="'/img/httpsms-github.png'"
            />
          </VCol>
        </VRow>

        <!-- Encryption -->
        <VRow class="mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="2">
            <div>
              <h3 class="text-display-medium mb-1">{{ $t('pages.index.encryptionTitle') }}</h3>
              <h5 class="text-title-large mb-3 font-weight-light my-2" v-html="$t('pages.index.encryptionDesc')">
              </h5>
              <VBtn
                to="/blog/end-to-end-encryption-to-sms-messages"
                color="primary"
              >
                <VIcon start :icon="mdiLockOutline" />
                {{ $t('pages.index.setupEncryptionBtn') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="1">
            <VImg
              class="mb-4"
              max-height="300"
              :src="'/img/mobile-encryption.svg'"
            />
          </VCol>
        </VRow>

        <!-- Multiple Phones -->
        <VRow class="mt-16">
          <VCol cols="12" md="6" class="d-flex align-center">
            <div>
              <h3 class="text-display-medium mb-1">{{ $t('pages.index.multiplePhonesTitle') }}</h3>
              <h5 class="text-title-large mb-3 font-weight-light my-2">
                {{ $t('pages.index.multiplePhonesDesc') }}
              </h5>
              <VBtn
                href="https://docs.httpsms.com/features/phone-api-keys"
                color="primary"
              >
                <VIcon start :icon="mdiCellphoneKey" />
                {{ $t('pages.index.documentationBtn') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6">
            <VImg
              class="mb-4"
              max-height="300"
              :src="'/img/manage-phones.svg'"
            />
          </VCol>
        </VRow>

        <!-- Schedule Messages -->
        <VRow class="mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="2">
            <div>
              <h3 class="text-display-medium mb-1">{{ $t('pages.index.scheduleTextMessagesTitle') }}</h3>
              <h5 class="text-headline-small my-2 font-weight-light">
                {{ $t('pages.index.scheduleTextMessagesDesc') }}
              </h5>
              <VBtn
                href="https://docs.httpsms.com/features/scheduling-sms-messages"
                color="primary"
              >
                <VIcon start :icon="mdiClockOutline" />
                {{ $t('pages.index.documentationBtn') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="1">
            <VImg
              class="mb-4"
              max-height="300"
              :src="'/img/schedule-messages.svg'"
            />
          </VCol>
        </VRow>
      </VContainer>
    </VSheet>

    <!-- Get Started Section -->
    <VContainer class="pb-16">
      <VRow>
        <VCol>
          <h2 class="text-display-large text-center mb-0">{{ $t('pages.index.getStartedSectionTitle') }}</h2>
        </VCol>
      </VRow>
      <VRow>
        <VCol cols="12">
          <VRow class="align-baseline">
            <VCol cols="12" md="5" class="pr-4">
              <VTimeline
                truncate-line="both"
                density="compact"
                class="mt-10 ml-n4"
              >
                <VTimelineItem dot-color="primary" :icon="mdiTallyMark1">
                  <VCard variant="elevated">
                    <VCardTitle class="text-headline-medium">{{ $t('pages.index.step1Title') }}</VCardTitle>
                    <VCardText class="text-body-large" v-html="$t('pages.index.step1Desc')">
                    </VCardText>
                  </VCard>
                </VTimelineItem>
                <VTimelineItem dot-color="primary" :icon="mdiTallyMark2">
                  <VCard variant="elevated">
                    <VCardTitle class="text-headline-medium">{{ $t('pages.index.step2Title') }}</VCardTitle>
                    <VCardText class="text-body-large" v-html="$t('pages.index.step2Desc').replace('{url}', config.public.appDownloadUrl)">
                    </VCardText>
                  </VCard>
                </VTimelineItem>
                <VTimelineItem dot-color="primary" :icon="mdiTallyMark3">
                  <VCard variant="elevated">
                    <VCardTitle class="text-headline-medium">{{ $t('pages.index.step3Title') }}</VCardTitle>
                    <VCardText class="text-body-large">
                      {{ $t('pages.index.step3Desc') }}
                      <a
                        class="text-decoration-none"
                        :href="config.public.appDocumentationUrl"
                      >
                        {{ config.public.appDocumentationUrl }}
                      </a>
                    </VCardText>
                  </VCard>
                </VTimelineItem>
              </VTimeline>
            </VCol>
            <VCol cols="12" md="7">
              <div class="w-100" :class="{ 'mt-n8': mdAndUp }">
                <VTabs
                  v-model="selectedTab"
                  color="primary"
                  bg-color="#212121"
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
                    <VIcon
                      color="#777bb3"
                      class="mr-2"
                      :icon="mdiLanguagePhp"
                    />
                    PHP
                  </VTab>
                  <VTab value="python">
                    <VIcon
                      color="#ffffff"
                      class="mr-2"
                      :icon="mdiLanguagePython"
                    />
                    Python
                  </VTab>
                  <VTab value="go">
                    <VIcon color="#00aed8" class="mr-2" :icon="mdiLanguageGo" />
                    Go
                  </VTab>
                  <VTab value="java">
                    <VIcon
                      color="#0c89c7"
                      class="mr-2"
                      :icon="mdiLanguageJava"
                    />
                    Java
                  </VTab>
                  <VTab value="curl">
                    <VIcon color="primary" class="mr-2" :icon="mdiPowershell" />
                    cURL
                  </VTab>
                  <VTab value="c-sharp">
                    <VIcon
                      color="#68217a"
                      class="mr-2"
                      :icon="mdiLanguageCsharp"
                    />
                    C#
                  </VTab>
                </VTabs>
                <VTabsWindow v-model="selectedTab" v-highlight>
                  <VTabsWindowItem value="javascript">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-javascript">import HttpSms from 'httpsms'

const client = new HttpSms('' /* Get the API Key from https://httpsms.com/settings */);

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
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-php">&lt;?php
$apiKey = "Get API Key from https://httpsms.com/settings";

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
$result = file_get_contents( "https://api.httpsms.com/v1/messages/send", false, $context );

echo $result;</code></pre>
                  </VTabsWindowItem>
                  <VTabsWindowItem value="python">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-python">import requests
import json

api_key = "Get API Key from https://httpsms.com/settings"

url = 'https://api.httpsms.com/v1/messages/send'

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
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-go">import "github.com/NdoleStudio/httpsms-go"

client := htpsms.New(htpsms.WithAPIKey(/* API Key from https://httpsms.com/settings */))

client.Messages.Send(context.Background(), &amp;httpsms.MessageSendParams{
    Content: "This is a sample text message",
    From:    "+18005550199",
    To:      "+18005550100",
})</code></pre>
                  </VTabsWindowItem>
                  <VTabsWindowItem value="java">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-java">var client = HttpClient.newHttpClient();
var apiKey = "Get API Key from https://httpsms.com/settings";

var payload = """
        {
            "content": "This is a sample text message",
            "from": "+18005550199",
            "to": "+18005550100"
        }
        """;

var request = HttpRequest.newBuilder()
        .uri(URI.create("https://api.httpsms.com/v1/messages/send"))
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
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-bash">curl --location --request POST 'https://api.httpsms.com/v1/messages/send' \
--header 'x-api-key: Get API Key from https://httpsms.com/settings' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from": "+18005550199",
    "to": "+18005550100",
    "content": "This is a sample text message"
}'</code></pre>
                  </VTabsWindowItem>
                  <VTabsWindowItem value="c-sharp">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-csharp">var client = new HttpClient();
client.DefaultRequestHeaders.Add("x-api-key", ""/* Get API Key from https://httpsms.com/settings */);

var response = await client.PostAsync(
    "https://api.httpsms.com/v1/messages/send",
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
              </div>
            </VCol>
          </VRow>
        </VCol>
      </VRow>
    </VContainer>

    <!-- Pricing Section -->
    <template v-if="config.public.enableBilling">
    <VSheet class="mt-16 pb-16">
      <VContainer>
        <VRow>
          <VCol md="6" offset-md="3">
            <h2
              id="pricing"
              style="text-decoration-color: #329ef4"
              class="text-center text-display-large mb-4 text-decoration-underline dark:text-white"
            >
              {{ $t('pages.index.pricingTitle') }}
            </h2>
            <h4 class="text-center text-headline-small text-medium-emphasis" v-html="$t('pages.index.pricingDesc')">
            </h4>
            <div class="d-flex justify-center mt-4 align-center">
              <p
                class="text-headline-small mr-3 mt-3"
                :class="{ 'text-medium-emphasis': yearlyPricing }"
              >
                {{ $t('pages.index.monthlyLabel') }}
              </p>
              <VSwitch
                v-model="yearlyPricing"
                color="primary"
                class="mt-n2"
                hide-details
              />
              <p
                class="text-headline-small ml-3 mt-3"
                :class="{ 'text-medium-emphasis': !yearlyPricing }"
              >
                {{ $t('pages.index.yearlyLabel') }}
                <VChip color="primary" size="small">
                  <VIcon start :icon="mdiGift" size="small" />
                  {{ $t('pages.index.twoMonthsFreeBadge') }}
                </VChip>
              </p>
            </div>
          </VCol>
        </VRow>
        <VRow>
          <VCol cols="12">
            <VSlider
              v-model="pricing"
              :tick-labels="lgAndUp ? pricingLabelsFull : pricingLabels"
              :max="4"
              step="1"
              color="primary"
              thumb-color="primary"
              thumb-label="always"
              thumb-size="16"
              tick-size="8"
              show-ticks="always"
            >
              <template #thumb-label>
                {{ pricingLabels[pricing] }}
              </template>
            </VSlider>
          </VCol>
        </VRow>
        <VRow>
          <!-- Free Plan -->
          <VCol cols="12" lg="4">
            <VCard elevation="4" color="#121212">
              <VCardText>
                <h1 class="text-center text-display-medium mt-0 mb-4">{{ $t('pages.index.freePlanTitle') }}</h1>
                <p
                  class="text-body-large text-center mt-0 text-medium-emphasis"
                >
                  {{ $t('pages.index.freePlanDesc') }}
                </p>
                <p class="text-center">
                  <span class="text-display-small">{{ $t('pages.index.freePlanCost') }}</span>
                </p>
                <p class="text-center mt-n3 text-medium-emphasis">
                  {{ $t('pages.index.noCreditCard') }}
                </p>
                <VBtn block to="/login" variant="tonal" size="large"
                  >{{ $t('pages.index.getStartedBtn') }}</VBtn
                >
                <p class="mt-6 text-md-body-large text-title-medium">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  Send or receive up to <b>200</b> SMS/month
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  Offline notifications for your phone
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  Forward received messages via webhook
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  Basic email support
                </p>
              </VCardText>
            </VCard>
          </VCol>
          <!-- Pro Plan -->
          <VCol cols="12" lg="4">
            <VCard elevation="4" color="#000000">
              <VCardText>
                <h1
                  class="text-center text-display-medium mt-0 mb-4 text-primary"
                >
                  {{ $t('pages.index.proPlanTitle') }}
                </h1>
                <p
                  class="text-body-large text-center mt-0 text-medium-emphasis"
                >
                  {{ $t('pages.index.proPlanDesc') }}
                </p>
                <p v-if="!yearlyPricing" class="text-center">
                  <span class="text-display-small">$10</span>/month
                </p>
                <p v-else class="text-center">
                  <span class="text-display-small">$100</span>/year
                </p>
                <p
                  v-if="!yearlyPricing"
                  class="text-center mt-n3 text-medium-emphasis"
                  v-html="$t('pages.index.orPerYear').replace('{price}', '100')"
                >
                </p>
                <p v-else class="text-center mt-n3 text-medium-emphasis" v-html="$t('pages.index.orPerMonth').replace('{price}', '8.33')">
                </p>
                <VBtn block color="primary" to="/login" size="large"
                  >{{ $t('pages.index.tryForFreeBtn') }}</VBtn
                >
                <p class="mt-6 text-md-body-large text-title-medium">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  <span v-html="$t('pages.index.proPlanFeature1')"></span>
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('pages.index.offlineNotifications') }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('pages.index.forwardMessagesWebhook') }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('pages.index.prioritySupport') }}
                </p>
              </VCardText>
            </VCard>
          </VCol>
          <!-- Custom Plan -->
          <VCol cols="12" lg="4">
            <VCard elevation="4" color="#121212">
              <VCardText>
                <h1 class="text-center text-display-medium mt-0 mb-4">
                  {{ $t('pages.index.customPlanTitle').replace('{plan}', pricingLabels[pricing]) }}
                </h1>
                <p
                  class="text-body-large text-center mt-0 text-medium-emphasis"
                >
                  {{ $t('pages.index.customPlanDesc').replace('{messages}', planMessages) }}
                </p>
                <p v-if="!yearlyPricing" class="text-center">
                  <span class="text-display-small">${{ planMonthlyPrice }}</span
                  >/month
                </p>
                <p v-else class="text-center">
                  <span class="text-display-small">${{ planYearlyPrice }}</span
                  >/year
                </p>
                <p
                  v-if="!yearlyPricing"
                  class="text-center mt-n3 text-medium-emphasis"
                  v-html="$t('pages.index.orPerYear').replace('{price}', planYearlyPrice)"
                >
                </p>
                <p v-else class="text-center mt-n3 text-medium-emphasis" v-html="$t('pages.index.orPerMonth').replace('{price}', planYearlyMonthlyPrice)">
                </p>
                <VBtn block variant="tonal" to="/login" size="large"
                  >{{ $t('pages.index.tryForFreeBtn') }}</VBtn
                >
                <p class="mt-6 text-md-body-large text-title-medium">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  <span v-html="$t('pages.index.customPlanFeature1').replace('{messages}', pricingLabels[pricing])"></span>
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('pages.index.offlineNotifications') }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('pages.index.forwardMessagesWebhook') }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('pages.index.prioritySupport') }}
                </p>
              </VCardText>
            </VCard>
          </VCol>
        </VRow>
        <VRow>
          <VCol>
            <VAlert
              color="info"
              :icon="mdAndUp ? mdiLightbulbOn60 : undefined"
              :prominent="mdAndUp"
              variant="tonal"
            >
              <span v-html="$t('pages.index.contactSupportAlert')"></span>
            </VAlert>
          </VCol>
        </VRow>
      </VContainer>
    </VSheet>
    </template>

    <!-- Testimonials Section -->
    <VContainer class="mt-16">
      <VRow>
        <VCol cols="12" md="6">
          <VCard
            href="https://www.g2.com/products/httpsms/reviews/httpsms-review-8589834"
          >
            <VCardText class="pt-0 pb-0">
              <div class="d-flex">
                <VAvatar class="mt-6">
                  <VImg
                    alt="Joysankar M."
                    src="https://images.g2crowd.com/uploads/avatar/image/1662077/thumb_square_d5706804d1b343744a8feb693827fe34.jpeg"
                  />
                </VAvatar>
                <div>
                  <p class="text-title-medium ml-3">Joysankar M.</p>
                  <VRating
                    class="mt-n7"
                    color="yellow-darken-3"
                    :model-value="4.5"
                    half-increments
                    readonly
                  />
                </div>
                <VSpacer />
                <div style="width: 30px" class="mt-4">
                  <VImg
                    max-height="30"
                    src="https://company.g2.com/hs-fs/hubfs/brand-guide/reversed-g2@2x.png"
                  />
                </div>
              </div>
              <p class="text-title-large font-weight-light mt-0" v-html="$t('pages.index.testimonial1')">
              </p>
            </VCardText>
          </VCard>
        </VCol>
        <VCol cols="12" md="6">
          <VCard href="https://www.uneed.best/tool/httpsmscom?tab=comments">
            <VCardText class="pb-0">
              <div class="d-flex">
                <VAvatar class="mt-2">
                  <VImg
                    alt="Edmund Ciego Profile Picture"
                    src="https://lh3.googleusercontent.com/a/ACg8ocJktUViyMcJvzkPNpza7SZ3ql_nwOAzYk0uJ27TF5L_z0bRoPKE=s96-c"
                  />
                </VAvatar>
                <div>
                  <p class="text-title-medium mt-0 ml-3">Edmund Ciego</p>
                  <VRating
                    class="mt-n7"
                    color="yellow-darken-3"
                    :model-value="5"
                    half-increments
                    readonly
                  />
                </div>
                <VSpacer />
                <div>
                  <v-img width="64" src="/img/logos/uneed.svg" />
                </div>
              </div>
              <p class="text-title-large font-weight-light mt-0" v-html="$t('pages.index.testimonial2')">
              </p>
            </VCardText>
          </VCard>
        </VCol>
      </VRow>
    </VContainer>

    <!-- FAQ Section -->
    <VContainer class="pb-16">
      <VRow>
        <VCol md="8" offset-md="2">
          <h2
            class="text-md-display-large mb-4 text-center text-display-medium"
          >
            {{ $t('pages.index.faqTitle') }}
          </h2>
          <p class="text-center text-title-large mt-4 text-medium-emphasis" v-html="$t('pages.index.faqDesc')">
          </p>
        </VCol>
      </VRow>
      <VRow>
        <VCol md="8" offset-md="2" class="mb-16">
          <VExpansionPanels v-model="faqPanel">
            <VExpansionPanel>
              <VExpansionPanelTitle
                class="text-title-large text-md-headline-small"
              >
                {{ $t('pages.index.faq1Question') }}
                <template #actions>
                  <VIcon :icon="faqPanel === 0 ? mdiMinus : mdiPlus" />
                </template>
              </VExpansionPanelTitle>
              <VExpansionPanelText>
                <p class="mt-4">
                  {{ $t('pages.index.faq1Answer') }}
                </p>
              </VExpansionPanelText>
            </VExpansionPanel>
            <VExpansionPanel>
              <VExpansionPanelTitle
                class="text-title-large text-md-headline-small"
              >
                {{ $t('pages.index.faq2Question') }}
                <template #actions>
                  <VIcon :icon="faqPanel === 1 ? mdiMinus : mdiPlus" />
                </template>
              </VExpansionPanelTitle>
              <VExpansionPanelText>
                <p class="mt-4">
                  {{ $t('pages.index.faq2Answer') }}
                </p>
              </VExpansionPanelText>
            </VExpansionPanel>
            <VExpansionPanel>
              <VExpansionPanelTitle
                class="text-title-large text-md-headline-small"
              >
                {{ $t('pages.index.faq3Question') }}
                <template #actions>
                  <VIcon :icon="faqPanel === 2 ? mdiMinus : mdiPlus" />
                </template>
              </VExpansionPanelTitle>
              <VExpansionPanelText>
                <p class="mt-4" v-html="$t('pages.index.faq3Answer')">
                </p>
              </VExpansionPanelText>
            </VExpansionPanel>
            <VExpansionPanel>
              <VExpansionPanelTitle
                class="text-title-large text-md-headline-small"
              >
                {{ $t('pages.index.faq4Question') }}
                <template #actions>
                  <VIcon :icon="faqPanel === 3 ? mdiMinus : mdiPlus" />
                </template>
              </VExpansionPanelTitle>
              <VExpansionPanelText>
                <p class="mt-4">
                  {{ $t('pages.index.faq4Answer') }}
                </p>
              </VExpansionPanelText>
            </VExpansionPanel>
          </VExpansionPanels>
        </VCol>
      </VRow>
    </VContainer>
  </div>
</template>

<style lang="scss">
.gradient-header {
  color: #1ad37f;
  background-image: -webkit-linear-gradient(0deg, #1ad37f 14%, #329ef4 55%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.font-size-45 {
  font-size: 4.5rem;
}

.gradient-underline {
  color: white;
}
</style>
