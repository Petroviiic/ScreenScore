package com.screenscore.app

import android.app.usage.UsageEvents
import android.app.usage.UsageStatsManager
import android.content.Context
import android.content.pm.PackageManager
import android.graphics.drawable.Drawable
import java.util.Calendar
import java.util.concurrent.TimeUnit

data class AppUsageData(
    val packageName: String,
    val appName: String,
    val icon: Drawable?,
    val totalTimeMs: Long,
    val category: AppCategory
) {
    val totalTimeFormatted: String get() {
        val hours = TimeUnit.MILLISECONDS.toHours(totalTimeMs)
        val minutes = TimeUnit.MILLISECONDS.toMinutes(totalTimeMs) % 60
        return when {
            hours > 0 -> "${hours}h ${minutes}m"
            minutes > 0 -> "${minutes}m"
            else -> "<1m"
        }
    }
}

enum class AppCategory(val label: String) {
    PRODUCTIVE("Productive"),
    SOCIAL("Social"),
    ENTERTAINMENT("Entertainment"),
    COMMUNICATION("Communication"),
    HEALTH("Health & Fitness"),
    OTHER("Other")
}

data class DailyStats(
    val totalTimeMs: Long,
    val goalMs: Long,
    val topApps: List<AppUsageData>,
    val weeklyData: List<Long>
) {
    val goalPercent: Float get() = (totalTimeMs.toFloat() / goalMs).coerceIn(0f, 1f)
    val totalTimeFormatted: String get() {
        val hours = TimeUnit.MILLISECONDS.toHours(totalTimeMs)
        val minutes = TimeUnit.MILLISECONDS.toMinutes(totalTimeMs) % 60
        return "${hours}h ${minutes}m"
    }
    val goalTimeFormatted: String get() = "${TimeUnit.MILLISECONDS.toHours(goalMs)}h"
    val remainingFormatted: String get() {
        val remaining = (goalMs - totalTimeMs).coerceAtLeast(0)
        val hours = TimeUnit.MILLISECONDS.toHours(remaining)
        val minutes = TimeUnit.MILLISECONDS.toMinutes(remaining) % 60
        return if (hours > 0) "${hours}h ${minutes}m left" else "${minutes}m left"
    }
}

object ScreenTimeManager {

    @Volatile var cachedStats: DailyStats? = null
    @Volatile var cacheTime: Long = 0

    fun prefetchStats(context: Context) {
        Thread {
            cachedStats = getRealStats(context)
            cacheTime = System.currentTimeMillis()
        }.start()
    }

    fun saveStatsToPrefs(context: Context) {
        val stats = cachedStats ?: return
        val today = java.text.SimpleDateFormat("yyyy-MM-dd", java.util.Locale.getDefault()).format(java.util.Date())
        context.getSharedPreferences("screenscore_cache", Context.MODE_PRIVATE).edit()
            .putLong("total_time_ms", stats.totalTimeMs)
            .putString("weekly_data", stats.weeklyData.joinToString(","))
            .putString("stats_date", today)
            .apply()
    }

    fun loadStatsFromPrefs(context: Context): DailyStats? {
        val prefs = context.getSharedPreferences("screenscore_cache", Context.MODE_PRIVATE)
        val today = java.text.SimpleDateFormat("yyyy-MM-dd", java.util.Locale.getDefault()).format(java.util.Date())
        val savedDate = prefs.getString("stats_date", null)
        if (savedDate != today) return null
        val totalTimeMs = prefs.getLong("total_time_ms", -1L)
        if (totalTimeMs < 0) return null
        val weekly = prefs.getString("weekly_data", null)
            ?.split(",")?.mapNotNull { it.toLongOrNull() }
            ?: List(7) { 0L }
        val dummy = DailyStats(totalTimeMs, TimeUnit.HOURS.toMillis(4), emptyList(), weekly)
        cachedStats = dummy
        cacheTime = System.currentTimeMillis()
        return dummy
    }

    var demoMode = false

