# Objective

Document and implement the correct authorization flow for the `DeleteUser` RPC so that:

1. The caller is authenticated via an access token (JWT) carried in request metadata — not via credentials in the request body.
2. The handler authorizes the request by verifying the caller's identity matches the `user_id` in the `DeleteUserRequest`.
3. Upon successful deletion, the backend invalidates all sessions for the deleted user.
4. The client performs purely local cleanup (discard tokens, navigate to login) with no additional RPCs required.

# Research & Context

## Design Decisions Already Made

- [`CreateUserResponse`](proto/fjarm/users/v1/create_user_response.proto) and [`DeleteUserResponse`](proto/fjarm/users/v1/delete_user_response.proto) are empty messages — `google.rpc.Status` has been removed from both. Errors are communicated exclusively via ConnectRPC's transport-level error mechanism (`connect.NewError`).
- [`DeleteUserRequest`](proto/fjarm/users/v1/delete_user_request.proto) identifies the resource by `UserId`, consistent with [`GetUserRequest`](proto/fjarm/users/v1/get_user_request.proto) and [`UpdateUserRequest`](proto/fjarm/users/v1/update_user_request.proto). Credentials (email+password) do NOT belong in the request body.
- The `UserId` field is `OUTPUT_ONLY` on the [`User`](proto/fjarm/users/v1/user.proto) message (server-generated per AIP-122), but clients learn their `user_id` after login via `GetUser` or from the access token claims.

## Existing Auth Infrastructure

- [`AuthenticationService.CreateSession`](proto/fjarm/authentication/v1/authentication_service.proto) accepts email+password and returns a [`Session`](proto/fjarm/authentication/v1/session.proto) containing an `AccessToken` (JWT) and `RefreshToken`.
- The interceptor chain already includes a [constant-timing interceptor](api/internal/obfuscation/v1/pkg/interceptor/connect_rpc_constant_timing_interceptor.go) and a [request-ID interceptor](api/internal/tracing/v1/pkg/interceptor/connect_rpc_request_id_interceptor.go). An auth interceptor does not yet exist.
- The Android client's [`ServerTransportModule`](android/library/serverTransport/src/main/java/xyz/fjarm/servertransport/ServerTransportModule.kt) has TODOs for adding an `Authenticator` and `Interceptor` to the `OkHttpClient`.
- The Android client's [`AttemptLoginUseCaseImpl`](android/feature/login/library/src/main/java/xyz/fjarm/loginlibrary/AttemptLoginUseCaseImpl.kt) has a TODO for persisting session tokens to an encrypted DataStore.

## Correct DeleteUser Authorization Flow

```
┌────────┐                              ┌──────────────────┐                        ┌────────────────┐
│ Client │                              │  Auth Interceptor │                        │ DeleteUser     │
│        │                              │  (ConnectRPC)     │                        │ Handler        │
└───┬────┘                              └────────┬─────────┘                        └───────┬────────┘
    │                                            │                                          │
    │  DeleteUserRequest + Bearer <jwt>          │                                          │
    │───────────────────────────────────────────>│                                          │
    │                                            │                                          │
    │                                            │ 1. Validate JWT signature & expiry       │
    │                                            │ 2. Extract caller identity (user_id)     │
    │                                            │ 3. Inject identity into context          │
    │                                            │                                          │
    │                                            │  If JWT invalid:                         │
    │  <── UNAUTHENTICATED ─────────────────────│                                          │
    │                                            │                                          │
    │                                            │  If JWT valid:                           │
    │                                            │─────────────────────────────────────────>│
    │                                            │                                          │
    │                                            │          4. Compare caller identity from │
    │                                            │             context with user_id in req  │
    │                                            │                                          │
    │                                            │          If mismatch:                    │
    │  <── PERMISSION_DENIED ──────────────────────────────────────────────────────────────│
    │                                            │                                          │
    │                                            │          If match:                       │
    │                                            │          5. Hard-delete user             │
    │                                            │          6. Invalidate all sessions      │
    │                                            │                                          │
    │  <── DeleteUserResponse {} ──────────────────────────────────────────────────────────│
    │                                            │                                          │
    │  7. Discard local tokens                   │                                          │
    │  8. Navigate to login screen               │                                          │
    │                                            │                                          │
```

