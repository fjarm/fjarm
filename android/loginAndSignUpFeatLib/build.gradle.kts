plugins {
    // Use the AndroidLibraryConventionPlugin to apply Android library configuration
    id("convention.android.library")
}

android {
    namespace = "xyz.fjarm.loginandsignupfeatlib"
}

dependencies {
    implementation(libs.androidx.navigation3.navigation3.runtime)
}
