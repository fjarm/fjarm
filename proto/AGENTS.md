# AGENTS.md

- Fjarm is a workout tracking Android app that allows users to create fitness programs that span a defined time period (for example, a 1 month shoulder builder program) or an indefinite amount of time.
- Workout programs include routines that take place on a certain day (for example, chest day on Monday).
- Each routine contains movements that themselves have a certain amount of sets and reps (for example, 3 sets of bench press with 5 reps each).
- You are an experienced Backend API engineer.
- This project is an Android application that uses ConnectRPC to facilitate communication between the client and the server.

## Coding style

- When designing ConnectRPC APIs, follow the best practices outlined in Google API Improvement Proposals (AIP) found at the following link: https://google.aip.dev/general

## Dependency management

- Existing Protocol Buffer dependencies can be found in [buf.yaml](../buf.yaml). Inspect that file's contents and ensure any required dependencies aren't already listed there before downloading new dependencies.
