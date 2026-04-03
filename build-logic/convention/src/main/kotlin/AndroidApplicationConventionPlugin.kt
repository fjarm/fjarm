import Dependencies.addCommonAndroidDependencies
import KotlinConfig.configureKotlinCompilerArgsAndJVMToolchain
import com.android.build.api.dsl.ApplicationExtension
import org.gradle.api.JavaVersion
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
class AndroidApplicationConventionPlugin: Plugin<Project> {
    override fun apply(target: Project) {
        with(target) {
            // Apply plugins
            pluginManager.apply("com.android.application")
            pluginManager.apply("org.jetbrains.kotlin.android")

            // Configure Android using shared configuration
            extensions.configure<ApplicationExtension> {
                compileSdk = targetSDK

                compileOptions {
                    sourceCompatibility = JavaVersion.VERSION_21
                    targetCompatibility = JavaVersion.VERSION_21
                }

                // Application-specific configuration
                defaultConfig {
                    minSdk = minSDK
                    targetSdk = targetSDK
                    testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
                    vectorDrawables {
                        useSupportLibrary = true
                    }
                }

                buildTypes {
                    getByName("release") {
                        // Override library's setting - apps should minify
                        isMinifyEnabled = true
                        proguardFiles(
                            getDefaultProguardFile("proguard-android-optimize.txt"),
                            "proguard-rules.pro",
                        )
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
            configureKotlinCompilerArgsAndJVMToolchain()

            // Add common dependencies using shared configuration
            addCommonAndroidDependencies()
        }
    }
}
