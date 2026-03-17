plugins {
    // Use the AndroidLibraryConventionPlugin to apply Android library configuration
    id("convention.android.library")

    // Use the HiltConventionPlugin to apply Dagger Hilt configuration
    id("convention.hilt")

    // Use the ComposeConventionPlugin to apply Jetpack Compose configuration
    id("convention.compose")

    // Use the ComposeCompilerMetricsConventionPlugin to enable Compose Compiler Metrics
    id("convention.compose.metrics")
}

android {
    namespace = "xyz.fjarm.loginandsignupfeat"
}

dependencies {
    implementation(project(":fjarmThemeLib"))
    implementation(project(":loginAndSignUpFeatLib"))
    implementation(project(":previewsLib"))

    implementation(libs.androidx.material3)
    implementation(libs.androidx.navigation3.navigation3.runtime)

    debugImplementation(project(":testActivityLib"))

    testImplementation(libs.androidx.ui.test.junit4)
    testImplementation(libs.kotlinx.coroutines.test)
    testImplementation(libs.robolectric)
}
