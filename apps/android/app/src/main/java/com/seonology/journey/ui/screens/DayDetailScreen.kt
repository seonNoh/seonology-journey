package com.seonology.journey.ui.screens

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.Restaurant
import androidx.compose.material.icons.filled.Hotel
import androidx.compose.material.icons.filled.Schedule
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ExposedDropdownMenuBox
import androidx.compose.material3.ExposedDropdownMenuDefaults
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TimePicker
import androidx.compose.material3.rememberTimePickerState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.seonology.journey.ui.theme.Spacing

data class ScheduleItem(
    val id: String,
    val title: String,
    val category: String,
    val startTime: String?,
    val endTime: String?,
)

data class MealSlot(
    val type: String, // breakfast, lunch, dinner
    val name: String,
    val rating: Int = 0,
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DayDetailScreen(
    dayId: String,
    schedules: List<ScheduleItem>,
    meals: List<MealSlot>,
    onAddSchedule: (title: String, category: String) -> Unit,
) {
    var showAddSheet by remember { mutableStateOf(false) }

    Scaffold(
        floatingActionButton = {
            FloatingActionButton(onClick = { showAddSheet = true }) {
                Icon(Icons.Default.Add, contentDescription = "일정 추가")
            }
        },
    ) { padding ->
        LazyColumn(
            modifier = Modifier.fillMaxSize().padding(padding).padding(Spacing.base),
            verticalArrangement = Arrangement.spacedBy(Spacing.sm),
        ) {
            // Meals section
            item {
                Text("식사", style = MaterialTheme.typography.titleMedium)
                Spacer(Modifier.height(Spacing.xs))
            }
            items(meals) { meal ->
                MealCard(meal)
            }

            // Schedule section
            item {
                Spacer(Modifier.height(Spacing.base))
                Text("일정", style = MaterialTheme.typography.titleMedium)
                Spacer(Modifier.height(Spacing.xs))
            }
            items(schedules) { schedule ->
                ScheduleCard(schedule)
            }
        }

        if (showAddSheet) {
            AddScheduleSheet(
                onDismiss = { showAddSheet = false },
                onAdd = { title, cat ->
                    onAddSchedule(title, cat)
                    showAddSheet = false
                },
            )
        }
    }
}

@Composable
private fun MealCard(meal: MealSlot) {
    val icon = when (meal.type) {
        "breakfast" -> Icons.Default.Restaurant
        "lunch" -> Icons.Default.Restaurant
        else -> Icons.Default.Restaurant
    }
    val label = when (meal.type) {
        "breakfast" -> "조식"
        "lunch" -> "중식"
        else -> "석식"
    }
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant),
    ) {
        Row(modifier = Modifier.padding(Spacing.md)) {
            Icon(icon, contentDescription = label)
            Column(modifier = Modifier.padding(start = Spacing.sm)) {
                Text("$label: ${meal.name}", style = MaterialTheme.typography.bodyLarge)
                if (meal.rating > 0) {
                    Text("${"♥".repeat(meal.rating)}${"♡".repeat(5 - meal.rating)}", style = MaterialTheme.typography.bodyMedium)
                }
            }
        }
    }
}

@Composable
private fun ScheduleCard(schedule: ScheduleItem) {
    val icon = when (schedule.category) {
        "accommodation" -> Icons.Default.Hotel
        "meal" -> Icons.Default.Restaurant
        else -> Icons.Default.Schedule
    }
    Card(modifier = Modifier.fillMaxWidth()) {
        Row(modifier = Modifier.padding(Spacing.md)) {
            Icon(icon, contentDescription = schedule.category)
            Column(modifier = Modifier.padding(start = Spacing.sm)) {
                Text(schedule.title, style = MaterialTheme.typography.bodyLarge)
                if (schedule.startTime != null) {
                    Text("${schedule.startTime} ~ ${schedule.endTime ?: ""}", style = MaterialTheme.typography.labelMedium)
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun AddScheduleSheet(onDismiss: () -> Unit, onAdd: (String, String) -> Unit) {
    var title by remember { mutableStateOf("") }
    var category by remember { mutableStateOf("sightseeing") }
    var expanded by remember { mutableStateOf(false) }
    val categories = listOf("sightseeing", "meal", "accommodation", "transport", "activity")

    ModalBottomSheet(onDismissRequest = onDismiss) {
        Column(modifier = Modifier.padding(Spacing.lg)) {
            Text("일정 추가", style = MaterialTheme.typography.headlineSmall)
            Spacer(Modifier.height(Spacing.base))
            OutlinedTextField(
                value = title,
                onValueChange = { title = it },
                label = { Text("제목") },
                modifier = Modifier.fillMaxWidth(),
            )
            Spacer(Modifier.height(Spacing.sm))
            ExposedDropdownMenuBox(expanded = expanded, onExpandedChange = { expanded = it }) {
                OutlinedTextField(
                    value = category,
                    onValueChange = {},
                    readOnly = true,
                    label = { Text("카테고리") },
                    trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded) },
                    modifier = Modifier.fillMaxWidth().menuAnchor(),
                )
                ExposedDropdownMenu(expanded = expanded, onDismissRequest = { expanded = false }) {
                    categories.forEach { cat ->
                        DropdownMenuItem(
                            text = { Text(cat) },
                            onClick = { category = cat; expanded = false },
                        )
                    }
                }
            }
            Spacer(Modifier.height(Spacing.lg))
            Row(horizontalArrangement = Arrangement.End, modifier = Modifier.fillMaxWidth()) {
                TextButton(onClick = onDismiss) { Text("취소") }
                TextButton(onClick = { onAdd(title, category) }, enabled = title.isNotBlank()) { Text("추가") }
            }
        }
    }
}
