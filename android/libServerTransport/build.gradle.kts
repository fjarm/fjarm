plugins {
    id("java-library")
    alias(libs.plugins.jetbrainsKotlinJvm)
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
}
