package com.seonology.journey.ui

import android.app.Activity
import android.widget.Toast
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.animation.core.RepeatMode
import androidx.compose.animation.core.animateFloat
import androidx.compose.animation.core.infiniteRepeatable
import androidx.compose.animation.core.rememberInfiniteTransition
import androidx.compose.animation.core.tween
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.horizontalScroll
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
import androidx.compose.foundation.layout.offset
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
import androidx.compose.foundation.text.BasicTextField
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.CalendarMonth
import androidx.compose.material.icons.filled.CameraAlt
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Favorite
import androidx.compose.material.icons.filled.Flight
import androidx.compose.material.icons.filled.Hotel
import androidx.compose.material.icons.filled.LocalAtm
import androidx.compose.material.icons.filled.Logout
import androidx.compose.material.icons.filled.Note
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material.icons.filled.Place
import androidx.compose.material.icons.filled.Restaurant
import androidx.compose.material.icons.filled.ShoppingBag
import androidx.compose.material.icons.filled.Star
import androidx.compose.material.icons.filled.Wallet
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.LinearProgressIndicator
import androidx.compose.material3.OutlinedButton
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
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.rotate
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.SolidColor
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.viewinterop.AndroidView
import android.webkit.WebView
import android.webkit.WebSettings
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import kotlinx.coroutines.launch
import com.seonology.journey.auth.AuthStore
import com.seonology.journey.auth.KeycloakAuth
import com.seonology.journey.data.Accommodation
import com.seonology.journey.data.CreateScheduleRequest
import com.seonology.journey.data.CreateTripRequest
import com.seonology.journey.data.Day
import com.seonology.journey.data.Expense
import com.seonology.journey.data.GeoPoint
import com.seonology.journey.data.JourneyApi
import com.seonology.journey.data.Meal
import com.seonology.journey.data.Network
import com.seonology.journey.data.Note
import com.seonology.journey.data.Schedule
import com.seonology.journey.data.Trip
import com.seonology.journey.ui.components.MascotAccessory
import com.seonology.journey.ui.components.MascotExpression
import com.seonology.journey.ui.components.PlaceSearchField
import com.seonology.journey.ui.components.SbBear
import com.seonology.journey.ui.components.SbChick
import com.seonology.journey.ui.components.SbChip
import com.seonology.journey.ui.components.SbMiniBear
import com.seonology.journey.ui.components.SbPaw
import com.seonology.journey.ui.components.SbPetal
import com.seonology.journey.ui.components.SbSection
import com.seonology.journey.ui.components.SbStatusPill
import com.seonology.journey.ui.components.categoryFor
import com.seonology.journey.ui.components.moodFor
import com.seonology.journey.ui.theme.Sakura100
import com.seonology.journey.ui.theme.Sakura200
import com.seonology.journey.ui.theme.Sakura300
import com.seonology.journey.ui.theme.Sakura400
import com.seonology.journey.ui.theme.Sakura50
import com.seonology.journey.ui.theme.Sakura500
import com.seonology.journey.ui.theme.Sakura600
import com.seonology.journey.ui.theme.Sakura700
import com.seonology.journey.ui.theme.Sakura900
import com.seonology.journey.ui.theme.SbBeige
import com.seonology.journey.ui.theme.SbCream
import com.seonology.journey.ui.theme.Spacing
import com.seonology.journey.ui.theme.Warm100
import com.seonology.journey.ui.theme.Warm400
import com.seonology.journey.ui.theme.Warm50
import com.seonology.journey.ui.theme.Warm500
import com.seonology.journey.ui.theme.Warm600
import com.seonology.journey.ui.theme.Warm700
import com.seonology.journey.ui.theme.Warm800
import com.seonology.journey.ui.theme.Warm900
import net.openid.appauth.AuthorizationException
import net.openid.appauth.AuthorizationResponse

/**
 * Sakura Bear 디자인 시스템을 적용한 앱 셸. 인증 상태 → 목록 → 상세 →
 * Day 상세 → 메모/지출 의 단순한 스택 네비게이션을 직접 구성한다.
 *
 * 디자인 원본은 `design/sakura-bear-screens.jsx`, `design/sakura-bear-mascot.jsx`,
 * `design/sakura-bear-system.jsx` 참조.
 */
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
            // 로그인 실패는 조용히 무시한다.
        }
    }

    LaunchedEffect(authed) {
        if (authed) {
            com.seonology.journey.sync.SyncWorker.enqueueOneShot(context)
            com.seonology.journey.sync.SyncWorker.enqueuePeriodic(context)
        }
    }

    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(
                Brush.verticalGradient(listOf(Sakura50, SbCream, Warm50)),
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
                    onOpenTrip = { nav.navigate("trips/$it") },
                    onCreateTrip = { nav.navigate("create-trip") },
                    onLogout = {
                        store.clear()
                        authed = false
                    },
                )
            }
            composable("create-trip") {
                CreateTripScreen(
                    api = api,
                    onBack = { nav.popBackStack() },
                    onCreated = { tripId ->
                        nav.popBackStack()
                        nav.navigate("trips/$tripId")
                    },
                )
            }
            composable("trips/{tripId}") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                TripDetailScreen(
                    api = api,
                    tripId = tripId,
                    onBack = { nav.popBackStack() },
                    onOpenDay = { nav.navigate("days/$it") },
                    onOpenNotes = { nav.navigate("trips/$tripId/notes") },
                    onOpenExpenses = { nav.navigate("trips/$tripId/expenses") },
                    onOpenMap = { nav.navigate("trips/$tripId/map") },
                    onOpenChecklist = { nav.navigate("trips/$tripId/checklist") },
                    onOpenReservations = { nav.navigate("trips/$tripId/reservations") },
                    onOpenTags = { nav.navigate("trips/$tripId/tags") },
                    onOpenCompanions = { nav.navigate("trips/$tripId/companions") },
                    onOpenShare = { nav.navigate("trips/$tripId/share") },
                    onOpenMedia = { nav.navigate("trips/$tripId/media") },
                    onOpenNearby = { nav.navigate("trips/$tripId/nearby") },
                    onOpenTransit = { nav.navigate("trips/$tripId/transit") },
                )
            }
            composable("trips/{tripId}/checklist") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                com.seonology.journey.ui.screens.ChecklistFullScreen(
                    api = api, tripId = tripId, onBack = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/reservations") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                com.seonology.journey.ui.screens.ReservationsFullScreen(
                    api = api, tripId = tripId, onBack = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/tags") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                com.seonology.journey.ui.screens.TagsFullScreen(
                    api = api, tripId = tripId, onBack = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/companions") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                com.seonology.journey.ui.screens.CompanionsFullScreen(
                    api = api, tripId = tripId, onBack = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/share") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                com.seonology.journey.ui.screens.ShareFullScreen(
                    api = api, tripId = tripId, onBack = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/media") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                com.seonology.journey.ui.screens.MediaFullScreen(
                    api = api, tripId = tripId, onBack = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/nearby") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                com.seonology.journey.ui.screens.NearbyFullScreen(
                    api = api, tripId = tripId, onBack = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/transit") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                com.seonology.journey.ui.screens.TransitFullScreen(
                    api = api, tripId = tripId, onBack = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/notes/add") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                com.seonology.journey.ui.screens.NoteAddScreen(
                    api = api, tripId = tripId,
                    onBack = { nav.popBackStack() },
                    onSaved = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/expenses/add") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                com.seonology.journey.ui.screens.ExpenseAddScreen(
                    api = api, tripId = tripId,
                    onBack = { nav.popBackStack() },
                    onSaved = { nav.popBackStack() },
                )
            }
            composable("days/{dayId}/meal") { entry ->
                val dayId = entry.arguments?.getString("dayId").orEmpty()
                com.seonology.journey.ui.screens.MealEditScreen(
                    api = api, dayId = dayId,
                    onBack = { nav.popBackStack() },
                    onSaved = { nav.popBackStack() },
                )
            }
            composable("days/{dayId}/accommodation") { entry ->
                val dayId = entry.arguments?.getString("dayId").orEmpty()
                com.seonology.journey.ui.screens.AccommodationEditScreen(
                    api = api, dayId = dayId,
                    onBack = { nav.popBackStack() },
                    onSaved = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/map") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                TripMapScreen(
                    api = api,
                    tripId = tripId,
                    onBack = { nav.popBackStack() },
                )
            }
            composable("days/{dayId}") { entry ->
                val dayId = entry.arguments?.getString("dayId").orEmpty()
                DayDetailScreen(
                    api = api,
                    dayId = dayId,
                    onBack = { nav.popBackStack() },
                    onAddSchedule = { nav.navigate("days/$dayId/add-schedule") },
                    onEditMeal = { nav.navigate("days/$dayId/meal") },
                    onEditAccommodation = { nav.navigate("days/$dayId/accommodation") },
                )
            }
            composable("days/{dayId}/add-schedule") { entry ->
                val dayId = entry.arguments?.getString("dayId").orEmpty()
                AddScheduleScreen(
                    api = api,
                    dayId = dayId,
                    onBack = { nav.popBackStack() },
                    onCreated = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/notes") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                NotesScreen(
                    api = api,
                    tripId = tripId,
                    onBack = { nav.popBackStack() },
                    onAdd = { nav.navigate("trips/$tripId/notes/add") },
                )
            }
            composable("trips/{tripId}/expenses") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                ExpensesScreen(
                    api = api,
                    tripId = tripId,
                    onBack = { nav.popBackStack() },
                    onAdd = { nav.navigate("trips/$tripId/expenses/add") },
                )
            }
        }
    }
}

