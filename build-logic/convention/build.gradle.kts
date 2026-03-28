plugins {
    `kotlin-dsl`
}

group = "xyz.fjarm.buildlogic"

dependencies {
    // Get versions from the version catalog
    compileOnly(libs.android.gradle.plugin)
    compileOnly(libs.kotlin.gradle.plugin)
    compileOnly(libs.compose.compiler.gradle.plugin)
    // RoborazziConventionPlugin needs to be able to access [RoborazziExtension]
    compileOnly(libs.io.github.takahirom.roborazzi.roborazzi.gradle.plugin)
}

gradlePlugin {
    plugins {
        register("androidApplication") {
            id = "convention.android.application"
            implementationClass = "AndroidApplicationConventionPlugin"
        }
        register("androidLibrary") {
            id = "convention.android.library"
            implementationClass = "AndroidLibraryConventionPlugin"
        }
        register("compose") {
            id = "convention.compose"
            implementationClass = "ComposeConventionPlugin"
        }
        register("composeMetrics") {
            id = "convention.compose.metrics"
            implementationClass = "ComposeCompilerMetricsConventionPlugin"
        }
        register("hilt") {
            id = "convention.hilt"
            implementationClass = "HiltConventionPlugin"
        }
        register("roborazzi") {
            id = "convention.roborazzi"
            implementationClass = "RoborazziConventionPlugin"
        }
    }
}
