import AndroidConfig.configureCommonAndroid
import Dependencies.addCommonAndroidDependencies
import KotlinConfig.configureKotlinCompilerArgsAndJVMToolchain
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

            configureCommonAndroid()
            // Configure Android
            extensions.configure<LibraryExtension> {
                buildTypes {
                    getByName("release") {
                        isMinifyEnabled = false // Default value for library modules
                    }
                }

                defaultConfig {
                    consumerProguardFiles("consumer-rules.pro")
                }

                testOptions {
                    targetSdk = targetSDK

                    // This is needed to allow Robolectric tests to access AndroidManifest.xml files
                    // that are in debug source sets like the one that points to [TestActivity.kt].
                    unitTests {
                        isIncludeAndroidResources = true
                    }
                }

                lint {
                    targetSdk = targetSDK
                }
            }

            // Configure Kotlin
            configureKotlinCompilerArgsAndJVMToolchain()

            // Add common dependencies
            addCommonAndroidDependencies()
        }
    }
}
