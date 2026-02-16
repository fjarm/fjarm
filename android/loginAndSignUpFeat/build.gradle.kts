plugins {
    // Use the AndroidLibraryConventionPlugin to apply Android library configuration
    id("convention.android.library")

    // Use the AndroidHiltConventionPlugin to apply Dagger Hilt configuration
    id("convention.android.hilt")

    // Use the AndroidLibraryComposeConventionPlugin to apply Jetpack Compose configuration
    id("convention.android.library.compose")

    // Use the ComposeCompilerMetricsConventionPlugin to enable Compose Compiler Metrics
    id("convention.compose.metrics")
}

android {
    namespace = "xyz.fjarm.loginandsignupfeat"
}

dependencies {
    implementation(project(":fjarmThemeLib"))
    implementation(project(":previewsLib"))

    implementation(libs.androidx.material3)
}
