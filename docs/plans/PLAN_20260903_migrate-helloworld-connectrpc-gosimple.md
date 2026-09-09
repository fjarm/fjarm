# Objective

Migrate the ConnectRPC handlers in `api/internal/helloworld/v1/` from `buf.build/gen/go/fjarm/fjarm/connectrpc/go` to `buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple`. This follows the same pattern as the completed users service migration ([PLAN_20260903_migrate-users-connectrpc-gosimple.md](file:///Users/jeremymuhia/development/181_create_session_handler_v0/docs/plans/PLAN_20260903_migrate-users-connectrpc-gosimple.md)), but requires additional care because the helloworld handler accesses **request headers** via `connect.CallInfoForHandlerContext(ctx)` and the test sets headers via `connect.NewRequest().Header().Set()`.

# Research & Context

## Files requiring changes

| File | Change summary |
|---|---|
| [`connect_rpc_handler.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/helloworld/v1/internal/helloworld/connect_rpc_handler.go) | Update `GetHelloWorld` signature to plain protobuf; add `UnimplementedHelloWorldServiceHandler` embed; replace `req.Msg` with `req`; import `helloworldv1connect` from `gosimple` |
| [`connect_rpc_handler_test.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/helloworld/v1/internal/helloworld/connect_rpc_handler_test.go) | Update `helloworldv1connect` import to `gosimple`; replace `connect.NewRequest` + header setting with `connect.NewClientContext` for header injection; update response access from `output.Msg` to `output` |
| [`main.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/helloworld/v1/cmd/helloworld/main.go) | Update `helloworldv1connect` import to `gosimple` |
| `BUILD.bazel` files | Will be auto-updated by Gazelle |

## Header access pattern (key difference from the users migration)

The helloworld handler uses `connect.CallInfoForHandlerContext(ctx)` to read request headers. This API is **context-based** and remains functional with `gosimple` because the `gosimple`-generated `NewHelloWorldServiceHandler` still populates the handler context with call info.

On the **client** side, the current test sets headers via:
```go
req := connect.NewRequest(msg)
req.Header().Set(key, value)
output, err := client.GetHelloWorld(ctx, req)
```

With `gosimple`, the client no longer accepts `*connect.Request`. Instead, use `connect.NewClientContext` to inject headers:
```go
ctx, callInfo := connect.NewClientContext(context.Background())
callInfo.RequestHeader().Set(key, value)
output, err := client.GetHelloWorld(ctx, msg)
```

The response is also returned as a plain protobuf message, so `output.Msg.GetOutput()` becomes `output.GetOutput()`.

## Important constraints

- **Do NOT edit `BUILD.bazel` files manually.** Run `bazel run //:gazelle` instead.
- **Do NOT run `go mod tidy` directly.** Use `bazel run @rules_go//go -- mod tidy` instead.
- The `"connectrpc.com/connect"` import must be **retained** in both `connect_rpc_handler.go` (for `connect.CallInfoForHandlerContext`, `connect.NewError`, etc.) and `connect_rpc_handler_test.go` (for `connect.WithGRPC`, `connect.NewClientContext`).

# Tasks Checklist

## 1. Update `connect_rpc_handler.go`

- [ ] **1a.** Change the `helloworldv1connect` import from `connectrpc/go` to `connectrpc/gosimple`:
  ```go
  // Before:
  "buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/helloworld/v1/helloworldv1connect"

  // After:
  "buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/helloworld/v1/helloworldv1connect"
  ```
- [ ] **1b.** Embed `helloworldv1connect.UnimplementedHelloWorldServiceHandler` in the `ConnectRPCHandler` struct:
  ```go
  type ConnectRPCHandler struct {
      helloworldv1connect.UnimplementedHelloWorldServiceHandler
      domain    getHelloWorlder
      logger    *slog.Logger
      validator protovalidate.Validator
  }
  ```
- [ ] **1c.** Update `GetHelloWorld` method signature:
  ```go
  // Before:
  func (h *ConnectRPCHandler) GetHelloWorld(
      ctx context.Context,
      req *connect.Request[pb.GetHelloWorldRequest],
  ) (*connect.Response[pb.GetHelloWorldResponse], error) {

  // After:
  func (h *ConnectRPCHandler) GetHelloWorld(
      ctx context.Context,
      req *pb.GetHelloWorldRequest,
  ) (*pb.GetHelloWorldResponse, error) {
  ```
- [ ] **1d.** Update `GetHelloWorld` method body:
  - Replace `req.Msg` with `req` in all usages (logging, validation, field access)
  - Replace the `connect.NewResponse(&pb.GetHelloWorldResponse{...})` return with just `&pb.GetHelloWorldResponse{...}`
  - Keep `connect.CallInfoForHandlerContext(ctx)` and `connect.NewError(...)` calls unchanged

## 2. Update `connect_rpc_handler_test.go`

- [ ] **2a.** Change the `helloworldv1connect` import from `connectrpc/go` to `connectrpc/gosimple`:
  ```go
  // Before:
  "buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/helloworld/v1/helloworldv1connect"

  // After:
  "buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/helloworld/v1/helloworldv1connect"
  ```
- [ ] **2b.** Update the test body in `TestConnectRPCHandler_GetHelloWorld_gRPCClient`. Replace the `connect.NewRequest` + header setting + wrapped response pattern:
  ```go
  // Before:
  msg := &pb.GetHelloWorldRequest{Input: &pb.HelloWorldInput{Input: &tc.input}}
  req := connect.NewRequest(msg)
  req.Header().Set(tc.header[0], tc.header[1])

  output, err := client.GetHelloWorld(context.Background(), req)
  if err != nil && !tc.err {
      t.Errorf("GetHelloWorld got an unexpected error: %v", err)
  }
  if err == nil && tc.err {
      t.Errorf("GetHelloWorld expected an error but got none: %v", output.Msg.GetOutput().GetOutput())
  }

  if !tc.err && output.Msg.GetOutput().GetOutput() != tc.expected {
      t.Errorf("GetHelloWorld got: %v, want: %v", output.Msg.GetOutput().GetOutput(), tc.expected)
  }

  // After:
  msg := &pb.GetHelloWorldRequest{Input: &pb.HelloWorldInput{Input: &tc.input}}
  ctx, callInfo := connect.NewClientContext(context.Background())
  callInfo.RequestHeader().Set(tc.header[0], tc.header[1])

  output, err := client.GetHelloWorld(ctx, msg)
  if err != nil && !tc.err {
      t.Errorf("GetHelloWorld got an unexpected error: %v", err)
  }
  if err == nil && tc.err {
      t.Errorf("GetHelloWorld expected an error but got none: %v", output.GetOutput().GetOutput())
  }

  if !tc.err && output.GetOutput().GetOutput() != tc.expected {
      t.Errorf("GetHelloWorld got: %v, want: %v", output.GetOutput().GetOutput(), tc.expected)
  }
  ```
- [ ] **2c.** Retain the `"connectrpc.com/connect"` import — it is still needed for `connect.WithGRPC()` and `connect.NewClientContext`.

## 3. Update `api/internal/helloworld/v1/cmd/helloworld/main.go`

- [ ] **3a.** Change the `helloworldv1connect` import from `connectrpc/go` to `connectrpc/gosimple`:
  ```go
  // Before:
  "buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/helloworld/v1/helloworldv1connect"

  // After:
  "buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/helloworld/v1/helloworldv1connect"
  ```

## 4. Update `go.mod` and Bazel dependencies

- [ ] **4a.** Run `bazel run @rules_go//go -- mod tidy` from the repository root.
- [ ] **4b.** Run `bazel run //:gazelle` from the repository root to regenerate all `BUILD.bazel` files.

## 5. Build and test

- [ ] **5a.** Run `bazel build //api/internal/helloworld/v1/...` to verify compilation.
- [ ] **5b.** Run `bazel test //api/internal/helloworld/v1/...` to verify all tests pass.
- [ ] **5c.** Run `bazel test //api/...` to verify no regressions across the full API.

# Verification

```bash
# Build the helloworld service
bazel build //api/internal/helloworld/v1/...

# Run all helloworld tests
bazel test //api/internal/helloworld/v1/...

# Run full API tests to check for regressions
bazel test //api/...

# Verify no stale connectrpc/go references remain in helloworld
grep -rn 'connectrpc/go/' api/internal/helloworld/ | grep -v 'connectrpc/gosimple' | grep -v 'BUILD.bazel'
# Expected: no output
```

---

## [Update: 2026-09-04 03:41] - ctx.Value vs connect.CallInfoForHandlerContext analysis

* **Context**: Question raised whether the helloworld handler could use `ctx.Value(tracing.RequestIDKey)` (like the users service) instead of `connect.CallInfoForHandlerContext(ctx)` to read the request ID.
* **Finding**: The two patterns serve different purposes and are **not interchangeable** in this case:
  - **Users handler** ([`connect_rpc_handler.go:39`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/connect_rpc_handler.go#L39)): Uses `ctx.Value(tracing.RequestIDKey)` purely for **log decoration**. The tracing interceptor ([`connect_rpc_request_id_interceptor.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/tracing/v1/pkg/interceptor/connect_rpc_request_id_interceptor.go)) validates the header and logs it, but **does not** call `context.WithValue` to propagate the value. So `ctx.Value(tracing.RequestIDKey)` always returns `nil` in the users handler — it's harmless but ineffective.
  - **Helloworld handler** ([`connect_rpc_handler.go:34-42`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/helloworld/v1/internal/helloworld/connect_rpc_handler.go#L34-L42)): Uses `connect.CallInfoForHandlerContext(ctx)` to read the `request-id` header for **control flow** — it returns an error if the header is missing and uses the value to enrich the logger. This pattern works correctly because `CallInfoForHandlerContext` accesses real request metadata from the handler context.
* **Decision**: Keep `connect.CallInfoForHandlerContext(ctx)` in the helloworld handler. It is the correct gosimple-compatible API for accessing request headers in a handler when the request wrapper is no longer available. No changes to tasks.

---

## [Update: 2026-09-04 03:48] - Dependency on interceptor context propagation plan

* **Context**: A separate plan has been created — [PLAN_20260903_interceptor-context-propagation.md](file:///Users/jeremymuhia/development/181_create_session_handler_v0/docs/plans/PLAN_20260903_interceptor-context-propagation.md) — to update the tracing interceptor to propagate the request ID via `context.WithValue` and simplify all handlers to use `tracing.RequestIDFromContext(ctx)`.
* **Impact**: If the interceptor plan is executed **before** this gosimple migration plan, the helloworld handler's `connect.CallInfoForHandlerContext` boilerplate (tasks 1c/1d) will already have been replaced with `tracing.RequestIDFromContext(ctx)`, making the migration simpler. If executed **after**, the interceptor plan will clean up the `CallInfoForHandlerContext` usage post-migration.
* **Recommendation**: Execute the interceptor context propagation plan first, then this gosimple migration plan.

---

## [Update: 2026-09-04 05:50] - Implementation and Verification Complete
* **Context**: All tasks to migrate the helloworld service (`api/internal/helloworld/v1/`) to `connectrpc/gosimple` have been completed. Handlers, tests, and entrypoints were updated, and BUILD.bazel files were regenerated via Gazelle. `go.mod` and `MODULE.bazel` automatically pruned the unused `buf.build/gen/go/fjarm/fjarm/connectrpc/go` dependency.
* **Objective**: Mark all migration tasks complete and record successful verification.
* **Changes to Tasks**: Completed all handler updates, test client refactoring to `connect.NewClientContext`, entrypoint updates, Gazelle dependency sync, and test suites.
* **Revised Task Checklist**:
    - [x] Task 1a: Change `helloworldv1connect` import to `gosimple` in `connect_rpc_handler.go`
    - [x] Task 1b: Embed `helloworldv1connect.UnimplementedHelloWorldServiceHandler` in `ConnectRPCHandler`
    - [x] Task 1c: Update `GetHelloWorld` method signature to plain protobuf types (`*pb.GetHelloWorldRequest`, `*pb.GetHelloWorldResponse`)
    - [x] Task 1d: Update `GetHelloWorld` method body to use `req` directly and return plain protobuf response
    - [x] Task 2a: Change `helloworldv1connect` import to `gosimple` in `connect_rpc_handler_test.go`
    - [x] Task 2b: Update test calls to use `connect.NewClientContext` and plain protobuf requests/responses
    - [x] Task 2c: Retain `"connectrpc.com/connect"` import in test file
    - [x] Task 3a: Change `helloworldv1connect` import to `gosimple` in `cmd/helloworld/main.go`
    - [x] Task 4a: Run `bazel run @rules_go//go -- mod tidy`
    - [x] Task 4b: Run `bazel run //:gazelle` to update `BUILD.bazel` files
    - [x] Task 5a: Run `bazel build //api/internal/helloworld/v1/...`
    - [x] Task 5b: Run `bazel test //api/internal/helloworld/v1/...`
    - [x] Task 5c: Run `bazel test //api/...`
