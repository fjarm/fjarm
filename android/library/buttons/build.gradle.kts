plugins {
    // Use the AndroidLibraryConventionPlugin to apply Android library configuration
    id("convention.android.library")

    // Use the HiltConventionPlugin to apply Dagger Hilt configuration
    id("convention.hilt")

    // Use the ComposeConventionPlugin to apply Jetpack Compose configuration
    id("convention.compose")

    // Use the ComposeCompilerMetricsConventionPlugin to enable Compose Compiler Metrics
    id("convention.compose.metrics")

    // Use the RoborazziConventionPlugin to enable Roborazzi tests
    id("convention.roborazzi")
}

android {
    namespace = "xyz.fjarm.buttons"
}

dependencies {
    implementation(libs.androidx.compose.ui)
    implementation(libs.androidx.compose.material3)

    debugImplementation(project(":android:library:testActivity"))

    testDebugImplementation(libs.org.robolectric.robolectric)
    testDebugImplementation(libs.kotlinx.coroutines.test)
}
