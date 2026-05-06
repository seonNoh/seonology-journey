package com.seonology.journey.ui

import android.app.Activity
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.animation.core.RepeatMode
import androidx.compose.animation.core.animateFloat
import androidx.compose.animation.core.infiniteRepeatable
import androidx.compose.animation.core.rememberInfiniteTransition
import androidx.compose.animation.core.tween
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.statusBars
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.windowInsetsPadding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.CalendarMonth
import androidx.compose.material.icons.filled.Flight
import androidx.compose.material.icons.filled.Logout
import androidx.compose.material.icons.filled.Place
import androidx.compose.material3.AssistChip
import androidx.compose.material3.AssistChipDefaults
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilledIconButton
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.IconButtonDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.rotate
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.seonology.journey.R
import com.seonology.journey.auth.AuthStore
import com.seonology.journey.auth.KeycloakAuth
import com.seonology.journey.data.JourneyApi
import com.seonology.journey.data.Network
import com.seonology.journey.data.Trip
import com.seonology.journey.data.Day
import com.seonology.journey.ui.theme.Sakura100
import com.seonology.journey.ui.theme.Sakura200
import com.seonology.journey.ui.theme.Sakura500
import com.seonology.journey.ui.theme.Sakura50
import com.seonology.journey.ui.theme.Sakura600
import com.seonology.journey.ui.theme.Sakura700
import com.seonology.journey.ui.theme.Spacing
import com.seonology.journey.ui.theme.Warm50
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

    LaunchedEffect(authed) {
        if (authed) {
            com.seonology.journey.sync.SyncWorker.enqueueOneShot(context)
            com.seonology.journey.sync.SyncWorker.enqueuePeriodic(context)
        }
    }

    // Sakura-tinted edge-to-edge background so the whole app always has the
    // theme vibe rather than the default white.
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(
                Brush.verticalGradient(
                    listOf(Sakura50, Warm50),
                ),
            ),
    ) {
        if (!authed) {
            LoginScreen(onLogin = {
                val cfg = KeycloakAuth.config()
                val req = KeycloakAuth.buildAuthRequest(cfg)
                launcher.launch(authService.getAuthorizationRequestIntent(req))
            })
            return@Box
        }

        NavHost(
            navController = nav,
            startDestination = "trips",
            modifier = Modifier.fillMaxSize(),
        ) {
            composable("trips") {
                TripListScreen(
                    api = api,
                    onOpenTrip = { tripId -> nav.navigate("trips/$tripId") },
                    onLogout = {
                        store.clear()
                        authed = false
                    },
                )
            }
            composable("trips/{tripId}") { backStackEntry ->
                val tripId = backStackEntry.arguments?.getString("tripId").orEmpty()
                TripDetailScreen(
                    api = api,
                    tripId = tripId,
                    onBack = { nav.popBackStack() },
                )
            }
        }
    }
}

// ---- Login ----

@Composable
private fun LoginScreen(onLogin: () -> Unit) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .windowInsetsPadding(WindowInsets.statusBars)
            .padding(horizontal = Spacing.xl, vertical = Spacing.lg),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        // Animated plane + drifting petals for cute loading vibe.
        CuteHero()
        Spacer(Modifier.height(Spacing.lg))
        Text(
            "旅の記録, ここから.",
            style = MaterialTheme.typography.headlineMedium.copy(fontWeight = FontWeight.Bold),
            color = Sakura700,
        )
        Spacer(Modifier.height(Spacing.sm))
        Text(
            "여행 계획, 일정, 사진까지 한 곳에서.",
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
        )
        Spacer(Modifier.height(Spacing.xl))
        Button(
            onClick = onLogin,
            shape = RoundedCornerShape(12.dp),
            modifier = Modifier
                .fillMaxWidth()
                .height(52.dp),
        ) {
            Text("Keycloak 로그인", fontSize = 16.sp, fontWeight = FontWeight.SemiBold)
        }
    }
}

/**
 * CuteHero — 로그인/빈 상태에 쓰는 sakura 테마 일러스트. 비행기가 살짝
 * 위아래로 흔들리고 벚꽃 한 장이 회전한다. 정적 SVG 대신 인라인
 * 애니메이션이라 리소스 추가 없이 분위기가 산다.
 */
@Composable
private fun CuteHero() {
    val infinite = rememberInfiniteTransition(label = "cute-hero")
    val bob by infinite.animateFloat(
        initialValue = 0f,
        targetValue = 8f,
        animationSpec = infiniteRepeatable(tween(1400), RepeatMode.Reverse),
        label = "bob",
    )
    val spin by infinite.animateFloat(
        initialValue = 0f,
        targetValue = 360f,
        animationSpec = infiniteRepeatable(tween(5200)),
        label = "spin",
    )

    Box(
        modifier = Modifier
            .size(140.dp)
            .background(
                Brush.radialGradient(listOf(Sakura100, Sakura50)),
                CircleShape,
            ),
        contentAlignment = Alignment.Center,
    ) {
        Icon(
            Icons.Default.Flight,
            contentDescription = null,
            tint = Sakura500,
            modifier = Modifier
                .size(64.dp)
                .padding(bottom = bob.dp)
                .rotate(-20f),
        )
        // Petal orbiting the plane.
        Box(
            modifier = Modifier
                .size(120.dp)
                .rotate(spin),
        ) {
            Box(
                modifier = Modifier
                    .align(Alignment.TopEnd)
                    .size(14.dp)
                    .background(Sakura200, CircleShape),
            )
        }
    }
}

