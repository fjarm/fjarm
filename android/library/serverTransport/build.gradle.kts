import org.jetbrains.kotlin.gradle.dsl.JvmTarget

plugins {
    id("java-library")

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
        jvmTarget = JvmTarget.JVM_21
    }
}

dependencies {
    implementation(libs.kotlinx.coroutines.core.jvm)

    implementation(libs.com.connectrpc.connect.kotlin)
    implementation(libs.com.connectrpc.connect.kotlin.okhttp)
    implementation(libs.com.connectrpc.connect.kotlin.google.javalite.ext)

    implementation(libs.com.squareup.okhttp3.okhttp)
}
