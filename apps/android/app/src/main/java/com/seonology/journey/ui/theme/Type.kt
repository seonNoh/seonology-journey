package com.seonology.journey.ui.theme

import androidx.compose.material3.Typography
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp

// gowun_dodum.ttf 리소스는 배포하지 않는다 (라이선스/용량). 기본 system
// sans 폰트를 쓰도록 두고, 디자인이 확정되면 res/font/ 에 ttf 를 넣은
// 뒤 FontFamily(Font(R.font.gowun_dodum)) 로 바꾼다.
val JourneyFontFamily = FontFamily.Default

val JourneyTypography = Typography(
    displayLarge = TextStyle(fontFamily = JourneyFontFamily, fontSize = 32.sp, fontWeight = FontWeight.Bold),
    headlineMedium = TextStyle(fontFamily = JourneyFontFamily, fontSize = 24.sp, fontWeight = FontWeight.Bold),
    headlineSmall = TextStyle(fontFamily = JourneyFontFamily, fontSize = 20.sp, fontWeight = FontWeight.SemiBold),
    titleLarge = TextStyle(fontFamily = JourneyFontFamily, fontSize = 20.sp, fontWeight = FontWeight.Medium),
    titleMedium = TextStyle(fontFamily = JourneyFontFamily, fontSize = 16.sp, fontWeight = FontWeight.Medium),
    bodyLarge = TextStyle(fontFamily = JourneyFontFamily, fontSize = 16.sp),
    bodyMedium = TextStyle(fontFamily = JourneyFontFamily, fontSize = 14.sp),
    labelLarge = TextStyle(fontFamily = JourneyFontFamily, fontSize = 14.sp, fontWeight = FontWeight.Medium),
    labelMedium = TextStyle(fontFamily = JourneyFontFamily, fontSize = 12.sp),
)
