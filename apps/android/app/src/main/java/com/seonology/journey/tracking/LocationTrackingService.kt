package com.seonology.journey.tracking

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.Service
import android.content.Intent
import android.os.Build
import android.os.IBinder
import androidx.core.app.NotificationCompat
import com.google.android.gms.location.FusedLocationProviderClient
import com.google.android.gms.location.LocationCallback
import com.google.android.gms.location.LocationRequest
import com.google.android.gms.location.LocationResult
import com.google.android.gms.location.LocationServices
import com.google.android.gms.location.Priority
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.cancel
import kotlinx.coroutines.launch
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.Response
import okhttp3.WebSocket
import okhttp3.WebSocketListener

class LocationTrackingService : Service() {

    companion object {
        const val CHANNEL_ID = "journey_tracking"
        const val NOTIFICATION_ID = 1001
        const val EXTRA_WS_URL = "ws_url"
        const val EXTRA_TOKEN = "token"
    }

    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.IO)
    private lateinit var fusedClient: FusedLocationProviderClient
    private var webSocket: WebSocket? = null
    private var locationCallback: LocationCallback? = null

    override fun onBind(intent: Intent?): IBinder? = null

    override fun onCreate() {
        super.onCreate()
        createNotificationChannel()
        fusedClient = LocationServices.getFusedLocationProviderClient(this)
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        startForeground(NOTIFICATION_ID, buildNotification())

        val wsUrl = intent?.getStringExtra(EXTRA_WS_URL) ?: return START_NOT_STICKY
        val token = intent.getStringExtra(EXTRA_TOKEN) ?: ""

        connectWebSocket(wsUrl, token)
        startLocationUpdates()

        return START_STICKY
    }

    override fun onDestroy() {
        locationCallback?.let { fusedClient.removeLocationUpdates(it) }
        webSocket?.close(1000, "service stopped")
        scope.cancel()
        super.onDestroy()
    }

    private fun connectWebSocket(url: String, token: String) {
        val client = OkHttpClient.Builder().build()
        val request = Request.Builder()
            .url(url)
            .addHeader("Authorization", "Bearer $token")
            .build()

        webSocket = client.newWebSocket(request, object : WebSocketListener() {
            override fun onFailure(webSocket: WebSocket, t: Throwable, response: Response?) {
                // Exponential backoff reconnect handled by caller
                scope.launch {
                    kotlinx.coroutines.delay(3000)
                    connectWebSocket(url, token)
                }
            }
        })
    }

    @Suppress("MissingPermission")
    private fun startLocationUpdates() {
        val request = LocationRequest.Builder(Priority.PRIORITY_HIGH_ACCURACY, 5000L)
            .setMinUpdateDistanceMeters(10f)
            .build()

        locationCallback = object : LocationCallback() {
            override fun onLocationResult(result: LocationResult) {
                val loc = result.lastLocation ?: return
                if (loc.accuracy > 100f) return // skip inaccurate
                val msg = """{"lat":${loc.latitude},"lng":${loc.longitude},"ts":${System.currentTimeMillis()}}"""
                webSocket?.send(msg)
            }
        }

        fusedClient.requestLocationUpdates(request, locationCallback!!, mainLooper)
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(CHANNEL_ID, "Location Tracking", NotificationManager.IMPORTANCE_LOW)
            getSystemService(NotificationManager::class.java).createNotificationChannel(channel)
        }
    }

    private fun buildNotification(): Notification =
        NotificationCompat.Builder(this, CHANNEL_ID)
            .setContentTitle("Journey")
            .setContentText("위치 공유 중")
            .setSmallIcon(android.R.drawable.ic_dialog_map)
            .build()
}
