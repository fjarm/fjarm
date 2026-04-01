# Summary

Multi-module Android app built following MVI and Clean Architecture.

## Project structure

The project contains the following top-level Gradle modules:

* `:android:app` - contains Android Application Gradle modules such as `:android:app:fjarm`
* `:android:feature` - contains Android Library Gradle modules that host UI features such as `:android:feature:loginAndSignUp`
    * `:android:feature:{Feature Name}:api` - contains navigation keys (i.e. `NavKey` implementations) for each feature such as `:android:feature:loginAndSignUp:api`
    * `:android:feature:{Feature Name}:impl` - contains the UI code (i.e. `@Composable` functions and `ViewModel` classes) for each feature such as `:android:feature:loginAndSignUp:impl`
    * `:android:feature:{Feature Name}:library` - contains pure Kotlin Gradle modules that host feature-supporting code such as use-cases and repositories in `:android:feature:login:library`
* `:android:library` - contains Android Library and pure Kotlin Gradle modules that host globally shared code such as `:android:library:navigation` or `:android:library:serverTransport`

`:build-logic` is a special top-level Gradle module that contains convention plugins.

## AGP and Gradle version compatibility

The required AGP version for each Gradle version can be found [at this link](https://developer.android.com/build/releases/about-agp#updating-gradle).

Upgrade Gradle in the CLI like so:

```bash
./gradlew --gradle-version 9.3.1
```

The command may need to be run twice.

## Android Studio and AGP version compatibility

The required Android Studio version for each AGP version can be found [at this link](https://developer.android.com/build/releases/about-agp#android_gradle_plugin_and_android_studio_compatibility).
