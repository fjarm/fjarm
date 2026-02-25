import Config.addCommonAndroidDependencies
import Config.configureAndroid
import Config.configureKotlin
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

                testOptions {
                    targetSdk = 36

                    // This is needed to allow Robolectric tests to access AndroidManifest.xml files
                    // that are in debug source sets like the one that points to [TestActivity.kt].
                    unitTests {
                        isIncludeAndroidResources = true
                    }
                }

                lint {
                    targetSdk = 36
                }
            }

            // Configure Kotlin
            configureKotlin()

            // Add common dependencies
            addCommonAndroidDependencies()
        }
    }
}
