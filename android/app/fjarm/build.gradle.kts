plugins {
    // Use the AndroidApplicationConventionPlugin to apply Android application configuration
    id("convention.android.application")

    // Use ComposeConventionPlugin to apply Jetpack Compose configuration
    id("convention.compose")

    // Use the ComposeCompilerMetricsConventionPlugin to enable Compose Compiler Metrics
    id("convention.compose.metrics")

    // Use the HiltConventionPlugin to apply Dagger Hilt configuration
    id("convention.hilt")
}

android {
    namespace = "xyz.fjarm.fjarm"

    defaultConfig {
        applicationId = "xyz.fjarm.fjarm"
        versionCode = 1
        versionName = "1.0.0"
    }
}

dependencies {
    implementation(project(":android:feature:loginAndSignUp:api"))
    implementation(project(":android:feature:loginAndSignUp:impl"))
    implementation(project(":android:feature:underConstruction:impl"))

    implementation(project(":android:library:fjarmTheme"))
    implementation(project(":android:library:navigation"))

    implementation(libs.androidx.activity.compose)
    implementation(libs.androidx.compose.material3)
    implementation(libs.androidx.navigation3.navigation3.runtime)
    implementation(libs.androidx.navigation3.navigation3.ui)
}
