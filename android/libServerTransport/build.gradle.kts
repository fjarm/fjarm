plugins {
    id("java-library")
    alias(libs.plugins.jetbrainsKotlinJvm)

    // Dagger and Hilt related plugins
    alias(libs.plugins.com.google.devtools.ksp)
}

java {
    sourceCompatibility = JavaVersion.VERSION_21
    targetCompatibility = JavaVersion.VERSION_21
}

dependencies {
    implementation(libs.connectKotlin)
    implementation(libs.coroutinesCoreJvm)
    implementation(libs.okHttp3)
    implementation(libs.connectKotlinGoogleJavaLiteExt)
    implementation(libs.connectKotlinOkHttp)

    // Dagger and Hilt related deps
    implementation(libs.com.google.dagger.hilt.core)
    ksp(libs.com.google.dagger.hilt.compiler)
}
