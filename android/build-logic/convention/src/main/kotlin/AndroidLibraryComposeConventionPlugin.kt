import AndroidConfigConventionPlugin.configureComposeWithDependencies
import com.android.build.api.dsl.LibraryExtension
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.kotlin.dsl.configure

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
 *       id("convention.android.library") // NOTE: Must be applied before the Compose convention plugin
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
                configureComposeWithDependencies(this)
            }
        }
    }
}
