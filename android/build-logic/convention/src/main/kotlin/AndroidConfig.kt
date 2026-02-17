import com.android.build.api.dsl.CommonExtension
import org.gradle.api.JavaVersion
import org.gradle.api.Project
import org.gradle.api.artifacts.VersionCatalog
import org.gradle.api.artifacts.VersionCatalogsExtension
import org.gradle.kotlin.dsl.configure
import org.gradle.kotlin.dsl.dependencies
import org.gradle.kotlin.dsl.getByType
import org.gradle.kotlin.dsl.withType
import org.jetbrains.kotlin.gradle.dsl.KotlinAndroidProjectExtension
import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

/**
 * Shared Android configuration for both library and application modules.
 *
 * This object provides common configuration functions that work on CommonExtension,
 * which is the parent of both LibraryExtension and ApplicationExtension.
 */
internal object AndroidConfig {

    /**
     * Configures common Android settings for both app and library modules.
     */
    fun configureAndroid(
        commonExtension: CommonExtension<*, *, *, *, *, *>
    ) {
        commonExtension.apply {
            compileSdk = 36

            defaultConfig {
                minSdk = 31
                testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
            }

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
        }
    }

    /**
     * Configures Kotlin for Android projects.
     */
    fun Project.configureKotlin() {
        // Configure Kotlin extension
        extensions.configure<KotlinAndroidProjectExtension> {
            jvmToolchain(21)
        }

        // Configure Kotlin compile tasks
        tasks.withType<KotlinCompile>().configureEach {
            compilerOptions {
                freeCompilerArgs.addAll(
                    "-XXLanguage:+PropertyParamAnnotationDefaultTargetMode",
                    "-XXLanguage:+WhenGuards",
                )
            }
        }
    }

    /**
     * Adds common Android dependencies (core, testing).
     */
    fun Project.addCommonAndroidDependencies() {
        dependencies {
            add("implementation", catalog.findLibrary("androidx.core.ktx").get())

            add("testImplementation", catalog.findLibrary("junit").get())
            add("androidTestImplementation", catalog.findLibrary("androidx.junit").get())
            add("androidTestImplementation", catalog.findLibrary("androidx.espresso.core").get())
        }
    }

    /**
     * Configures Jetpack Compose for Android modules.
     */
    fun Project.configureComposeWithDependencies(
        commonExtension: CommonExtension<*, *, *, *, *, *>
    ) {
        commonExtension.apply {
            buildFeatures {
                compose = true
            }
        }

        dependencies {
            val bom = catalog.findLibrary("androidx.compose.bom").get()
            add("implementation", platform(bom))
            add("androidTestImplementation", platform(bom))

            // Core Compose dependencies
            add("implementation", catalog.findLibrary("androidx.compose.runtime").get())
            add("implementation", catalog.findLibrary("androidx.ui.tooling.preview").get())

            // Debug/Testing
            add("debugImplementation", catalog.findLibrary("androidx.ui.tooling").get())
            add("androidTestImplementation", catalog.findLibrary("androidx.ui.test.junit4").get())
        }
    }

    private val Project.catalog: VersionCatalog
        get() = extensions.getByType<VersionCatalogsExtension>().named("libs")
}