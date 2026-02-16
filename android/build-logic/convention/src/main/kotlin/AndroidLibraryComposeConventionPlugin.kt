import com.android.build.api.dsl.LibraryExtension
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.api.artifacts.VersionCatalogsExtension
import org.gradle.kotlin.dsl.configure
import org.gradle.kotlin.dsl.dependencies
import org.gradle.kotlin.dsl.getByType

/**
 * Convention plugin for Android library modules with Jetpack Compose.
 *
 * Applies:
 * - Android library conventions (via convention.android.library)
 * - Compose Compiler plugin
 * - Compose dependencies (BOM, Material3, UI, etc.)
 *
 * Usage:
 *   plugins {
 *       id("convention.android.library") // Typically need to apply this
 *       id("convention.android.library.compose")
 *   }
 *
 *   android {
 *       namespace = "com.example.yourmodule"
 *   }
 */
class AndroidLibraryComposeConventionPlugin : Plugin<Project> {
    override fun apply(target: Project) {
        with(target) {
            // Apply Compose Compiler plugin
            pluginManager.apply("org.jetbrains.kotlin.plugin.compose")

            // Enable Compose in Android
            extensions.configure<LibraryExtension> {
                buildFeatures {
                    compose = true
                }
            }

            // Add Compose dependencies
            dependencies {
                val bom = catalog.findLibrary("androidx.compose.bom").get()
                add("implementation", platform(bom))
                add("androidTestImplementation", platform(bom))

                // Core Compose dependencies
                add("implementation", catalog.findLibrary("androidx.compose.runtime").get())
                add("implementation", catalog.findLibrary("androidx.ui.tooling.preview").get())

                // Debug/Testing
                add("debugImplementation", catalog.findLibrary("androidx.ui.tooling").get())
                add("androidTestImplementation", catalog.findLibrary("androidx.ui.test.junit4").get())
            }
        }
    }

    private val Project.catalog
        get() = extensions.getByType<VersionCatalogsExtension>()
            .named("libs")
}
