plugins {
    // Use the AndroidLibraryConventionPlugin to apply Android library configuration
    id("convention.android.library")

    // Use the AndroidLibraryComposeConventionPlugin to apply Jetpack Compose configuration
    id("convention.android.library.compose")

    // Use the ComposeCompilerMetricsConventionPlugin to enable Compose Compiler Metrics
    id("convention.compose.metrics")
}

android {
    namespace = "xyz.fjarm.previewslib"
}

dependencies {
}
