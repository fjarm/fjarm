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

    implementation(project(":serverTransportLib"))
    implementation(libs.fjarmConnectSdk)
    implementation(libs.fjarmProtobufLiteSdk)

    testImplementation(libs.junit)
}
