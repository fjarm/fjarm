@file:Suppress("UnstableApiUsage")

pluginManagement {
    includeBuild("build-logic")
    repositories {
        google {
            content {
                includeGroupByRegex("com\\.android.*")
                includeGroupByRegex("com\\.google.*")
                includeGroupByRegex("androidx.*")
            }
        }
        mavenCentral()
        gradlePluginPortal()
    }
}

dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        google()
        mavenCentral()
        maven {
            name = "buf"
            url = uri("https://buf.build/gen/maven")
        }
    }
}

rootProject.name = "fjarm"
include(":android:app:fjarm")
include(":android:feature:login:api")
include(":android:feature:loginAndSignUp:api")
include(":android:feature:loginAndSignUp:impl")
include(":android:feature:underConstruction:impl")
include(":android:library:buttons")
include(":android:library:fjarmTheme")
include(":android:library:navigation")
include(":android:library:previews")
include(":android:library:serverTransport")
include(":android:library:testActivity")
