import AndroidConfig.addCommonHiltAndroidDependencies
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
 *       id("convention.android.hilt")
 *   }
 */
class AndroidHiltConventionPlugin : Plugin<Project> {
    override fun apply(target: Project) {
        with(target) {
            // Apply plugins
            pluginManager.apply("com.google.devtools.ksp")
            pluginManager.apply("com.google.dagger.hilt.android")

            // Add Hilt dependencies
            addCommonHiltAndroidDependencies()
        }
    }
}
