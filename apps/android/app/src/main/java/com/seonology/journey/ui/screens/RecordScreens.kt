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
import androidx.compose.material.icons.filled.AttachMoney
import androidx.compose.material.icons.filled.Note
import androidx.compose.material.icons.filled.PhotoCamera
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.Checkbox
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SmallFloatingActionButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import com.seonology.journey.ui.theme.Spacing

// --- Expense Screen ---

data class ExpenseItem(
    val id: String,
    val category: String,
    val amount: Double,
    val currency: String,
    val description: String?,
)

@Composable
fun ExpenseScreen(
    expenses: List<ExpenseItem>,
    onAdd: () -> Unit,
) {
    Scaffold(
        floatingActionButton = {
            FloatingActionButton(onClick = onAdd) {
                Icon(Icons.Default.AttachMoney, contentDescription = "지출 추가")
            }
        },
    ) { padding ->
        LazyColumn(
            modifier = Modifier.fillMaxSize().padding(padding).padding(Spacing.base),
            verticalArrangement = Arrangement.spacedBy(Spacing.sm),
        ) {
            item {
                Text("지출", style = MaterialTheme.typography.headlineSmall)
                Spacer(Modifier.height(Spacing.sm))
            }
            items(expenses) { expense ->
                Card(modifier = Modifier.fillMaxWidth()) {
                    Row(
                        modifier = Modifier.padding(Spacing.md).fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween,
                    ) {
                        Column {
                            Text(expense.category, style = MaterialTheme.typography.labelMedium)
                            expense.description?.let {
                                Text(it, style = MaterialTheme.typography.bodyMedium)
                            }
                        }
                        Text(
                            "${expense.currency} ${expense.amount}",
                            style = MaterialTheme.typography.titleMedium,
                        )
                    }
                }
            }
        }
    }
}

// --- Note Screen ---

data class NoteItem(val id: String, val title: String, val content: String?)

@Composable
fun NoteScreen(notes: List<NoteItem>, onAdd: () -> Unit) {
    Scaffold(
        floatingActionButton = {
            FloatingActionButton(onClick = onAdd) {
                Icon(Icons.Default.Note, contentDescription = "메모 추가")
            }
        },
    ) { padding ->
        LazyColumn(
            modifier = Modifier.fillMaxSize().padding(padding).padding(Spacing.base),
            verticalArrangement = Arrangement.spacedBy(Spacing.sm),
        ) {
            items(notes) { note ->
                Card(modifier = Modifier.fillMaxWidth()) {
                    Column(modifier = Modifier.padding(Spacing.md)) {
                        Text(note.title, style = MaterialTheme.typography.titleMedium)
                        note.content?.let {
                            Text(it, style = MaterialTheme.typography.bodyMedium, maxLines = 3)
                        }
                    }
                }
            }
        }
    }
}

// --- Checklist Screen ---

data class CheckItem(val id: String, val text: String, val checked: Boolean)

@Composable
fun ChecklistScreen(items: List<CheckItem>, onToggle: (String) -> Unit) {
    LazyColumn(
        modifier = Modifier.fillMaxSize().padding(Spacing.base),
        verticalArrangement = Arrangement.spacedBy(Spacing.xs),
    ) {
        item { Text("체크리스트", style = MaterialTheme.typography.headlineSmall) }
        items(items) { item ->
            Card(
                modifier = Modifier.fillMaxWidth(),
                colors = CardDefaults.cardColors(
                    containerColor = if (item.checked) MaterialTheme.colorScheme.surfaceVariant
                    else MaterialTheme.colorScheme.surface,
                ),
            ) {
                Row(
                    modifier = Modifier.padding(Spacing.sm),
                    verticalAlignment = Alignment.CenterVertically,
                ) {
                    Checkbox(checked = item.checked, onCheckedChange = { onToggle(item.id) })
                    Text(item.text, style = MaterialTheme.typography.bodyLarge)
                }
            }
        }
    }
}

// --- FAB Quick Record ---

@Composable
fun QuickRecordFab(
    onPhoto: () -> Unit,
    onNote: () -> Unit,
    onExpense: () -> Unit,
) {
    var expanded by remember { mutableStateOf(false) }
    Column(horizontalAlignment = Alignment.End) {
        if (expanded) {
            SmallFloatingActionButton(onClick = { onPhoto(); expanded = false }) {
                Icon(Icons.Default.PhotoCamera, contentDescription = "사진")
            }
            Spacer(Modifier.height(Spacing.xs))
            SmallFloatingActionButton(onClick = { onNote(); expanded = false }) {
                Icon(Icons.Default.Note, contentDescription = "메모")
            }
            Spacer(Modifier.height(Spacing.xs))
            SmallFloatingActionButton(onClick = { onExpense(); expanded = false }) {
                Icon(Icons.Default.AttachMoney, contentDescription = "지출")
            }
            Spacer(Modifier.height(Spacing.sm))
        }
        FloatingActionButton(onClick = { expanded = !expanded }) {
            Icon(Icons.Default.Add, contentDescription = "퀵 기록")
        }
    }
}
