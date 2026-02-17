import AndroidConfig.addCommonHiltAndroidDependencies
import AndroidConfig.addCommonHiltDependencies
import org.gradle.api.Plugin
import org.gradle.api.Project

/**
 * Convention plugin for modules using Dagger Hilt.
 *
 * Applies:
 * - KSP plugin
 * - Hilt Android plugin
 * - Hilt dependencies (compiler, runtime, navigation)
 *
 * Usage:
 *   plugins {
 *       id("convention.android.library") // NOTE: Must be applied before the Hilt convention plugin
 *       id("convention.hilt")
 *   }
 */
class HiltConventionPlugin : Plugin<Project> {
    override fun apply(target: Project) {
        with(target) {
            // Always apply KSP plugin
            pluginManager.apply("com.google.devtools.ksp")

            pluginManager.withPlugin("org.jetbrains.kotlin.jvm") {
                addCommonHiltDependencies()
            }

            pluginManager.withPlugin("com.android.library") {
                // Apply Android specific Hilt plugin
                pluginManager.apply("com.google.dagger.hilt.android")
                // Add Hilt dependencies
                addCommonHiltAndroidDependencies()
            }

            pluginManager.withPlugin("com.android.application") {
                // Apply Android specific Hilt plugin
                pluginManager.apply("com.google.dagger.hilt.android")
                // Add Hilt dependencies
                addCommonHiltAndroidDependencies()
            }
        }
    }
}
