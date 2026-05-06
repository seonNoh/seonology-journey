package com.seonology.journey.ui.components

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.layout.size
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.rotate
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.StrokeCap
import androidx.compose.ui.graphics.drawscope.DrawScope
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.graphics.drawscope.withTransform
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import com.seonology.journey.ui.theme.MascotBearCheek
import com.seonology.journey.ui.theme.MascotBearEar
import com.seonology.journey.ui.theme.MascotBearFur
import com.seonology.journey.ui.theme.MascotBearOutline
import com.seonology.journey.ui.theme.MascotBearSnout
import com.seonology.journey.ui.theme.MascotChickBeak
import com.seonology.journey.ui.theme.MascotChickBody
import com.seonology.journey.ui.theme.MascotChickOutline
import com.seonology.journey.ui.theme.MascotMiniCheek
import com.seonology.journey.ui.theme.MascotMiniEar
import com.seonology.journey.ui.theme.MascotMiniFur
import com.seonology.journey.ui.theme.MascotMiniOutline
import com.seonology.journey.ui.theme.MascotMiniSnout
import com.seonology.journey.ui.theme.Sakura300
import com.seonology.journey.ui.theme.Sakura500

/**
 * Sakura Bear mascot expressions / accessories.
 *
 * Compose 캔버스로 SVG (100x100 viewBox) 디자인을 그대로 재현한다.
 * 디자인 원본은 `design/sakura-bear-mascot.jsx` 참조.
 */
enum class MascotExpression { Plain, Happy, Sleep, Wink, Surprise }

enum class MascotAccessory { None, Flower, Camera, SleepCap }

data class BearPalette(
    val fur: Color,
    val ear: Color,
    val snout: Color,
    val outline: Color,
    val cheek: Color,
)

val MainBearPalette = BearPalette(
    fur = MascotBearFur,
    ear = MascotBearEar,
    snout = MascotBearSnout,
    outline = MascotBearOutline,
    cheek = MascotBearCheek,
)

val MiniBearPalette = BearPalette(
    fur = MascotMiniFur,
    ear = MascotMiniEar,
    snout = MascotMiniSnout,
    outline = MascotMiniOutline,
    cheek = MascotMiniCheek,
)

@Composable
fun SbBear(
    size: Dp = 80.dp,
    expression: MascotExpression = MascotExpression.Plain,
    accessory: MascotAccessory = MascotAccessory.None,
    palette: BearPalette = MainBearPalette,
) {
    Canvas(modifier = Modifier.size(size)) {
        val s = this.size.width / 100f
        withTransform({ scale(s, s, pivot = Offset.Zero) }) {
            drawBear(palette, expression, accessory)
        }
    }
}

@Composable
fun SbMiniBear(
    size: Dp = 40.dp,
    expression: MascotExpression = MascotExpression.Plain,
) = SbBear(size = size, expression = expression, palette = MiniBearPalette)

@Composable
fun SbChick(
    size: Dp = 40.dp,
    expression: MascotExpression = MascotExpression.Plain,
) {
    Canvas(modifier = Modifier.size(size)) {
        val s = this.size.width / 100f
        withTransform({ scale(s, s, pivot = Offset.Zero) }) {
            drawChick(expression)
        }
    }
}

@Composable
fun SbPetal(
    size: Dp = 14.dp,
    color: Color = Sakura500,
    rotateDeg: Float = 0f,
    opacity: Float = 0.7f,
) {
    Canvas(modifier = Modifier.size(size).rotate(rotateDeg)) {
        val s = this.size.width / 20f
        withTransform({ scale(s, s, pivot = Offset.Zero) }) {
            drawPetal(color, opacity)
        }
    }
}

@Composable
fun SbPaw(size: Dp = 16.dp, color: Color = Sakura300) {
    Canvas(modifier = Modifier.size(size)) {
        val s = this.size.width / 20f
        withTransform({ scale(s, s, pivot = Offset.Zero) }) {
            drawPaw(color)
        }
    }
}

// ───────────────────────────────────────────────────────────────────
// Drawing primitives (in 100x100 unit space).
// ───────────────────────────────────────────────────────────────────

private fun DrawScope.drawCircleAt(c: Color, cx: Float, cy: Float, r: Float) {
    drawCircle(color = c, radius = r, center = Offset(cx, cy))
}

private fun DrawScope.strokeCircleAt(c: Color, cx: Float, cy: Float, r: Float, width: Float) {
    drawCircle(color = c, radius = r, center = Offset(cx, cy), style = Stroke(width = width))
}

private fun DrawScope.drawOvalAt(c: Color, cx: Float, cy: Float, rx: Float, ry: Float, alpha: Float = 1f) {
    drawOval(
        color = c,
        topLeft = Offset(cx - rx, cy - ry),
        size = Size(rx * 2f, ry * 2f),
        alpha = alpha,
    )
}

