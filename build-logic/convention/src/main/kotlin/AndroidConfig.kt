//import com.android.build.api.dsl.CommonExtension
//import org.gradle.api.JavaVersion

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
//    fun configureAndroid(
//        commonExtension: CommonExtension<*, *, *, *, *, *>
//    ) {
//        commonExtension.apply {
//            compileSdk = targetSDK
//
//            defaultConfig {
//                minSdk = minSDK
//                testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
//            }
//
//            compileOptions {
//                sourceCompatibility = JavaVersion.VERSION_21
//                targetCompatibility = JavaVersion.VERSION_21
//            }
//
//            buildTypes {
//                getByName("release") {
//                    isMinifyEnabled = false
//                    proguardFiles(
//                        getDefaultProguardFile("proguard-android-optimize.txt"),
//                        "proguard-rules.pro",
//                    )
//                }
//            }
//        }
//    }
}
