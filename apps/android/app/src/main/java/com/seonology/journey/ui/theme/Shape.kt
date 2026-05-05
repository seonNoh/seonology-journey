package com.seonology.journey.ui.theme

import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Shapes
import androidx.compose.ui.unit.dp

val JourneyShapes = Shapes(
    small = RoundedCornerShape(8.dp),
    medium = RoundedCornerShape(12.dp),   // Web rounded-card (12px)
    large = RoundedCornerShape(16.dp),    // Web rounded-2xl
    extraLarge = RoundedCornerShape(24.dp),
)
