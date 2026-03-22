plugins {
    // Use the AndroidLibraryConventionPlugin to apply Android library configuration
    id("convention.android.library")

    // Apply the parcelize plugin to support @Parcelize annotations
    alias(libs.plugins.kotlin.parcelize)
}

android {
    namespace = "xyz.fjarm.loginandsignupapi"
}

dependencies {
    implementation(libs.androidx.navigation3.navigation3.runtime)
}
