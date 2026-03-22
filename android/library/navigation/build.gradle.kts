plugins {
    // Use the AndroidLibraryConventionPlugin to apply Android library configuration
    id("convention.android.library")

    // Use the HiltConventionPlugin to apply Dagger Hilt configuration
    id("convention.hilt")

    alias(libs.plugins.kotlin.parcelize)
}

android {
    namespace = "xyz.fjarm.navigation"
}

dependencies {
    implementation(libs.androidx.navigation3.navigation3.runtime)

    testImplementation(libs.kotlinx.coroutines.test)
}
