import Dependencies.addCommonRoborazziDependencies
import io.github.takahirom.roborazzi.RoborazziExtension
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.kotlin.dsl.configure

/**
 * Convention plugin that enables Roborazzi tests for Android library modules with UI.
 *
 * Usage:
 *   plugins {
 *       id("convention.roborazzi")
 *   }
 *
 * To test and generate screenshots, run Gradle with the flag(s):
 *   ./gradlew testDebugUnitTest --Proborazzi.test.record=true -Proborazzi.test.compare=true -Proborazzi.test.verify=true
 */
class RoborazziConventionPlugin: Plugin<Project> {

    override fun apply(target: Project) {
        with(target) {
            // Apply the Roborazzi plugin
            pluginManager.apply("io.github.takahirom.roborazzi")

            // Ensure the Android library plugin is applied
            pluginManager.withPlugin("com.android.library") {
                configureRoborazziOutputDir()
                addCommonRoborazziDependencies()
            }
        }
    }

    private fun Project.configureRoborazziOutputDir() {
        // Check if Roborazzi plugin is applied
        pluginManager.withPlugin("io.github.takahirom.roborazzi") {
            extensions.configure<RoborazziExtension> {
                outputDir.set(layout.projectDirectory.dir("src/screenshots"))
            }
        }
    }
}
