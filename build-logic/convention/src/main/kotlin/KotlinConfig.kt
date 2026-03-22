import org.gradle.api.Project
import org.gradle.kotlin.dsl.configure
import org.gradle.kotlin.dsl.withType
import org.jetbrains.kotlin.gradle.dsl.KotlinAndroidProjectExtension
import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

internal object KotlinConfig {
    /**
     * Configures Kotlin for Android projects.
     */
    fun Project.configureKotlin() {
        // Configure Kotlin extension
        extensions.configure<KotlinAndroidProjectExtension> {
            jvmToolchain(jDKVersion)
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
}
