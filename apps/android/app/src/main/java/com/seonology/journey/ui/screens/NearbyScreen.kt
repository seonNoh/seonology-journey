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
import androidx.compose.material.icons.filled.Map
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilterChip
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import com.seonology.journey.ui.theme.Spacing

data class NearbyPlace(
    val name: String,
    val address: String,
    val rating: Float?,
    val type: String,
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun NearbyScreen(
    places: List<NearbyPlace>,
    onCategoryToggle: (String) -> Unit,
    onNavigate: (NearbyPlace) -> Unit,
) {
    val categories = listOf("restaurant", "cafe", "hotel", "convenience_store")
    var selectedCategories by remember { mutableStateOf(setOf("restaurant")) }

    Column(modifier = Modifier.fillMaxSize().padding(Spacing.base)) {
        Text("주변 추천", style = MaterialTheme.typography.headlineSmall)
        Spacer(Modifier.height(Spacing.sm))

        // Category chips
        Row(horizontalArrangement = Arrangement.spacedBy(Spacing.xs)) {
            categories.forEach { cat ->
                val label = when (cat) {
                    "restaurant" -> "식당"
                    "cafe" -> "카페"
                    "hotel" -> "숙소"
                    else -> "편의점"
                }
                FilterChip(
                    selected = cat in selectedCategories,
                    onClick = {
                        selectedCategories = if (cat in selectedCategories) selectedCategories - cat else selectedCategories + cat
                        onCategoryToggle(cat)
                    },
                    label = { Text(label) },
                )
            }
        }

        Spacer(Modifier.height(Spacing.base))

        LazyColumn(verticalArrangement = Arrangement.spacedBy(Spacing.sm)) {
            items(places) { place ->
                Card(modifier = Modifier.fillMaxWidth(), onClick = { onNavigate(place) }) {
                    Column(modifier = Modifier.padding(Spacing.md)) {
                        Text(place.name, style = MaterialTheme.typography.titleMedium)
                        Text(place.address, style = MaterialTheme.typography.bodyMedium)
                        place.rating?.let {
                            Text("평점 $it", style = MaterialTheme.typography.labelMedium)
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun TransitScreen(
    onSearch: (originLat: Double, originLng: Double, destLat: Double, destLng: Double) -> Unit,
) {
    // Transit search screen - simplified for skeleton
    Column(modifier = Modifier.fillMaxSize().padding(Spacing.base)) {
        Text("교통 검색", style = MaterialTheme.typography.headlineSmall)
        Spacer(Modifier.height(Spacing.base))
        Icon(Icons.Default.Map, contentDescription = null, modifier = Modifier.fillMaxWidth())
        Text("출발지와 도착지를 입력하세요", style = MaterialTheme.typography.bodyLarge)
    }
}
