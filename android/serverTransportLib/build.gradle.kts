plugins {
    id("java-library")
    id("convention.hilt")
    alias(libs.plugins.jetbrainsKotlinJvm)
}

java {
    sourceCompatibility = JavaVersion.VERSION_21
    targetCompatibility = JavaVersion.VERSION_21
}

dependencies {
    implementation(libs.connectKotlin)
    implementation(libs.connectKotlinGoogleJavaLiteExt)
    implementation(libs.connectKotlinOkHttp)
    implementation(libs.coroutinesCoreJvm)
    implementation(libs.okHttp3)
}
