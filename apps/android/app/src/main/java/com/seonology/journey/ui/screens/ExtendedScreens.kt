package com.seonology.journey.ui.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.BasicTextField
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.Check
import androidx.compose.material.icons.filled.ContentCopy
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material.icons.filled.Hotel
import androidx.compose.material.icons.filled.Image
import androidx.compose.material.icons.filled.Link
import androidx.compose.material.icons.filled.LocationOn
import androidx.compose.material.icons.filled.MyLocation
import androidx.compose.material.icons.filled.People
import androidx.compose.material.icons.filled.Place
import androidx.compose.material.icons.filled.Restaurant
import androidx.compose.material.icons.filled.Search
import androidx.compose.material.icons.filled.Tag
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.Checkbox
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Slider
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.SolidColor
import androidx.compose.ui.platform.ClipboardManager
import androidx.compose.ui.platform.LocalClipboardManager
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.AnnotatedString
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.seonology.journey.data.AddCompanionRequest
import com.seonology.journey.data.Companion as TripCompanion
import com.seonology.journey.data.ChecklistItem as Chk
import com.seonology.journey.data.CreateChecklistRequest
import com.seonology.journey.data.CreateExpenseRequest
import com.seonology.journey.data.CreateNoteRequest
import com.seonology.journey.data.CreateReservationRequest
import com.seonology.journey.data.CreateShareRequest
import com.seonology.journey.data.CreateTagRequest
import com.seonology.journey.data.JourneyApi
import com.seonology.journey.data.Money
import com.seonology.journey.data.NearbyPlace
import com.seonology.journey.data.Reservation
import com.seonology.journey.data.Share
import com.seonology.journey.data.Tag
import com.seonology.journey.data.TransitRoute
import com.seonology.journey.data.UpdateChecklistRequest
import com.seonology.journey.data.UpdateCompanionRequest
import com.seonology.journey.data.UpsertAccommodationRequest
import com.seonology.journey.data.UpsertMealRequest
import com.seonology.journey.ui.CenteredLoader
import com.seonology.journey.ui.EmptyCard
import com.seonology.journey.ui.ErrorState
import com.seonology.journey.ui.SakuraScaffold
import com.seonology.journey.ui.SbField
import com.seonology.journey.ui.theme.Sakura100
import com.seonology.journey.ui.theme.Sakura400
import com.seonology.journey.ui.theme.Sakura50
import com.seonology.journey.ui.theme.Sakura500
import com.seonology.journey.ui.theme.Sakura600
import com.seonology.journey.ui.theme.Sakura700
import com.seonology.journey.ui.theme.Spacing
import com.seonology.journey.ui.theme.Warm100
import com.seonology.journey.ui.theme.Warm400
import com.seonology.journey.ui.theme.Warm500
import com.seonology.journey.ui.theme.Warm700
import com.seonology.journey.ui.theme.Warm800
import com.seonology.journey.ui.theme.Warm900
import kotlinx.coroutines.launch

/**
 * 본 파일은 웹과 안드로이드 기능 동일성을 위해 추가된 화면들을 모은 모듈이다.
 * - 체크리스트 / 예약 / 태그 / 동행자 / 공유 링크
 * - 주변 검색 / 교통 검색
 * - 식사·숙소 편집 / 지출 추가 / 메모 추가
 * 디자인 시스템(SakuraScaffold, EmptyCard 등)은 JourneyApp.kt 의 internal API
 * 를 재사용한다.
 */

// =======================================================================
// Common helpers
// =======================================================================

@Composable
private fun PrimaryButton(text: String, enabled: Boolean = true, onClick: () -> Unit) {
    Box(
        modifier = Modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(12.dp))
            .background(if (enabled) Sakura500 else Sakura100)
            .clickable(enabled = enabled, onClick = onClick)
            .padding(vertical = 12.dp),
        contentAlignment = Alignment.Center,
    ) {
        Text(text, color = Color.White, fontSize = 14.sp, fontWeight = FontWeight.Bold)
    }
}

@Composable
private fun SmallChip(text: String, bg: Color = Sakura100, fg: Color = Sakura700) {
    Box(
        modifier = Modifier
            .clip(RoundedCornerShape(99.dp))
            .background(bg)
            .padding(horizontal = 10.dp, vertical = 4.dp),
    ) {
        Text(text, color = fg, fontSize = 11.sp, fontWeight = FontWeight.Bold)
    }
}

@Composable
private fun SectionTitle(text: String) {
    Text(
        text,
        modifier = Modifier.padding(horizontal = Spacing.base),
        fontSize = 13.sp,
        fontWeight = FontWeight.Bold,
        color = Warm900,
    )
}

// =======================================================================
// Checklist
// =======================================================================

private val CHECKLIST_CATEGORIES = listOf(
    "CHECKLIST_CATEGORY_PACKING" to "짐",
    "CHECKLIST_CATEGORY_TODO" to "할일",
    "CHECKLIST_CATEGORY_BOOKING" to "예약",
)