// ──────────────────────────────────────────────────────────────────────
// Login (design: ScreenLogin)
// ──────────────────────────────────────────────────────────────────────

@Composable
private fun LoginScreen(onLogin: () -> Unit) {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(
                Brush.radialGradient(listOf(Sakura100, Sakura50, SbCream)),
            )
            .windowInsetsPadding(WindowInsets.statusBars),
    ) {
        FloatingPetals(
            entries = listOf(
                PetalEntry(30.dp, 80.dp, -20f, 14.dp, Sakura300),
                PetalEntry(340.dp, 60.dp, 30f, 12.dp, Sakura400),
                PetalEntry(60.dp, 200.dp, 60f, 16.dp, Sakura300),
                PetalEntry(310.dp, 300.dp, -30f, 14.dp, Sakura400),
                PetalEntry(40.dp, 380.dp, 10f, 10.dp, Sakura300),
                PetalEntry(330.dp, 460.dp, -50f, 12.dp, Sakura400),
            ),
        )

        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(horizontal = Spacing.xl, vertical = Spacing.lg),
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            Spacer(Modifier.height(60.dp))
            HeroMascotScene()
            Spacer(Modifier.height(Spacing.lg))
            Text(
                "SEONOLOGY · JOURNEY",
                fontSize = 13.sp,
                fontWeight = FontWeight.Bold,
                letterSpacing = 2.sp,
                color = Sakura600,
            )
            Spacer(Modifier.height(Spacing.sm))
            Text("旅の記録,", fontSize = 30.sp, fontWeight = FontWeight.Bold, color = Sakura900)
            Text("곰돌이와 함께.", fontSize = 30.sp, fontWeight = FontWeight.Bold, color = Sakura900)
            Spacer(Modifier.height(Spacing.md))
            Text("여행 계획·일정·사진·메모·지출까지", fontSize = 13.sp, color = Warm600)
            Text("한 권의 폭신폭신한 노트에.", fontSize = 13.sp, color = Warm600)

            Spacer(Modifier.weight(1f))

            Button(
                onClick = onLogin,
                shape = RoundedCornerShape(18.dp),
                colors = ButtonDefaults.buttonColors(
                    containerColor = Sakura500,
                    contentColor = Color.White,
                ),
                modifier = Modifier
                    .fillMaxWidth()
                    .height(54.dp),
            ) {
                Icon(Icons.Default.CheckCircle, contentDescription = null, modifier = Modifier.size(18.dp))
                Spacer(Modifier.width(Spacing.sm))
                Text("Keycloak으로 로그인", fontSize = 15.sp, fontWeight = FontWeight.Bold)
            }
            Spacer(Modifier.height(10.dp))
            OutlinedButton(
                onClick = {},
                enabled = false,
                shape = RoundedCornerShape(16.dp),
                modifier = Modifier
                    .fillMaxWidth()
                    .height(48.dp),
                colors = ButtonDefaults.outlinedButtonColors(
                    containerColor = Color.White,
                    contentColor = Sakura700,
                    disabledContentColor = Warm400,
                    disabledContainerColor = Color.White,
                ),
                border = BorderStroke(1.dp, Sakura100),
            ) {
                Text("둘러보기 (게스트)", fontSize = 13.sp, fontWeight = FontWeight.Bold)
            }
            Spacer(Modifier.height(Spacing.md))
            Text(
                "계속하시면 이용약관 및 개인정보처리방침에 동의하게 됩니다.",
                fontSize = 10.sp,
                color = Warm400,
            )
            Spacer(Modifier.height(Spacing.sm))
        }
    }
}

@Composable
private fun HeroMascotScene() {
    Box(
        modifier = Modifier
            .width(220.dp)
            .height(180.dp),
    ) {
        // Ground band behind the mascots.
        Box(
            modifier = Modifier
                .align(Alignment.BottomCenter)
                .fillMaxWidth()
                .height(14.dp)
                .clip(RoundedCornerShape(topStart = 50.dp, topEnd = 50.dp))
                .background(Sakura100),
        )
        // Center: main bear.
        Box(modifier = Modifier.align(Alignment.BottomCenter)) {
            SbBear(
                size = 148.dp,
                expression = MascotExpression.Happy,
                accessory = MascotAccessory.Flower,
            )
        }
        // Left: chick.
        Box(modifier = Modifier.align(Alignment.BottomStart)) {
            SbChick(size = 56.dp, expression = MascotExpression.Happy)
        }
        // Right: mini bear, lifted slightly.
        Box(
            modifier = Modifier
                .align(Alignment.BottomEnd)
                .padding(bottom = 24.dp),
        ) {
            SbMiniBear(size = 70.dp, expression = MascotExpression.Happy)
        }
    }
}

private data class PetalEntry(
    val x: Dp,
    val y: Dp,
    val rotateDeg: Float,
    val size: Dp,
    val color: Color,
)

@Composable
private fun FloatingPetals(entries: List<PetalEntry>) {
    Box(modifier = Modifier.fillMaxSize()) {
        entries.forEach { p ->
            Box(modifier = Modifier.offset(x = p.x, y = p.y)) {
                SbPetal(size = p.size, color = p.color, rotateDeg = p.rotateDeg, opacity = 0.55f)
            }
        }
    }
}

