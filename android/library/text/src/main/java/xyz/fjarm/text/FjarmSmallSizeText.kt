package xyz.fjarm.text

import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.tooling.preview.PreviewDynamicColors
import androidx.compose.ui.tooling.preview.PreviewFontScale
import androidx.compose.ui.tooling.preview.PreviewLightDark
import androidx.compose.ui.tooling.preview.PreviewScreenSizes
import androidx.compose.ui.unit.sp

@Composable
fun FjarmSmallSizeText(
    text: String,
    modifier: Modifier = Modifier,
    fontWeight: FontWeight = FontWeight.Normal,
    textDecoration: TextDecoration = TextDecoration.None,
) {
    Text(
        text = text,
        modifier = modifier,
        fontSize = 8.sp,
        fontWeight = fontWeight,
        textDecoration = textDecoration,
    )
}

@PreviewLightDark
@PreviewFontScale
@PreviewDynamicColors
@PreviewScreenSizes
@Composable
fun FjarmSmallSizeTextPreview() {
    FjarmSmallSizeText("Fjarm")
}