@Composable
fun ChecklistFullScreen(api: JourneyApi, tripId: String, onBack: () -> Unit) {
    var items by remember { mutableStateOf<List<Chk>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    var input by remember { mutableStateOf("") }
    var category by remember { mutableStateOf(CHECKLIST_CATEGORIES[0].first) }
    val scope = rememberCoroutineScope()

    suspend fun reload() {
        runCatching { api.listChecklist(tripId).items }
            .onSuccess { items = it; loading = false }
            .onFailure { error = it.message; loading = false }
    }

    LaunchedEffect(tripId) { reload() }

    SakuraScaffold(title = "체크리스트", onBack = onBack) { padding ->
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
                Column(modifier = Modifier.padding(horizontal = Spacing.base)) {
                    CategoryDropdown(
                        options = CHECKLIST_CATEGORIES,
                        selected = category,
                        onSelected = { category = it },
                    )
                    Spacer(Modifier.height(Spacing.sm))
                    SbField(label = "항목", value = input, onChange = { input = it }, placeholder = "추가할 항목")
                    Spacer(Modifier.height(Spacing.sm))
                    PrimaryButton(
                        text = "추가",
                        enabled = input.isNotBlank(),
                        onClick = {
                            scope.launch {
                                runCatching {
                                    api.createChecklistItem(tripId, CreateChecklistRequest(input.trim(), category))
                                }.onSuccess { input = ""; reload() }
                            }
                        },
                    )
                }

                val grouped = items.groupBy { it.category.ifBlank { "CHECKLIST_CATEGORY_UNSPECIFIED" } }
                if (items.isEmpty()) {
                    EmptyCard("항목이 없습니다.")
                } else {
                    grouped.forEach { (cat, list) ->
                        SectionTitle(CHECKLIST_CATEGORIES.firstOrNull { it.first == cat }?.second ?: "기타")
                        Column(
                            modifier = Modifier.padding(horizontal = Spacing.base),
                            verticalArrangement = Arrangement.spacedBy(6.dp),
                        ) {
                            list.forEach { it ->
                                ChecklistRow(
                                    item = it,
                                    onToggle = {
                                        scope.launch {
                                            runCatching {
                                                api.updateChecklistItem(it.id, UpdateChecklistRequest(isChecked = !(it.isChecked ?: false)))
                                            }.onSuccess { reload() }
                                        }
                                    },
                                    onDelete = {
                                        scope.launch {
                                            runCatching { api.deleteChecklistItem(it.id) }.onSuccess { reload() }
                                        }
                                    },
                                )
                            }
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun ChecklistRow(item: Chk, onToggle: () -> Unit, onDelete: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(12.dp))
            .background(Color.White)
            .border(1.dp, Sakura100, RoundedCornerShape(12.dp))
            .padding(horizontal = 8.dp, vertical = 4.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Checkbox(checked = item.isChecked ?: false, onCheckedChange = { onToggle() })
        Text(
            item.item,
            modifier = Modifier.weight(1f),
            fontSize = 13.sp,
            color = if (item.isChecked == true) Warm400 else Warm800,
        )
        IconButton(onClick = onDelete) {
            Icon(Icons.Default.Delete, contentDescription = "삭제", tint = Warm400)
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun CategoryDropdown(
    options: List<Pair<String, String>>,
    selected: String,
    onSelected: (String) -> Unit,
) {
    var expanded by remember { mutableStateOf(false) }
    val current = options.firstOrNull { it.first == selected }?.second ?: selected
    Box {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .clip(RoundedCornerShape(12.dp))
                .background(Color.White)
                .border(1.dp, Sakura100, RoundedCornerShape(12.dp))
                .clickable { expanded = true }
                .padding(horizontal = 12.dp, vertical = 12.dp),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Text(current, fontSize = 13.sp, color = Warm800, modifier = Modifier.weight(1f))
            Text("▾", color = Warm500, fontSize = 12.sp)
        }
        DropdownMenu(expanded = expanded, onDismissRequest = { expanded = false }) {
            options.forEach { (k, v) ->
                DropdownMenuItem(text = { Text(v) }, onClick = {
                    onSelected(k); expanded = false
                })
            }
        }
    }
}

// =======================================================================
// Reservations
// =======================================================================

private val RESERVATION_TYPES = listOf(
    "RESERVATION_TYPE_FLIGHT" to "항공",
    "RESERVATION_TYPE_HOTEL" to "호텔",
    "RESERVATION_TYPE_ACTIVITY" to "액티비티",
    "RESERVATION_TYPE_RESTAURANT" to "식당",
    "RESERVATION_TYPE_TRANSPORT" to "교통",
)

@Composable
fun ReservationsFullScreen(api: JourneyApi, tripId: String, onBack: () -> Unit) {
    var items by remember { mutableStateOf<List<Reservation>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    var showAdd by remember { mutableStateOf(false) }
    val scope = rememberCoroutineScope()

    suspend fun reload() {
        runCatching { api.listReservations(tripId).reservations }
            .onSuccess { items = it; loading = false }
            .onFailure { error = it.message; loading = false }
    }
    LaunchedEffect(tripId) { reload() }

    SakuraScaffold(title = "예약", onBack = onBack) { padding ->
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
                Column(modifier = Modifier.padding(horizontal = Spacing.base)) {
                    PrimaryButton("예약 추가", onClick = { showAdd = true })
                }
                if (items.isEmpty()) EmptyCard("아직 예약이 없습니다.")
                else Column(
                    modifier = Modifier.padding(horizontal = Spacing.base),
                    verticalArrangement = Arrangement.spacedBy(8.dp),
                ) {
                    items.forEach { r ->
                        ReservationCard(r) {
                            scope.launch {
                                runCatching { api.deleteReservation(r.id) }.onSuccess { reload() }
                            }
                        }
                    }
                }
            }
        }
    }

    if (showAdd) {
        ReservationAddDialog(
            onDismiss = { showAdd = false },
            onSubmit = { req ->
                scope.launch {
                    runCatching { api.createReservation(tripId, req) }.onSuccess {
                        showAdd = false; reload()
                    }
                }
            },
        )
    }
}

@Composable
private fun ReservationCard(r: Reservation, onDelete: () -> Unit) {
    val label = RESERVATION_TYPES.firstOrNull { it.first == r.type }?.second ?: "기타"
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(14.dp))
            .background(Color.White)
            .border(1.dp, Sakura100, RoundedCornerShape(14.dp))
            .padding(Spacing.md),
    ) {
        Column(modifier = Modifier.weight(1f)) {
            Text(label, fontSize = 11.sp, color = Warm500, fontWeight = FontWeight.Bold)
            Text(r.vendor ?: "(공급자 없음)", fontSize = 14.sp, fontWeight = FontWeight.Bold, color = Warm900)
            if (!r.confirmNumber.isNullOrBlank()) Text("예약번호 ${r.confirmNumber}", fontSize = 11.sp, color = Warm500)
            if (!r.reservedAt.isNullOrBlank()) Text(r.reservedAt, fontSize = 11.sp, color = Warm500)
            r.cost?.let { Text("${it.amount} ${it.currency}", fontSize = 12.sp, color = Sakura600, fontWeight = FontWeight.Bold) }
            if (!r.notes.isNullOrBlank()) Text(r.notes, fontSize = 12.sp, color = Warm700, maxLines = 2, overflow = TextOverflow.Ellipsis)
        }
        IconButton(onClick = onDelete) {
            Icon(Icons.Default.Delete, contentDescription = "삭제", tint = Warm400)
        }
    }
}

@Composable
private fun ReservationAddDialog(onDismiss: () -> Unit, onSubmit: (CreateReservationRequest) -> Unit) {
    var type by remember { mutableStateOf(RESERVATION_TYPES[0].first) }
    var vendor by remember { mutableStateOf("") }
    var confirmNumber by remember { mutableStateOf("") }
    var reservedAt by remember { mutableStateOf("") }
    var notes by remember { mutableStateOf("") }
    var costAmount by remember { mutableStateOf("") }
    var costCurrency by remember { mutableStateOf("JPY") }

    DialogBox(title = "예약 추가", onDismiss = onDismiss, onConfirm = {
        val cost = costAmount.toLongOrNull()?.let { Money(costCurrency, it) }
        onSubmit(CreateReservationRequest(
            type = type,
            vendor = vendor.ifBlank { null },
            confirmNumber = confirmNumber.ifBlank { null },
            reservedAt = reservedAt.ifBlank { null },
            cost = cost,
            notes = notes.ifBlank { null },
        ))
    }) {
        CategoryDropdown(options = RESERVATION_TYPES, selected = type, onSelected = { type = it })
        Spacer(Modifier.height(Spacing.sm))
        SbField(label = "공급자", value = vendor, onChange = { vendor = it })
        Spacer(Modifier.height(Spacing.sm))
        SbField(label = "예약번호", value = confirmNumber, onChange = { confirmNumber = it })
        Spacer(Modifier.height(Spacing.sm))
        SbField(label = "예약 시각 (ISO)", value = reservedAt, onChange = { reservedAt = it }, placeholder = "2025-09-01T18:00:00Z")
        Spacer(Modifier.height(Spacing.sm))
        Row {
            Box(Modifier.weight(1f)) { SbField(label = "비용", value = costAmount, onChange = { costAmount = it }) }
            Spacer(Modifier.size(Spacing.sm))
            Box(Modifier.weight(1f)) { SbField(label = "통화", value = costCurrency, onChange = { costCurrency = it }) }
        }
        Spacer(Modifier.height(Spacing.sm))
        SbField(label = "메모", value = notes, onChange = { notes = it })
    }
}

// =======================================================================
// Tags
// =======================================================================

@Composable
fun TagsFullScreen(api: JourneyApi, tripId: String, onBack: () -> Unit) {
    var allTags by remember { mutableStateOf<List<Tag>>(emptyList()) }
    var linked by remember { mutableStateOf<Set<String>>(emptySet()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    var name by remember { mutableStateOf("") }
    var color by remember { mutableStateOf("#f9a8d4") }
    val scope = rememberCoroutineScope()

    suspend fun reload() {
        runCatching {
            val all = api.listTags().tags
            val mine = api.listTripTags(tripId).tags.map { it.id }.toSet()
            all to mine
        }.onSuccess { (a, m) -> allTags = a; linked = m; loading = false }
            .onFailure { error = it.message; loading = false }
    }
    LaunchedEffect(tripId) { reload() }

    SakuraScaffold(title = "태그", onBack = onBack) { padding ->
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
                Column(modifier = Modifier.padding(horizontal = Spacing.base)) {
                    SbField(label = "이름", value = name, onChange = { name = it })
                    Spacer(Modifier.height(Spacing.sm))
                    SbField(label = "색상 (#hex)", value = color, onChange = { color = it })
                    Spacer(Modifier.height(Spacing.sm))
                    PrimaryButton("태그 생성", enabled = name.isNotBlank()) {
                        scope.launch {
                            runCatching { api.createTag(CreateTagRequest(name.trim(), color)) }
                                .onSuccess { name = ""; reload() }
                        }
                    }
                }

                SectionTitle("이 여행에 붙은 태그")
                val attached = allTags.filter { it.id in linked }
                if (attached.isEmpty()) EmptyCard("아직 붙은 태그가 없습니다.")
                else TagWrap(tags = attached, attached = linked, onClick = { t ->
                    scope.launch {
                        runCatching { api.detachTripTag(tripId, t.id) }.onSuccess { reload() }
                    }
                }, onDelete = { t ->
                    scope.launch { runCatching { api.deleteTag(t.id) }.onSuccess { reload() } }
                })

                SectionTitle("전체 태그")
                if (allTags.isEmpty()) EmptyCard("등록된 태그가 없습니다.")
                else TagWrap(tags = allTags, attached = linked, onClick = { t ->
                    scope.launch {
                        runCatching {
                            if (t.id in linked) api.detachTripTag(tripId, t.id)
                            else api.attachTripTag(tripId, t.id)
                        }.onSuccess { reload() }
                    }
                }, onDelete = { t ->
                    scope.launch { runCatching { api.deleteTag(t.id) }.onSuccess { reload() } }
                })
            }
        }
    }
}

@OptIn(androidx.compose.foundation.layout.ExperimentalLayoutApi::class)
@Composable
private fun TagWrap(
    tags: List<Tag>,
    attached: Set<String>,
    onClick: (Tag) -> Unit,
    onDelete: (Tag) -> Unit,
) {
    androidx.compose.foundation.layout.FlowRow(
        modifier = Modifier
            .padding(horizontal = Spacing.base)
            .fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(6.dp),
        verticalArrangement = Arrangement.spacedBy(6.dp),
    ) {
        tags.forEach { t ->
            val isAttached = t.id in attached
            val bg = parseColor(t.color) ?: Sakura400
            Row(
                modifier = Modifier
                    .clip(RoundedCornerShape(99.dp))
                    .background(if (isAttached) bg else Warm100)
                    .clickable { onClick(t) }
                    .padding(horizontal = 10.dp, vertical = 6.dp),
                verticalAlignment = Alignment.CenterVertically,
            ) {
                Text(t.name, color = if (isAttached) Color.White else Warm800, fontSize = 11.sp, fontWeight = FontWeight.Bold)
                Spacer(Modifier.size(4.dp))
                Icon(Icons.Default.Delete, contentDescription = "삭제", tint = if (isAttached) Color.White else Warm400, modifier = Modifier
                    .size(12.dp)
                    .clickable { onDelete(t) })
            }
        }
    }
}

private fun parseColor(hex: String?): Color? {
    if (hex.isNullOrBlank()) return null
    val s = hex.removePrefix("#")
    return runCatching {
        when (s.length) {
            6 -> Color(0xFF000000L or s.toLong(16))
            8 -> Color(s.toLong(16))
            else -> null
        }
    }.getOrNull()
}

// =======================================================================
// Companions
// =======================================================================

private val COMPANION_ROLES = listOf(
    "COMPANION_ROLE_EDITOR" to "편집자",
    "COMPANION_ROLE_VIEWER" to "뷰어",
)
private fun roleLabel(role: String): String = when (role) {
    "COMPANION_ROLE_OWNER" -> "소유자"
    "COMPANION_ROLE_EDITOR" -> "편집자"
    "COMPANION_ROLE_VIEWER" -> "뷰어"
    else -> role
}

@Composable
fun CompanionsFullScreen(api: JourneyApi, tripId: String, onBack: () -> Unit) {
    var items by remember { mutableStateOf<List<TripCompanion>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    var memberId by remember { mutableStateOf("") }
    var role by remember { mutableStateOf(COMPANION_ROLES[0].first) }
    val scope = rememberCoroutineScope()

    suspend fun reload() {
        runCatching { api.listCompanions(tripId).companions }
            .onSuccess { items = it; loading = false }
            .onFailure { error = it.message; loading = false }
    }
    LaunchedEffect(tripId) { reload() }

    SakuraScaffold(title = "동행자", onBack = onBack) { padding ->
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
                Column(modifier = Modifier.padding(horizontal = Spacing.base)) {
                    SbField(label = "멤버 ID (Keycloak sub)", value = memberId, onChange = { memberId = it })
                    Spacer(Modifier.height(Spacing.sm))
                    CategoryDropdown(options = COMPANION_ROLES, selected = role, onSelected = { role = it })
                    Spacer(Modifier.height(Spacing.sm))
                    PrimaryButton("동행자 추가", enabled = memberId.isNotBlank()) {
                        scope.launch {
                            runCatching { api.addCompanion(tripId, AddCompanionRequest(memberId.trim(), role)) }
                                .onSuccess { memberId = ""; reload() }
                        }
                    }
                }

                if (items.isEmpty()) EmptyCard("동행자가 없습니다.")
                else Column(
                    modifier = Modifier.padding(horizontal = Spacing.base),
                    verticalArrangement = Arrangement.spacedBy(8.dp),
                ) {
                    items.forEach { c ->
                        Row(
                            modifier = Modifier
                                .fillMaxWidth()
                                .clip(RoundedCornerShape(12.dp))
                                .background(Color.White)
                                .border(1.dp, Sakura100, RoundedCornerShape(12.dp))
                                .padding(Spacing.md),
                            verticalAlignment = Alignment.CenterVertically,
                        ) {
                            Box(
                                modifier = Modifier
                                    .size(32.dp)
                                    .clip(RoundedCornerShape(99.dp))
                                    .background(Sakura100),
                                contentAlignment = Alignment.Center,
                            ) {
                                Text((c.displayName ?: c.memberId).take(1).uppercase(), color = Sakura700, fontWeight = FontWeight.Bold)
                            }
                            Spacer(Modifier.size(Spacing.sm))
                            Column(modifier = Modifier.weight(1f)) {
                                Text(c.displayName ?: c.memberId, fontSize = 13.sp, fontWeight = FontWeight.Bold, color = Warm900)
                                Text(c.memberId, fontSize = 10.sp, color = Warm500)
                            }
                            SmallChip(roleLabel(c.role))
                            if (c.role != "COMPANION_ROLE_OWNER") {
                                IconButton(onClick = {
                                    scope.launch {
                                        runCatching { api.removeCompanion(tripId, c.memberId) }.onSuccess { reload() }
                                    }
                                }) {
                                    Icon(Icons.Default.Delete, contentDescription = "삭제", tint = Warm400)
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

// =======================================================================
// Share Links
// =======================================================================

@Composable
fun ShareFullScreen(api: JourneyApi, tripId: String, onBack: () -> Unit) {
    var shares by remember { mutableStateOf<List<Share>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    var permission by remember { mutableStateOf("COMPANION_ROLE_VIEWER") }
    var hours by remember { mutableStateOf("72") }
    val scope = rememberCoroutineScope()
    val clipboard: ClipboardManager = LocalClipboardManager.current

    suspend fun reload() {
        runCatching { api.listShares(tripId).shares }
            .onSuccess { shares = it; loading = false }
            .onFailure { error = it.message; loading = false }
    }
    LaunchedEffect(tripId) { reload() }

    SakuraScaffold(title = "공유 링크", onBack = onBack) { padding ->
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
                Column(modifier = Modifier.padding(horizontal = Spacing.base)) {
                    CategoryDropdown(
                        options = listOf(
                            "COMPANION_ROLE_VIEWER" to "조회 전용",
                            "COMPANION_ROLE_EDITOR" to "편집 가능",
                        ),
                        selected = permission,
                        onSelected = { permission = it },
                    )
                    Spacer(Modifier.height(Spacing.sm))
                    SbField(label = "만료 (시간)", value = hours, onChange = { hours = it })
                    Spacer(Modifier.height(Spacing.sm))
                    PrimaryButton("링크 생성") {
                        scope.launch {
                            runCatching {
                                api.createShare(tripId, CreateShareRequest(permission, hours.toIntOrNull() ?: 72))
                            }.onSuccess { reload() }
                        }
                    }
                }

                if (shares.isEmpty()) EmptyCard("아직 공유 링크가 없습니다.")
                else Column(
                    modifier = Modifier.padding(horizontal = Spacing.base),
                    verticalArrangement = Arrangement.spacedBy(8.dp),
                ) {
                    shares.forEach { s ->
                        val url = "https://journey.seonology.com/join/${s.code}"
                        Column(
                            modifier = Modifier
                                .fillMaxWidth()
                                .clip(RoundedCornerShape(14.dp))
                                .background(Color.White)
                                .border(1.dp, Sakura100, RoundedCornerShape(14.dp))
                                .padding(Spacing.md),
                        ) {
                            Text(url, fontSize = 12.sp, color = Warm800)
                            Spacer(Modifier.height(4.dp))
                            Text(
                                "권한: ${roleLabel(s.permission)}" + (s.expiresAt?.let { " · 만료 $it" } ?: ""),
                                fontSize = 11.sp,
                                color = Warm500,
                            )
                            Spacer(Modifier.height(8.dp))
                            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                                OutlinedButton(onClick = { clipboard.setText(AnnotatedString(url)) }) {
                                    Icon(Icons.Default.ContentCopy, contentDescription = null, modifier = Modifier.size(14.dp))
                                    Spacer(Modifier.size(4.dp))
                                    Text("복사", fontSize = 11.sp)
                                }
                                OutlinedButton(onClick = {
                                    scope.launch {
                                        runCatching { api.deleteShare(s.code) }.onSuccess { reload() }
                                    }
                                }) {
                                    Icon(Icons.Default.Delete, contentDescription = null, modifier = Modifier.size(14.dp))
                                    Spacer(Modifier.size(4.dp))
                                    Text("삭제", fontSize = 11.sp)
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

// =======================================================================
// Nearby
// =======================================================================

@Composable
fun NearbyFullScreen(api: JourneyApi, tripId: String, onBack: () -> Unit) {
    var lat by remember { mutableStateOf("35.6895") }
    var lng by remember { mutableStateOf("139.6917") }
    var radius by remember { mutableStateOf(1000f) }
    var type by remember { mutableStateOf("restaurant") }
    var results by remember { mutableStateOf<List<NearbyPlace>>(emptyList()) }
    var loading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }
    val scope = rememberCoroutineScope()

    val types = listOf(
        "restaurant" to "식당",
        "cafe" to "카페",
        "tourist_attraction" to "관광지",
        "convenience_store" to "편의점",
        "hotel" to "호텔",
    )

    SakuraScaffold(title = "주변 검색", onBack = onBack) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .verticalScroll(rememberScrollState())
                .padding(bottom = Spacing.xl),
            verticalArrangement = Arrangement.spacedBy(Spacing.md),
        ) {
            Column(modifier = Modifier.padding(horizontal = Spacing.base)) {
                Row {
                    Box(Modifier.weight(1f)) { SbField(label = "위도", value = lat, onChange = { lat = it }) }
                    Spacer(Modifier.size(Spacing.sm))
                    Box(Modifier.weight(1f)) { SbField(label = "경도", value = lng, onChange = { lng = it }) }
                }
                Spacer(Modifier.height(Spacing.sm))
                Text("반경: ${radius.toInt()} m", fontSize = 11.sp, color = Warm500)
                Slider(
                    value = radius,
                    onValueChange = { radius = it },
                    valueRange = 200f..5000f,
                )
                CategoryDropdown(options = types, selected = type, onSelected = { type = it })
                Spacer(Modifier.height(Spacing.sm))
                PrimaryButton(
                    text = if (loading) "검색 중…" else "검색",
                    enabled = !loading,
                ) {
                    val la = lat.toDoubleOrNull(); val lo = lng.toDoubleOrNull()
                    if (la == null || lo == null) return@PrimaryButton
                    loading = true; error = null
                    scope.launch {
                        runCatching { api.nearby(la, lo, radius.toInt(), type).results }
                            .onSuccess { results = it }
                            .onFailure { error = it.message }
                        loading = false
                    }
                }
            }

            error?.let { ErrorState(it) }

            if (results.isEmpty() && !loading) EmptyCard("결과가 없습니다.")
            else Column(
                modifier = Modifier.padding(horizontal = Spacing.base),
                verticalArrangement = Arrangement.spacedBy(8.dp),
            ) {
                results.forEach { p ->
                    Column(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clip(RoundedCornerShape(12.dp))
                            .background(Color.White)
                            .border(1.dp, Sakura100, RoundedCornerShape(12.dp))
                            .padding(Spacing.md),
                    ) {
                        Text(p.name, fontSize = 13.sp, fontWeight = FontWeight.Bold, color = Warm900)
                        Text(p.address, fontSize = 11.sp, color = Warm500)
                        p.rating?.let { Text("평점 ${"%.1f".format(it)}", fontSize = 11.sp, color = Sakura600) }
                    }
                }
            }
        }
    }
}

// =======================================================================
// Transit
// =======================================================================

@Composable
fun TransitFullScreen(api: JourneyApi, tripId: String, onBack: () -> Unit) {
    var origin by remember { mutableStateOf("") }
    var destination by remember { mutableStateOf("") }
    var departure by remember { mutableStateOf("") }
    var routes by remember { mutableStateOf<List<TransitRoute>>(emptyList()) }
    var loading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }
    val scope = rememberCoroutineScope()

    SakuraScaffold(title = "교통 검색", onBack = onBack) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .verticalScroll(rememberScrollState())
                .padding(bottom = Spacing.xl),
            verticalArrangement = Arrangement.spacedBy(Spacing.md),
        ) {
            Column(modifier = Modifier.padding(horizontal = Spacing.base)) {
                SbField(label = "출발지 (lat,lng)", value = origin, onChange = { origin = it }, placeholder = "35.6895,139.6917")
                Spacer(Modifier.height(Spacing.sm))
                SbField(label = "도착지 (lat,lng)", value = destination, onChange = { destination = it }, placeholder = "35.6762,139.6503")
                Spacer(Modifier.height(Spacing.sm))
                SbField(label = "출발 시간 (선택, ISO)", value = departure, onChange = { departure = it })
                Spacer(Modifier.height(Spacing.sm))
                PrimaryButton(if (loading) "검색 중…" else "경로 검색", enabled = !loading) {
                    val (oLat, oLng) = origin.split(",").mapNotNull { it.trim().toDoubleOrNull() }.takeIf { it.size >= 2 }?.let { it[0] to it[1] } ?: return@PrimaryButton
                    val (dLat, dLng) = destination.split(",").mapNotNull { it.trim().toDoubleOrNull() }.takeIf { it.size >= 2 }?.let { it[0] to it[1] } ?: return@PrimaryButton
                    loading = true; error = null
                    scope.launch {
                        runCatching { api.transit(oLat, oLng, dLat, dLng, departure.ifBlank { null }).routes }
                            .onSuccess { routes = it }
                            .onFailure { error = it.message }
                        loading = false
                    }
                }
            }
            error?.let { ErrorState(it) }
            if (routes.isEmpty() && !loading) EmptyCard("경로가 없습니다.")
            else Column(
                modifier = Modifier.padding(horizontal = Spacing.base),
                verticalArrangement = Arrangement.spacedBy(8.dp),
            ) {
                routes.forEachIndexed { idx, r ->
                    Column(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clip(RoundedCornerShape(14.dp))
                            .background(Color.White)
                            .border(1.dp, Sakura100, RoundedCornerShape(14.dp))
                            .padding(Spacing.md),
                    ) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Text(r.summary ?: "경로 ${idx + 1}", modifier = Modifier.weight(1f), fontSize = 13.sp, fontWeight = FontWeight.Bold)
                            Text(r.duration ?: "", fontSize = 12.sp, color = Sakura600)
                        }
                        r.distance?.let { Text(it, fontSize = 11.sp, color = Warm500) }
                        Spacer(Modifier.height(6.dp))
                        r.steps.forEachIndexed { i, s ->
                            Text("${i + 1}. ${s.instruction} (${s.distance})", fontSize = 11.sp, color = Warm700)
                        }
                    }
                }
            }
        }
    }
}

// =======================================================================
// Meal Edit
// =======================================================================

private val MEAL_TYPES = listOf(
    "MEAL_TYPE_BREAKFAST" to "조식",
    "MEAL_TYPE_LUNCH" to "중식",
    "MEAL_TYPE_DINNER" to "석식",
)
private val MEAL_SOURCES = listOf(
    "MEAL_SOURCE_LOCAL" to "현지 식당",
    "MEAL_SOURCE_HOTEL" to "호텔 조식",
    "MEAL_SOURCE_CONVENIENCE" to "편의점",
    "MEAL_SOURCE_SKIP" to "스킵",
)

@Composable
fun MealEditScreen(api: JourneyApi, dayId: String, onBack: () -> Unit, onSaved: () -> Unit) {
    var mealType by remember { mutableStateOf(MEAL_TYPES[0].first) }
    var source by remember { mutableStateOf(MEAL_SOURCES[0].first) }
    var restaurantName by remember { mutableStateOf("") }
    var menu by remember { mutableStateOf("") }
    var costAmount by remember { mutableStateOf("") }
    var costCurrency by remember { mutableStateOf("JPY") }
    var rating by remember { mutableStateOf("") }
    var review by remember { mutableStateOf("") }
    var saving by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }
    val scope = rememberCoroutineScope()

    SakuraScaffold(title = "식사 입력", onBack = onBack) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .verticalScroll(rememberScrollState())
                .padding(horizontal = Spacing.base, vertical = Spacing.md),
            verticalArrangement = Arrangement.spacedBy(Spacing.sm),
        ) {
            CategoryDropdown(MEAL_TYPES, mealType) { mealType = it }
            CategoryDropdown(MEAL_SOURCES, source) { source = it }
            SbField(label = "식당", value = restaurantName, onChange = { restaurantName = it })
            SbField(label = "메뉴", value = menu, onChange = { menu = it })
            Row {
                Box(Modifier.weight(1f)) { SbField(label = "비용", value = costAmount, onChange = { costAmount = it }) }
                Spacer(Modifier.size(Spacing.sm))
                Box(Modifier.weight(1f)) { SbField(label = "통화", value = costCurrency, onChange = { costCurrency = it }) }
            }
            SbField(label = "평점 (1-5)", value = rating, onChange = { rating = it })
            SbField(label = "리뷰", value = review, onChange = { review = it })
            error?.let { Text(it, color = Color.Red, fontSize = 12.sp) }
            PrimaryButton(if (saving) "저장 중…" else "저장", enabled = !saving) {
                saving = true; error = null
                scope.launch {
                    runCatching {
                        api.upsertMeal(dayId, UpsertMealRequest(
                            mealType = mealType,
                            source = source,
                            restaurantName = restaurantName.ifBlank { null },
                            menu = menu.ifBlank { null },
                            cost = costAmount.toLongOrNull()?.let { Money(costCurrency, it) },
                            rating = rating.toIntOrNull(),
                            review = review.ifBlank { null },
                        ))
                    }.onSuccess { onSaved() }
                        .onFailure { error = it.message }
                    saving = false
                }
            }
        }
    }
}

// =======================================================================
// Accommodation Edit
// =======================================================================

@Composable
fun AccommodationEditScreen(api: JourneyApi, dayId: String, onBack: () -> Unit, onSaved: () -> Unit) {
    var name by remember { mutableStateOf("") }
    var checkIn by remember { mutableStateOf("") }
    var checkOut by remember { mutableStateOf("") }
    var costAmount by remember { mutableStateOf("") }
    var costCurrency by remember { mutableStateOf("JPY") }
    var amenities by remember { mutableStateOf("") }
    var address by remember { mutableStateOf("") }
    var saving by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(dayId) {
        runCatching { api.getAccommodation(dayId).accommodation }.onSuccess { a ->
            if (a != null) {
                name = a.name
                checkIn = a.checkInTime ?: ""
                checkOut = a.checkOutTime ?: ""
                costAmount = a.cost?.amount?.toString() ?: ""
                costCurrency = a.cost?.currency ?: "JPY"
                amenities = a.amenities ?: ""
                address = a.address ?: ""
            }
        }
    }

    SakuraScaffold(title = "숙소 입력", onBack = onBack) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .verticalScroll(rememberScrollState())
                .padding(horizontal = Spacing.base, vertical = Spacing.md),
            verticalArrangement = Arrangement.spacedBy(Spacing.sm),
        ) {
            SbField(label = "숙소명", value = name, onChange = { name = it })
            SbField(label = "체크인 (ISO)", value = checkIn, onChange = { checkIn = it })
            SbField(label = "체크아웃 (ISO)", value = checkOut, onChange = { checkOut = it })
            Row {
                Box(Modifier.weight(1f)) { SbField(label = "비용", value = costAmount, onChange = { costAmount = it }) }
                Spacer(Modifier.size(Spacing.sm))
                Box(Modifier.weight(1f)) { SbField(label = "통화", value = costCurrency, onChange = { costCurrency = it }) }
            }
            SbField(label = "편의시설", value = amenities, onChange = { amenities = it })
            SbField(label = "주소", value = address, onChange = { address = it })
            error?.let { Text(it, color = Color.Red, fontSize = 12.sp) }
            PrimaryButton(if (saving) "저장 중…" else "저장", enabled = !saving && name.isNotBlank()) {
                saving = true; error = null
                scope.launch {
                    runCatching {
                        api.upsertAccommodation(dayId, UpsertAccommodationRequest(
                            name = name,
                            checkInTime = checkIn.ifBlank { null },
                            checkOutTime = checkOut.ifBlank { null },
                            cost = costAmount.toLongOrNull()?.let { Money(costCurrency, it) },
                            amenities = amenities.ifBlank { null },
                            address = address.ifBlank { null },
                        ))
                    }.onSuccess { onSaved() }
                        .onFailure { error = it.message }
                    saving = false
                }
            }
        }
    }
}

// =======================================================================
// Expense Add
// =======================================================================

private val EXPENSE_CATEGORIES = listOf(
    "EXPENSE_CATEGORY_TRANSPORT" to "교통",
    "EXPENSE_CATEGORY_FOOD" to "식사",
    "EXPENSE_CATEGORY_LODGING" to "숙박",
    "EXPENSE_CATEGORY_ACTIVITY" to "체험",
    "EXPENSE_CATEGORY_SHOPPING" to "쇼핑",
    "EXPENSE_CATEGORY_OTHER" to "기타",
)

@Composable
fun ExpenseAddScreen(api: JourneyApi, tripId: String, onBack: () -> Unit, onSaved: () -> Unit) {
    var category by remember { mutableStateOf(EXPENSE_CATEGORIES[0].first) }
    var amount by remember { mutableStateOf("") }
    var currency by remember { mutableStateOf("JPY") }
    var description by remember { mutableStateOf("") }
    var saving by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }
    val scope = rememberCoroutineScope()

    SakuraScaffold(title = "지출 추가", onBack = onBack) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .verticalScroll(rememberScrollState())
                .padding(horizontal = Spacing.base, vertical = Spacing.md),
            verticalArrangement = Arrangement.spacedBy(Spacing.sm),
        ) {
            CategoryDropdown(EXPENSE_CATEGORIES, category) { category = it }
            Row {
                Box(Modifier.weight(1f)) { SbField(label = "금액", value = amount, onChange = { amount = it }) }
                Spacer(Modifier.size(Spacing.sm))
                Box(Modifier.weight(1f)) { SbField(label = "통화", value = currency, onChange = { currency = it }) }
            }
            SbField(label = "설명", value = description, onChange = { description = it })
            error?.let { Text(it, color = Color.Red, fontSize = 12.sp) }
            PrimaryButton(if (saving) "저장 중…" else "추가", enabled = !saving && amount.toLongOrNull() != null) {
                saving = true; error = null
                scope.launch {
                    runCatching {
                        api.createExpense(tripId, CreateExpenseRequest(
                            category = category,
                            amount = Money(currency, amount.toLong()),
                            description = description.ifBlank { null },
                        ))
                    }.onSuccess { onSaved() }
                        .onFailure { error = it.message }
                    saving = false
                }
            }
        }
    }
}

// =======================================================================
// Note Add
// =======================================================================

private val NOTE_MOODS = listOf("설렘", "맛있음", "평온", "피곤")

@Composable
fun NoteAddScreen(api: JourneyApi, tripId: String, onBack: () -> Unit, onSaved: () -> Unit) {
    var content by remember { mutableStateOf("") }
    var mood by remember { mutableStateOf<String?>(null) }
    var saving by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }
    val scope = rememberCoroutineScope()

    SakuraScaffold(title = "메모 추가", onBack = onBack) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(horizontal = Spacing.base, vertical = Spacing.md),
            verticalArrangement = Arrangement.spacedBy(Spacing.sm),
        ) {
            SbField(label = "내용", value = content, onChange = { content = it }, placeholder = "여행 메모를 적어보세요")
            Row(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                NOTE_MOODS.forEach { m ->
                    val active = mood == m
                    Box(
                        modifier = Modifier
                            .clip(RoundedCornerShape(99.dp))
                            .background(if (active) Sakura500 else Color.White)
                            .border(1.dp, Sakura100, RoundedCornerShape(99.dp))
                            .clickable { mood = if (active) null else m }
                            .padding(horizontal = 12.dp, vertical = 6.dp),
                    ) {
                        Text("#$m", color = if (active) Color.White else Warm700, fontSize = 11.sp, fontWeight = FontWeight.Bold)
                    }
                }
            }
            error?.let { Text(it, color = Color.Red, fontSize = 12.sp) }
            PrimaryButton(if (saving) "저장 중…" else "저장", enabled = !saving && content.isNotBlank()) {
                saving = true; error = null
                scope.launch {
                    runCatching {
                        api.createNote(tripId, CreateNoteRequest(content = content, mood = mood))
                    }.onSuccess { onSaved() }
                        .onFailure { error = it.message }
                    saving = false
                }
            }
        }
    }
}

