plugins {
    // Use the AndroidApplicationConventionPlugin to apply Android application configuration
    id("convention.android.application")

    // Use ComposeConventionPlugin to apply Jetpack Compose configuration
    id("convention.compose")

    // Use the AndroidHiltConventionPlugin to apply Dagger Hilt configuration
    id("convention.android.hilt")

    // Use the ComposeCompilerMetricsConventionPlugin to enable Compose Compiler Metrics
    id("convention.compose.metrics")
}

android {
    namespace = "xyz.fjarm.helloworld"

    defaultConfig {
        applicationId = "xyz.fjarm.helloworld"
        versionCode = Config.VERSION_CODE
        versionName = Config.VERSION_NAME
    }

    buildFeatures {
        buildConfig = true
    }

    testOptions {
        unitTests {
            isIncludeAndroidResources = true
        }
    }
}

dependencies {
    implementation(project(":helloWorldLib"))
    implementation(libs.fjarmProtobufLiteSdk)

    implementation(libs.androidx.lifecycle.runtime.compose)
    implementation(libs.androidx.lifecycle.runtime.ktx)
    implementation(libs.androidx.lifecycle.viewmodel.compose)
    implementation(libs.androidx.lifecycle.viewmodel.ktx)
    implementation(libs.androidx.activity.compose)
    implementation(libs.androidx.ui)
    implementation(libs.androidx.ui.graphics)
    implementation(libs.androidx.material3)

    testImplementation(libs.androidx.ui.test.junit4)
    testImplementation(libs.kotlinx.coroutines.test)
    testImplementation(libs.robolectric)

    debugImplementation(libs.androidx.ui.test.manifest)
}
