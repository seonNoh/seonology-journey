package com.seonology.journey.ui.screens

import android.Manifest
import android.content.Context
import android.content.Intent
import android.os.Build
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import com.seonology.journey.tracking.LocationTrackingService
import com.seonology.journey.ui.theme.Spacing

@Composable
fun TrackingScreen(
    wsUrl: String,
    token: String,
    partnerConsented: Boolean,
) {
    val context = LocalContext.current
    var myConsent by remember { mutableStateOf(false) }
    var showConsentDialog by remember { mutableStateOf(false) }
    var permissionsGranted by remember { mutableStateOf(false) }

    val permissions = buildList {
        add(Manifest.permission.ACCESS_FINE_LOCATION)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.UPSIDE_DOWN_CAKE) {
            add(Manifest.permission.FOREGROUND_SERVICE_LOCATION)
        }
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            add(Manifest.permission.POST_NOTIFICATIONS)
        }
    }.toTypedArray()

    val permissionLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.RequestMultiplePermissions()
    ) { results ->
        permissionsGranted = results.values.all { it }
        if (permissionsGranted) {
            showConsentDialog = true
        }
    }

    Column(
        modifier = Modifier.fillMaxSize().padding(Spacing.lg),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        Text("실시간 위치 공유", style = MaterialTheme.typography.headlineSmall)
        Spacer(Modifier.height(Spacing.base))

        when {
            myConsent && partnerConsented -> {
                Text("위치 공유 활성 중", color = MaterialTheme.colorScheme.primary)
                Spacer(Modifier.height(Spacing.base))
                Button(onClick = {
                    myConsent = false
                    stopTracking(context)
                }) {
                    Text("공유 중지")
                }
            }
            else -> {
                if (!partnerConsented) {
                    Text("상대방의 동의를 기다리고 있습니다", color = MaterialTheme.colorScheme.outline)
                }
                Spacer(Modifier.height(Spacing.base))
                Button(onClick = {
                    permissionLauncher.launch(permissions)
                }) {
                    Text("위치 공유 시작")
                }
            }
        }
    }

    if (showConsentDialog) {
        AlertDialog(
            onDismissRequest = { showConsentDialog = false },
            title = { Text("위치 공유 동의") },
            text = { Text("상대방과 실시간으로 위치를 공유합니다. 동의하시겠습니까?") },
            confirmButton = {
                TextButton(onClick = {
                    myConsent = true
                    showConsentDialog = false
                    if (partnerConsented) startTracking(context, wsUrl, token)
                }) { Text("동의") }
            },
            dismissButton = {
                TextButton(onClick = { showConsentDialog = false }) { Text("취소") }
            },
        )
    }
}

private fun startTracking(context: Context, wsUrl: String, token: String) {
    val intent = Intent(context, LocationTrackingService::class.java).apply {
        putExtra(LocationTrackingService.EXTRA_WS_URL, wsUrl)
        putExtra(LocationTrackingService.EXTRA_TOKEN, token)
    }
    context.startForegroundService(intent)
}

private fun stopTracking(context: Context) {
    context.stopService(Intent(context, LocationTrackingService::class.java))
}
