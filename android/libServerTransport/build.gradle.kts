plugins {
    id("java-library")
    alias(libs.plugins.jetbrainsKotlinJvm)
}

java {
    sourceCompatibility = JavaVersion.VERSION_17
    targetCompatibility = JavaVersion.VERSION_17
}

dependencies {
    implementation(libs.coroutinesCoreJvm)
    implementation(libs.grpcAndroid)
    implementation(libs.grpcOkHttp)
    implementation(libs.connectKotlin)
    implementation(libs.connectKotlinGoogleJavaLiteExt)
    implementation(libs.connectKotlinOkHttp)
}