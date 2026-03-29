package xyz.fjarm.text

import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.PreviewDynamicColors
import androidx.compose.ui.tooling.preview.PreviewFontScale
import androidx.compose.ui.tooling.preview.PreviewLightDark
import androidx.compose.ui.tooling.preview.PreviewScreenSizes
import xyz.fjarm.fjarmtheme.FjarmTypography

@Composable
fun FjarmHeaderText(
    text: String,
    modifier: Modifier = Modifier,
) {
    Text(
        text = text,
        modifier = modifier,
        style = FjarmTypography.headlineLarge.copy(
            fontWeight = FontWeight.Bold,
        ),
    )
}

@PreviewLightDark
@PreviewFontScale
@PreviewDynamicColors
@PreviewScreenSizes
@Composable
fun FjarmHeaderTextPreview() {
    FjarmHeaderText("Fjarm")
}
