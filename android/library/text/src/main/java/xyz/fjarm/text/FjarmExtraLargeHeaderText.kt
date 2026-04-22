package xyz.fjarm.text

import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.PreviewDynamicColors
import androidx.compose.ui.tooling.preview.PreviewFontScale
import androidx.compose.ui.tooling.preview.PreviewLightDark
import androidx.compose.ui.tooling.preview.PreviewScreenSizes
import androidx.compose.ui.unit.sp
import xyz.fjarm.fjarmtheme.FjarmTypography

@Composable
fun FjarmExtraLargeHeaderText(
    text: String,
    modifier: Modifier = Modifier,
) {
    Text(
        text = text,
        modifier = modifier,
        style = FjarmTypography.displayLarge.copy(
            fontWeight = FontWeight.Black,
        ),
        fontSize = 72.sp,
        lineHeight = 72.sp,
        letterSpacing = (-2).sp,
    )
}

@PreviewLightDark
@PreviewFontScale
@PreviewDynamicColors
@PreviewScreenSizes
@Composable
fun FjarmExtraLargeHeaderTextPreview() {
    FjarmExtraLargeHeaderText("Fjarm")
}
