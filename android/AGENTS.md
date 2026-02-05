# AGENTS.md

- You are an experienced Android developer.
- This project an Android application built with Jetpack Compose. The programming languages and tools used include Kotlin, Gradle, and ConnectRPC. The project follows a layered architecture with unidirectional data flow.

## Testing instructions

- Do not use a mocking library.
- Do not use backticks for test names.
- Follow Behavior-Driven-Design patterns for writing tests in a Given-When-Then format.

## Coding style

- Use interfaces and prefer composition over inheritance.
- Clean up any unused imports.
- Follow official architecture recommendations, including use of a layered architecture.
- For example, use a unidirectional data flow (UDF), ViewModels, lifecycle-aware UI state collection, and other recommendations.
- As it specifically pertains to UDF, follow Model-View-Intent architecture pattern.
- As it specifically pertains to layered architecture, place business logic in uses cases to separate concerns.
- Review any suggested code for memory leaks.
- Review any suggested code for leaked coroutines.

## Dependency management

- Existing dependencies can be found in [gradle/libs.versions.toml](gradle/libs.versions.toml). Inspect that file's contents and ensure any required dependencies aren't already listed there before downloading new dependencies.
