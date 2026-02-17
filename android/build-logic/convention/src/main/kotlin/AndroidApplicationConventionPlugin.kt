import AndroidConfig.addCommonAndroidDependencies
import AndroidConfig.configureAndroid
import AndroidConfig.configureKotlin
import com.android.build.api.dsl.ApplicationExtension
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.kotlin.dsl.configure

/**
 * Convention plugin for Android application modules.
 *
 * Applies:
 * - Android application plugin
 * - Kotlin Android plugin
 * - Standard Android configuration (SDK versions, Java version, etc.)
 *
 * Usage:
 *   plugins {
 *       id("convention.android.application")
 *   }
 *
 *   android {
 *       namespace = "com.example.yourapp"
 *       defaultConfig {
 *           applicationId = "com.example.yourapp"
 *           versionCode = 1
 *           versionName = "1.0"
 *       }
 *   }
 */
class AndroidApplicationConventionPlugin : Plugin<Project> {
    override fun apply(target: Project) {
        with(target) {
            // Apply plugins
            pluginManager.apply("com.android.application")
            pluginManager.apply("org.jetbrains.kotlin.android")

            // Configure Android using shared configuration
            extensions.configure<ApplicationExtension> {
                configureAndroid(this)

                // Application-specific configuration
                defaultConfig {
                    targetSdk = 36
                    vectorDrawables {
                        useSupportLibrary = true
                    }
                }

                buildTypes {
                    getByName("release") {
                        // Override library's setting - apps should minify
                        isMinifyEnabled = true
                    }
                }

                packaging {
                    resources {
                        excludes += "/META-INF/{AL2.0,LGPL2.1}"
                    }
                }

                testOptions {
                    unitTests {
                        isIncludeAndroidResources = true
                    }
                }
            }

            // Configure Kotlin using shared configuration
            configureKotlin()

            // Add common dependencies using shared configuration
            addCommonAndroidDependencies()
        }
    }
}