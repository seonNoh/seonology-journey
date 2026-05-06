package com.seonology.journey.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.BasicTextField
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.LocationOn
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.SolidColor
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.TextFieldValue
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.seonology.journey.data.GeocodePlace
import com.seonology.journey.data.JourneyApi
import com.seonology.journey.ui.theme.Sakura100
import com.seonology.journey.ui.theme.Sakura50
import com.seonology.journey.ui.theme.Sakura500
import com.seonology.journey.ui.theme.Sakura600
import com.seonology.journey.ui.theme.Sakura700
import com.seonology.journey.ui.theme.Sakura900
import com.seonology.journey.ui.theme.Warm500
import kotlinx.coroutines.delay

/**
 * 지도 기반 장소 검색 입력 필드. 백엔드 `/external/geocode` 를 호출해
 * 식당/호텔/관광지 등 POI 후보를 자동완성으로 보여주고, 사용자가 선택하면
 * 이름과 좌표를 함께 폼에 채워준다. 300ms 디바운스로 호출 빈도를 제한한다.
 */
@Composable
fun PlaceSearchField(
    api: JourneyApi,
    label: String,
    value: String,
    onChange: (String) -> Unit,
    onSelect: (GeocodePlace) -> Unit,
    modifier: Modifier = Modifier,
    placeholder: String = "식당 / 호텔 / 관광지 검색",
) {
    var input by remember(value) { mutableStateOf(TextFieldValue(value)) }
    var results by remember { mutableStateOf<List<GeocodePlace>>(emptyList()) }
    var loading by remember { mutableStateOf(false) }
    var open by remember { mutableStateOf(false) }

    LaunchedEffect(input.text) {
        val q = input.text.trim()
        if (q.length < 2) {
            results = emptyList()
            loading = false
            return@LaunchedEffect
        }
        loading = true
        delay(300)
        runCatching { api.geocode(q) }
            .onSuccess {
                results = it.places
                open = it.places.isNotEmpty()
            }
            .onFailure { results = emptyList() }
        loading = false
    }

    Column(modifier = modifier.fillMaxWidth()) {
        Text(
            text = label,
            fontSize = 12.sp,
            fontWeight = FontWeight.Bold,
            color = Sakura700,
        )
        Spacer(Modifier.height(4.dp))
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .clip(RoundedCornerShape(12.dp))
                .background(Color.White)
                .border(1.dp, Sakura100, RoundedCornerShape(12.dp))
                .padding(horizontal = 12.dp, vertical = 10.dp),
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Icon(
                    imageVector = Icons.Default.Search,
                    contentDescription = null,
                    tint = Sakura600,
                    modifier = Modifier.size(16.dp),
                )
                Spacer(Modifier.width(8.dp))
                Box(modifier = Modifier.weight(1f)) {
                    if (input.text.isEmpty()) {
                        Text(
                            text = placeholder,
                            color = Warm500,
                            fontSize = 13.sp,
                        )
                    }
                    BasicTextField(
                        value = input,
                        onValueChange = {
                            input = it
                            onChange(it.text)
                        },
                        singleLine = true,
                        textStyle = TextStyle(color = Sakura900, fontSize = 13.sp),
                        cursorBrush = SolidColor(Sakura500),
                        modifier = Modifier.fillMaxWidth(),
                    )
                }
                if (loading) {
                    Spacer(Modifier.width(8.dp))
                    CircularProgressIndicator(
                        modifier = Modifier.size(14.dp),
                        color = Sakura500,
                        strokeWidth = 2.dp,
                    )
                }
            }
        }
        if (open && results.isNotEmpty()) {
            Spacer(Modifier.height(6.dp))
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .clip(RoundedCornerShape(12.dp))
                    .background(Color.White)
                    .border(1.dp, Sakura100, RoundedCornerShape(12.dp)),
            ) {
                LazyColumn(modifier = Modifier.heightIn(max = 220.dp)) {
                    items(results) { place ->
                        Row(
                            verticalAlignment = Alignment.Top,
                            modifier = Modifier
                                .fillMaxWidth()
                                .clickable {
                                    input = TextFieldValue(place.name)
                                    onChange(place.name)
                                    onSelect(place)
                                    open = false
                                }
                                .padding(horizontal = 12.dp, vertical = 10.dp),
                        ) {
                            Icon(
                                imageVector = Icons.Default.LocationOn,
                                contentDescription = null,
                                tint = Sakura500,
                                modifier = Modifier.size(16.dp),
                            )
                            Spacer(Modifier.width(8.dp))
                            Column(modifier = Modifier.weight(1f)) {
                                Text(
                                    text = place.name,
                                    fontSize = 13.sp,
                                    fontWeight = FontWeight.Bold,
                                    color = Sakura900,
                                )
                                if (place.address.isNotEmpty() && place.address != place.name) {
                                    Spacer(Modifier.height(2.dp))
                                    Text(
                                        text = place.address,
                                        fontSize = 11.sp,
                                        color = Warm500,
                                    )
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

/**
 * 동일 파일 내에서만 사용되는 모디파이어 helper 가 필요한 경우 여기에.
 */
