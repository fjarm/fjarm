# Objective

Migrate the ConnectRPC handlers in `api/internal/users/v1/internal/users/` (and the dependent wiring in `api/internal/users/v1/cmd/users/main.go`) from the `buf.build/gen/go/fjarm/fjarm/connectrpc/go` generated code to `buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple`. The `gosimple` plugin generates service handler interfaces with simplified method signatures that accept and return plain protobuf messages instead of `connect.Request[T]`/`connect.Response[T]` wrappers. An existing migration for the authentication service serves as the reference pattern.

# Research & Context

## What `gosimple` changes

The standard `connectrpc/go` plugin generates handler interfaces like:

```go
func (h *Handler) CreateUser(
    ctx context.Context,
    req *connect.Request[userspb.CreateUserRequest],
) (*connect.Response[userspb.CreateUserResponse], error)
```

The `gosimple` plugin generates simplified interfaces like:

```go
func (h *Handler) CreateUser(
    ctx context.Context,
    req *userspb.CreateUserRequest,
) (*userspb.CreateUserResponse, error)
```

Key differences:
- Method parameters use `*userspb.XxxRequest` instead of `*connect.Request[userspb.XxxRequest]`
- Method return types use `*userspb.XxxResponse` instead of `*connect.Response[userspb.XxxResponse]`
- No need to access `req.Msg` — the request proto is passed directly
- No need to call `connect.NewResponse()` — the response proto is returned directly
- The handler struct should embed `UnimplementedUserServiceHandler` from the `connectrpc/go` package (this provides default no-op implementations for unimplemented methods)

## Dependency availability

