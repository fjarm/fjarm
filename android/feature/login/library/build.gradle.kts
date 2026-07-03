plugins {
    id("java-library")
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
    implementation(project(":android:library:serverTransport"))

    implementation(libs.build.buf.gen.fjarm.connectrpc.kotlin.lite)
    implementation(libs.build.buf.gen.fjarm.protocolbuffers.kotlin.lite)
}
