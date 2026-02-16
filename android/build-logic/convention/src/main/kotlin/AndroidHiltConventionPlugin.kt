import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.api.artifacts.VersionCatalogsExtension
import org.gradle.kotlin.dsl.dependencies
import org.gradle.kotlin.dsl.getByType

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
            dependencies {
                add("ksp", catalog.findLibrary("com.google.dagger.hilt.android.compiler").get())
                add("kspTest", catalog.findLibrary("com.google.dagger.hilt.android.compiler").get())

                add("implementation", catalog.findLibrary("com.google.dagger.hilt.android").get())
                add("implementation", catalog.findLibrary("androidx.hilt.hilt.navigation.compose").get())

                add("testImplementation", catalog.findLibrary("com.google.dagger.hilt.android.testing").get())
            }
        }
    }

    private val Project.catalog
        get() = extensions.getByType<VersionCatalogsExtension>()
            .named("libs")
}