// ──────────────────────────────────────────────────────────────────────
// Trip list (design: ScreenTripList)
// ──────────────────────────────────────────────────────────────────────

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun TripListScreen(
    api: JourneyApi,
    onOpenTrip: (String) -> Unit,
    onCreateTrip: () -> Unit,
    onLogout: () -> Unit,
) {
    var trips by remember { mutableStateOf<List<Trip>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(Unit) {
        loading = true
        runCatching { api.listTrips() }
            .onSuccess { trips = it.trips; loading = false }
            .onFailure { error = it.message; loading = false }
    }

    Box(modifier = Modifier.fillMaxSize()) {
        // Decorative background petals.
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(240.dp),
        ) {
            FloatingPetals(
                entries = listOf(
                    PetalEntry(16.dp, 24.dp, -10f, 12.dp, Sakura500),
                    PetalEntry(320.dp, 16.dp, 20f, 10.dp, Sakura400),
                    PetalEntry(60.dp, 84.dp, 45f, 14.dp, Sakura300),
                    PetalEntry(280.dp, 120.dp, -30f, 10.dp, Sakura500),
                    PetalEntry(200.dp, 60.dp, 70f, 8.dp, Sakura400),
                ),
            )
        }

        Column(modifier = Modifier.fillMaxSize()) {
            TripListHeader(onLogout = onLogout)
            when {
                loading -> CenteredLoader()
                error != null -> ErrorState(error!!)
                trips.isEmpty() -> EmptyTrips()
                else -> {
                    val upcoming = trips.firstOrNull { it.status != "TRIP_STATUS_COMPLETED" }
                        ?: trips.first()
                    val rest = trips.filter { it.id != upcoming.id }
                    LazyColumn(
                        modifier = Modifier.fillMaxSize(),
                        contentPadding = PaddingValues(bottom = Spacing.xl),
                        verticalArrangement = Arrangement.spacedBy(Spacing.md),
                    ) {
                        item { NextTripHeroCard(upcoming, onClick = { onOpenTrip(upcoming.id) }) }
                        item {
                            QuickActionRow(
                                onCreateTrip = onCreateTrip,
                                onSchedule = { onOpenTrip(upcoming.id) },
                                onNotes = { onOpenTrip(upcoming.id) },
                            )
                        }
                        item {
                            SbSection(
                                title = "다가오는 여행",
                                count = "(${trips.count { it.status != "TRIP_STATUS_COMPLETED" }})",
                            )
                        }
                        items(rest, key = { it.id }) { t ->
                            TripRow(trip = t, indexHint = trips.indexOf(t)) { onOpenTrip(t.id) }
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun TripListHeader(onLogout: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .windowInsetsPadding(WindowInsets.statusBars)
            .padding(horizontal = Spacing.base, vertical = Spacing.sm),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        SbBear(size = 48.dp, expression = MascotExpression.Happy, accessory = MascotAccessory.Flower)
        Spacer(Modifier.width(Spacing.sm))
        Column(modifier = Modifier.weight(1f)) {
            Text(
                "HELLO ♡ ようこそ",
                fontSize = 11.sp,
                fontWeight = FontWeight.Bold,
                letterSpacing = 1.sp,
                color = Sakura600,
            )
            Text(
                "오늘도 살랑살랑 여행 ♪",
                fontSize = 18.sp,
                fontWeight = FontWeight.Bold,
                color = Sakura900,
            )
        }
        IconRoundButton(icon = Icons.Default.Notifications, onClick = {})
        Spacer(Modifier.width(Spacing.xs))
        IconRoundButton(icon = Icons.Default.Logout, onClick = onLogout)
    }
}

@Composable
private fun IconRoundButton(icon: ImageVector, onClick: () -> Unit) {
    Box(
        modifier = Modifier
            .size(40.dp)
            .clip(RoundedCornerShape(12.dp))
            .background(Color.White)
            .border(1.dp, Sakura100, RoundedCornerShape(12.dp))
            .clickable(onClick = onClick),
        contentAlignment = Alignment.Center,
    ) {
        Icon(icon, contentDescription = null, tint = Sakura600, modifier = Modifier.size(18.dp))
    }
}

@Composable
private fun NextTripHeroCard(trip: Trip, onClick: () -> Unit) {
    Box(
        modifier = Modifier
            .padding(horizontal = Spacing.base)
            .fillMaxWidth()
            .clip(RoundedCornerShape(24.dp))
            .background(Brush.linearGradient(listOf(Sakura100, Color.White)))
            .border(1.5.dp, Sakura200, RoundedCornerShape(24.dp))
            .clickable(onClick = onClick),
    ) {
        // Big bear in lower-right corner.
        Box(
            modifier = Modifier
                .align(Alignment.BottomEnd)
                .offset(x = 4.dp, y = 8.dp),
        ) {
            SbBear(size = 92.dp, expression = MascotExpression.Happy, accessory = MascotAccessory.Camera)
        }
        Column(
            modifier = Modifier
                .padding(Spacing.base)
                .fillMaxWidth(0.72f),
            verticalArrangement = Arrangement.spacedBy(Spacing.xs),
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Box(
                    modifier = Modifier
                        .size(6.dp)
                        .clip(CircleShape)
                        .background(Sakura500),
                )
                Spacer(Modifier.width(Spacing.xs))
                Text(
                    "NEXT TRIP",
                    fontSize = 11.sp,
                    fontWeight = FontWeight.Bold,
                    letterSpacing = 1.sp,
                    color = Sakura600,
                )
            }
            Text(
                trip.title.ifBlank { "이름 없는 여행" },
                fontSize = 22.sp,
                fontWeight = FontWeight.Bold,
                color = Sakura900,
            )
            if (!trip.startDate.isNullOrBlank() || !trip.endDate.isNullOrBlank()) {
                Text(
                    "${trip.startDate ?: "-"} → ${trip.endDate ?: "-"}",
                    fontSize = 12.sp,
                    color = Warm500,
                )
            }
            Spacer(Modifier.height(Spacing.sm))
            Row(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                if (!trip.destination.isNullOrBlank()) {
                    SbChip(text = trip.destination, leadingIcon = Icons.Default.Place)
                }
                SbChip(text = "★ HONEY POT", bg = Sakura500, fg = Color.White)
            }
            Spacer(Modifier.height(Spacing.sm))
            Row(verticalAlignment = Alignment.CenterVertically) {
                LinearProgressIndicator(
                    progress = { 0.35f },
                    color = Sakura500,
                    trackColor = Sakura100,
                    modifier = Modifier
                        .height(5.dp)
                        .weight(1f)
                        .clip(RoundedCornerShape(99.dp)),
                )
                Spacer(Modifier.width(Spacing.sm))
                Text(
                    "준비 35%",
                    fontSize = 11.sp,
                    fontWeight = FontWeight.Bold,
                    color = Sakura700,
                )
            }
        }
    }
}

@Composable
private fun QuickActionRow(
    onCreateTrip: () -> Unit,
    onSchedule: () -> Unit,
    onNotes: () -> Unit,
) {
    val actions = listOf(
        QuickAction("일정", Icons.Default.CalendarMonth, onSchedule),
        QuickAction("메모", Icons.Default.Note, onNotes),
    )
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .horizontalScroll(rememberScrollState())
            .padding(horizontal = Spacing.base),
        horizontalArrangement = Arrangement.spacedBy(Spacing.sm),
    ) {
        Row(
            verticalAlignment = Alignment.CenterVertically,
            modifier = Modifier
                .clip(RoundedCornerShape(99.dp))
                .background(Sakura500)
                .clickable {
                    onCreateTrip()
                }
                .padding(horizontal = 14.dp, vertical = 8.dp),
        ) {
            Icon(Icons.Default.Add, contentDescription = null, tint = Color.White, modifier = Modifier.size(14.dp))
            Spacer(Modifier.width(4.dp))
            Text("새 여행", color = Color.White, fontSize = 12.sp, fontWeight = FontWeight.Bold)
        }
        actions.forEach { a ->
            Row(
                verticalAlignment = Alignment.CenterVertically,
                modifier = Modifier
                    .clip(RoundedCornerShape(99.dp))
                    .background(Color.White)
                    .border(1.dp, Sakura100, RoundedCornerShape(99.dp))
                    .clickable(onClick = a.onClick)
                    .padding(horizontal = 14.dp, vertical = 8.dp),
            ) {
                Icon(a.icon, contentDescription = null, tint = Sakura600, modifier = Modifier.size(14.dp))
                Spacer(Modifier.width(4.dp))
                Text(a.label, color = Sakura700, fontSize = 12.sp, fontWeight = FontWeight.Bold)
            }
        }
    }
}

private data class QuickAction(
    val label: String,
    val icon: ImageVector,
    val onClick: () -> Unit,
)

@Composable
private fun TripRow(trip: Trip, indexHint: Int, onClick: () -> Unit) {
    Card(
        onClick = onClick,
        shape = RoundedCornerShape(18.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp, pressedElevation = 4.dp),
        border = BorderStroke(1.dp, Sakura100),
        modifier = Modifier
            .padding(horizontal = Spacing.base)
            .fillMaxWidth(),
    ) {
        Row(
            modifier = Modifier.padding(Spacing.base),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Box(
                modifier = Modifier
                    .size(56.dp)
                    .clip(RoundedCornerShape(16.dp))
                    .background(Sakura100),
                contentAlignment = Alignment.Center,
            ) {
                when (indexHint % 3) {
                    0 -> SbMiniBear(size = 42.dp, expression = MascotExpression.Happy)
                    1 -> SbChick(size = 42.dp, expression = MascotExpression.Happy)
                    else -> SbBear(size = 42.dp, expression = MascotExpression.Plain)
                }
            }
            Spacer(Modifier.width(Spacing.md))
            Column(modifier = Modifier.weight(1f)) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Text(
                        trip.title.ifBlank { "이름 없는 여행" },
                        fontSize = 14.sp,
                        fontWeight = FontWeight.Bold,
                        color = Sakura900,
                        modifier = Modifier.weight(1f),
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis,
                    )
                    if (!trip.status.isNullOrBlank()) SbStatusPill(trip.status, small = true)
                }
                if (!trip.destination.isNullOrBlank()) {
                    Spacer(Modifier.height(3.dp))
                    InlineIconText(Icons.Default.Place, trip.destination)
                }
                if (!trip.startDate.isNullOrBlank()) {
                    Spacer(Modifier.height(2.dp))
                    InlineIconText(
                        Icons.Default.CalendarMonth,
                        "${trip.startDate} → ${trip.endDate ?: ""}",
                    )
                }
            }
        }
    }
}

@Composable
private fun InlineIconText(icon: ImageVector, text: String, color: Color = Warm500) {
    Row(verticalAlignment = Alignment.CenterVertically) {
        Icon(icon, contentDescription = null, tint = Sakura500, modifier = Modifier.size(11.dp))
        Spacer(Modifier.width(4.dp))
        Text(text, fontSize = 11.sp, color = color)
    }
}

// ──────────────────────────────────────────────────────────────────────
// Trip detail (design: ScreenTripDetail)
// ──────────────────────────────────────────────────────────────────────

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun TripDetailScreen(
    api: JourneyApi,
    tripId: String,
    onBack: () -> Unit,
    onOpenDay: (String) -> Unit,
    onOpenNotes: (String) -> Unit,
    onOpenExpenses: (String) -> Unit,
    onOpenMap: () -> Unit,
    onOpenChecklist: () -> Unit,
    onOpenReservations: () -> Unit,
    onOpenTags: () -> Unit,
    onOpenCompanions: () -> Unit,
    onOpenShare: () -> Unit,
    onOpenMedia: () -> Unit,
    onOpenNearby: () -> Unit,
    onOpenTransit: () -> Unit,
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
            .onSuccess { (t, d) -> trip = t; days = d; loading = false }
            .onFailure { error = it.message; loading = false }
    }

    SakuraScaffold(title = trip?.title ?: "여행 상세", onBack = onBack) { padding ->
        when {
            loading -> CenteredLoader()
            error != null -> ErrorState(error!!)
            else -> Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding)
                    .verticalScroll(rememberScrollState())
                    .padding(bottom = Spacing.xl),
                verticalArrangement = Arrangement.spacedBy(Spacing.md),
            ) {
                trip?.let { TripHeroHeader(it) }
                TripPrepStrip()
                TripDetailNavRow(
                    onSchedule = { /* 같은 화면의 일정 섹션을 노출 — 별도 동작 없음 */ },
                    onMap = onOpenMap,
                    onNotes = { onOpenNotes(tripId) },
                    onExpenses = { onOpenExpenses(tripId) },
                    onChecklist = onOpenChecklist,
                    onReservations = onOpenReservations,
                    onTags = onOpenTags,
                    onCompanions = onOpenCompanions,
                    onShare = onOpenShare,
                    onMedia = onOpenMedia,
                    onNearby = onOpenNearby,
                    onTransit = onOpenTransit,
                )
                SbSection(title = "일정", icon = Icons.Default.CalendarMonth, count = "${days.size}일")
                if (days.isEmpty()) {
                    EmptyCard("등록된 일정이 없습니다.")
                } else {
                    Column(
                        modifier = Modifier.padding(horizontal = Spacing.base),
                        verticalArrangement = Arrangement.spacedBy(Spacing.sm),
                    ) {
                        days.forEach { d -> DayRow(d, onClick = { onOpenDay(d.id) }) }
                    }
                }
            }
        }
    }
}

