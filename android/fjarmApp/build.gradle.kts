plugins {
    // Use the AndroidApplicationConventionPlugin to apply Android application configuration
    id("convention.android.application")

    // Use ComposeConventionPlugin to apply Jetpack Compose configuration
    id("convention.compose")

    // Use the HiltConventionPlugin to apply Dagger Hilt configuration
    id("convention.hilt")

    // Use the ComposeCompilerMetricsConventionPlugin to enable Compose Compiler Metrics
    id("convention.compose.metrics")
}

android {
    namespace = "xyz.fjarm.fjarm"

    defaultConfig {
        applicationId = "xyz.fjarm.fjarm"
        versionCode = 1
        versionName = "1.0.0"
    }

    buildFeatures {
        buildConfig = true
    }
}

dependencies {
    implementation(project(":fjarmThemeLib"))

    implementation(libs.androidx.activity.compose)

    implementation(libs.androidx.lifecycle.runtime.ktx)

    implementation(libs.androidx.material3)

    implementation(libs.androidx.ui)
    implementation(libs.androidx.ui.graphics)

    debugImplementation(libs.androidx.ui.test.manifest)
}
