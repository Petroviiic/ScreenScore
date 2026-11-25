package com.screenscore.app

import android.app.AppOpsManager
import android.content.Context
import android.graphics.Bitmap
import android.graphics.Canvas
import android.os.Process
import android.util.Base64
import com.facebook.react.bridge.*
import java.io.ByteArrayOutputStream

class ScreenTimeModule(reactContext: ReactApplicationContext) :
    ReactContextBaseJavaModule(reactContext) {

    override fun getName() = "ScreenTimeModule"

    // packageName -> "data:image/png;base64,..." (null ako aplikacija nema ikonicu)
    private val iconCache = HashMap<String, String?>()

    @ReactMethod
    fun getAppIcon(packageName: String, promise: Promise) {
        if (iconCache.containsKey(packageName)) {
            promise.resolve(iconCache[packageName])
            return
        }
        val uri = try {
            val drawable = reactApplicationContext.packageManager.getApplicationIcon(packageName)
            val size = 96
            val bitmap = Bitmap.createBitmap(size, size, Bitmap.Config.ARGB_8888)
            drawable.setBounds(0, 0, size, size)
            drawable.draw(Canvas(bitmap))
            val out = ByteArrayOutputStream()
            bitmap.compress(Bitmap.CompressFormat.PNG, 100, out)
            "data:image/png;base64," + Base64.encodeToString(out.toByteArray(), Base64.NO_WRAP)
        } catch (e: Exception) {
            null
        }
        iconCache[packageName] = uri
        promise.resolve(uri)
    }

    @ReactMethod
    fun checkPermission(promise: Promise) {
        try {
            val appOps = reactApplicationContext
                .getSystemService(Context.APP_OPS_SERVICE) as AppOpsManager
            val mode = appOps.checkOpNoThrow(
                AppOpsManager.OPSTR_GET_USAGE_STATS,
                Process.myUid(),
                reactApplicationContext.packageName
            )
            promise.resolve(mode == AppOpsManager.MODE_ALLOWED)
        } catch (e: Exception) {
            promise.resolve(false)
        }
    }

    @ReactMethod
    fun getWeeklyStats(promise: Promise) {
        try {
            val stats = ScreenTimeManager.getRealStats(reactApplicationContext)
            if (stats == null) {
                promise.resolve(null)
                return
            }

            val result = Arguments.createMap()

            val weeklyArray = Arguments.createArray()
            stats.weeklyData.forEach { weeklyArray.pushDouble(it.toDouble()) }
            result.putArray("weeklyData", weeklyArray)

            result.putDouble("totalTimeMs", stats.totalTimeMs.toDouble())
            result.putDouble("goalMs", stats.goalMs.toDouble())

            val appsArray = Arguments.createArray()
            stats.topApps.forEach { app ->
                val appMap = Arguments.createMap()
                appMap.putString("packageName", app.packageName)
                appMap.putString("appName", app.appName)
                appMap.putDouble("totalTimeMs", app.totalTimeMs.toDouble())
                appMap.putString("category", app.category.label)
                appsArray.pushMap(appMap)
            }
            result.putArray("topApps", appsArray)

            promise.resolve(result)
        } catch (e: Exception) {
            promise.reject("ERROR", e.message)
        }
    }

    @ReactMethod
    fun getAppsForDay(daysAgo: Int, promise: Promise) {
        try {
            val apps = ScreenTimeManager.getAppsForDay(reactApplicationContext, daysAgo)
            val appsArray = Arguments.createArray()
            apps.forEach { app ->
                val appMap = Arguments.createMap()
                appMap.putString("packageName", app.packageName)
                appMap.putString("appName", app.appName)
                appMap.putDouble("totalTimeMs", app.totalTimeMs.toDouble())
                appMap.putString("category", app.category.label)
                appsArray.pushMap(appMap)
            }
            promise.resolve(appsArray)
        } catch (e: Exception) {
            promise.reject("ERROR", e.message)
        }
    }
}