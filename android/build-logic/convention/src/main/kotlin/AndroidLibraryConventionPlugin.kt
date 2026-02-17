import AndroidConfigConventionPlugin.addCommonAndroidDependencies
import AndroidConfigConventionPlugin.configureAndroid
import AndroidConfigConventionPlugin.configureKotlin
import com.android.build.api.dsl.LibraryExtension
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.kotlin.dsl.configure

/**
 * Convention plugin for Android library modules.
 *
 * Applies:
 * - Android library plugin
 * - Kotlin Android plugin
 * - Standard Android configuration (SDK versions, Java version, etc.)
 *
 * Usage:
 *   plugins {
 *       id("convention.android.library")
 *   }
 *
 *   android {
 *       namespace = "com.example.yourmodule"  // Still need to set this
 *   }
 */
class AndroidLibraryConventionPlugin : Plugin<Project> {
    override fun apply(target: Project) {
        with(target) {
            // Apply plugins
            pluginManager.apply("com.android.library")
            pluginManager.apply("org.jetbrains.kotlin.android")

            // Configure Android
            extensions.configure<LibraryExtension> {
                configureAndroid(this)

                defaultConfig {
                    consumerProguardFiles("consumer-rules.pro")
                }
            }

            // Configure Kotlin
            configureKotlin()

            // Add common dependencies
            addCommonAndroidDependencies()
        }
    }
}
