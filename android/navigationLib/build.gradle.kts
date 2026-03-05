plugins {
    // Use the AndroidLibraryConventionPlugin to apply Android library configuration
    id("convention.android.library")

    // Use the HiltConventionPlugin to apply Dagger Hilt configuration
    id("convention.hilt")
}

android {
    namespace = "xyz.fjarm.navigationlib"
}

dependencies {
    implementation(project(":loginAndSignUpFeatLib"))

    implementation(libs.androidx.navigation3.navigation3.runtime)

    testImplementation(libs.androidx.ui.test.junit4)
    testImplementation(libs.kotlinx.coroutines.test)
}
