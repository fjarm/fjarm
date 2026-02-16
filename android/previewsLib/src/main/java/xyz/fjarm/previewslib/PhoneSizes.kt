package xyz.fjarm.previewslib

import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview

/**
 * Multi-preview annotation that represents various phone sizes.
 */
@Preview(name = "Phone - Small", device = Devices.PIXEL_3, showSystemUi = true)
@Preview(name = "Phone - Medium", device = Devices.NEXUS_5, showSystemUi = true)
@Preview(name = "Phone - Large", device = Devices.PIXEL_4_XL, showSystemUi = true)
@Preview(name = "Foldable - Unfolded", device = Devices.FOLDABLE, showSystemUi = true)
annotation class PreviewPhoneSizes
