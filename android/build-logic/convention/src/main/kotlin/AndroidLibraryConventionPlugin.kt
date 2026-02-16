import com.android.build.api.dsl.LibraryExtension
import org.gradle.api.JavaVersion
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.api.artifacts.VersionCatalogsExtension
import org.gradle.kotlin.dsl.configure
import org.gradle.kotlin.dsl.dependencies
import org.gradle.kotlin.dsl.getByType
import org.jetbrains.kotlin.gradle.dsl.KotlinAndroidProjectExtension

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
                compileSdk = 36  // Or use: findProperty("compileSdk")?.toString()?.toInt() ?: 35

                defaultConfig {
                    minSdk = 31

                    testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
                    consumerProguardFiles("consumer-rules.pro")
                }

                testOptions {
                    targetSdk = 36
                }

                lint {
                    targetSdk = 36
                }

                buildTypes {
                    release {
                        isMinifyEnabled = false
                        proguardFiles(
                            getDefaultProguardFile("proguard-android-optimize.txt"),
                            "proguard-rules.pro"
                        )
                    }
                }

                compileOptions {
                    sourceCompatibility = JavaVersion.VERSION_21
                    targetCompatibility = JavaVersion.VERSION_21
                }
            }

            // Configure Kotlin
            extensions.configure<KotlinAndroidProjectExtension> {
                compilerOptions {
                    freeCompilerArgs.addAll(
                        "-XXLanguage:+PropertyParamAnnotationDefaultTargetMode",
                        "-XXLanguage:+WhenGuards",
                    )
                }
                jvmToolchain(21)
            }

            // Add common dependencies
            dependencies {
                add("implementation", catalog.findLibrary("androidx.core.ktx").get())

                add("testImplementation", catalog.findLibrary("junit").get())
                add("androidTestImplementation", catalog.findLibrary("androidx.junit").get())
                add("androidTestImplementation", catalog.findLibrary("androidx.espresso.core").get())
            }
        }
    }

    private val Project.catalog
        get() = extensions.getByType<VersionCatalogsExtension>()
            .named("libs")
}