    private val CATEGORY_MAP = mapOf(
        "com.instagram.android"         to AppCategory.SOCIAL,
        "com.facebook.katana"           to AppCategory.SOCIAL,
        "com.twitter.android"           to AppCategory.SOCIAL,
        "com.zhiliaoapp.musically"      to AppCategory.ENTERTAINMENT,
        "com.ss.android.ugc.trill"      to AppCategory.ENTERTAINMENT,
        "com.snapchat.android"          to AppCategory.SOCIAL,
        "com.reddit.frontpage"          to AppCategory.SOCIAL,
        "com.google.android.youtube"    to AppCategory.ENTERTAINMENT,
        "com.spotify.music"             to AppCategory.ENTERTAINMENT,
        "com.netflix.mediaclient"       to AppCategory.ENTERTAINMENT,
        "com.whatsapp"                  to AppCategory.COMMUNICATION,
        "com.google.android.gm"         to AppCategory.COMMUNICATION,
        "org.telegram.messenger"        to AppCategory.COMMUNICATION,
        "com.discord"                   to AppCategory.COMMUNICATION,
        "com.duolingo"                  to AppCategory.PRODUCTIVE,
        "com.google.android.apps.docs"  to AppCategory.PRODUCTIVE,
        "com.microsoft.office.word"     to AppCategory.PRODUCTIVE,
        "com.google.android.calendar"   to AppCategory.PRODUCTIVE,
        "com.nike.plusgps"              to AppCategory.HEALTH,
        "com.strava"                    to AppCategory.HEALTH,
        "com.brave.browser"             to AppCategory.OTHER,
        "org.mozilla.firefox"           to AppCategory.OTHER,
        "com.android.chrome"            to AppCategory.OTHER,
    )

    val KNOWN_APP_NAMES_PUBLIC = mapOf(
        "com.zhiliaoapp.musically"   to "TikTok",
        "com.ss.android.ugc.trill"   to "TikTok",
        "com.instagram.android"      to "Instagram",
        "com.facebook.katana"        to "Facebook",
        "com.snapchat.android"       to "Snapchat",
        "com.twitter.android"        to "Twitter/X",
        "com.google.android.youtube" to "YouTube",
        "com.spotify.music"          to "Spotify",
        "com.whatsapp"               to "WhatsApp",
        "org.telegram.messenger"     to "Telegram",
        "com.discord"                to "Discord",
        "com.netflix.mediaclient"    to "Netflix",
        "com.brave.browser"          to "Brave",
        "com.android.chrome"         to "Chrome",
        "org.mozilla.firefox"        to "Firefox",
        "com.reddit.frontpage"       to "Reddit",
        "com.duolingo"               to "Duolingo",
        "com.google.android.gm"      to "Gmail",
        "com.mi.android.globalFileexplorer" to "File Manager",
        "com.android.fileexplorer"   to "File Manager",
        "com.android.settings"       to "Settings",
        "com.google.android.gms"     to "Google Play Services",
        "com.android.vending"        to "Play Store",
        "com.coloros.filemanager"    to "File Manager",
        "com.asus.filemanager"       to "File Manager",
        "com.sec.android.app.myfiles" to "My Files",
        "com.google.android.googlequicksearchbox" to "Google",
        "com.miui.securitycenter"    to "Security",
        "com.miui.systemui"          to "MIUI System",
        "com.miui.notification"      to "Notifications",
        "com.android.packageinstaller" to "Package Installer",
        "com.google.android.packageinstaller" to "Package Installer",
        "com.miui.packageinstaller"  to "Package Installer",
        "com.android.permissioncontroller" to "Permission Controller",
        "com.google.android.permissioncontroller" to "Permission Controller",
        "com.xiaomi.misettings"      to "Settings",
    )

