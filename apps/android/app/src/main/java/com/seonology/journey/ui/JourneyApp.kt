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
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.CalendarMonth
import androidx.compose.material.icons.filled.Flight
import androidx.compose.material.icons.filled.Hotel
import androidx.compose.material.icons.filled.Image
import androidx.compose.material.icons.filled.ListAlt
import androidx.compose.material.icons.filled.Logout
import androidx.compose.material.icons.filled.Note
import androidx.compose.material.icons.filled.Place
import androidx.compose.material.icons.filled.Restaurant
import androidx.compose.material.icons.filled.Wallet
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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.seonology.journey.R
import com.seonology.journey.auth.AuthStore
import com.seonology.journey.auth.KeycloakAuth
import com.seonology.journey.data.Accommodation
import com.seonology.journey.data.Day
import com.seonology.journey.data.Expense
import com.seonology.journey.data.JourneyApi
import com.seonology.journey.data.Meal
import com.seonology.journey.data.Network
import com.seonology.journey.data.Note
import com.seonology.journey.data.Schedule
import com.seonology.journey.data.Trip
import com.seonology.journey.ui.theme.Sakura100
import com.seonology.journey.ui.theme.Sakura200
import com.seonology.journey.ui.theme.Sakura50
import com.seonology.journey.ui.theme.Sakura500
import com.seonology.journey.ui.theme.Sakura600
import com.seonology.journey.ui.theme.Sakura700
import com.seonology.journey.ui.theme.Spacing
import com.seonology.journey.ui.theme.Warm50
import net.openid.appauth.AuthorizationException
import net.openid.appauth.AuthorizationResponse

