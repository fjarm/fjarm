import com.android.build.api.dsl.CommonExtension
import org.gradle.api.Project
import org.gradle.kotlin.dsl.configure

internal object ComposeConfig {

    /**
     * Configures Jetpack Compose for Android modules.
     */
    fun Project.configureComposeBuildFeaturesEnabled() {
        extensions.configure<CommonExtension> {
            buildFeatures.apply {
                compose = true
            }
        }
    }
}
