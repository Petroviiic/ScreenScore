package com.screenscore.app

import android.app.AppOpsManager
import android.content.Context
import android.os.Process
import com.facebook.react.bridge.*

class ScreenTimeModule(reactContext: ReactApplicationContext) :
    ReactContextBaseJavaModule(reactContext) {

    override fun getName() = "ScreenTimeModule"

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