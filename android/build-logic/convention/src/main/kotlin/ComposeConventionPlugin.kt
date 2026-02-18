import Config.configureComposeWithDependencies
import com.android.build.api.dsl.ApplicationExtension
import com.android.build.api.dsl.LibraryExtension
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.kotlin.dsl.configure

/**
 * Convention plugin for modules with Jetpack Compose.
 * Works for both Android applications and libraries.
 *
 * Applies:
 * - Base Android conventions (via convention.android.library or convention.android.application)
 * - Compose Compiler plugin
 * - Common Compose dependencies (BOM, Material3, runtime, preview)
 *
 * Usage in library modules:
 *   plugins {
 *       id("convention.android.library")
 *       id("convention.compose")
 *   }
 *
 * Usage in application modules:
 *   plugins {
 *       id("convention.android.application")
 *       id("convention.compose")
 *   }
 *
 *   android {
 *       namespace = "com.your.module"
 *       buildFeatures {
 *           buildConfig = true  // If needed in app modules
 *       }
 *   }
 *
 *   dependencies {
 *       // Add app-specific dependencies:
 *       implementation(libs.androidx.activity.compose)
 *       implementation(libs.androidx.lifecycle.viewmodel.compose)
 *   }
 */
class ComposeConventionPlugin : Plugin<Project> {
    override fun apply(target: Project) {
        with(target) {
            // Apply Compose Compiler plugin
            pluginManager.apply("org.jetbrains.kotlin.plugin.compose")

            // Configure for library modules
            pluginManager.withPlugin("com.android.library") {
                extensions.configure<LibraryExtension> {
                    configureComposeWithDependencies(this)
                }
            }

            // Configure for application modules
            pluginManager.withPlugin("com.android.application") {
                extensions.configure<ApplicationExtension> {
                    configureComposeWithDependencies(this)
                }
            }
        }
    }
}
