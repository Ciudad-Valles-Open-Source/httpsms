package com.nerus.httpsms

import android.content.ComponentName
import android.content.Context
import android.content.Intent
import android.net.Uri
import android.os.Build
import android.provider.Settings as AndroidSettings
import timber.log.Timber

/**
 * Some manufacturers (Xiaomi, Huawei, Oppo, Vivo, ...) kill background apps even when the standard
 * Android battery optimization is disabled. This helper detects them and opens their
 * "auto-start / background activity" screens so old devices keep receiving messages.
 */
object BatteryGuide {
    private val aggressiveManufacturers = setOf(
        "xiaomi", "redmi", "poco", "huawei", "honor", "oppo", "realme",
        "vivo", "oneplus", "samsung", "asus", "meizu", "letv", "tecno", "infinix", "itel"
    )

    private val vendorIntents = listOf(
        ComponentName("com.miui.securitycenter", "com.miui.permcenter.autostart.AutoStartManagementActivity"),
        ComponentName("com.huawei.systemmanager", "com.huawei.systemmanager.startupmgr.ui.StartupNormalAppListActivity"),
        ComponentName("com.huawei.systemmanager", "com.huawei.systemmanager.optimize.process.ProtectActivity"),
        ComponentName("com.coloros.safecenter", "com.coloros.safecenter.permission.startup.StartupAppListActivity"),
        ComponentName("com.oppo.safe", "com.oppo.safe.permission.startup.StartupAppListActivity"),
        ComponentName("com.vivo.permissionmanager", "com.vivo.permissionmanager.activity.BgStartUpManagerActivity"),
        ComponentName("com.iqoo.secure", "com.iqoo.secure.ui.phoneoptimize.AddWhiteListActivity"),
        ComponentName("com.samsung.android.lool", "com.samsung.android.sm.ui.battery.BatteryActivity"),
        ComponentName("com.oneplus.security", "com.oneplus.security.chainlaunch.view.ChainLaunchAppListActivity"),
        ComponentName("com.asus.mobilemanager", "com.asus.mobilemanager.autostart.AutoStartActivity"),
        ComponentName("com.transsion.phonemaster", "com.cyin.himgr.autostart.AutoStartActivity"),
    )

    fun manufacturerName(): String {
        return Build.MANUFACTURER.replaceFirstChar { it.uppercase() }
    }

    fun isAggressiveManufacturer(): Boolean {
        return aggressiveManufacturers.contains(Build.MANUFACTURER.lowercase())
    }

    /** Opens the vendor-specific auto-start screen, falling back to the app details page. */
    fun openSettings(context: Context) {
        for (component in vendorIntents) {
            try {
                val intent = Intent().setComponent(component).addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                context.startActivity(intent)
                return
            } catch (e: Exception) {
                Timber.d("vendor screen [${component.className}] not available: ${e.javaClass.simpleName}")
            }
        }

        val fallback = Intent(AndroidSettings.ACTION_APPLICATION_DETAILS_SETTINGS)
            .setData(Uri.fromParts("package", context.packageName, null))
            .addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
        context.startActivity(fallback)
    }
}
