package com.nerus.httpsms

import android.app.Application
import androidx.appcompat.app.AppCompatDelegate

class HttpSmsApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        AppCompatDelegate.setDefaultNightMode(Settings.getThemeMode(this))
    }
}
