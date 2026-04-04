import com.android.build.api.dsl.CommonExtension
import org.gradle.api.JavaVersion
import org.gradle.api.Project
import org.gradle.kotlin.dsl.configure

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
    fun Project.configureCommonAndroid() {
        extensions.configure<CommonExtension> {
            compileSdk = targetSDK

            defaultConfig.apply {
                minSdk = minSDK
                testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
            }

            compileOptions.apply {
                sourceCompatibility = JavaVersion.VERSION_21
                targetCompatibility = JavaVersion.VERSION_21
            }

            buildTypes.apply {
                getByName("release") {
                    proguardFiles(
                        getDefaultProguardFile("proguard-android-optimize.txt"),
                        "proguard-rules.pro",
                    )
                }
            }
        }
    }
}