    private val EXCLUDE_PACKAGES = setOf(
        "android",
        "com.android.systemui",
        "com.miui.home", "com.android.launcher", "com.android.launcher3",
        "com.google.android.apps.nexuslauncher", "com.sec.android.app.launcher",
        "com.huawei.android.launcher", "com.oneplus.launcher",
        "com.xiaomi.xmsf", "com.xiaomi.account",
        "com.miui.daemon", "com.miui.systemAdSolution",
        "com.miui.contentcatcher", "com.miui.mishare.connectivity",
        "le.android.apps.wellbeing", "com.google.android.apps.wellbeing",
        "com.samsung.android.digitalwellbeing",
        "com.android.packageinstaller", "com.google.android.packageinstaller",
        "com.miui.packageinstaller", "i.global.packageinstaller",
        "bal.packageinstaller", "com.samsung.android.packageinstaller",
        "com.android.vending.billing.InAppBillingService.COIN",
        "com.miui.global.packageinstaller", "com.lbe.security.miui",
        "com.xiaomi.packageinstaller", "packageinstaller",
        "com.android.permissioncontroller", "com.google.android.permissioncontroller",
        "com.google.android.inputmethod.latin",
        "com.samsung.android.honeyboard",
    )

    fun midnightDaysAgo(daysAgo: Int): Long {
        val cal = Calendar.getInstance()
        cal.add(Calendar.DAY_OF_YEAR, -daysAgo)
        cal.set(Calendar.HOUR_OF_DAY, 0)
        cal.set(Calendar.MINUTE, 0)
        cal.set(Calendar.SECOND, 0)
        cal.set(Calendar.MILLISECOND, 0)
        return cal.timeInMillis
    }

    fun calcUsageFromEvents(context: Context, startMs: Long, endMs: Long): Map<String, Long> {
        val usm = context.getSystemService(Context.USAGE_STATS_SERVICE) as UsageStatsManager
        val result = mutableMapOf<String, Long>()

        val isToday = endMs > System.currentTimeMillis() - TimeUnit.MINUTES.toMillis(1)
        val queryEnd = if (isToday) endMs else endMs + TimeUnit.MINUTES.toMillis(30)

        val events = usm.queryEvents(startMs, queryEnd) ?: return result

        val foregroundStart = mutableMapOf<String, Long>()
        val event = UsageEvents.Event()

        while (events.hasNextEvent()) {
            events.getNextEvent(event)
            val pkg = event.packageName
            if (pkg in EXCLUDE_PACKAGES) continue
            if (pkg.contains("packageinstaller", ignoreCase = true)) continue
            if (pkg.contains("wellbeing", ignoreCase = true)) continue
            if (pkg.contains("permissioncontroller", ignoreCase = true)) continue
            if (pkg.contains("faceservice", ignoreCase = true)) continue
            if (pkg.contains("inputmethod", ignoreCase = true)) continue

            when (event.eventType) {
                UsageEvents.Event.MOVE_TO_FOREGROUND -> {
                    foregroundStart[pkg] = event.timeStamp
                }
                UsageEvents.Event.MOVE_TO_BACKGROUND -> {
                    val start = foregroundStart.remove(pkg)
                    if (start != null) {
                        val clampedEnd = event.timeStamp.coerceAtMost(endMs)
                        val dur = clampedEnd - start
                        if (dur > 0) result[pkg] = (result[pkg] ?: 0L) + dur
                    } else {
                        val clampedEnd = event.timeStamp.coerceAtMost(endMs)
                        if (clampedEnd - startMs <= TimeUnit.MINUTES.toMillis(30)) {
                            val dur = clampedEnd - startMs
                            if (dur > 0) result[pkg] = (result[pkg] ?: 0L) + dur
                        }
                    }
                }
            }
        }

        foregroundStart.forEach { (pkg, start) ->
            val dur = endMs - start
            if (dur > 0) result[pkg] = (result[pkg] ?: 0L) + dur
        }

        return result.filter { it.value >= 5_000 }
    }

    private fun buildAppList(context: Context, usage: Map<String, Long>): List<AppUsageData> {
        val pm = context.packageManager
        return usage.entries
            .filter { it.value >= 5_000 }
            .sortedByDescending { it.value }
            .take(30)
            .map { (packageName, timeMs) ->
                val appName = try {
                    pm.getApplicationLabel(pm.getApplicationInfo(packageName, 0)).toString()
                } catch (e: PackageManager.NameNotFoundException) {
                    KNOWN_APP_NAMES_PUBLIC[packageName]
                        ?: packageName.substringAfterLast(".").replaceFirstChar { it.uppercase() }
                }
                val icon = try { pm.getApplicationIcon(packageName) } catch (e: Exception) { null }
                val category = CATEGORY_MAP[packageName] ?: AppCategory.OTHER
                AppUsageData(packageName, appName, icon, timeMs, category)
            }
    }

