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
        register("composeMetrics") {
            id = "convention.compose.metrics"
            implementationClass = "ComposeCompilerMetricsConventionPlugin"
        }
    }
}
