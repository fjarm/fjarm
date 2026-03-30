// Top-level build file where you can add configuration options common to all sub-projects/modules.
plugins {
    alias(libs.plugins.android.application) apply false
    alias(libs.plugins.kotlin.android) apply false
    alias(libs.plugins.kotlin.compose) apply false

    // Our plugins below

    // Dagger/Hilt Android
    alias(libs.plugins.com.google.dagger.hilt.android) apply false

    // Dagger
    alias(libs.plugins.com.google.devtools.ksp) apply false

    // Android Library
    alias(libs.plugins.android.library) apply false

    // Parcelize plugin
    alias(libs.plugins.kotlin.parcelize) apply false

    // Kotlin JVM plugin
    alias(libs.plugins.jetbrains.kotlin.jvm) apply false

    // Roborazzi plugin
    alias(libs.plugins.io.github.takahirom.roborazzi) apply false
}
