plugins {
    `kotlin-dsl`
}

group = "xyz.fjarm.buildlogic"

dependencies {
    // Get versions from the version catalog
    compileOnly(libs.android.gradlePlugin)
    compileOnly(libs.kotlin.gradlePlugin)
    compileOnly(libs.compose.compiler.gradlePlugin)
}

gradlePlugin {
    plugins {
        register("androidHilt") {
            id = "convention.android.hilt"
            implementationClass = "AndroidHiltConventionPlugin"
        }
        register("androidLibrary") {
            id = "convention.android.library"
            implementationClass = "AndroidLibraryConventionPlugin"
        }
        register("androidLibraryCompose") {
            id = "convention.android.library.compose"
            implementationClass = "AndroidLibraryComposeConventionPlugin"
        }
        register("composeMetrics") {
            id = "convention.compose.metrics"
            implementationClass = "ComposeCompilerMetricsConventionPlugin"
        }
    }
}