// ---- Trip list ----

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun TripListScreen(
    api: JourneyApi,
    onOpenTrip: (String) -> Unit,
    onLogout: () -> Unit,
) {
    var trips by remember { mutableStateOf<List<Trip>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(Unit) {
        loading = true
        runCatching { api.listTrips() }
            .onSuccess {
                trips = it.trips
                loading = false
            }
            .onFailure {
                error = it.message
                loading = false
            }
    }

    Scaffold(
        containerColor = Warm50.copy(alpha = 0f),
        topBar = {
            TopAppBar(
                title = {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Icon(
                            Icons.Default.Flight,
                            contentDescription = null,
                            tint = Sakura600,
                            modifier = Modifier
                                .size(20.dp)
                                .rotate(-20f),
                        )
                        Spacer(Modifier.width(Spacing.sm))
                        Text(
                            "Seonology Journey",
                            fontWeight = FontWeight.Bold,
                            color = Sakura700,
                        )
                    }
                },
                actions = {
                    IconButton(onClick = onLogout) {
                        Icon(Icons.Default.Logout, contentDescription = "logout", tint = Sakura600)
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = Warm50.copy(alpha = 0.75f),
                ),
            )
        },
    ) { padding ->
        when {
            loading -> CenteredLoader(padding)
            error != null -> ErrorState(padding, error!!)
            trips.isEmpty() -> EmptyTrips(padding)
            else -> {
                LazyColumn(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(padding)
                        .padding(horizontal = Spacing.base),
                    verticalArrangement = Arrangement.spacedBy(Spacing.md),
                    contentPadding = PaddingValues(vertical = Spacing.md),
                ) {
                    items(trips, key = { it.id }) { t ->
                        TripCard(trip = t, onClick = { onOpenTrip(t.id) })
                    }
                }
            }
        }
    }
}

@Composable
private fun TripCard(trip: Trip, onClick: () -> Unit) {
    Card(
        onClick = onClick,
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp, pressedElevation = 6.dp),
        modifier = Modifier.fillMaxWidth(),
    ) {
        Column(
            modifier = Modifier.padding(Spacing.base),
            verticalArrangement = Arrangement.spacedBy(Spacing.xs),
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(
                    trip.title.ifBlank { "이름 없는 여행" },
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold,
                    color = MaterialTheme.colorScheme.onSurface,
                    modifier = Modifier.weight(1f),
                )
                if (!trip.status.isNullOrBlank()) {
                    StatusPill(trip.status)
                }
            }
            if (!trip.destination.isNullOrBlank()) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(
                        Icons.Default.Place,
                        contentDescription = null,
                        tint = Sakura500,
                        modifier = Modifier.size(16.dp),
                    )
                    Spacer(Modifier.width(Spacing.xs))
                    Text(
                        trip.destination,
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
                }
            }
            if (!trip.startDate.isNullOrBlank()) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(
                        Icons.Default.CalendarMonth,
                        contentDescription = null,
                        tint = Sakura500,
                        modifier = Modifier.size(16.dp),
                    )
                    Spacer(Modifier.width(Spacing.xs))
                    Text(
                        "${trip.startDate} → ${trip.endDate ?: ""}",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
                }
            }
        }
    }
}

@Composable
private fun StatusPill(status: String) {
    val label = when (status) {
        "TRIP_STATUS_PLANNING" -> "계획중"
        "TRIP_STATUS_ONGOING" -> "여행중"
        "TRIP_STATUS_COMPLETED" -> "완료"
        "TRIP_STATUS_ARCHIVED" -> "보관"
        else -> status
    }
    AssistChip(
        onClick = {},
        label = { Text(label, fontSize = 11.sp) },
        colors = AssistChipDefaults.assistChipColors(
            containerColor = Sakura100,
            labelColor = Sakura700,
        ),
        border = null,
    )
}

