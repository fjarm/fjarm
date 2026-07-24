plugins {
    id("java-library")
    id("java-test-fixtures")

    // Use the HiltConventionPlugin to apply Dagger Hilt configuration
    id("convention.hilt")

    alias(libs.plugins.jetbrains.kotlin.jvm)
}

java {
    sourceCompatibility = JavaVersion.VERSION_21
    targetCompatibility = JavaVersion.VERSION_21
}

kotlin {
    compilerOptions {
        jvmTarget = org.jetbrains.kotlin.gradle.dsl.JvmTarget.JVM_21
    }
}

dependencies {
    implementation(libs.kotlinx.coroutines.core.jvm)
    testImplementation(libs.kotlinx.coroutines.test)
    testImplementation(libs.junit)

    testFixturesApi(libs.kotlinx.coroutines.test)
    testFixturesApi(libs.com.google.dagger.hilt.core)
    testFixturesApi(libs.junit)
    kspTestFixtures(libs.com.google.dagger.hilt.compiler)
}