private fun DrawScope.strokeOvalAt(c: Color, cx: Float, cy: Float, rx: Float, ry: Float, width: Float) {
    drawOval(
        color = c,
        topLeft = Offset(cx - rx, cy - ry),
        size = Size(rx * 2f, ry * 2f),
        style = Stroke(width = width),
    )
}

private fun DrawScope.drawBear(p: BearPalette, expr: MascotExpression, acc: MascotAccessory) {
    // Ears
    drawCircleAt(p.fur, 22f, 28f, 13f)
    strokeCircleAt(p.outline, 22f, 28f, 13f, 1.5f)
    drawCircleAt(p.fur, 78f, 28f, 13f)
    strokeCircleAt(p.outline, 78f, 28f, 13f, 1.5f)
    drawCircleAt(p.ear, 22f, 28f, 7f)
    drawCircleAt(p.ear, 78f, 28f, 7f)

    // Head
    drawOvalAt(p.fur, 50f, 56f, 34f, 32f)
    strokeOvalAt(p.outline, 50f, 56f, 34f, 32f, 1.5f)

    // Snout
    drawOvalAt(p.snout, 50f, 65f, 18f, 14f)

    // Cheeks
    drawOvalAt(p.cheek, 22f, 62f, 5f, 5f, alpha = 0.65f)
    drawOvalAt(p.cheek, 78f, 62f, 5f, 5f, alpha = 0.65f)

    // Eyes
    drawBearEyes(expr, p.outline)

    // Nose
    drawOvalAt(p.outline, 50f, 60f, 3f, 2.2f)

    // Mouth
    drawBearMouth(expr, p.outline)

    // Accessory
    when (acc) {
        MascotAccessory.Flower -> drawFlower(72f, 8f)
        MascotAccessory.Camera -> drawCamera(64f, 72f, p.outline)
        MascotAccessory.SleepCap -> drawSleepCap(p.outline)
        MascotAccessory.None -> Unit
    }
}

private fun DrawScope.drawBearEyes(expr: MascotExpression, outline: Color) {
    when (expr) {
        MascotExpression.Sleep -> {
            // Closed-eye smile arcs
            arc(34f, 50f, 42f, 50f, dy = 3f, color = outline, width = 1.6f)
            arc(58f, 50f, 66f, 50f, dy = 3f, color = outline, width = 1.6f)
        }
        MascotExpression.Happy -> {
            arc(34f, 51f, 42f, 51f, dy = -4f, color = outline, width = 1.8f)
            arc(58f, 51f, 66f, 51f, dy = -4f, color = outline, width = 1.8f)
        }
        MascotExpression.Wink -> {
            arc(34f, 51f, 42f, 51f, dy = -4f, color = outline, width = 1.8f)
            drawOvalAt(outline, 62f, 50f, 2.6f, 3.4f)
            drawCircleAt(Color.White, 62.7f, 49f, 0.8f)
        }
        MascotExpression.Surprise -> {
            drawCircleAt(outline, 38f, 50f, 3f)
            drawCircleAt(outline, 62f, 50f, 3f)
            drawCircleAt(Color.White, 38.7f, 49f, 1f)
            drawCircleAt(Color.White, 62.7f, 49f, 1f)
        }
        MascotExpression.Plain -> {
            drawOvalAt(outline, 38f, 50f, 2.6f, 3.4f)
            drawOvalAt(outline, 62f, 50f, 2.6f, 3.4f)
            drawCircleAt(Color.White, 38.7f, 49f, 0.8f)
            drawCircleAt(Color.White, 62.7f, 49f, 0.8f)
        }
    }
}

private fun DrawScope.drawBearMouth(expr: MascotExpression, outline: Color) {
    when (expr) {
        MascotExpression.Surprise -> drawOvalAt(outline, 50f, 66f, 2.5f, 3f)
        MascotExpression.Happy -> arc(44f, 64f, 56f, 64f, dy = 6f, color = outline, width = 1.6f)
        else -> {
            // "w" shape: M 50 62 Q 46 68 42 65   and   M 50 62 Q 54 68 58 65
            drawQuad(50f, 62f, 46f, 68f, 42f, 65f, outline, 1.4f)
            drawQuad(50f, 62f, 54f, 68f, 58f, 65f, outline, 1.4f)
        }
    }
}

private fun DrawScope.drawFlower(cx: Float, cy: Float) {
    // 5 petals + center, around (cx, cy)
    drawCircleAt(Color(0xFFFFA3B8), cx - 4f, cy - 2f, 3f)
    drawCircleAt(Color(0xFFFFA3B8), cx + 4f, cy - 2f, 3f)
    drawCircleAt(Color(0xFFFFA3B8), cx - 2f, cy + 3f, 3f)
    drawCircleAt(Color(0xFFFFA3B8), cx + 2f, cy + 3f, 3f)
    drawCircleAt(Sakura500, cx, cy, 4f)
    drawCircleAt(Color.White, cx, cy, 1.5f)
}

