package com.seonology.journey.ui

import android.app.Activity
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Logout
import androidx.compose.material.icons.filled.Place
import androidx.compose.material3.Button
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.OutlinedCard
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.seonology.journey.auth.AuthStore
import com.seonology.journey.auth.KeycloakAuth
import com.seonology.journey.data.JourneyApi
import com.seonology.journey.data.Network
import com.seonology.journey.data.Trip
import kotlinx.coroutines.launch
import net.openid.appauth.AuthorizationException
import net.openid.appauth.AuthorizationResponse

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun JourneyApp() {
    val context = LocalContext.current
    val store = remember { AuthStore(context) }
    val api = remember { Network.journeyApi(store) }
    val nav = rememberNavController()
    var authed by remember { mutableStateOf(store.isAuthenticated) }

    val scope = rememberCoroutineScope()
    val authService = remember { KeycloakAuth.newService(context as Activity) }
    val launcher = rememberLauncherForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
        val data = result.data ?: return@rememberLauncherForActivityResult
        val resp = AuthorizationResponse.fromIntent(data)
        val ex = AuthorizationException.fromIntent(data)
        if (resp != null) {
            authService.performTokenRequest(resp.createTokenExchangeRequest()) { tokenResp, _ ->
                KeycloakAuth.handleTokenResponse(context, tokenResp)
                authed = store.isAuthenticated
            }
        } else if (ex != null) {
            // ignore for MVP.
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Seonology Journey") },
                actions = {
                    if (authed) {
                        IconButton(onClick = {
                            store.clear()
                            authed = false
                        }) { Icon(Icons.Default.Logout, contentDescription = "logout") }
                    }
                },
            )
        },
    ) { padding ->
        if (!authed) {
            LoginScreen(padding) {
                val cfg = KeycloakAuth.config()
                val req = KeycloakAuth.buildAuthRequest(cfg)
                launcher.launch(authService.getAuthorizationRequestIntent(req))
            }
            return@Scaffold
        }
        NavHost(navController = nav, startDestination = "trips", modifier = Modifier.padding(padding)) {
            composable("trips") { TripListScreen(api) }
        }
    }

    LaunchedEffect(Unit) {
        scope.launch { /* future: token refresh */ }
    }
}

@Composable
private fun LoginScreen(padding: PaddingValues, onLogin: () -> Unit) {
    Column(
        modifier = Modifier.fillMaxSize().padding(padding).padding(24.dp),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        Text("旅の記録, ここから.", style = androidx.compose.material3.MaterialTheme.typography.headlineSmall)
        Button(onClick = onLogin, modifier = Modifier.padding(top = 16.dp)) {
            Text("Keycloak 로그인")
        }
    }
}

@Composable
private fun TripListScreen(api: JourneyApi) {
    var trips by remember { mutableStateOf<List<Trip>>(emptyList()) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(Unit) {
        runCatching { api.listTrips() }
            .onSuccess { trips = it.trips }
            .onFailure { error = it.message }
    }

    LazyColumn(
        modifier = Modifier.fillMaxSize().padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp),
    ) {
        error?.let { item { Text("Error: $it") } }
        items(trips) { t ->
            OutlinedCard {
                Column(Modifier.padding(16.dp)) {
                    Text(t.title, style = androidx.compose.material3.MaterialTheme.typography.titleMedium)
                    if (t.destination != null) {
                        Row(t.destination)
                    }
                    if (t.startDate != null) {
                        Text("${t.startDate} → ${t.endDate ?: ""}")
                    }
                }
            }
        }
        if (trips.isEmpty() && error == null) {
            item { Text("아직 여행이 없습니다.") }
        }
    }
}

@Composable
private fun Row(text: String) {
    Column {
        Icon(Icons.Default.Place, contentDescription = null)
        Text(text)
    }
}