- `gosimple` is **already** present in [`go.mod`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/go.mod#L7) at version `v1.20.0-20260811010424-ad8df5ebf10d.1`
- The Bazel external repo `build_buf_gen_go_fjarm_fjarm_connectrpc_gosimple` is **already** declared in [`MODULE.bazel`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/MODULE.bazel#L22)
- No new dependencies need to be added to `go.mod` or `MODULE.bazel`

## Reference implementation (authentication service)

The authentication service has already been migrated and serves as the canonical pattern:

- **Handler struct** ([`connect_rpc_handler.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/authentication/v1/internal/authentication/connect_rpc_handler.go#L15-L19)):
  Embeds `authenticationv1connect.UnimplementedAuthenticationServiceHandler` from `connectrpc/go`
- **Handler method** ([`connect_rpc_handler_create_session.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/authentication/v1/internal/authentication/connect_rpc_handler_create_session.go#L9-L14)):
  Uses simplified `(ctx, *proto.Request) (*proto.Response, error)` signature
- **Wiring / main.go** ([`main.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/authentication/v1/cmd/authentication/main.go#L13)):
  Imports `NewAuthenticationServiceHandler` from `connectrpc/gosimple`

## Files requiring changes

| File | Change summary |
|---|---|
| [`connect_rpc_handler.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/connect_rpc_handler.go) | Update handler method signatures; embed `UnimplementedUserServiceHandler`; remove `connect` import |
| [`connect_rpc_handler_test.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/connect_rpc_handler_test.go) | Update import from `connectrpc/go` → `connectrpc/gosimple`; update test client calls |
| [`main.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/cmd/users/main.go) | Update import from `connectrpc/go` → `connectrpc/gosimple` |
| `BUILD.bazel` files | Will be auto-updated by Gazelle |

## Important constraints

- **Do NOT edit `BUILD.bazel` files manually.** Run `bazel run //:gazelle` instead.
- **Do NOT run `go mod tidy` directly.** Use `bazel run @rules_go//go -- mod tidy` instead.
- The `"connectrpc.com/connect"` import should be **removed** from `connect_rpc_handler.go` since the handler methods will no longer reference `connect.Request`, `connect.Response`, `connect.NewResponse`, or `connect.NewError`. However, `connect` error codes **are still needed** for mapping domain errors to RPC error codes — use `connect.NewError(...)` to return errors. Therefore, the `"connectrpc.com/connect"` import must be **retained** in `connect_rpc_handler.go`.
- The `userDomain` internal interface (defined in `connect_rpc_handler.go`) already works with plain protobuf types — **no changes** are needed to `domain.go` or any other domain/repository files.

# Tasks Checklist

## 1. Update `connect_rpc_handler.go`

- [ ] **1a.** Add import for `usersv1connect "buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/users/v1/usersv1connect"`. This provides the `UnimplementedUserServiceHandler` type. (Note: this import may already be implicitly available; verify the generated package contents.)
- [ ] **1b.** Embed `usersv1connect.UnimplementedUserServiceHandler` in the `ConnectRPCHandler` struct, following the authentication service pattern:
  ```go
  type ConnectRPCHandler struct {
      usersv1connect.UnimplementedUserServiceHandler
      domain    userDomain
      logger    *slog.Logger
      validator protovalidate.Validator
  }
  ```
- [ ] **1c.** Update `CreateUser` method signature — remove `connect.Request`/`connect.Response` wrappers:
  ```go
  // Before:
  func (h *ConnectRPCHandler) CreateUser(
      ctx context.Context,
      req *connect.Request[userspb.CreateUserRequest],
  ) (*connect.Response[userspb.CreateUserResponse], error) {

  // After:
  func (h *ConnectRPCHandler) CreateUser(
      ctx context.Context,
      req *userspb.CreateUserRequest,
  ) (*userspb.CreateUserResponse, error) {
  ```
- [ ] **1d.** Update `CreateUser` method body:
  - Replace all occurrences of `req.Msg` with `req` (the request proto is now passed directly)
  - Replace `connect.NewResponse(res)` return with just `res`
  - Replace `connect.NewError(code, err)` error returns. Since the `gosimple` handler wrapper translates returned errors into Connect errors, use `connect.NewError(...)` as before — these still work correctly. Retain the `"connectrpc.com/connect"` import.
- [ ] **1e.** Update `GetUser` method signature:
  ```go
  // Before:
  func (h *ConnectRPCHandler) GetUser(
      ctx context.Context,
      req *connect.Request[userspb.GetUserRequest],
  ) (*connect.Response[userspb.GetUserResponse], error) {

  // After:
  func (h *ConnectRPCHandler) GetUser(
      ctx context.Context,
      req *userspb.GetUserRequest,
  ) (*userspb.GetUserResponse, error) {
  ```
- [ ] **1f.** Update `UpdateUser` method signature:
  ```go
  // Before:
  func (h *ConnectRPCHandler) UpdateUser(
      ctx context.Context,
      req *connect.Request[userspb.UpdateUserRequest],
  ) (*connect.Response[userspb.UpdateUserResponse], error) {

  // After:
  func (h *ConnectRPCHandler) UpdateUser(
      ctx context.Context,
      req *userspb.UpdateUserRequest,
  ) (*userspb.UpdateUserResponse, error) {
  ```
- [ ] **1g.** Update `DeleteUser` method signature:
  ```go
  // Before:
  func (h *ConnectRPCHandler) DeleteUser(
      ctx context.Context,
      req *connect.Request[userspb.DeleteUserRequest],
  ) (*connect.Response[userspb.DeleteUserResponse], error) {

  // After:
  func (h *ConnectRPCHandler) DeleteUser(
      ctx context.Context,
      req *userspb.DeleteUserRequest,
  ) (*userspb.DeleteUserResponse, error) {
  ```

## 2. Update `connect_rpc_handler_test.go`

- [ ] **2a.** Change the `usersv1connect` import from the old to the new package:
  ```go
  // Before:
  "buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/users/v1/usersv1connect"

  // After:
  "buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/users/v1/usersv1connect"
  ```
  The `NewUserServiceHandler` and `NewUserServiceClient` functions will now come from the `gosimple` package. Their call signatures at the call sites (line 46 and line 182) remain the same.
- [ ] **2b.** Update the test client calls. In `TestConnectRPCHandler_CreateUser_gRPCClient`, the `client.CreateUser` call currently sends a `connect.NewRequest(req)` wrapped request and receives a `connect.Response`. With `gosimple`, the client API also uses plain protobuf messages:
  ```go
  // Before:
  _, err := client.CreateUser(context.Background(), connect.NewRequest(req))

  // After:
  _, err := client.CreateUser(context.Background(), req)
  ```
  The `connect.NewRequest` wrapper is no longer needed.
- [ ] **2c.** If the `"connectrpc.com/connect"` import is no longer used in the test file after removing `connect.NewRequest`, check whether it is still needed for `connect.WithGRPC()`, `connect.Code`, and `connect.CodeOf()`. If it is, keep the import. If not, remove it. (It **is** still needed for `connect.WithGRPC()`, `connect.Code`, `connect.CodeOf`, `connect.CodeInvalidArgument`, etc.)

## 3. Update `api/internal/users/v1/cmd/users/main.go`

- [ ] **3a.** Change the blank import from old to new:
  ```go
  // Before:
  _ "buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/users/v1/usersv1connect"

  // After:
  _ "buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/users/v1/usersv1connect"
  ```

## 4. Update `go.mod` and Bazel dependencies

- [ ] **4a.** Run `bazel run @rules_go//go -- mod tidy` from the repository root to clean up `go.mod`/`go.sum`. This may remove `buf.build/gen/go/fjarm/fjarm/connectrpc/go` if no other packages in the module still import it, or it may keep it if the `connect_rpc_handler.go` still imports the `usersv1connect` package from it for `UnimplementedUserServiceHandler`.
- [ ] **4b.** Run `bazel run //:gazelle` from the repository root to regenerate all `BUILD.bazel` files. This will update the `deps` in:
  - `api/internal/users/v1/internal/users/BUILD.bazel` (the `go_test` target will switch from `@build_buf_gen_go_fjarm_fjarm_connectrpc_go` to `@build_buf_gen_go_fjarm_fjarm_connectrpc_gosimple`, and the `go_library` target will gain `@build_buf_gen_go_fjarm_fjarm_connectrpc_go//fjarm/users/v1/usersv1connect` for the embedded `Unimplemented` type)
  - `api/internal/users/v1/cmd/users/BUILD.bazel` (will switch to `@build_buf_gen_go_fjarm_fjarm_connectrpc_gosimple`)

## 5. Build and test

- [ ] **5a.** Run `bazel build //api/internal/users/v1/...` to verify compilation.
- [ ] **5b.** Run `bazel test //api/internal/users/v1/...` to verify all tests pass.

# Verification

```bash
# Build the users service and its internal packages
bazel build //api/internal/users/v1/...

# Run all tests in the users service
bazel test //api/internal/users/v1/...

# Verify that no stale references to the old connectrpc/go import remain in the users directory
grep -rn 'connectrpc/go/' api/internal/users/v1/ | grep -v 'connectrpc/gosimple' | grep -v 'BUILD.bazel'
# Expected: no output (all references should now be gosimple)
```

---

## [Update: 2026-09-03 06:17] - UnimplementedUserServiceHandler available from gosimple

* **Context**: The authentication module's [`connect_rpc_handler.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/authentication/v1/internal/authentication/connect_rpc_handler.go) has been updated so that `UnimplementedAuthenticationServiceHandler` is imported from `connectrpc/gosimple` instead of `connectrpc/go`. This means the `gosimple` generated package also exports the `Unimplemented*` types.
* **Objective**: Simplify the users migration so that **all** generated code references — including `UnimplementedUserServiceHandler` — come from the single `connectrpc/gosimple` package. No `connectrpc/go` import is needed in any users file.
* **Changes to Tasks**:
  - **Task 1a is superseded.** Do NOT add an import of `connectrpc/go` for `UnimplementedUserServiceHandler`. Instead, add an import of `connectrpc/gosimple`:
    ```go
    usersv1connect "buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/users/v1/usersv1connect"
    ```
  - **Task 1b is unchanged** in intent but the import source changes to `gosimple` as described above.
  - **Task 4a** (`go mod tidy`): After this migration, `buf.build/gen/go/fjarm/fjarm/connectrpc/go` may be removed from `go.mod` entirely if no other packages in the module still import it.
  - **Task 4b** (Gazelle): The `go_library` target in `api/internal/users/v1/internal/users/BUILD.bazel` will gain `@build_buf_gen_go_fjarm_fjarm_connectrpc_gosimple//fjarm/users/v1/usersv1connect` (not `@build_buf_gen_go_fjarm_fjarm_connectrpc_go`).
  - **Verification grep** is simplified: no need to exclude `UnimplementedUserServiceHandler` from the stale-reference check since it now also comes from `gosimple`.
* **Revised Task Checklist**:
    - [ ] In task 1a, import `usersv1connect` from `buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/users/v1/usersv1connect` (not `connectrpc/go`)
    - [ ] In task 1b, embed `usersv1connect.UnimplementedUserServiceHandler` (sourced from `gosimple`)
    - [ ] All other tasks (1c–1g, 2a–2c, 3a, 4a–4b, 5a–5b) remain unchanged

---

## [Update: 2026-09-04 03:25] - Implementation and Verification Complete
* **Context**: All tasks to migrate the ConnectRPC handlers in `api/internal/users/v1/internal/users/` and the wiring in `api/internal/users/v1/cmd/users/main.go` from `buf.build/gen/go/fjarm/fjarm/connectrpc/go` to `buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple` have been implemented. `go mod tidy` and Gazelle were executed via Bazel, and all tests pass.
* **Objective**: Mark all migration tasks as completed and record successful verification.
* **Changes to Tasks**: Completed all handler updates, test updates, cmd entrypoint updates, dependency syncs, and verification steps.
* **Revised Task Checklist**:
    - [x] In task 1a, import `usersv1connect` from `buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/users/v1/usersv1connect` (not `connectrpc/go`)
    - [x] In task 1b, embed `usersv1connect.UnimplementedUserServiceHandler` (sourced from `gosimple`)
    - [x] Task 1c: Update `CreateUser` method signature to plain protobuf types (`*userspb.CreateUserRequest`, `*userspb.CreateUserResponse`)
    - [x] Task 1d: Update `CreateUser` method body (`req` instead of `req.Msg`, return `res` directly, retain `connect.NewError`)
    - [x] Task 1e: Update `GetUser` method signature to plain protobuf types
    - [x] Task 1f: Update `UpdateUser` method signature to plain protobuf types
    - [x] Task 1g: Update `DeleteUser` method signature to plain protobuf types
    - [x] Task 2a: Update `usersv1connect` import in `connect_rpc_handler_test.go` to `connectrpc/gosimple`
    - [x] Task 2b: Update test client calls to pass `req` directly instead of wrapping with `connect.NewRequest(req)`
    - [x] Task 2c: Retain `"connectrpc.com/connect"` import in test file for `connect.WithGRPC()` and `connect.CodeOf()`
    - [x] Task 3a: Update blank import in `api/internal/users/v1/cmd/users/main.go` to `connectrpc/gosimple`
    - [x] Task 4a: Run `bazel run @rules_go//go -- mod tidy`
    - [x] Task 4b: Run `bazel run //:gazelle` to update `BUILD.bazel` files
    - [x] Task 5a: Run `bazel build //api/internal/users/v1/...`
    - [x] Task 5b: Run `bazel test //api/internal/users/v1/...`
