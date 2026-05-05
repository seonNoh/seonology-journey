package com.seonology.journey.ui.theme

import androidx.compose.material3.Typography
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.Font
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp
import com.seonology.journey.R

val GowunDodum = FontFamily(
    Font(R.font.gowun_dodum, FontWeight.Normal),
)

val JourneyTypography = Typography(
    displayLarge = TextStyle(fontFamily = GowunDodum, fontSize = 32.sp, fontWeight = FontWeight.Bold),
    headlineMedium = TextStyle(fontFamily = GowunDodum, fontSize = 24.sp, fontWeight = FontWeight.Bold),
    headlineSmall = TextStyle(fontFamily = GowunDodum, fontSize = 20.sp, fontWeight = FontWeight.SemiBold),
    titleLarge = TextStyle(fontFamily = GowunDodum, fontSize = 20.sp, fontWeight = FontWeight.Medium),
    titleMedium = TextStyle(fontFamily = GowunDodum, fontSize = 16.sp, fontWeight = FontWeight.Medium),
    bodyLarge = TextStyle(fontFamily = GowunDodum, fontSize = 16.sp),
    bodyMedium = TextStyle(fontFamily = GowunDodum, fontSize = 14.sp),
    labelLarge = TextStyle(fontFamily = GowunDodum, fontSize = 14.sp, fontWeight = FontWeight.Medium),
    labelMedium = TextStyle(fontFamily = GowunDodum, fontSize = 12.sp),
)
