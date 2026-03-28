plugins {
    // Use the AndroidLibraryConventionPlugin to apply Android library configuration
    id("convention.android.library")

    // Use the ComposeConventionPlugin to apply Jetpack Compose configuration
    id("convention.compose")

    // Use the ComposeCompilerMetricsConventionPlugin to enable Compose Compiler Metrics
    id("convention.compose.metrics")
}

android {
    namespace = "xyz.fjarm.buttons"
}

dependencies {
    implementation(project(":android:library:fjarmTheme"))

    implementation(libs.androidx.compose.ui)
    implementation(libs.androidx.compose.material3)
}
