package xyz.fjarm.previewslib


import androidx.compose.ui.tooling.preview.AndroidUiModes.UI_MODE_NIGHT_NO
import androidx.compose.ui.tooling.preview.AndroidUiModes.UI_MODE_NIGHT_YES
import androidx.compose.ui.tooling.preview.Preview

/**
 * Multi-preview for light and dark themes.
 */
@Preview(name = "Light Mode", uiMode = UI_MODE_NIGHT_NO, showSystemUi = true)
@Preview(name = "Dark Mode", uiMode = UI_MODE_NIGHT_YES, showSystemUi = true)
annotation class PreviewLightDarkTheme