## Key Principles

- **Authentication is a cross-cutting concern.** It belongs in an interceptor, not in individual handler request messages. The interceptor validates the JWT and injects the caller identity into the request context. Every protected RPC benefits from this without duplicating logic.
- **Authorization is a handler-level concern.** The handler compares the caller identity (from context) against the resource being acted upon. For `DeleteUser`, the caller's `user_id` (extracted from the JWT by the interceptor) must match the `user_id` in the `DeleteUserRequest`.
- **Credentials never appear in resource operation messages.** `UserEmailAddress` and `UserPassword` are used only in [`CreateUserRequest`](proto/fjarm/users/v1/create_user_request.proto) (registration) and [`CreateSessionRequest`](proto/fjarm/authentication/v1/create_session_request.proto) (login). They do not appear in Get, Update, or Delete request messages.
- **No additional RPCs after DeleteUser.** The backend invalidates sessions server-side. The client discards local tokens and navigates to the login screen. If the client crashes before cleanup, the dead tokens will be rejected on next use (`UNAUTHENTICATED`), which triggers the same login redirect.

## References

- [AIP-135: Standard methods: Delete](https://google.aip.dev/135) — Delete requests identify the resource by name/ID; response should be empty for hard deletes.
- [AIP-122: Resource names](https://google.aip.dev/122) — Resource IDs are server-generated (`OUTPUT_ONLY`) but known to clients after creation.
- [AIP-193: Errors](https://google.aip.dev/193) — Error signaling via `google.rpc.Status` at the transport level, not in response payloads.

# Tasks Checklist

## Server-Side

- [ ] **1. Create a ConnectRPC auth interceptor** that validates JWTs from the `Authorization` header, extracts the caller's user identity, and injects it into the Go `context.Context`. On failure, return `connect.CodeUnauthenticated`.
- [ ] **2. Register the auth interceptor** in the users service's interceptor chain (in `main.go`), positioned after the tracing interceptor. Determine which RPCs are protected (e.g., `GetUser`, `UpdateUser`, `DeleteUser`) vs. unprotected (`CreateUser`).
- [ ] **3. Implement the `DeleteUser` handler** with an authorization check: extract the caller identity from the context and compare it to the `user_id` in the `DeleteUserRequest`. Return `connect.CodePermissionDenied` on mismatch.
- [ ] **4. Implement session invalidation on user deletion.** After the user record is hard-deleted, revoke/delete all sessions associated with that `user_id` so that any outstanding access and refresh tokens become invalid.

## Client-Side (Android)

- [ ] **5. Persist session tokens** to an encrypted DataStore after login (resolves the existing TODO in `AttemptLoginUseCaseImpl`).
- [ ] **6. Add an OkHttp interceptor** that attaches the `Authorization: Bearer <jwt>` header to every outgoing request for protected RPCs (resolves the existing TODO in `ServerTransportModule`).
- [ ] **7. Implement the delete-user flow.** On successful `DeleteUserResponse`:
  - Clear tokens from the encrypted DataStore.
  - Navigate to the login screen.
- [ ] **8. Handle `UNAUTHENTICATED` globally.** If any RPC returns `UNAUTHENTICATED` (e.g., expired token, server-revoked session), clear local tokens and navigate to the login screen as a fallback.

# Verification

- [ ] Unit test the auth interceptor: valid JWT → context populated; expired/malformed JWT → `UNAUTHENTICATED`.
- [ ] Unit test `DeleteUser` handler: caller identity matches `user_id` → success; mismatch → `PERMISSION_DENIED`.
- [ ] Integration test: after `DeleteUser` succeeds, subsequent RPCs with the old access token return `UNAUTHENTICATED`.
- [ ] Integration test: after `DeleteUser` succeeds, attempting to refresh the session with the old refresh token returns `UNAUTHENTICATED`.