private fun DrawScope.drawCamera(cx: Float, cy: Float, outline: Color) {
    // Body rect (cx-10, cy-6, 20x14)
    val bodyColor = Color(0xFF5A2230)
    drawRect(
        color = bodyColor,
        topLeft = Offset(cx - 10f, cy - 6f),
        size = Size(20f, 14f),
    )
    drawCircleAt(Color(0xFFFFB3C7), cx, cy, 4f)
    drawCircleAt(bodyColor, cx, cy, 2f)
}

private fun DrawScope.drawSleepCap(outline: Color) {
    // Stylized: triangle-ish path 18,14 -> 50,-8 -> 82,14 closed.
    val path = Path().apply {
        moveTo(18f, 14f)
        quadraticBezierTo(50f, -8f, 82f, 14f)
        lineTo(78f, 22f)
        quadraticBezierTo(50f, 6f, 22f, 22f)
        close()
    }
    drawPath(path, color = Sakura500)
    drawPath(path, color = outline, style = Stroke(width = 1.2f))
    drawCircleAt(Color.White, 86f, 10f, 4f)
    strokeCircleAt(outline, 86f, 10f, 4f, 1.2f)
}

private fun DrawScope.drawChick(expr: MascotExpression) {
    val body = MascotChickBody
    val beak = MascotChickBeak
    val outline = MascotChickOutline

    // Body
    drawOvalAt(body, 50f, 55f, 34f, 30f)
    strokeOvalAt(outline, 50f, 55f, 34f, 30f, 1.5f)

    // Tuft (head feather): line from 42,22 -> 50,14 -> 58,22
    val tuft = Path().apply {
        moveTo(42f, 22f)
        lineTo(50f, 14f)
        lineTo(58f, 22f)
    }
    drawPath(tuft, color = outline, style = Stroke(width = 1.5f))

    // Eyes
    if (expr == MascotExpression.Happy) {
        arc(38f, 49f, 46f, 49f, dy = -3f, color = outline, width = 1.6f)
        arc(54f, 49f, 62f, 49f, dy = -3f, color = outline, width = 1.6f)
    } else {
        drawCircleAt(outline, 42f, 50f, 2.2f)
        drawCircleAt(outline, 58f, 50f, 2.2f)
    }

    // Beak (triangle)
    val beakPath = Path().apply {
        moveTo(46f, 60f)
        lineTo(54f, 60f)
        lineTo(50f, 66f)
        close()
    }
    drawPath(beakPath, color = beak)
    drawPath(beakPath, color = outline, style = Stroke(width = 1f))

    // Cheeks
    drawCircleAt(Color(0xFFFFB3C7).copy(alpha = 0.7f), 32f, 58f, 3f)
    drawCircleAt(Color(0xFFFFB3C7).copy(alpha = 0.7f), 68f, 58f, 3f)
}

private fun DrawScope.drawPetal(color: Color, opacity: Float) {
    // Path equivalent of: M10 2 C 14 4 14 10 10 12 C 6 10 6 4 10 2 Z
    val path = Path().apply {
        moveTo(10f, 2f)
        cubicTo(14f, 4f, 14f, 10f, 10f, 12f)
        cubicTo(6f, 10f, 6f, 4f, 10f, 2f)
        close()
    }
    drawPath(path, color = color, alpha = opacity)
    drawCircleAt(Color.White.copy(alpha = 0.6f * opacity), 10f, 11f, 1.2f)
}

private fun DrawScope.drawPaw(color: Color) {
    drawOvalAt(color, 10f, 13f, 5f, 4f)
    drawCircleAt(color, 5f, 7f, 2f)
    drawCircleAt(color, 10f, 5f, 2f)
    drawCircleAt(color, 15f, 7f, 2f)
    drawCircleAt(color, 3f, 11f, 1.5f)
    drawCircleAt(color, 17f, 11f, 1.5f)
}

// ───────────────────────────────────────────────────────────────────
// Path helpers.
// ───────────────────────────────────────────────────────────────────

/**
 * Quadratic-bezier "smile" arc. (x1,y1) → (x2,y2), control auto at midpoint
 * shifted by `dy` (negative dy = curve upward).
 */
private fun DrawScope.arc(x1: Float, y1: Float, x2: Float, y2: Float, dy: Float, color: Color, width: Float) {
    val cx = (x1 + x2) / 2f
    val cy = (y1 + y2) / 2f + dy
    drawQuad(x1, y1, cx, cy, x2, y2, color, width)
}

private fun DrawScope.drawQuad(
    x1: Float, y1: Float,
    cx: Float, cy: Float,
    x2: Float, y2: Float,
    color: Color, width: Float,
) {
    val path = Path().apply {
        moveTo(x1, y1)
        quadraticBezierTo(cx, cy, x2, y2)
    }
    drawPath(
        path = path,
        color = color,
        style = Stroke(width = width, cap = StrokeCap.Round),
    )
}