    fun getRealStats(context: Context): DailyStats? {
        val now = System.currentTimeMillis()
        val todayStart = midnightDaysAgo(0)

        val todayUsage = calcUsageFromEvents(context, todayStart, now)
        if (todayUsage.isEmpty()) return null

        val totalToday = todayUsage.values.sum().coerceAtMost(TimeUnit.HOURS.toMillis(16))
        val appDataList = buildAppList(context, todayUsage)

        val weeklyData = (6 downTo 0).map { daysAgo ->
            if (daysAgo == 0) totalToday
            else {
                val dayStart = midnightDaysAgo(daysAgo)
                val dayEnd = midnightDaysAgo(daysAgo - 1)
                calcUsageFromEvents(context, dayStart, dayEnd)
                    .values.sum().coerceAtMost(TimeUnit.HOURS.toMillis(16))
            }
        }

        return DailyStats(totalToday, TimeUnit.HOURS.toMillis(4), appDataList, weeklyData)
    }

    fun getAppsForDay(context: Context, daysAgo: Int): List<AppUsageData> {
        val startMs = midnightDaysAgo(daysAgo)
        val endMs = if (daysAgo == 0) System.currentTimeMillis() else midnightDaysAgo(daysAgo - 1)
        val usage = calcUsageFromEvents(context, startMs, endMs)
        return buildAppList(context, usage)
    }

    fun getDemoStats(): DailyStats {
        val mockApps = listOf(
            Triple("Instagram",   AppCategory.SOCIAL,        107 * 60_000L),
            Triple("Brave",       AppCategory.OTHER,          96 * 60_000L),
            Triple("WhatsApp",    AppCategory.COMMUNICATION,  74 * 60_000L),
            Triple("YouTube",     AppCategory.ENTERTAINMENT,  58 * 60_000L),
            Triple("TikTok",      AppCategory.ENTERTAINMENT,  34 * 60_000L),
            Triple("Gmail",       AppCategory.COMMUNICATION,  12 * 60_000L),
            Triple("Spotify",     AppCategory.ENTERTAINMENT,  11 * 60_000L),
            Triple("Google Maps", AppCategory.OTHER,           8 * 60_000L),
            Triple("Duolingo",    AppCategory.PRODUCTIVE,      6 * 60_000L),
        ).map { (name, cat, time) ->
            AppUsageData("demo.$name", name, null, time, cat)
        }
        val todayTotal = mockApps.sumOf { it.totalTimeMs }
        return DailyStats(
            totalTimeMs = todayTotal,
            goalMs = TimeUnit.HOURS.toMillis(4),
            topApps = mockApps,
            weeklyData = listOf(
                TimeUnit.HOURS.toMillis(5) + TimeUnit.MINUTES.toMillis(43),
                TimeUnit.HOURS.toMillis(7) + TimeUnit.MINUTES.toMillis(36),
                TimeUnit.HOURS.toMillis(6) + TimeUnit.MINUTES.toMillis(30),
                TimeUnit.HOURS.toMillis(5) + TimeUnit.MINUTES.toMillis(41),
                todayTotal, 0L, 0L,
            )
        )
    }

    fun getStats(context: Context): DailyStats {
        if (demoMode) return getDemoStats()
        val age = System.currentTimeMillis() - cacheTime
        if (cachedStats != null && age < 30_000) return cachedStats!!
        val fresh = getRealStats(context) ?: getDemoStats()
        cachedStats = fresh
        cacheTime = System.currentTimeMillis()
        return fresh
    }

    fun getFreshStats(context: Context): DailyStats {
        if (demoMode) return getDemoStats()
        val fresh = getRealStats(context) ?: getDemoStats()
        cachedStats = fresh
        cacheTime = System.currentTimeMillis()
        return fresh
    }
}