@Composable
private fun TripHeroHeader(trip: Trip) {
    Box(
        modifier = Modifier
            .padding(horizontal = Spacing.base)
            .fillMaxWidth()
            .clip(RoundedCornerShape(22.dp))
            .background(Brush.linearGradient(listOf(Sakura100, Sakura50, Color.White)))
            .border(1.5.dp, Sakura200, RoundedCornerShape(22.dp)),
    ) {
        Column {
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(120.dp)
                    .background(
                        Brush.linearGradient(listOf(Sakura200, Sakura100, SbBeige)),
                    ),
                contentAlignment = Alignment.Center,
            ) {
                Text(
                    "COVER · SAKURA",
                    fontSize = 10.sp,
                    fontWeight = FontWeight.SemiBold,
                    letterSpacing = 1.sp,
                    color = Sakura700,
                )
            }
            Column(
                modifier = Modifier.padding(Spacing.base),
                verticalArrangement = Arrangement.spacedBy(Spacing.xs),
            ) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Text(
                        trip.title.ifBlank { "이름 없는 여행" },
                        fontSize = 22.sp,
                        fontWeight = FontWeight.Bold,
                        color = Sakura900,
                        modifier = Modifier.weight(1f),
                    )
                    if (!trip.status.isNullOrBlank()) SbStatusPill(trip.status)
                }
                if (!trip.destination.isNullOrBlank()) {
                    Row(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                        SbChip(text = trip.destination, leadingIcon = Icons.Default.Place)
                    }
                }
                if (!trip.startDate.isNullOrBlank()) {
                    Spacer(Modifier.height(Spacing.xs))
                    Text(
                        "${trip.startDate} → ${trip.endDate ?: "-"}",
                        fontSize = 12.sp,
                        color = Warm600,
                    )
                }
                if (!trip.description.isNullOrBlank()) {
                    Spacer(Modifier.height(Spacing.xs))
                    Text(trip.description, fontSize = 13.sp, color = Warm800)
                }
            }
        }
        Box(
            modifier = Modifier
                .align(Alignment.TopEnd)
                .padding(top = 64.dp, end = Spacing.sm),
        ) {
            SbBear(size = 84.dp, expression = MascotExpression.Wink, accessory = MascotAccessory.Flower)
        }
    }
}

@Composable
private fun TripPrepStrip() {
    Row(
        modifier = Modifier
            .padding(horizontal = Spacing.base)
            .fillMaxWidth()
            .clip(RoundedCornerShape(14.dp))
            .background(Color.White)
            .border(1.dp, Sakura100, RoundedCornerShape(14.dp))
            .padding(Spacing.md),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        SbChick(size = 36.dp, expression = MascotExpression.Happy)
        Spacer(Modifier.width(Spacing.sm))
        Column(modifier = Modifier.weight(1f)) {
            Text("여행 준비 35%", fontSize = 12.sp, fontWeight = FontWeight.Bold, color = Sakura900)
            Spacer(Modifier.height(4.dp))
            LinearProgressIndicator(
                progress = { 0.35f },
                color = Sakura500,
                trackColor = Sakura100,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(5.dp)
                    .clip(RoundedCornerShape(99.dp)),
            )
            Spacer(Modifier.height(4.dp))
            Text("항공권 ✓ · 숙소 ✓ · 환전 · 짐싸기", fontSize = 10.sp, color = Warm500)
        }
    }
}

/**
 * 여행 상세 화면 상단의 가로 스크롤 탭 행. 일정/지도/메모/지출 4개 탭을
 * 칩 모양으로 노출한다.
 */
@Composable
private fun TripDetailNavRow(
    onSchedule: () -> Unit,
    onMap: () -> Unit,
    onNotes: () -> Unit,
    onExpenses: () -> Unit,
    onChecklist: () -> Unit,
    onReservations: () -> Unit,
    onTags: () -> Unit,
    onCompanions: () -> Unit,
    onShare: () -> Unit,
    onMedia: () -> Unit,
    onNearby: () -> Unit,
    onTransit: () -> Unit,
) {
    data class Tab(val label: String, val icon: ImageVector, val onClick: () -> Unit)
    val tabs = listOf(
        Tab("일정", Icons.Default.CalendarMonth, onSchedule),
        Tab("지도", Icons.Default.Place, onMap),
        Tab("메모", Icons.Default.Note, onNotes),
        Tab("지출", Icons.Default.Wallet, onExpenses),
        Tab("체크", Icons.Default.CheckCircle, onChecklist),
        Tab("예약", Icons.Default.Flight, onReservations),
        Tab("사진", Icons.Default.CameraAlt, onMedia),
        Tab("동행", Icons.Default.Notifications, onCompanions),
        Tab("태그", Icons.Default.Star, onTags),
        Tab("공유", Icons.Default.Favorite, onShare),
        Tab("주변", Icons.Default.Place, onNearby),
        Tab("교통", Icons.Default.Flight, onTransit),
    )
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .horizontalScroll(rememberScrollState())
            .padding(horizontal = Spacing.base),
        horizontalArrangement = Arrangement.spacedBy(Spacing.sm),
    ) {
        tabs.forEach { t ->
            Row(
                verticalAlignment = Alignment.CenterVertically,
                modifier = Modifier
                    .clip(RoundedCornerShape(99.dp))
                    .background(Color.White)
                    .border(1.dp, Sakura100, RoundedCornerShape(99.dp))
                    .clickable(onClick = t.onClick)
                    .padding(horizontal = 14.dp, vertical = 10.dp),
            ) {
                Icon(t.icon, contentDescription = null, tint = Sakura600, modifier = Modifier.size(14.dp))
                Spacer(Modifier.width(6.dp))
                Text(t.label, color = Sakura700, fontSize = 12.sp, fontWeight = FontWeight.Bold)
            }
        }
    }
}

@Composable
private fun DayRow(d: Day, onClick: () -> Unit) {
    Card(
        onClick = onClick,
        shape = RoundedCornerShape(14.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 0.dp),
        border = BorderStroke(1.dp, Sakura100),
        modifier = Modifier.fillMaxWidth(),
    ) {
        Row(
            modifier = Modifier.padding(Spacing.md),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Box(
                modifier = Modifier
                    .size(44.dp)
                    .clip(CircleShape)
                    .background(Sakura100)
                    .border(1.5.dp, Sakura300, CircleShape),
                contentAlignment = Alignment.Center,
            ) {
                Column(horizontalAlignment = Alignment.CenterHorizontally) {
                    Text("DAY", fontSize = 9.sp, fontWeight = FontWeight.Bold, color = Sakura700)
                    Text("${d.dayNumber}", fontSize = 14.sp, fontWeight = FontWeight.Bold, color = Sakura700)
                }
            }
            Spacer(Modifier.width(Spacing.md))
            Column(modifier = Modifier.weight(1f)) {
                Row(verticalAlignment = Alignment.Bottom) {
                    Text(d.date, fontSize = 13.sp, fontWeight = FontWeight.Bold, color = Warm900)
                    if (!d.region.isNullOrBlank()) {
                        Spacer(Modifier.width(Spacing.xs))
                        Text("· ${d.region}", fontSize = 11.sp, color = Warm500)
                    }
                }
                if (!d.dailySummary.isNullOrBlank()) {
                    Spacer(Modifier.height(2.dp))
                    Text(
                        d.dailySummary,
                        fontSize = 11.sp,
                        color = Warm600,
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis,
                    )
                }
            }
            Icon(
                Icons.AutoMirrored.Filled.ArrowBack,
                contentDescription = null,
                tint = Warm400,
                modifier = Modifier
                    .size(16.dp)
                    .rotate(180f),
            )
        }
    }
}

// ──────────────────────────────────────────────────────────────────────
// Day detail (design: ScreenDayDetail)
// ──────────────────────────────────────────────────────────────────────

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun DayDetailScreen(
    api: JourneyApi,
    dayId: String,
    onBack: () -> Unit,
    onAddSchedule: () -> Unit,
    onEditMeal: () -> Unit,
    onEditAccommodation: () -> Unit,
) {
    var schedules by remember { mutableStateOf<List<Schedule>>(emptyList()) }
    var meals by remember { mutableStateOf<List<Meal>>(emptyList()) }
    var accommodation by remember { mutableStateOf<Accommodation?>(null) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(dayId) {
        loading = true
        val scheduleResult = runCatching { api.listSchedules(dayId).schedules }
        val mealResult = runCatching { api.listMeals(dayId).meals }
        val accommodationResult = runCatching { api.getAccommodation(dayId).accommodation }

        scheduleResult.onSuccess { schedules = it }
        mealResult.onSuccess { meals = it }
        accommodationResult.onSuccess { accommodation = it }

        val firstError = listOfNotNull(
            scheduleResult.exceptionOrNull(),
            mealResult.exceptionOrNull(),
        ).firstOrNull()
        if (schedules.isEmpty() && meals.isEmpty() && accommodation == null && firstError != null) {
            error = firstError.message
        }
        loading = false
    }

    SakuraScaffold(title = "일정 상세", onBack = onBack) { padding ->
        when {
            loading -> CenteredLoader()
            error != null -> ErrorState(error!!)
            else -> Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding)
                    .verticalScroll(rememberScrollState())
                    .padding(bottom = Spacing.xl),
                verticalArrangement = Arrangement.spacedBy(Spacing.md),
            ) {
                DayHeaderCard()

                SbSection(
                    title = "일정",
                    icon = Icons.Default.CalendarMonth,
                    count = "${schedules.size}건",
                    action = {
                        Row(
                            verticalAlignment = Alignment.CenterVertically,
                            modifier = Modifier
                                .clip(RoundedCornerShape(99.dp))
                                .background(Sakura500)
                                .clickable(onClick = onAddSchedule)
                                .padding(horizontal = 12.dp, vertical = 6.dp),
                        ) {
                            Icon(Icons.Default.Add, contentDescription = null, tint = Color.White, modifier = Modifier.size(12.dp))
                            Spacer(Modifier.width(4.dp))
                            Text("일정 추가", color = Color.White, fontSize = 11.sp, fontWeight = FontWeight.Bold)
                        }
                    },
                )
                if (schedules.isEmpty()) EmptyCard("등록된 일정이 없습니다.")
                else ScheduleTimeline(schedules)

                SbSection(title = "식사", icon = Icons.Default.Restaurant, count = "${meals.size}건", action = {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        modifier = Modifier
                            .clip(RoundedCornerShape(99.dp))
                            .background(Sakura500)
                            .clickable(onClick = onEditMeal)
                            .padding(horizontal = 12.dp, vertical = 6.dp),
                    ) {
                        Icon(Icons.Default.Add, contentDescription = null, tint = Color.White, modifier = Modifier.size(12.dp))
                        Spacer(Modifier.width(4.dp))
                        Text("식사 입력", color = Color.White, fontSize = 11.sp, fontWeight = FontWeight.Bold)
                    }
                })
                if (meals.isEmpty()) EmptyCard("식사 기록이 없습니다.")
                else Column(
                    modifier = Modifier.padding(horizontal = Spacing.base),
                    verticalArrangement = Arrangement.spacedBy(Spacing.sm),
                ) {
                    meals.forEach { MealCard(it) }
                }

                SbSection(title = "숙소", icon = Icons.Default.Hotel, action = {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        modifier = Modifier
                            .clip(RoundedCornerShape(99.dp))
                            .background(Sakura500)
                            .clickable(onClick = onEditAccommodation)
                            .padding(horizontal = 12.dp, vertical = 6.dp),
                    ) {
                        Icon(Icons.Default.Add, contentDescription = null, tint = Color.White, modifier = Modifier.size(12.dp))
                        Spacer(Modifier.width(4.dp))
                        Text("숙소 입력", color = Color.White, fontSize = 11.sp, fontWeight = FontWeight.Bold)
                    }
                })
                if (accommodation == null) EmptyCard("등록된 숙소가 없습니다.")
                else AccommodationCard(accommodation!!)
            }
        }
    }
}

