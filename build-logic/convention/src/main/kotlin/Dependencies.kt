import org.gradle.api.Project
import org.gradle.kotlin.dsl.dependencies

internal object Dependencies {

    /**
     * Adds common Android dependencies (core, testing).
     */
    fun Project.addCommonAndroidDependencies() {
        dependencies {
            add("implementation", catalog.findLibrary("androidx.core.ktx").get())

            add("testImplementation", catalog.findLibrary("junit").get())
            add("androidTestImplementation", catalog.findLibrary("androidx.junit").get())
            add("androidTestImplementation", catalog.findLibrary("androidx.espresso.core").get())
        }
    }

    // Compose dependencies
    fun Project.addCommonComposeDependencies() {
        dependencies {
            val bom = catalog.findLibrary("androidx.compose.bom").get()
            add("implementation", platform(bom))
            add("androidTestImplementation", platform(bom))

            // Core Compose dependencies
            add("implementation", catalog.findLibrary("androidx.compose.runtime").get())
            add("implementation", catalog.findLibrary("androidx.compose.ui.tooling.preview").get())

            // Needed for Previews to work correctly
            add("debugImplementation", catalog.findLibrary("androidx.compose.ui.tooling").get())
            // Needed to provide createAndroidComposeRule
            add("testDebugImplementation", catalog.findLibrary("androidx.compose.ui.test.junit4").get())
        }
    }

    // Dagger/Hilt dependencies
    fun Project.addCommonHiltDependencies() {
        // Add Hilt dependencies
        dependencies {
            add("ksp", catalog.findLibrary("com.google.dagger.hilt.compiler").get())
            add("implementation", catalog.findLibrary("com.google.dagger.hilt.core").get())
        }
    }

    // Dagger/Hilt Android dependencies
    fun Project.addCommonHiltAndroidDependencies() {
        dependencies {
            add("ksp", catalog.findLibrary("com.google.dagger.hilt.android.compiler").get())
            add("implementation", catalog.findLibrary("com.google.dagger.hilt.android").get())

            add("kspTest", catalog.findLibrary("com.google.dagger.hilt.android.compiler").get())
            add("testImplementation", catalog.findLibrary("com.google.dagger.hilt.android.testing").get())
        }
    }

    // Roborazzi dependencies
    fun Project.addCommonRoborazziDependencies() {
        dependencies {
            add("testImplementation", catalog.findLibrary("io.github.takahirom.roborazzi.roborazzi").get())
        }
    }
}
