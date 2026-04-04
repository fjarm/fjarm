plugins {
    // Use the AndroidLibraryConventionPlugin to apply Android library configuration
    id("convention.android.library")

    // Use the HiltConventionPlugin to apply Dagger Hilt configuration
    id("convention.hilt")
}

tasks.withType<Test> {
    failOnNoDiscoveredTests = false
}

android {
    namespace = "xyz.fjarm.testactivity"
}

dependencies {
}