@Composable
private fun DayHeaderCard() {
    Box(
        modifier = Modifier
            .padding(horizontal = Spacing.base)
            .fillMaxWidth()
            .clip(RoundedCornerShape(18.dp))
            .background(Brush.linearGradient(listOf(Sakura100, Sakura50)))
            .border(1.5.dp, Sakura200, RoundedCornerShape(18.dp))
            .padding(Spacing.base),
    ) {
        Column {
            Row(verticalAlignment = Alignment.CenterVertically) {
                SbBear(size = 48.dp, expression = MascotExpression.Happy, accessory = MascotAccessory.Camera)
                Spacer(Modifier.width(Spacing.sm))
                Column {
                    Text(
                        "TODAY · 오늘의 여정",
                        fontSize = 11.sp,
                        fontWeight = FontWeight.Bold,
                        letterSpacing = 1.sp,
                        color = Sakura600,
                    )
                    Text(
                        "한 걸음 한 걸음 폭신폭신",
                        fontSize = 16.sp,
                        fontWeight = FontWeight.Bold,
                        color = Sakura900,
                    )
                }
            }
            Spacer(Modifier.height(Spacing.sm))
            Row(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                SbChip(text = "기록 모음", bg = Color.White, fg = Sakura700)
                SbChip(text = "사진 노트", bg = Color.White, fg = Sakura700)
                SbChip(text = "Sakura Bear", bg = Sakura500, fg = Color.White)
            }
        }
    }
}

@Composable
private fun ScheduleTimeline(schedules: List<Schedule>) {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = Spacing.base),
        verticalArrangement = Arrangement.spacedBy(Spacing.sm),
    ) {
        schedules.forEach { ScheduleTimelineRow(it) }
    }
}

@Composable
private fun ScheduleTimelineRow(s: Schedule) {
    Row(modifier = Modifier.fillMaxWidth()) {
        Column(
            modifier = Modifier
                .width(50.dp)
                .padding(top = 12.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            Text(
                s.startTime ?: "—",
                fontSize = 10.sp,
                fontWeight = FontWeight.Bold,
                color = Sakura700,
            )
        }
        Box(
            modifier = Modifier
                .width(20.dp)
                .padding(top = 14.dp),
            contentAlignment = Alignment.TopCenter,
        ) {
            Box(
                modifier = Modifier
                    .size(10.dp)
                    .clip(CircleShape)
                    .background(Sakura500)
                    .border(2.dp, Sakura200, CircleShape),
            )
        }
        Spacer(Modifier.width(Spacing.sm))
        Card(
            shape = RoundedCornerShape(14.dp),
            colors = CardDefaults.cardColors(containerColor = Color.White),
            elevation = CardDefaults.cardElevation(defaultElevation = 0.dp),
            border = BorderStroke(1.dp, Sakura100),
            modifier = Modifier.weight(1f),
        ) {
            Column(
                modifier = Modifier.padding(Spacing.md),
                verticalArrangement = Arrangement.spacedBy(Spacing.xs),
            ) {
                Text(
                    s.title.ifBlank { "무제" },
                    fontSize = 13.sp,
                    fontWeight = FontWeight.Bold,
                    color = Warm900,
                )
                if (!s.placeName.isNullOrBlank()) {
                    InlineIconText(Icons.Default.Place, s.placeName)
                }
                if (!s.notes.isNullOrBlank()) {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clip(RoundedCornerShape(8.dp))
                            .background(Sakura50)
                            .padding(Spacing.sm),
                    ) {
                        Text(s.notes, fontSize = 11.sp, color = Warm700)
                    }
                }
            }
        }
    }
}

@Composable
private fun MealCard(m: Meal) {
    val label = when (m.mealType) {
        "MEAL_TYPE_BREAKFAST" -> "아침"
        "MEAL_TYPE_LUNCH" -> "점심"
        "MEAL_TYPE_DINNER" -> "저녁"
        "MEAL_TYPE_SNACK" -> "간식"
        else -> m.mealType?.ifBlank { null } ?: "식사"
    }
    Card(
        shape = RoundedCornerShape(14.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 0.dp),
        border = BorderStroke(1.dp, Sakura100),
        modifier = Modifier.fillMaxWidth(),
    ) {
        Row(
            modifier = Modifier.padding(Spacing.md),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Box(
                modifier = Modifier
                    .size(40.dp)
                    .clip(RoundedCornerShape(12.dp))
                    .background(Sakura50),
                contentAlignment = Alignment.Center,
            ) {
                Icon(Icons.Default.Restaurant, contentDescription = null, tint = Sakura600, modifier = Modifier.size(18.dp))
            }
            Spacer(Modifier.width(Spacing.md))
            Column(modifier = Modifier.weight(1f)) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    SbChip(text = label, bg = Sakura100, fg = Sakura700, small = true)
                    Spacer(Modifier.width(6.dp))
                    Text(
                        m.restaurantName?.ifBlank { null } ?: m.menu?.ifBlank { null } ?: "(미기록)",
                        fontSize = 13.sp,
                        fontWeight = FontWeight.Bold,
                        color = Warm900,
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis,
                    )
                }
                if (!m.review.isNullOrBlank()) {
                    Spacer(Modifier.height(4.dp))
                    Text(m.review, fontSize = 11.sp, color = Warm500)
                }
            }
            if (m.cost != null) {
                Spacer(Modifier.width(Spacing.sm))
                Text(
                    "${m.cost.amount} ${m.cost.currency}",
                    fontSize = 11.sp,
                    fontWeight = FontWeight.Bold,
                    color = Sakura700,
                )
            }
        }
    }
}

@Composable
private fun AccommodationCard(a: Accommodation) {
    Card(
        shape = RoundedCornerShape(14.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 0.dp),
        border = BorderStroke(1.dp, Sakura100),
        modifier = Modifier
            .padding(horizontal = Spacing.base)
            .fillMaxWidth(),
    ) {
        Column(modifier = Modifier.padding(Spacing.base)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Box(
                    modifier = Modifier
                        .size(56.dp)
                        .clip(RoundedCornerShape(12.dp))
                        .background(Brush.linearGradient(listOf(Sakura100, SbBeige))),
                    contentAlignment = Alignment.Center,
                ) {
                    Icon(Icons.Default.Hotel, contentDescription = null, tint = Sakura700)
                }
                Spacer(Modifier.width(Spacing.md))
                Column(modifier = Modifier.weight(1f)) {
                    Text(a.name, fontSize = 14.sp, fontWeight = FontWeight.Bold, color = Warm900)
                    if (!a.address.isNullOrBlank()) {
                        Spacer(Modifier.height(2.dp))
                        InlineIconText(Icons.Default.Place, a.address)
                    }
                }
            }
            Spacer(Modifier.height(Spacing.sm))
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .clip(RoundedCornerShape(10.dp))
                    .background(Sakura50)
                    .padding(Spacing.sm),
                horizontalArrangement = Arrangement.SpaceBetween,
            ) {
                Text("체크인 ${a.checkInTime ?: "-"}", fontSize = 11.sp, color = Warm700)
                Text("체크아웃 ${a.checkOutTime ?: "-"}", fontSize = 11.sp, color = Warm700)
                if (a.cost != null) {
                    Text(
                        "${a.cost.amount} ${a.cost.currency}",
                        fontSize = 11.sp,
                        fontWeight = FontWeight.Bold,
                        color = Sakura700,
                    )
                }
            }
        }
    }
}

