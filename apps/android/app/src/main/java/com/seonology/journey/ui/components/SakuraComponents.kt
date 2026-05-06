package com.seonology.journey.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.seonology.journey.ui.theme.CatActivity
import com.seonology.journey.ui.theme.CatActivityBg
import com.seonology.journey.ui.theme.CatFood
import com.seonology.journey.ui.theme.CatFoodBg
import com.seonology.journey.ui.theme.CatLodging
import com.seonology.journey.ui.theme.CatLodgingBg
import com.seonology.journey.ui.theme.CatOther
import com.seonology.journey.ui.theme.CatOtherBg
import com.seonology.journey.ui.theme.CatShopping
import com.seonology.journey.ui.theme.CatShoppingBg
import com.seonology.journey.ui.theme.CatTransport
import com.seonology.journey.ui.theme.CatTransportBg
import com.seonology.journey.ui.theme.MoodCalm
import com.seonology.journey.ui.theme.MoodCalmBg
import com.seonology.journey.ui.theme.MoodExcited
import com.seonology.journey.ui.theme.MoodExcitedBg
import com.seonology.journey.ui.theme.MoodTasty
import com.seonology.journey.ui.theme.MoodTastyBg
import com.seonology.journey.ui.theme.MoodTired
import com.seonology.journey.ui.theme.MoodTiredBg
import com.seonology.journey.ui.theme.Sakura100
import com.seonology.journey.ui.theme.Sakura600
import com.seonology.journey.ui.theme.Sakura700
import com.seonology.journey.ui.theme.Sakura900
import com.seonology.journey.ui.theme.Sky100
import com.seonology.journey.ui.theme.Sky700
import com.seonology.journey.ui.theme.Spacing
import com.seonology.journey.ui.theme.Warm100
import com.seonology.journey.ui.theme.Warm200
import com.seonology.journey.ui.theme.Warm500
import com.seonology.journey.ui.theme.Warm700

/** Pill chip used by Sakura Bear theme. */
@Composable
fun SbChip(
    text: String,
    bg: Color = Sakura100,
    fg: Color = Sakura700,
    leadingIcon: ImageVector? = null,
    small: Boolean = false,
) {
    Row(
        verticalAlignment = Alignment.CenterVertically,
        modifier = Modifier
            .background(bg, RoundedCornerShape(999.dp))
            .padding(
                horizontal = if (small) 8.dp else 10.dp,
                vertical = if (small) 2.dp else 4.dp,
            ),
    ) {
        if (leadingIcon != null) {
            Icon(
                imageVector = leadingIcon,
                contentDescription = null,
                tint = fg,
                modifier = Modifier.size(if (small) 10.dp else 12.dp),
            )
            Spacer(Modifier.width(4.dp))
        }
        Text(
            text = text,
            color = fg,
            fontSize = if (small) 10.sp else 11.sp,
            fontWeight = FontWeight.Bold,
        )
    }
}

/** Status pill for trip status. */
@Composable
fun SbStatusPill(status: String, small: Boolean = false) {
    val (label, bg, fg) = when (status) {
        "TRIP_STATUS_PLANNING", "계획중" -> Triple("계획중", Sakura100, Sakura700)
        "TRIP_STATUS_ONGOING", "여행중" -> Triple("여행중", Sky100, Sky700)
        "TRIP_STATUS_COMPLETED", "완료" -> Triple("완료", Warm200, Warm700)
        "TRIP_STATUS_ARCHIVED", "보관" -> Triple("보관", Warm100, Warm500)
        else -> Triple(status, Sakura100, Sakura700)
    }
    SbChip(text = label, bg = bg, fg = fg, small = small)
}

/** Section header used between cards/lists. */
@Composable
fun SbSection(
    title: String,
    icon: ImageVector? = null,
    count: String? = null,
    action: (@Composable () -> Unit)? = null,
) {
    Row(
        verticalAlignment = Alignment.CenterVertically,
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = Spacing.base, vertical = Spacing.sm),
    ) {
        if (icon != null) {
            Box(
                modifier = Modifier
                    .size(28.dp)
                    .background(Sakura100, RoundedCornerShape(10.dp)),
                contentAlignment = Alignment.Center,
            ) {
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    tint = Sakura600,
                    modifier = Modifier.size(16.dp),
                )
            }
            Spacer(Modifier.width(Spacing.sm))
        }
        Text(
            title,
            fontSize = 14.sp,
            fontWeight = FontWeight.Bold,
            color = Sakura900,
        )
        if (count != null) {
            Spacer(Modifier.width(6.dp))
            Text(
                count,
                fontSize = 12.sp,
                color = Warm500,
                fontWeight = FontWeight.SemiBold,
            )
        }
        Spacer(Modifier.weight(1f))
        action?.invoke()
    }
}

/** Mood metadata for Notes screen. */
data class MoodMeta(val label: String, val color: Color, val bg: Color, val emoji: String)

val MoodTable: Map<String, MoodMeta> = mapOf(
    "설렘" to MoodMeta("설렘", MoodExcited, MoodExcitedBg, "✨"),
    "맛있음" to MoodMeta("맛있음", MoodTasty, MoodTastyBg, "🍙"),
    "평온" to MoodMeta("평온", MoodCalm, MoodCalmBg, "🌿"),
    "피곤" to MoodMeta("피곤", MoodTired, MoodTiredBg, "💤"),
)

fun moodFor(name: String?): MoodMeta? = name?.let { MoodTable[it] }

/** Expense category metadata. */
data class CategoryMeta(val label: String, val color: Color, val bg: Color)

val CategoryTable: Map<String, CategoryMeta> = mapOf(
    "EXPENSE_CATEGORY_TRANSPORT" to CategoryMeta("교통", CatTransport, CatTransportBg),
    "EXPENSE_CATEGORY_FOOD" to CategoryMeta("식사", CatFood, CatFoodBg),
    "EXPENSE_CATEGORY_LODGING" to CategoryMeta("숙박", CatLodging, CatLodgingBg),
    "EXPENSE_CATEGORY_ACTIVITY" to CategoryMeta("체험", CatActivity, CatActivityBg),
    "EXPENSE_CATEGORY_SHOPPING" to CategoryMeta("쇼핑", CatShopping, CatShoppingBg),
    "EXPENSE_CATEGORY_OTHER" to CategoryMeta("기타", CatOther, CatOtherBg),
)

fun categoryFor(c: String?): CategoryMeta =
    c?.let { CategoryTable[it] } ?: CategoryMeta(c ?: "기타", CatOther, CatOtherBg)