// =======================================================================
// Media (basic gallery — no upload yet)
// =======================================================================

@Composable
fun MediaFullScreen(api: JourneyApi, tripId: String, onBack: () -> Unit) {
    var items by remember { mutableStateOf<List<com.seonology.journey.data.Media>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(tripId) {
        runCatching { api.listMedia(tripId).items }
            .onSuccess { items = it; loading = false }
            .onFailure { error = it.message; loading = false }
    }

    SakuraScaffold(title = "사진", onBack = onBack) { padding ->
        when {
            loading -> CenteredLoader()
            error != null -> ErrorState(error!!)
            items.isEmpty() -> EmptyCard("아직 사진이 없습니다.")
            else -> LazyColumn(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding)
                    .padding(horizontal = Spacing.base, vertical = Spacing.md),
                verticalArrangement = Arrangement.spacedBy(8.dp),
            ) {
                items(items) { m ->
                    Column(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clip(RoundedCornerShape(12.dp))
                            .background(Color.White)
                            .border(1.dp, Sakura100, RoundedCornerShape(12.dp))
                            .padding(Spacing.md),
                    ) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Icon(Icons.Default.Image, contentDescription = null, tint = Sakura500)
                            Spacer(Modifier.size(8.dp))
                            Text(m.caption ?: m.s3Key, fontSize = 12.sp, color = Warm800, modifier = Modifier.weight(1f), maxLines = 1, overflow = TextOverflow.Ellipsis)
                        }
                        m.takenAt?.let { Text(it, fontSize = 10.sp, color = Warm500) }
                    }
                }
            }
        }
    }
}

// =======================================================================
// Generic dialog
// =======================================================================

@Composable
private fun DialogBox(
    title: String,
    onDismiss: () -> Unit,
    onConfirm: () -> Unit,
    content: @Composable () -> Unit,
) {
    androidx.compose.ui.window.Dialog(onDismissRequest = onDismiss) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .clip(RoundedCornerShape(16.dp))
                .background(Color.White)
                .padding(Spacing.base),
        ) {
            Text(title, fontSize = 16.sp, fontWeight = FontWeight.Bold, color = Warm900)
            Spacer(Modifier.height(Spacing.sm))
            content()
            Spacer(Modifier.height(Spacing.md))
            Row(horizontalArrangement = Arrangement.End, modifier = Modifier.fillMaxWidth()) {
                OutlinedButton(onClick = onDismiss) { Text("취소") }
                Spacer(Modifier.size(8.dp))
                Box(
                    modifier = Modifier
                        .clip(RoundedCornerShape(8.dp))
                        .background(Sakura500)
                        .clickable(onClick = onConfirm)
                        .padding(horizontal = 16.dp, vertical = 10.dp),
                ) {
                    Text("확인", color = Color.White, fontWeight = FontWeight.Bold)
                }
            }
        }
    }
}
