package com.seonology.journey.ui.theme

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable

private val LightColors = lightColorScheme(
    primary = Sakura500,
    onPrimary = Sakura50,
    primaryContainer = Sakura100,
    onPrimaryContainer = Sakura900,
    secondary = Sky500,
    onSecondary = Sky50,
    secondaryContainer = Sky100,
    onSecondaryContainer = Sky900,
    surface = Warm50,
    onSurface = Warm900,
    surfaceVariant = Warm100,
    onSurfaceVariant = Warm700,
    background = Warm50,
    onBackground = Warm900,
    error = ErrorRed,
    outline = Warm300,
)

private val DarkColors = darkColorScheme(
    primary = Sakura300,
    onPrimary = Sakura900,
    primaryContainer = Sakura800,
    onPrimaryContainer = Sakura100,
    secondary = Sky300,
    onSecondary = Sky900,
    secondaryContainer = Sky800,
    onSecondaryContainer = Sky100,
    surface = Warm900,
    onSurface = Warm100,
    surfaceVariant = Warm800,
    onSurfaceVariant = Warm300,
    background = Warm900,
    onBackground = Warm100,
    error = ErrorRed,
    outline = Warm600,
)

@Composable
fun JourneyTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    content: @Composable () -> Unit,
) {
    MaterialTheme(
        colorScheme = if (darkTheme) DarkColors else LightColors,
        typography = JourneyTypography,
        shapes = JourneyShapes,
        content = content,
    )
}
