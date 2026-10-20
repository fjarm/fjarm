# Objective

Analyze whether user sign-up/registration should be a separate `RegistrationService` (as drafted in `registration.md`) or folded into the existing `UserService` in `proto/fjarm/users/v1`. Determine if the latter would create bloat.

# Research & Context

## Existing `UserService` ([user_service.proto](../../../proto/fjarm/users/v1/user_service.proto))
- Follows AIP standard CRUD pattern: `CreateUser`, `GetUser`, `UpdateUser`, `DeleteUser` (4 RPCs).
- Each RPC operates directly on the `User` resource.
- `CreateUserRequest` already accepts a full `User` message and an idempotency key.
- `CreateUserResponse` is intentionally empty to prevent account enumeration.

## Existing `AuthenticationService` ([authentication_service.proto](../../../proto/fjarm/authentication/v1/authentication_service.proto))
- Separate package: `fjarm.authentication.v1`.
- Contains `CreateSession` — a single RPC focused on session lifecycle, not user CRUD.
- Precedent: the project already separates concerns by domain boundary, not by grouping all user-adjacent operations together.

## Draft `RegistrationService` ([registration.md](../../../registration.md))
- Proposed package: `fjarm.registration.v1`.
- 3 RPCs: `BeginUserRegistration`, `CompleteUserRegistration`, `RefreshRegistrationToken`.
- Introduces new domain concepts: `RegistrationCredentials`, `RegistrationToken`.
- Multi-step flow with token issuance and expiry — fundamentally different lifecycle from CRUD.
- TODO comments suggest the design is still evolving (handle removal, credential vs. user distinction).
- `CompleteUserRegistration` would internally call `CreateUser` — an orchestration concern.

## `User` entity ([user.proto](../../../proto/fjarm/users/v1/user.proto))
- Resource with: `UserId`, `UserHandle`, `UserEmailAddress`, `UserAvatar`, `UserPassword`, `EntityTag`.
- `UserPassword` is `INPUT_ONLY`, `UserId` is `OUTPUT_ONLY`.

## Key observations
1. The registration draft already imports `fjarm.users.v1.UserEmailAddress` and `fjarm.users.v1.UserPassword` — it depends on user types but introduces its own domain concepts (tokens, credentials).
2. Registration is a multi-step workflow with ephemeral state (tokens with expiry). CRUD is stateless request/response.
3. The project already follows a "one service per bounded context" pattern (users vs. authentication).

# Tasks Checklist

- [x] Read `registration.md` draft
- [x] Read all `proto/fjarm/users/v1/*.proto` files
- [x] Read `proto/fjarm/authentication/v1/authentication_service.proto` for precedent
- [x] Analyze domain boundaries and service responsibilities
- [x] Write architectural recommendation

# Verification

- N/A — this is an architectural analysis, not a code change.
