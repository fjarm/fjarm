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
    namespace = "xyz.fjarm.loginandsignupimpl"
}

dependencies {
    implementation(project(":android:feature:loginAndSignUp:api"))
    implementation(project(":android:library:fjarmTheme"))
    implementation(project(":android:library:previews"))
    implementation(libs.androidx.compose.material3)
    implementation(libs.androidx.hilt.hilt.navigation.compose)
    implementation(libs.androidx.navigation3.navigation3.runtime)

    debugImplementation(project(":android:library:testActivity"))

    testDebugImplementation(libs.androidx.compose.ui.test.junit4)
    testDebugImplementation(libs.org.robolectric.robolectric)
    testImplementation(libs.kotlinx.coroutines.test)
}