// ──────────────────────────────────────────────────────────────────────
// Notes (design: ScreenNotes)
// ──────────────────────────────────────────────────────────────────────

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun NotesScreen(api: JourneyApi, tripId: String, onBack: () -> Unit, onAdd: () -> Unit) {
    var notes by remember { mutableStateOf<List<Note>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    var moodFilter by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(tripId) {
        loading = true
        runCatching { api.listNotes(tripId).notes }
            .onSuccess { notes = it; loading = false }
            .onFailure { error = it.message; loading = false }
    }

    SakuraScaffold(title = "여행 메모", onBack = onBack) { padding ->
        when {
            loading -> CenteredLoader()
            error != null -> ErrorState(error!!)
            else -> Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding)
                    .verticalScroll(rememberScrollState())
                    .padding(bottom = Spacing.xl),
                verticalArrangement = Arrangement.spacedBy(Spacing.md),
            ) {
                MoodFilterRow(moodFilter, onSelect = { moodFilter = it })

                Row(
                    modifier = Modifier
                        .padding(horizontal = Spacing.base)
                        .fillMaxWidth(),
                    horizontalArrangement = Arrangement.End,
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        modifier = Modifier
                            .clip(RoundedCornerShape(99.dp))
                            .background(Sakura500)
                            .clickable(onClick = onAdd)
                            .padding(horizontal = 12.dp, vertical = 6.dp),
                    ) {
                        Icon(Icons.Default.Add, contentDescription = null, tint = Color.White, modifier = Modifier.size(12.dp))
                        Spacer(Modifier.width(4.dp))
                        Text("메모 추가", color = Color.White, fontSize = 11.sp, fontWeight = FontWeight.Bold)
                    }
                }

                val filtered = if (moodFilter == null) notes
                else notes.filter { it.mood == moodFilter }

                if (filtered.isEmpty()) {
                    EmptyCard("메모가 없습니다.")
                } else {
                    Column(
                        modifier = Modifier.padding(horizontal = Spacing.base),
                        verticalArrangement = Arrangement.spacedBy(Spacing.md),
                    ) {
                        filtered.forEachIndexed { i, n -> StickyNoteCard(n, i) }
                    }
                }

                Spacer(Modifier.height(Spacing.md))
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.Center,
                ) {
                    SbPaw(size = 14.dp, color = Sakura400)
                    Spacer(Modifier.width(8.dp))
                    SbPaw(size = 14.dp, color = Sakura300)
                    Spacer(Modifier.width(8.dp))
                    SbPaw(size = 14.dp, color = Sakura400)
                }
            }
        }
    }
}

@Composable
private fun MoodFilterRow(selected: String?, onSelect: (String?) -> Unit) {
    val moods = listOf("설렘", "맛있음", "평온", "피곤")
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .horizontalScroll(rememberScrollState())
            .padding(horizontal = Spacing.base),
        horizontalArrangement = Arrangement.spacedBy(6.dp),
    ) {
        MoodPill(label = "전체", emoji = null, isActive = selected == null, onClick = { onSelect(null) })
        moods.forEach { m ->
            val meta = moodFor(m)
            MoodPill(
                label = m,
                emoji = meta?.emoji,
                isActive = selected == m,
                onClick = { onSelect(m) },
            )
        }
    }
}

@Composable
private fun MoodPill(label: String, emoji: String?, isActive: Boolean, onClick: () -> Unit) {
    val bg = if (isActive) Sakura500 else Color.White
    val fg = if (isActive) Color.White else Warm700
    Row(
        verticalAlignment = Alignment.CenterVertically,
        modifier = Modifier
            .clip(RoundedCornerShape(99.dp))
            .background(bg)
            .then(
                if (isActive) Modifier
                else Modifier.border(1.dp, Sakura100, RoundedCornerShape(99.dp))
            )
            .clickable(onClick = onClick)
            .padding(horizontal = 12.dp, vertical = 6.dp),
    ) {
        if (emoji != null) {
            Text(emoji, fontSize = 11.sp)
            Spacer(Modifier.width(4.dp))
        }
        Text("#$label", color = fg, fontSize = 11.sp, fontWeight = FontWeight.Bold)
    }
}

@Composable
private fun StickyNoteCard(n: Note, index: Int) {
    val meta = moodFor(n.mood)
    val accent = meta?.color ?: Sakura400
    val rotateDeg = if (index % 2 == 0) -0.6f else 0.4f
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .rotate(rotateDeg)
            .clip(RoundedCornerShape(18.dp))
            .background(Color.White)
            .border(1.dp, Sakura100, RoundedCornerShape(18.dp)),
    ) {
        // Top accent stripe (mood color).
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(4.dp)
                .background(accent),
        )
        Column(modifier = Modifier.padding(Spacing.base)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                if (meta != null) {
                    SbChip(text = "${meta.emoji} ${meta.label}", bg = meta.bg, fg = meta.color, small = true)
                } else if (!n.mood.isNullOrBlank()) {
                    SbChip(text = "#${n.mood}", small = true)
                }
                Spacer(Modifier.weight(1f))
                Icon(Icons.Default.Favorite, contentDescription = null, tint = Sakura400, modifier = Modifier.size(14.dp))
            }
            Spacer(Modifier.height(Spacing.sm))
            Text(n.content, fontSize = 13.sp, color = Warm800, lineHeight = 20.sp)
        }
    }
}

// ──────────────────────────────────────────────────────────────────────
// Expenses (design: ScreenExpenses)
// ──────────────────────────────────────────────────────────────────────

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun ExpensesScreen(api: JourneyApi, tripId: String, onBack: () -> Unit, onAdd: () -> Unit) {
    var expenses by remember { mutableStateOf<List<Expense>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(tripId) {
        loading = true
        runCatching { api.listExpenses(tripId).expenses }
            .onSuccess { expenses = it; loading = false }
            .onFailure { error = it.message; loading = false }
    }

    SakuraScaffold(title = "지출", onBack = onBack) { padding ->
        when {
            loading -> CenteredLoader()
            error != null -> ErrorState(error!!)
            else -> Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding)
                    .verticalScroll(rememberScrollState())
                    .padding(bottom = Spacing.xl),
                verticalArrangement = Arrangement.spacedBy(Spacing.md),
            ) {
                ExpensesTotalCard(expenses)
                Row(
                    modifier = Modifier
                        .padding(horizontal = Spacing.base)
                        .fillMaxWidth(),
                    horizontalArrangement = Arrangement.End,
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        modifier = Modifier
                            .clip(RoundedCornerShape(99.dp))
                            .background(Sakura500)
                            .clickable(onClick = onAdd)
                            .padding(horizontal = 12.dp, vertical = 6.dp),
                    ) {
                        Icon(Icons.Default.Add, contentDescription = null, tint = Color.White, modifier = Modifier.size(12.dp))
                        Spacer(Modifier.width(4.dp))
                        Text("지출 추가", color = Color.White, fontSize = 11.sp, fontWeight = FontWeight.Bold)
                    }
                }
                SbSection(title = "카테고리별")
                CategoryBreakdown(expenses)
                SbSection(title = "기록", count = "${expenses.size}건")
                if (expenses.isEmpty()) EmptyCard("지출 기록이 없습니다.")
                else Column(
                    modifier = Modifier.padding(horizontal = Spacing.base),
                    verticalArrangement = Arrangement.spacedBy(6.dp),
                ) {
                    expenses.forEach { ExpenseRow(it) }
                }
            }
        }
    }
}

@Composable
private fun ExpensesTotalCard(expenses: List<Expense>) {
    val currency = expenses.mapNotNull { it.amount?.currency }
        .groupingBy { it }.eachCount()
        .maxByOrNull { it.value }?.key ?: "JPY"
    val total = expenses.mapNotNull { it.amount }
        .filter { it.currency == currency }
        .sumOf { it.amount }

    Box(
        modifier = Modifier
            .padding(horizontal = Spacing.base)
            .fillMaxWidth()
            .clip(RoundedCornerShape(22.dp))
            .background(Brush.linearGradient(listOf(Sakura500, Sakura400)))
            .padding(Spacing.base),
    ) {
        Column(modifier = Modifier.fillMaxWidth(0.7f)) {
            Text(
                "총 지출",
                color = Color.White.copy(alpha = 0.85f),
                fontSize = 11.sp,
                fontWeight = FontWeight.Bold,
                letterSpacing = 1.sp,
            )
            Spacer(Modifier.height(4.dp))
            Row(verticalAlignment = Alignment.Bottom) {
                Text(
                    "%,d".format(total),
                    color = Color.White,
                    fontSize = 32.sp,
                    fontWeight = FontWeight.Bold,
                )
                Spacer(Modifier.width(6.dp))
                Text(
                    currency,
                    color = Color.White.copy(alpha = 0.85f),
                    fontSize = 14.sp,
                )
            }
            Spacer(Modifier.height(2.dp))
            Text(
                "기록 ${expenses.size}건",
                color = Color.White.copy(alpha = 0.85f),
                fontSize = 12.sp,
            )
        }
        Box(
            modifier = Modifier
                .align(Alignment.BottomEnd)
                .offset(x = 6.dp, y = 6.dp),
        ) {
            SbBear(size = 68.dp, expression = MascotExpression.Surprise)
        }
    }
}

