# Summary

Multi-module Android app built following MVI and Clean Architecture.

# Project structure

The project contains the following top-level Gradle modules:

* `:app` - contains Android Application Gradle modules such as `:app:fjarm`
* `:feature` - contains Android Library Gradle modules that host UI features such as `:feature:loginAndSignUp`
    * `:feature:{Feature Name}:api` - contains navigation keys (i.e. `NavKey` implementations) for each feature such as `:feature:loginAndSignUp:api`
    * `:feature:{Feature Name}:impl` - contains the UI code (i.e. `@Composable` functions and `ViewModel` classes) for each feature such as `:feature:loginAndSignUp:impl`
* `:library` - contains Android Library and pure Kotlin Gradle modules that host feature-supporting code such as `:library:navigation` or `:library:serverTransport`
    * `:library:loginAndSignUp` - pure Kotlin Gradle module where each use-case is implemented

`:build-logic` is a special top-level Gradle module that contains convention plugins.