// ---- Trip detail ----

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun TripDetailScreen(
    api: JourneyApi,
    tripId: String,
    onBack: () -> Unit,
) {
    var trip by remember { mutableStateOf<Trip?>(null) }
    var days by remember { mutableStateOf<List<Day>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(tripId) {
        loading = true
        runCatching {
            val t = api.getTrip(tripId).trip
            val d = api.listDays(tripId).days
            t to d
        }
            .onSuccess { (t, d) ->
                trip = t
                days = d
                loading = false
            }
            .onFailure {
                error = it.message
                loading = false
            }
    }

    Scaffold(
        containerColor = Warm50.copy(alpha = 0f),
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        trip?.title ?: "여행 상세",
                        fontWeight = FontWeight.Bold,
                        color = Sakura700,
                    )
                },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(
                            Icons.Default.ArrowBack,
                            contentDescription = "뒤로",
                            tint = Sakura600,
                        )
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = Warm50.copy(alpha = 0.75f),
                ),
            )
        },
    ) { padding ->
        when {
            loading -> CenteredLoader(padding)
            error != null -> ErrorState(padding, error!!)
            else -> Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding)
                    .verticalScroll(rememberScrollState())
                    .padding(horizontal = Spacing.base, vertical = Spacing.md),
                verticalArrangement = Arrangement.spacedBy(Spacing.md),
            ) {
                trip?.let { t -> TripHeaderCard(t) }

                Text(
                    "일정 (${days.size}일)",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold,
                    color = MaterialTheme.colorScheme.onSurface,
                )

                if (days.isEmpty()) {
                    Text(
                        "등록된 일정이 없습니다.",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
                } else {
                    days.forEach { d -> DayRow(d) }
                }
                Spacer(Modifier.height(Spacing.xl))
            }
        }
    }
}

@Composable
private fun TripHeaderCard(trip: Trip) {
    Card(
        shape = RoundedCornerShape(20.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp),
        modifier = Modifier.fillMaxWidth(),
    ) {
        Column(
            modifier = Modifier.padding(Spacing.base),
            verticalArrangement = Arrangement.spacedBy(Spacing.xs),
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(
                    trip.title,
                    style = MaterialTheme.typography.titleLarge,
                    fontWeight = FontWeight.Bold,
                    modifier = Modifier.weight(1f),
                )
                if (!trip.status.isNullOrBlank()) StatusPill(trip.status)
            }
            if (!trip.destination.isNullOrBlank()) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(
                        Icons.Default.Place,
                        contentDescription = null,
                        tint = Sakura500,
                        modifier = Modifier.size(16.dp),
                    )
                    Spacer(Modifier.width(Spacing.xs))
                    Text(trip.destination, color = MaterialTheme.colorScheme.onSurfaceVariant)
                }
            }
            if (!trip.startDate.isNullOrBlank()) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(
                        Icons.Default.CalendarMonth,
                        contentDescription = null,
                        tint = Sakura500,
                        modifier = Modifier.size(16.dp),
                    )
                    Spacer(Modifier.width(Spacing.xs))
                    Text(
                        "${trip.startDate} → ${trip.endDate ?: ""}",
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
                }
            }
            if (!trip.description.isNullOrBlank()) {
                Spacer(Modifier.height(Spacing.xs))
                Text(
                    trip.description,
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onSurface,
                )
            }
        }
    }
}

@Composable
private fun DayRow(d: Day) {
    Card(
        shape = RoundedCornerShape(14.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp),
        modifier = Modifier.fillMaxWidth(),
    ) {
        Row(
            modifier = Modifier.padding(Spacing.base),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            FilledIconButton(
                onClick = {},
                shape = CircleShape,
                colors = IconButtonDefaults.filledIconButtonColors(
                    containerColor = Sakura100,
                    contentColor = Sakura700,
                ),
                modifier = Modifier.size(44.dp),
            ) {
                Text(
                    "${d.dayNumber}",
                    fontWeight = FontWeight.Bold,
                    fontSize = 14.sp,
                )
            }
            Spacer(Modifier.width(Spacing.md))
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    d.date,
                    style = MaterialTheme.typography.bodyLarge,
                    fontWeight = FontWeight.SemiBold,
                )
                if (!d.region.isNullOrBlank()) {
                    Text(
                        d.region,
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
                }
            }
        }
    }
}

// ---- Shared empty / loading / error states ----

@Composable
private fun CenteredLoader(padding: PaddingValues) {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .padding(padding),
        contentAlignment = Alignment.Center,
    ) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            CuteHero()
            Spacer(Modifier.height(Spacing.md))
            Text(
                "여행을 꺼내오고 있어요",
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant,
            )
            Spacer(Modifier.height(Spacing.sm))
            CircularProgressIndicator(
                color = Sakura500,
                modifier = Modifier.size(24.dp),
                strokeWidth = 3.dp,
            )
        }
    }
}

@Composable
private fun EmptyTrips(padding: PaddingValues) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(padding)
            .padding(Spacing.xl),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        Image(
            painter = painterResource(R.drawable.il_empty_trips),
            contentDescription = null,
            modifier = Modifier.size(200.dp),
        )
        Spacer(Modifier.height(Spacing.md))
        Text(
            "아직 여행이 없어요",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.SemiBold,
        )
        Spacer(Modifier.height(Spacing.xs))
        Text(
            "웹에서 첫 번째 여행을 만들면 여기에 보여요.",
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
        )
    }
}

@Composable
private fun ErrorState(padding: PaddingValues, message: String) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(padding)
            .padding(Spacing.xl),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        Text(
            "잠깐만요",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.SemiBold,
            color = Sakura700,
        )
        Spacer(Modifier.height(Spacing.xs))
        Text(
            message,
            style = MaterialTheme.typography.bodySmall,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
        )
    }
}