@Composable
private fun CategoryBreakdown(expenses: List<Expense>) {
    val groups = expenses.groupBy { it.category }
    val totals = groups.mapValues { (_, list) ->
        list.mapNotNull { it.amount?.amount?.toDouble() }.sum()
    }.filter { it.value > 0 }
    if (totals.isEmpty()) {
        EmptyCard("아직 카테고리 데이터가 없습니다.")
        return
    }
    val max = totals.values.max()
    Column(
        modifier = Modifier.padding(horizontal = Spacing.base),
        verticalArrangement = Arrangement.spacedBy(Spacing.sm),
    ) {
        totals.entries.sortedByDescending { it.value }.forEach { (cat, value) ->
            val meta = categoryFor(cat)
            Row(verticalAlignment = Alignment.CenterVertically) {
                Box(
                    modifier = Modifier
                        .size(28.dp)
                        .clip(RoundedCornerShape(8.dp))
                        .background(meta.bg),
                    contentAlignment = Alignment.Center,
                ) {
                    Icon(
                        iconForCategory(cat),
                        contentDescription = null,
                        tint = meta.color,
                        modifier = Modifier.size(14.dp),
                    )
                }
                Spacer(Modifier.width(Spacing.sm))
                Column(modifier = Modifier.weight(1f)) {
                    Row {
                        Text(meta.label, fontSize = 12.sp, fontWeight = FontWeight.Bold, modifier = Modifier.weight(1f))
                        Text(
                            "%,d".format(value.toLong()),
                            fontSize = 12.sp,
                            fontWeight = FontWeight.Bold,
                            color = meta.color,
                        )
                    }
                    Spacer(Modifier.height(4.dp))
                    LinearProgressIndicator(
                        progress = { (value / max).toFloat() },
                        color = meta.color,
                        trackColor = Warm100,
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(5.dp)
                            .clip(RoundedCornerShape(99.dp)),
                    )
                }
            }
        }
    }
}

private fun iconForCategory(c: String?): ImageVector = when (c) {
    "EXPENSE_CATEGORY_TRANSPORT" -> Icons.Default.Flight
    "EXPENSE_CATEGORY_FOOD" -> Icons.Default.Restaurant
    "EXPENSE_CATEGORY_LODGING" -> Icons.Default.Hotel
    "EXPENSE_CATEGORY_ACTIVITY" -> Icons.Default.Star
    "EXPENSE_CATEGORY_SHOPPING" -> Icons.Default.ShoppingBag
    else -> Icons.Default.LocalAtm
}

@Composable
private fun ExpenseRow(e: Expense) {
    val meta = categoryFor(e.category)
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(14.dp))
            .background(Color.White)
            .border(1.dp, Sakura100, RoundedCornerShape(14.dp))
            .padding(Spacing.md),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Box(
            modifier = Modifier
                .size(36.dp)
                .clip(RoundedCornerShape(10.dp))
                .background(meta.bg),
            contentAlignment = Alignment.Center,
        ) {
            Icon(
                iconForCategory(e.category),
                contentDescription = null,
                tint = meta.color,
                modifier = Modifier.size(16.dp),
            )
        }
        Spacer(Modifier.width(Spacing.sm))
        Column(modifier = Modifier.weight(1f)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                SbChip(text = meta.label, bg = meta.bg, fg = meta.color, small = true)
            }
            Spacer(Modifier.height(3.dp))
            Text(
                e.description?.ifBlank { null } ?: "(메모 없음)",
                fontSize = 12.sp,
                color = Warm800,
                maxLines = 1,
                overflow = TextOverflow.Ellipsis,
            )
        }
        if (e.amount != null) {
            Spacer(Modifier.width(Spacing.sm))
            Text(
                "${e.amount.amount} ${e.amount.currency}",
                fontSize = 13.sp,
                fontWeight = FontWeight.Bold,
                color = Sakura700,
            )
        }
    }
}

// ──────────────────────────────────────────────────────────────────────
// Shared scaffold / states
// ──────────────────────────────────────────────────────────────────────

@OptIn(ExperimentalMaterial3Api::class)
@Composable
internal fun SakuraScaffold(
    title: String,
    onBack: () -> Unit,
    content: @Composable (PaddingValues) -> Unit,
) {
    Scaffold(
        containerColor = Color.Transparent,
        topBar = {
            TopAppBar(
                title = {
                    Text(title, fontWeight = FontWeight.Bold, color = Sakura700, fontSize = 17.sp)
                },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(
                            Icons.AutoMirrored.Filled.ArrowBack,
                            contentDescription = "뒤로",
                            tint = Sakura600,
                        )
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(containerColor = Color.Transparent),
            )
        },
        content = content,
    )
}

@Composable
internal fun CenteredLoader() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center,
    ) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            BobbingBear()
            Spacer(Modifier.height(Spacing.md))
            Text("꺼내오고 있어요", fontSize = 13.sp, color = Warm600)
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
private fun BobbingBear() {
    val infinite = rememberInfiniteTransition(label = "bobbing-bear")
    val bob by infinite.animateFloat(
        initialValue = 0f,
        targetValue = 6f,
        animationSpec = infiniteRepeatable(tween(1400), RepeatMode.Reverse),
        label = "bob",
    )
    Box(
        modifier = Modifier
            .size(140.dp)
            .clip(CircleShape)
            .background(Brush.radialGradient(listOf(Sakura100, Sakura50))),
        contentAlignment = Alignment.Center,
    ) {
        Box(modifier = Modifier.padding(bottom = bob.dp)) {
            SbBear(size = 96.dp, expression = MascotExpression.Happy, accessory = MascotAccessory.Flower)
        }
    }
}

@Composable
private fun EmptyTrips() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(Spacing.xl),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        SbBear(size = 120.dp, expression = MascotExpression.Surprise, accessory = MascotAccessory.Flower)
        Spacer(Modifier.height(Spacing.md))
        Text(
            "아직 여행이 없어요",
            fontSize = 16.sp,
            fontWeight = FontWeight.Bold,
            color = Sakura900,
        )
        Spacer(Modifier.height(Spacing.xs))
        Text(
            "웹에서 첫 번째 여행을 만들면 여기에 보여요.",
            fontSize = 12.sp,
            color = Warm500,
        )
    }
}

@Composable
internal fun EmptyCard(text: String) {
    Box(
        modifier = Modifier
            .padding(horizontal = Spacing.base)
            .fillMaxWidth()
            .clip(RoundedCornerShape(12.dp))
            .background(Color.White)
            .border(1.dp, Sakura100, RoundedCornerShape(12.dp))
            .padding(Spacing.base),
    ) {
        Text(text, fontSize = 13.sp, color = Warm500)
    }
}

@Composable
internal fun ErrorState(message: String) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(Spacing.xl),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        SbBear(size = 96.dp, expression = MascotExpression.Sleep)
        Spacer(Modifier.height(Spacing.sm))
        Text(
            "잠깐만요",
            fontSize = 16.sp,
            fontWeight = FontWeight.Bold,
            color = Sakura700,
        )
        Spacer(Modifier.height(Spacing.xs))
        Text(message, fontSize = 12.sp, color = Warm500)
    }
}

// ──────────────────────────────────────────────────────────────────────
// Create Trip / Add Schedule (place search 연동)
// ──────────────────────────────────────────────────────────────────────

@Composable
private fun CreateTripScreen(
    api: JourneyApi,
    onBack: () -> Unit,
    onCreated: (String) -> Unit,
) {
    val context = LocalContext.current
    val scope = androidx.compose.runtime.rememberCoroutineScope()
    var title by remember { mutableStateOf("") }
    var destination by remember { mutableStateOf("") }
    var startDate by remember { mutableStateOf("") }
    var endDate by remember { mutableStateOf("") }
    var currency by remember { mutableStateOf("JPY") }
    var submitting by remember { mutableStateOf(false) }

    SakuraScaffold(title = "새 여행", onBack = onBack) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .verticalScroll(rememberScrollState())
                .padding(horizontal = Spacing.base, vertical = Spacing.md),
            verticalArrangement = Arrangement.spacedBy(Spacing.md),
        ) {
            SbField(label = "여행 제목", value = title, onChange = { title = it }, placeholder = "예: 도쿄 사쿠라 5일")
            PlaceSearchField(
                api = api,
                label = "목적지",
                value = destination,
                onChange = { destination = it },
                onSelect = { p -> destination = p.address.ifEmpty { p.name } },
                placeholder = "도시 / 지역 검색 (예: 도쿄, 삿포로)",
            )
            Row(horizontalArrangement = Arrangement.spacedBy(Spacing.sm)) {
                SbField(label = "시작일", value = startDate, onChange = { startDate = it }, placeholder = "YYYY-MM-DD", modifier = Modifier.weight(1f))
                SbField(label = "종료일", value = endDate, onChange = { endDate = it }, placeholder = "YYYY-MM-DD", modifier = Modifier.weight(1f))
            }
            SbField(label = "기준 통화", value = currency, onChange = { currency = it }, placeholder = "JPY / KRW / USD")

            Spacer(Modifier.height(Spacing.sm))
            Button(
                onClick = {
                    if (title.isBlank()) {
                        Toast.makeText(context, "여행 제목을 입력해주세요", Toast.LENGTH_SHORT).show()
                        return@Button
                    }
                    submitting = true
                    scope.launch {
                        runCatching {
                            api.createTrip(
                                CreateTripRequest(
                                    title = title.trim(),
                                    destination = destination.ifBlank { null },
                                    startDate = startDate.ifBlank { null },
                                    endDate = endDate.ifBlank { null },
                                    budgetCurrency = currency.ifBlank { null },
                                ),
                            )
                        }
                            .onSuccess {
                                Toast.makeText(context, "여행이 만들어졌어요", Toast.LENGTH_SHORT).show()
                                onCreated(it.trip.id)
                            }
                            .onFailure { e ->
                                Toast.makeText(context, "생성 실패: ${e.message}", Toast.LENGTH_LONG).show()
                                submitting = false
                            }
                    }
                },
                enabled = !submitting,
                colors = ButtonDefaults.buttonColors(containerColor = Sakura500, contentColor = Color.White),
                modifier = Modifier.fillMaxWidth(),
            ) {
                Text(if (submitting) "만드는 중…" else "여행 만들기", fontWeight = FontWeight.Bold)
            }
        }
    }
}

