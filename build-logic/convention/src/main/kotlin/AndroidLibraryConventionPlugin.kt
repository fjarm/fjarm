import Dependencies.addCommonAndroidDependencies
import KotlinConfig.configureKotlinCompilerArgsAndJVMToolchain
import com.android.build.api.dsl.LibraryExtension
import org.gradle.api.JavaVersion
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

            // Configure Android
            extensions.configure<LibraryExtension> {
                compileSdk = targetSDK

                compileOptions {
                    sourceCompatibility = JavaVersion.VERSION_21
                    targetCompatibility = JavaVersion.VERSION_21
                }

                buildTypes {
                    getByName("release") {
                        isMinifyEnabled = false
                        proguardFiles(
                            getDefaultProguardFile("proguard-android-optimize.txt"),
                            "proguard-rules.pro",
                        )
                    }
                }

                defaultConfig {
                    minSdk = minSDK
                    testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
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
