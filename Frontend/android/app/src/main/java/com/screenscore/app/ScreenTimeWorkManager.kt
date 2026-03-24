package com.screenscore.app

import android.app.usage.UsageStatsManager
import android.content.Context
import androidx.work.*
import org.json.JSONObject
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL
import java.text.SimpleDateFormat
import java.util.*
import java.util.concurrent.TimeUnit

class ScreenTimeWorkManager(context: Context, workerParams: WorkerParameters) :
    Worker(context, workerParams) {

    override fun doWork(): Result {
        return try {
            val stats = ScreenTimeManager.getRealStats(applicationContext)
                ?: return Result.retry()

            val screenTimeMinutes = (stats.totalTimeMs / 60000).toInt()

            val sdf = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSS'Z'", Locale.US)
            sdf.timeZone = TimeZone.getTimeZone("UTC")
            val recordedAt = sdf.format(Date())

            val json = JSONObject()
            json.put("screen_time", screenTimeMinutes)
            json.put("recorded_at", recordedAt)

            // Same endpoint as syncToBackend() in Home.jsx
            val url = URL("http://https://shenika-ovarian-unpiratically.ngrok-free.dev/v1/stats/sync-stats")
            val connection = url.openConnection() as HttpURLConnection
            connection.requestMethod = "POST"
            connection.setRequestProperty("Content-Type", "application/json")
            connection.doOutput = true
            connection.connectTimeout = 10000
            connection.readTimeout = 10000

            val writer = OutputStreamWriter(connection.outputStream)
            writer.write(json.toString())
            writer.flush()
            writer.close()

            val responseCode = connection.responseCode
            connection.disconnect()

            if (responseCode in 200..299) Result.success() else Result.retry()

        } catch (e: Exception) {
            Result.retry()
        }
    }

    companion object {
        private const val WORK_NAME = "ScreenTimePeriodicSync"

        fun schedulePeriodicSync(context: Context) {
            val workRequest = PeriodicWorkRequestBuilder<ScreenTimeWorkManager>(
                15, TimeUnit.MINUTES
            ).build()

            WorkManager.getInstance(context).enqueueUniquePeriodicWork(
                WORK_NAME,
                ExistingPeriodicWorkPolicy.KEEP,
                workRequest
            )
        }
    }
}