@Composable
private fun AddScheduleScreen(
    api: JourneyApi,
    dayId: String,
    onBack: () -> Unit,
    onCreated: () -> Unit,
) {
    val context = LocalContext.current
    val scope = androidx.compose.runtime.rememberCoroutineScope()
    var title by remember { mutableStateOf("") }
    var startTime by remember { mutableStateOf("") }
    var endTime by remember { mutableStateOf("") }
    var placeName by remember { mutableStateOf("") }
    var notes by remember { mutableStateOf("") }
    var location by remember { mutableStateOf<GeoPoint?>(null) }
    var submitting by remember { mutableStateOf(false) }

    SakuraScaffold(title = "일정 추가", onBack = onBack) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .verticalScroll(rememberScrollState())
                .padding(horizontal = Spacing.base, vertical = Spacing.md),
            verticalArrangement = Arrangement.spacedBy(Spacing.md),
        ) {
            SbField(label = "제목", value = title, onChange = { title = it }, placeholder = "예: 센소지 참배")
            Row(horizontalArrangement = Arrangement.spacedBy(Spacing.sm)) {
                SbField(label = "시작", value = startTime, onChange = { startTime = it }, placeholder = "HH:mm", modifier = Modifier.weight(1f))
                SbField(label = "종료", value = endTime, onChange = { endTime = it }, placeholder = "HH:mm", modifier = Modifier.weight(1f))
            }
            PlaceSearchField(
                api = api,
                label = "장소",
                value = placeName,
                onChange = {
                    placeName = it
                    // 사용자가 직접 텍스트를 바꾸는 동안에는 좌표 무효화.
                    location = null
                },
                onSelect = { p ->
                    placeName = p.name
                    location = p.location
                },
                placeholder = "식당 / 호텔 / 관광지 검색",
            )
            if (location != null) {
                Text(
                    "위치: %.5f, %.5f".format(location!!.latitude, location!!.longitude),
                    fontSize = 11.sp,
                    color = Warm500,
                )
            }
            SbField(label = "메모", value = notes, onChange = { notes = it }, placeholder = "준비물, 주의사항 등")

            Spacer(Modifier.height(Spacing.sm))
            Button(
                onClick = {
                    if (title.isBlank()) {
                        Toast.makeText(context, "제목을 입력해주세요", Toast.LENGTH_SHORT).show()
                        return@Button
                    }
                    submitting = true
                    scope.launch {
                        runCatching {
                            api.createSchedule(
                                dayId,
                                CreateScheduleRequest(
                                    title = title.trim(),
                                    startTime = startTime.ifBlank { null },
                                    endTime = endTime.ifBlank { null },
                                    placeName = placeName.ifBlank { null },
                                    location = location,
                                    notes = notes.ifBlank { null },
                                ),
                            )
                        }
                            .onSuccess {
                                Toast.makeText(context, "일정이 추가됐어요", Toast.LENGTH_SHORT).show()
                                onCreated()
                            }
                            .onFailure { e ->
                                Toast.makeText(context, "추가 실패: ${e.message}", Toast.LENGTH_LONG).show()
                                submitting = false
                            }
                    }
                },
                enabled = !submitting,
                colors = ButtonDefaults.buttonColors(containerColor = Sakura500, contentColor = Color.White),
                modifier = Modifier.fillMaxWidth(),
            ) {
                Text(if (submitting) "저장 중…" else "일정 저장", fontWeight = FontWeight.Bold)
            }
        }
    }
}

@Composable
internal fun SbField(
    label: String,
    value: String,
    onChange: (String) -> Unit,
    placeholder: String = "",
    modifier: Modifier = Modifier,
) {
    Column(modifier = modifier.fillMaxWidth()) {
        Text(label, fontSize = 12.sp, fontWeight = FontWeight.Bold, color = Sakura700)
        Spacer(Modifier.height(4.dp))
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .clip(RoundedCornerShape(12.dp))
                .background(Color.White)
                .border(1.dp, Sakura100, RoundedCornerShape(12.dp))
                .padding(horizontal = 12.dp, vertical = 10.dp),
        ) {
            if (value.isEmpty() && placeholder.isNotEmpty()) {
                Text(placeholder, color = Warm500, fontSize = 13.sp)
            }
            BasicTextField(
                value = value,
                onValueChange = onChange,
                singleLine = true,
                textStyle = TextStyle(color = Sakura900, fontSize = 13.sp),
                cursorBrush = SolidColor(Sakura500),
                modifier = Modifier.fillMaxWidth(),
            )
        }
    }
}

// ──────────────────────────────────────────────────────────────────────
// Trip map (OpenStreetMap via Leaflet WebView)
// ──────────────────────────────────────────────────────────────────────

/**
 * tripId 의 모든 day → schedule.location 을 모아서 OSM(Leaflet) 지도에
 * 마커로 표시한다. 좌표가 없으면 trip.destination 을 geocode 로 조회해서
 * 중심 좌표로 사용한다. 모두 실패하면 일본 중부를 디폴트로 보여준다.
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun TripMapScreen(
    api: JourneyApi,
    tripId: String,
    onBack: () -> Unit,
) {
    data class MapMarker(val lat: Double, val lng: Double, val title: String)

    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    var markers by remember { mutableStateOf<List<MapMarker>>(emptyList()) }
    var center by remember { mutableStateOf(35.0 to 138.0) }
    var zoom by remember { mutableStateOf(5) }

    LaunchedEffect(tripId) {
        loading = true
        runCatching {
            val trip = api.getTrip(tripId).trip
            val days = api.listDays(tripId).days
            val collected = mutableListOf<MapMarker>()
            for (d in days) {
                val schedules = runCatching { api.listSchedules(d.id).schedules }.getOrDefault(emptyList())
                for (s in schedules) {
                    val loc = s.location
                    if (loc != null && (loc.latitude != 0.0 || loc.longitude != 0.0)) {
                        collected += MapMarker(loc.latitude, loc.longitude, s.placeName ?: s.title)
                    }
                }
            }
            collected to trip
        }
            .onSuccess { (list, trip) ->
                markers = list
                if (list.isNotEmpty()) {
                    center = list.first().lat to list.first().lng
                    zoom = if (list.size > 1) 10 else 13
                } else if (!trip.destination.isNullOrBlank()) {
                    runCatching { api.geocode(trip.destination).places.firstOrNull()?.location }
                        .getOrNull()
                        ?.let { p ->
                            center = p.latitude to p.longitude
                            zoom = 11
                        }
                }
                loading = false
            }
            .onFailure { error = it.message; loading = false }
    }

    SakuraScaffold(title = "지도", onBack = onBack) { padding ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding),
        ) {
            when {
                loading -> CenteredLoader()
                error != null -> ErrorState(error!!)
                else -> {
                    val html = remember(markers, center, zoom) {
                        buildLeafletHtml(
                            centerLat = center.first,
                            centerLng = center.second,
                            zoom = zoom,
                            markers = markers.map { Triple(it.lat, it.lng, it.title) },
                        )
                    }
                    AndroidView(
                        factory = { ctx ->
                            WebView(ctx).apply {
                                settings.javaScriptEnabled = true
                                settings.domStorageEnabled = true
                                settings.cacheMode = WebSettings.LOAD_DEFAULT
                                setBackgroundColor(android.graphics.Color.WHITE)
                            }
                        },
                        update = { wv ->
                            wv.loadDataWithBaseURL(
                                "https://www.openstreetmap.org/",
                                html,
                                "text/html",
                                "utf-8",
                                null,
                            )
                        },
                        modifier = Modifier.fillMaxSize(),
                    )
                    if (markers.isEmpty()) {
                        Text(
                            "좌표가 등록된 일정이 없습니다.\n일정 추가 시 장소 검색에서 선택해주세요.",
                            modifier = Modifier
                                .align(Alignment.TopCenter)
                                .padding(top = Spacing.md, start = Spacing.base, end = Spacing.base)
                                .clip(RoundedCornerShape(12.dp))
                                .background(Color.White.copy(alpha = 0.95f))
                                .border(1.dp, Sakura100, RoundedCornerShape(12.dp))
                                .padding(horizontal = 14.dp, vertical = 10.dp),
                            color = Sakura700,
                            fontSize = 12.sp,
                            fontWeight = FontWeight.SemiBold,
                        )
                    }
                }
            }
        }
    }
}

private fun buildLeafletHtml(
    centerLat: Double,
    centerLng: Double,
    zoom: Int,
    markers: List<Triple<Double, Double, String>>,
): String {
    val markersJs = markers.joinToString(",\n") { (lat, lng, title) ->
        val safe = title.replace("\\", "\\\\").replace("'", "\\'").replace("\n", " ")
        "[${lat}, ${lng}, '${safe}']"
    }
    return """
        <!doctype html>
        <html><head>
        <meta charset="utf-8"/>
        <meta name="viewport" content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=no"/>
        <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"
              integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY=" crossorigin=""/>
        <style>html,body,#map{height:100%;margin:0;padding:0;background:#fff}</style>
        </head><body>
        <div id="map"></div>
        <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"
                integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo=" crossorigin=""></script>
        <script>
          var map = L.map('map').setView([${centerLat}, ${centerLng}], ${zoom});
          L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            maxZoom: 19,
            attribution: '&copy; OpenStreetMap'
          }).addTo(map);
          var pts = [
            ${markersJs}
          ];
          var bounds = [];
          pts.forEach(function(p){
            var m = L.marker([p[0], p[1]]).addTo(map).bindPopup(p[2]);
            bounds.push([p[0], p[1]]);
          });
          if (bounds.length > 1) { map.fitBounds(bounds, {padding:[24,24]}); }
        </script>
        </body></html>
    """.trimIndent()
}
