import com.android.build.api.dsl.LibraryExtension
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.kotlin.dsl.configure
import org.jetbrains.kotlin.compose.compiler.gradle.ComposeCompilerGradlePluginExtension

/**
 * Convention plugin that enables Compose Compiler Metrics for Android library modules.
 *
 * Usage:
 *   plugins {
 *       id("convention.compose.metrics")
 *   }
 *
 * To generate metrics, run Gradle with the flag:
 *   ./gradlew assembleRelease -PcomposeCompilerReports=true
 */
class ComposeCompilerMetricsConventionPlugin : Plugin<Project> {
    override fun apply(target: Project) {
        with(target) {
            // Ensure the Android library plugin is applied
            pluginManager.withPlugin("com.android.library") {
                configureComposeMetrics()
            }

            // Also support Android application modules
            pluginManager.withPlugin("com.android.application") {
                configureComposeMetrics()
            }
        }
    }

    private fun Project.configureComposeMetrics() {
        // Check if Compose Compiler plugin is applied
        pluginManager.withPlugin("org.jetbrains.kotlin.plugin.compose") {

            // Configure the Compose Compiler extension
            extensions.configure<ComposeCompilerGradlePluginExtension> {
                // Only generate reports if the flag is passed
                if (project.findProperty("composeCompilerReports") == "true") {
                    val outputDir = layout.buildDirectory.dir("compose_compiler")

                    reportsDestination.set(outputDir)
                    metricsDestination.set(outputDir)

                    // Optional: Enable stability configuration file generation
                    // This helps understand why classes are unstable
                    stabilityConfigurationFiles.add(
                        layout.projectDirectory.file("compose_compiler_config.conf")
                    )

                    logger.lifecycle("Compose Compiler Metrics enabled for ${project.name}")
                    logger.lifecycle("Output directory: ${outputDir.get().asFile.absolutePath}")
                }
            }
        }
    }
}
