package xyz.fjarm.previewslib

import androidx.compose.ui.tooling.preview.Preview

/**
 * Multi-preview for testing different font scales (accessibility).
 */

@Preview(name = "Font Scale - Small", fontScale = 0.85f, showSystemUi = true)
@Preview(name = "Font Scale - Large", fontScale = 1.3f, showSystemUi = true)
@Preview(name = "Font Scale - Huge", fontScale = 2f, showSystemUi = true)
annotation class PreviewFontScales
