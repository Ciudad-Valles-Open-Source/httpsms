import java.util.Properties
import java.io.FileInputStream

plugins {
    id("com.android.application")
    id("com.google.gms.google-services")
    id("io.sentry.android.gradle") version "6.14.0"
    id("org.jetbrains.kotlin.plugin.compose")
}

val gitHash = providers.exec {
    commandLine("git", "rev-parse", "--short", "HEAD")
}.standardOutput.asText.map { it.trim() }

// Optional build properties: `-PNAME=value`, gradle.properties / ~/.gradle/gradle.properties,
// or an environment variable with the same name. All default to an empty string.
fun optionalBuildProperty(name: String): String = providers.gradleProperty(name)
    .orElse(providers.environmentVariable(name))
    .getOrElse("")
    .trim()

fun String.asBuildConfigString(): String =
    "\"" + replace("\\", "\\\\").replace("\"", "\\\"") + "\""

// Optional server URL pre-filled on the login screen. Nothing is hard-coded: when it is
// empty the user must type (or scan) the URL of their own httpSMS server.
val defaultServerUrl: String = optionalBuildProperty("HTTPSMS_DEFAULT_SERVER_URL")
    .trimEnd('/')
    .let { if (it.endsWith("/v1", ignoreCase = true)) it.dropLast(3).trimEnd('/') else it }

if (defaultServerUrl.isNotEmpty() &&
    !defaultServerUrl.startsWith("https://", ignoreCase = true) &&
    !defaultServerUrl.startsWith("http://", ignoreCase = true)
) {
    throw GradleException("HTTPSMS_DEFAULT_SERVER_URL must start with http:// or https:// (got [$defaultServerUrl])")
}

// Optional crash reporting and remote debug logs. Both are OFF unless the person building the
// app provides their OWN credentials, so a default build sends nothing to any third party.
//  - HTTPSMS_SENTRY_DSN: Sentry DSN; when empty the Sentry SDK is not initialised at all.
//  - HTTPSMS_LOGS_URL + HTTPSMS_LOGS_TOKEN: endpoint receiving debug logs (JSON array, Bearer
//    token, e.g. an Axiom ingest URL); when either is empty no logs leave the phone.
val sentryDsn: String = optionalBuildProperty("HTTPSMS_SENTRY_DSN")
val sentryEnabled: Boolean = sentryDsn.isNotEmpty()
val logsUrl: String = optionalBuildProperty("HTTPSMS_LOGS_URL")
val logsToken: String = optionalBuildProperty("HTTPSMS_LOGS_TOKEN")

if (logsUrl.isNotEmpty() && !logsUrl.startsWith("https://", ignoreCase = true) &&
    !logsUrl.startsWith("http://", ignoreCase = true)
) {
    throw GradleException("HTTPSMS_LOGS_URL must start with http:// or https:// (got [$logsUrl])")
}

android {
    compileSdk = 37

    defaultConfig {
        applicationId = "com.nerus.httpsms"
        minSdk = 23
        targetSdk = 37
        versionCode = 1
        versionName = gitHash.getOrElse("unknown")
        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
        buildConfigField("String", "DEFAULT_SERVER_URL", defaultServerUrl.asBuildConfigString())
        buildConfigField("String", "LOGS_URL", logsUrl.asBuildConfigString())
        buildConfigField("String", "LOGS_TOKEN", logsToken.asBuildConfigString())
        manifestPlaceholders["sentryDsn"] = sentryDsn
        manifestPlaceholders["sentryEnabled"] = sentryEnabled.toString()
    }

    // Load keystore properties
    val keystorePropertiesFile = rootProject.file("key.properties")
    val keystoreProperties = Properties()
    if (keystorePropertiesFile.exists()) {
        keystoreProperties.load(FileInputStream(keystorePropertiesFile))
    }

    signingConfigs {
        create("release") {
            if (keystorePropertiesFile.exists()) {
                keyAlias = keystoreProperties.getProperty("keyAlias")
                keyPassword = keystoreProperties.getProperty("keyPassword")
                storeFile = keystoreProperties.getProperty("storeFile")?.let { file(it) }
                storePassword = keystoreProperties.getProperty("storePassword")
            }
        }
    }

    buildTypes {
        getByName("debug") {
            manifestPlaceholders["sentryEnvironment"] = "development"
        }
        getByName("release") {
            signingConfig = signingConfigs.getByName("release")
            manifestPlaceholders["sentryEnvironment"] = "production"
            isMinifyEnabled = false
            proguardFiles(getDefaultProguardFile("proguard-android-optimize.txt"), "proguard-rules.pro")
        }
    }
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_1_8
        targetCompatibility = JavaVersion.VERSION_1_8
    }
    namespace = "com.nerus.httpsms"

    buildFeatures {
        buildConfig = true
        compose = true
    }
}

// The Sentry Gradle plugin must not report build telemetry or upload anything on its own.
sentry {
    telemetry.set(false)
    autoUploadProguardMapping.set(false)
    autoUploadNativeSymbols.set(false)
    includeSourceContext.set(false)
}

dependencies {
    val composeBom = platform("androidx.compose:compose-bom:2026.06.01")
    implementation(composeBom)
    androidTestImplementation(composeBom)

    implementation("androidx.compose.ui:ui")
    implementation("androidx.compose.ui:ui-graphics")
    implementation("androidx.compose.ui:ui-tooling-preview")
    implementation("androidx.compose.material3:material3")
    implementation("androidx.compose.material:material-icons-extended")
    implementation("androidx.activity:activity-compose:1.13.0")
    implementation("androidx.lifecycle:lifecycle-viewmodel-compose:2.11.0")

    implementation(platform("com.google.firebase:firebase-bom:34.16.0"))
    implementation("com.journeyapps:zxing-android-embedded:4.3.0")
    implementation("com.google.firebase:firebase-analytics")
    implementation("com.google.firebase:firebase-messaging")
    implementation("com.squareup.okhttp3:okhttp:5.4.0")
    implementation("com.jakewharton.timber:timber:5.0.1")
    implementation("androidx.preference:preference-ktx:1.2.1")
    implementation("androidx.work:work-runtime-ktx:2.11.2")
    implementation("androidx.core:core-ktx:1.19.0")
    implementation("androidx.cardview:cardview:1.0.0")
    implementation("com.beust:klaxon:5.6")
    implementation("androidx.appcompat:appcompat:1.7.1")
    implementation("org.apache.commons:commons-text:1.15.0")
    implementation("com.google.android.material:material:1.14.0")
    implementation("androidx.constraintlayout:constraintlayout:2.2.1")
    implementation("com.googlecode.libphonenumber:libphonenumber:9.0.34")
    implementation("com.klinkerapps:android-smsmms:5.2.6")
    testImplementation("junit:junit:4.13.2")
    androidTestImplementation("androidx.test.ext:junit:1.3.0")
    androidTestImplementation("androidx.test.espresso:espresso-core:3.7.0")
}
