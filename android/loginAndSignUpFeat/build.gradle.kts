plugins {
    alias(libs.plugins.android.library)
    alias(libs.plugins.composeCompiler)
    alias(libs.plugins.jetbrainsKotlinAndroid)

    // Dagger and Hilt related plugins
    alias(libs.plugins.com.google.devtools.ksp)
    alias(libs.plugins.com.google.dagger.hilt.android)
}

android {
    namespace = "xyz.fjarm.loginandsignupfeat"
    compileSdk = Config.COMPILE_SDK

    defaultConfig {
        minSdk = Config.MIN_SDK

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
        consumerProguardFiles("consumer-rules.pro")
    }
    testOptions {
        targetSdk = Config.TARGET_SDK
    }
    lint {
        targetSdk = Config.TARGET_SDK
    }

    buildFeatures {
        compose = true
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
    kotlin {
        compilerOptions {
            freeCompilerArgs = listOf(
                "-XXLanguage:+PropertyParamAnnotationDefaultTargetMode",
                "-XXLanguage:+WhenGuards",
            )
        }
        jvmToolchain(JavaVersion.VERSION_21.majorVersion.toInt())
    }
}

composeCompiler {
    // Only generate reports if we pass the flag -PcomposeCompilerReports=true
    if (project.findProperty("composeCompilerReports") == "true") {
        reportsDestination = layout.buildDirectory.dir("compose_compiler")
        metricsDestination = layout.buildDirectory.dir("compose_compiler")
    }
}

dependencies {
    // Dagger and Hilt deps
    ksp(libs.com.google.dagger.hilt.android.compiler)
    kspTest(libs.com.google.dagger.hilt.android.compiler)
    implementation(libs.com.google.dagger.hilt.android)
    implementation(libs.androidx.hilt.hilt.navigation.compose)

    implementation(project(":fjarmThemeLib"))
    implementation(project(":previewsLib"))
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.material3)
    implementation(libs.androidx.compose.runtime)
    implementation(libs.androidx.ui.tooling.preview)
    implementation(platform(libs.androidx.compose.bom))
    testImplementation(libs.junit)
    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(platform(libs.androidx.compose.bom))
    debugImplementation(libs.androidx.ui.tooling)
}