/**
 * 앱 전체를 감싸는 엔트리. 인증 상태 → 목록 → 상세 → Day 상세 → 서브
 * 리스트(일정/식사/숙소/메모/지출) 순서의 단순한 스택 네비게이션을
 * 직접 구성한다. 향후 한 화면 안에서 여러 탭으로 전환해야 한다면
 * Scaffold bottom bar 로 쪼개면 된다.
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
            // 로그인 실패는 조용히 무시하고 사용자가 다시 시도하도록 둔다.
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
            .background(Brush.verticalGradient(listOf(Sakura50, Warm50))),
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
                    onLogout = {
                        store.clear()
                        authed = false
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
                )
            }
            composable("days/{dayId}") { entry ->
                val dayId = entry.arguments?.getString("dayId").orEmpty()
                DayDetailScreen(
                    api = api,
                    dayId = dayId,
                    onBack = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/notes") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                NotesScreen(
                    api = api,
                    tripId = tripId,
                    onBack = { nav.popBackStack() },
                )
            }
            composable("trips/{tripId}/expenses") { entry ->
                val tripId = entry.arguments?.getString("tripId").orEmpty()
                ExpensesScreen(
                    api = api,
                    tripId = tripId,
                    onBack = { nav.popBackStack() },
                )
            }
        }
    }
}

// ──────────────────────────────────────────────────────────────────────
// Login
// ──────────────────────────────────────────────────────────────────────

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
            .background(Brush.radialGradient(listOf(Sakura100, Sakura50)), CircleShape),
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

// ──────────────────────────────────────────────────────────────────────
// Trip list
// ──────────────────────────────────────────────────────────────────────

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
            .onSuccess { trips = it.trips; loading = false }
            .onFailure { error = it.message; loading = false }
    }

    Scaffold(
        containerColor = Color.Transparent,
        topBar = {
            TopAppBar(
                title = {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Icon(
                            Icons.Default.Flight,
                            contentDescription = null,
                            tint = Sakura600,
                            modifier = Modifier.size(20.dp).rotate(-20f),
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
                colors = TopAppBarDefaults.topAppBarColors(containerColor = Color.Transparent),
            )
        },
    ) { padding ->
        when {
            loading -> CenteredLoader(padding)
            error != null -> ErrorState(padding, error!!)
            trips.isEmpty() -> EmptyTrips(padding)
            else -> LazyColumn(
                modifier = Modifier.fillMaxSize().padding(padding).padding(horizontal = Spacing.base),
                verticalArrangement = Arrangement.spacedBy(Spacing.md),
                contentPadding = PaddingValues(vertical = Spacing.md),
            ) {
                items(trips, key = { it.id }) { t ->
                    TripCard(
                        trip = t,
                        onClick = { onOpenTrip(t.id) },
                    )
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
                if (!trip.status.isNullOrBlank()) StatusPill(trip.status)
            }
            if (!trip.destination.isNullOrBlank()) InlineIconText(Icons.Default.Place, trip.destination)
            if (!trip.startDate.isNullOrBlank()) {
                InlineIconText(
                    Icons.Default.CalendarMonth,
                    "${trip.startDate} → ${trip.endDate ?: ""}",
                )
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

@Composable
private fun InlineIconText(icon: ImageVector, text: String) {
    Row(verticalAlignment = Alignment.CenterVertically) {
        Icon(icon, contentDescription = null, tint = Sakura500, modifier = Modifier.size(16.dp))
        Spacer(Modifier.width(Spacing.xs))
        Text(
            text,
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
        )
    }
}

// ──────────────────────────────────────────────────────────────────────
// Trip detail
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

    TopBarScaffold(
        title = trip?.title ?: "여행 상세",
        onBack = onBack,
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
                trip?.let { TripHeaderCard(it) }

                SectionNavRow(
                    QuickNav("메모", Icons.Default.Note) { onOpenNotes(tripId) },
                    QuickNav("지출", Icons.Default.Wallet) { onOpenExpenses(tripId) },
                )

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
                    days.forEach { d ->
                        DayRow(d, onClick = { onOpenDay(d.id) })
                    }
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
            if (!trip.destination.isNullOrBlank()) InlineIconText(Icons.Default.Place, trip.destination)
            if (!trip.startDate.isNullOrBlank()) {
                InlineIconText(
                    Icons.Default.CalendarMonth,
                    "${trip.startDate} → ${trip.endDate ?: ""}",
                )
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

private data class QuickNav(val label: String, val icon: ImageVector, val onClick: () -> Unit)

@Composable
private fun SectionNavRow(vararg items: QuickNav) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(Spacing.sm),
    ) {
        items.forEach { nav ->
            Card(
                onClick = nav.onClick,
                shape = RoundedCornerShape(14.dp),
                colors = CardDefaults.cardColors(containerColor = Sakura100),
                modifier = Modifier.weight(1f).height(72.dp),
            ) {
                Column(
                    modifier = Modifier.fillMaxSize().padding(Spacing.sm),
                    verticalArrangement = Arrangement.Center,
                    horizontalAlignment = Alignment.CenterHorizontally,
                ) {
                    Icon(nav.icon, contentDescription = null, tint = Sakura700)
                    Spacer(Modifier.height(Spacing.xs))
                    Text(
                        nav.label,
                        style = MaterialTheme.typography.labelLarge,
                        fontWeight = FontWeight.SemiBold,
                        color = Sakura700,
                    )
                }
            }
        }
    }
}

@Composable
private fun DayRow(d: Day, onClick: () -> Unit) {
    Card(
        onClick = onClick,
        shape = RoundedCornerShape(14.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp, pressedElevation = 4.dp),
        modifier = Modifier.fillMaxWidth(),
    ) {
        Row(
            modifier = Modifier.padding(Spacing.base),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            FilledIconButton(
                onClick = onClick,
                shape = CircleShape,
                colors = IconButtonDefaults.filledIconButtonColors(
                    containerColor = Sakura100,
                    contentColor = Sakura700,
                ),
                modifier = Modifier.size(44.dp),
            ) {
                Text("${d.dayNumber}", fontWeight = FontWeight.Bold, fontSize = 14.sp)
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
                if (!d.dailySummary.isNullOrBlank()) {
                    Text(
                        d.dailySummary,
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                        maxLines = 1,
                    )
                }
            }
        }
    }
}

// ──────────────────────────────────────────────────────────────────────
// Day detail
// ──────────────────────────────────────────────────────────────────────

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun DayDetailScreen(
    api: JourneyApi,
    dayId: String,
    onBack: () -> Unit,
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

        // 모두 실패했을 때만 에러 화면. 일부 404 (예: 숙소 미등록) 는 정상.
        val firstError = listOfNotNull(
            scheduleResult.exceptionOrNull(),
            mealResult.exceptionOrNull(),
        ).firstOrNull()
        if (schedules.isEmpty() && meals.isEmpty() && accommodation == null && firstError != null) {
            error = firstError.message
        }
        loading = false
    }

    TopBarScaffold(title = "일정 상세", onBack = onBack) { padding ->
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
                SectionTitle(Icons.Default.ListAlt, "일정 (${schedules.size}건)")
                if (schedules.isEmpty()) EmptyCard("등록된 일정이 없습니다.")
                else schedules.forEach { ScheduleCard(it) }

                SectionTitle(Icons.Default.Restaurant, "식사 (${meals.size}건)")
                if (meals.isEmpty()) EmptyCard("식사 기록이 없습니다.")
                else meals.forEach { MealCard(it) }

                SectionTitle(Icons.Default.Hotel, "숙소")
                if (accommodation == null) EmptyCard("등록된 숙소가 없습니다.")
                else AccommodationCard(accommodation!!)

                Spacer(Modifier.height(Spacing.xl))
            }
        }
    }
}

@Composable
private fun SectionTitle(icon: ImageVector, text: String) {
    Row(verticalAlignment = Alignment.CenterVertically) {
        Icon(icon, contentDescription = null, tint = Sakura600, modifier = Modifier.size(20.dp))
        Spacer(Modifier.width(Spacing.sm))
        Text(
            text,
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.onSurface,
        )
    }
}

@Composable
private fun ScheduleCard(s: Schedule) {
    Card(
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp),
        modifier = Modifier.fillMaxWidth(),
    ) {
        Column(
            modifier = Modifier.padding(Spacing.base),
            verticalArrangement = Arrangement.spacedBy(Spacing.xs),
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                val time = listOfNotNull(s.startTime, s.endTime).joinToString(" - ")
                if (time.isNotBlank()) {
                    Text(
                        time,
                        style = MaterialTheme.typography.labelSmall,
                        color = Sakura600,
                        fontWeight = FontWeight.SemiBold,
                    )
                    Spacer(Modifier.width(Spacing.sm))
                }
                Text(
                    s.title.ifBlank { "무제" },
                    style = MaterialTheme.typography.bodyLarge,
                    fontWeight = FontWeight.SemiBold,
                    modifier = Modifier.weight(1f),
                )
            }
            if (!s.placeName.isNullOrBlank()) {
                InlineIconText(Icons.Default.Place, s.placeName)
            }
            if (!s.notes.isNullOrBlank()) {
                Text(
                    s.notes,
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant,
                )
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
        else -> m.mealType
    }
    Card(
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp),
        modifier = Modifier.fillMaxWidth(),
    ) {
        Column(
            modifier = Modifier.padding(Spacing.base),
            verticalArrangement = Arrangement.spacedBy(Spacing.xs),
        ) {
            Text(
                label,
                style = MaterialTheme.typography.labelMedium,
                color = Sakura600,
                fontWeight = FontWeight.SemiBold,
            )
            Text(
                m.restaurantName?.ifBlank { null } ?: m.menu?.ifBlank { null } ?: "(미기록)",
                style = MaterialTheme.typography.bodyLarge,
                fontWeight = FontWeight.SemiBold,
            )
            if (m.cost != null) {
                Text(
                    "${m.cost.amount} ${m.cost.currency}",
                    style = MaterialTheme.typography.bodySmall,
                    color = Sakura700,
                )
            }
            if (!m.review.isNullOrBlank()) {
                Text(
                    m.review,
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant,
                )
            }
        }
    }
}

@Composable
private fun AccommodationCard(a: Accommodation) {
    Card(
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp),
        modifier = Modifier.fillMaxWidth(),
    ) {
        Column(
            modifier = Modifier.padding(Spacing.base),
            verticalArrangement = Arrangement.spacedBy(Spacing.xs),
        ) {
            Text(a.name, style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.Bold)
            if (!a.address.isNullOrBlank()) InlineIconText(Icons.Default.Place, a.address)
            val checkin = a.checkInTime.orEmpty()
            val checkout = a.checkOutTime.orEmpty()
            if (checkin.isNotBlank() || checkout.isNotBlank()) {
                Text(
                    "체크인 ${checkin.ifBlank { "-" }} / 체크아웃 ${checkout.ifBlank { "-" }}",
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant,
                )
            }
            if (a.cost != null) {
                Text(
                    "${a.cost.amount} ${a.cost.currency}",
                    style = MaterialTheme.typography.bodyMedium,
                    color = Sakura700,
                )
            }
        }
    }
}

// ──────────────────────────────────────────────────────────────────────
// Notes / Expenses lists
// ──────────────────────────────────────────────────────────────────────

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun NotesScreen(api: JourneyApi, tripId: String, onBack: () -> Unit) {
    var notes by remember { mutableStateOf<List<Note>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(tripId) {
        loading = true
        runCatching { api.listNotes(tripId).notes }
            .onSuccess { notes = it; loading = false }
            .onFailure { error = it.message; loading = false }
    }

    TopBarScaffold(title = "메모", onBack = onBack) { padding ->
        when {
            loading -> CenteredLoader(padding)
            error != null -> ErrorState(padding, error!!)
            notes.isEmpty() -> EmptyCenter("메모가 없습니다.", padding)
            else -> LazyColumn(
                modifier = Modifier.fillMaxSize().padding(padding).padding(horizontal = Spacing.base),
                contentPadding = PaddingValues(vertical = Spacing.md),
                verticalArrangement = Arrangement.spacedBy(Spacing.sm),
            ) {
                items(notes, key = { it.id }) { n ->
                    Card(
                        shape = RoundedCornerShape(12.dp),
                        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
                        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp),
                        modifier = Modifier.fillMaxWidth(),
                    ) {
                        Column(
                            modifier = Modifier.padding(Spacing.base),
                            verticalArrangement = Arrangement.spacedBy(Spacing.xs),
                        ) {
                            if (!n.mood.isNullOrBlank()) {
                                Text(
                                    "#${n.mood}",
                                    style = MaterialTheme.typography.labelSmall,
                                    color = Sakura600,
                                )
                            }
                            Text(
                                n.content,
                                style = MaterialTheme.typography.bodyMedium,
                                color = MaterialTheme.colorScheme.onSurface,
                            )
                        }
                    }
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun ExpensesScreen(api: JourneyApi, tripId: String, onBack: () -> Unit) {
    var expenses by remember { mutableStateOf<List<Expense>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(tripId) {
        loading = true
        runCatching { api.listExpenses(tripId).expenses }
            .onSuccess { expenses = it; loading = false }
            .onFailure { error = it.message; loading = false }
    }

    TopBarScaffold(title = "지출", onBack = onBack) { padding ->
        when {
            loading -> CenteredLoader(padding)
            error != null -> ErrorState(padding, error!!)
            expenses.isEmpty() -> EmptyCenter("지출 기록이 없습니다.", padding)
            else -> LazyColumn(
                modifier = Modifier.fillMaxSize().padding(padding).padding(horizontal = Spacing.base),
                contentPadding = PaddingValues(vertical = Spacing.md),
                verticalArrangement = Arrangement.spacedBy(Spacing.sm),
            ) {
                items(expenses, key = { it.id }) { e ->
                    Card(
                        shape = RoundedCornerShape(12.dp),
                        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
                        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp),
                        modifier = Modifier.fillMaxWidth(),
                    ) {
                        Row(
                            modifier = Modifier.padding(Spacing.base),
                            verticalAlignment = Alignment.CenterVertically,
                        ) {
                            Column(modifier = Modifier.weight(1f)) {
                                Text(
                                    expenseCategoryLabel(e.category),
                                    style = MaterialTheme.typography.labelSmall,
                                    color = Sakura600,
                                )
                                Text(
                                    e.description.orEmpty().ifBlank { "(메모 없음)" },
                                    style = MaterialTheme.typography.bodyMedium,
                                    color = MaterialTheme.colorScheme.onSurface,
                                )
                            }
                            if (e.amount != null) {
                                Text(
                                    "${e.amount.amount} ${e.amount.currency}",
                                    style = MaterialTheme.typography.titleMedium,
                                    fontWeight = FontWeight.Bold,
                                    color = Sakura700,
                                )
                            }
                        }
                    }
                }
            }
        }
    }
}

private fun expenseCategoryLabel(c: String): String = when (c) {
    "EXPENSE_CATEGORY_TRANSPORT" -> "교통"
    "EXPENSE_CATEGORY_FOOD" -> "식사"
    "EXPENSE_CATEGORY_LODGING" -> "숙박"
    "EXPENSE_CATEGORY_ACTIVITY" -> "체험"
    "EXPENSE_CATEGORY_SHOPPING" -> "쇼핑"
    "EXPENSE_CATEGORY_OTHER" -> "기타"
    else -> c
}

// ──────────────────────────────────────────────────────────────────────
// Scaffold / shared states
// ──────────────────────────────────────────────────────────────────────

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun TopBarScaffold(
    title: String,
    onBack: () -> Unit,
    content: @Composable (PaddingValues) -> Unit,
) {
    Scaffold(
        containerColor = Color.Transparent,
        topBar = {
            TopAppBar(
                title = {
                    Text(title, fontWeight = FontWeight.Bold, color = Sakura700)
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
private fun CenteredLoader(padding: PaddingValues) {
    Box(
        modifier = Modifier.fillMaxSize().padding(padding),
        contentAlignment = Alignment.Center,
    ) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            CuteHero()
            Spacer(Modifier.height(Spacing.md))
            Text(
                "꺼내오고 있어요",
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
        modifier = Modifier.fillMaxSize().padding(padding).padding(Spacing.xl),
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
private fun EmptyCard(text: String) {
    Card(
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 0.dp),
        modifier = Modifier.fillMaxWidth(),
    ) {
        Text(
            text,
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            modifier = Modifier.padding(Spacing.base),
        )
    }
}

@Composable
private fun EmptyCenter(text: String, padding: PaddingValues) {
    Box(
        modifier = Modifier.fillMaxSize().padding(padding),
        contentAlignment = Alignment.Center,
    ) {
        Text(
            text,
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
        )
    }
}

@Composable
private fun ErrorState(padding: PaddingValues, message: String) {
    Column(
        modifier = Modifier.fillMaxSize().padding(padding).padding(Spacing.xl),
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
