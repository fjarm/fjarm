# Objective

Update the ConnectRPC request-ID tracing interceptor to propagate the validated `request-id` into the request context via `context.WithValue`. Then simplify the helloworld handler to read the request ID from the context (matching the pattern already used by the users service), removing the `connect.CallInfoForHandlerContext` boilerplate. Additionally, introduce a proper unexported context key type to avoid the antipattern of using a bare string as a context key.

# Research & Context

## Current state

The [tracing interceptor](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/tracing/v1/pkg/interceptor/connect_rpc_request_id_interceptor.go) already:
1. Extracts the `request-id` header from every incoming request (line 29)
2. Validates it's non-empty (line 30)
3. Rejects the request with `tracing.ErrRequestIDNotFound` if missing (line 36)
4. Logs the validated value (lines 39–43)
5. Calls `next(ctx, req)` — but **discards** the validated request ID instead of propagating it via context

## Consumers of the request ID

Four call sites read the request ID from context via `ctx.Value(tracing.RequestIDKey)`:

| File | Line | Purpose |
|---|---|---|
| [`connect_rpc_handler.go` (users)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/connect_rpc_handler.go#L41) | 41 | Log decoration — currently returns `nil` |
| [`domain.go` (users)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/domain.go#L87) | 87 | Log decoration — currently returns `nil` |
| [`in_memory_repository.go` (users)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/in_memory_repository.go#L26) | 26 | Log decoration — currently returns `nil` |
| [`interactor.go` (helloworld)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/helloworld/v1/internal/helloworld/interactor.go#L33) | 33 | Log decoration — currently returns `nil` |

All four currently return `nil` because the interceptor never calls `context.WithValue`.

One additional call site uses `connect.CallInfoForHandlerContext` to read the request ID directly from headers:

| File | Lines | Purpose |
|---|---|---|
| [`connect_rpc_handler.go` (helloworld)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/helloworld/v1/internal/helloworld/connect_rpc_handler.go#L34-L42) | 34–42 | Control flow — returns error if missing; enriches logger |

## Context key antipattern

[`RequestIDKey`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/tracing/request_id.go) is defined as `const RequestIDKey = "request-id"` — a bare `string`. Using plain strings as context keys risks collisions with other packages that might use the same string. The Go standard practice is to use an unexported type:

```go
type contextKey string
const requestIDContextKey contextKey = "request-id"
```

Since `RequestIDKey` is also used as the **HTTP header name** (by the interceptor and tests), it should remain as a public string constant for that purpose. A separate, unexported context key should be introduced for `context.WithValue`/`ctx.Value`, with exported helper functions `RequestIDFromContext(ctx)` and `ContextWithRequestID(ctx, id)` to encapsulate the key.

## Files requiring changes

| File | Change summary |
|---|---|
| [`request_id.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/tracing/request_id.go) | Add unexported `contextKey` type, `requestIDContextKey` constant, and exported helper functions `RequestIDFromContext` and `ContextWithRequestID` |
| [`connect_rpc_request_id_interceptor.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/tracing/v1/pkg/interceptor/connect_rpc_request_id_interceptor.go) | Add `context.WithValue` to propagate request ID before calling `next` (using the new `ContextWithRequestID` helper) |
| [`connect_rpc_request_id_interceptor_test.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/tracing/v1/pkg/interceptor/connect_rpc_request_id_interceptor_test.go) | Add test case verifying the request ID is present in the context passed to `next` |
| [`connect_rpc_handler.go` (helloworld)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/helloworld/v1/internal/helloworld/connect_rpc_handler.go) | Remove `connect.CallInfoForHandlerContext` boilerplate; read request ID via `tracing.RequestIDFromContext(ctx)` |
| [`connect_rpc_handler.go` (users)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/connect_rpc_handler.go) | Replace `ctx.Value(tracing.RequestIDKey)` with `tracing.RequestIDFromContext(ctx)` |
| [`domain.go` (users)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/domain.go) | Replace `ctx.Value(tracing.RequestIDKey)` with `tracing.RequestIDFromContext(ctx)` |
| [`in_memory_repository.go` (users)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/in_memory_repository.go) | Replace `ctx.Value(tracing.RequestIDKey)` with `tracing.RequestIDFromContext(ctx)` |
| [`interactor.go` (helloworld)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/helloworld/v1/internal/helloworld/interactor.go) | Replace `ctx.Value(tracing.RequestIDKey)` with `tracing.RequestIDFromContext(ctx)` |
| `BUILD.bazel` files | Will be auto-updated by Gazelle |

## Important constraints

- **Do NOT edit `BUILD.bazel` files manually.** Run `bazel run //:gazelle` instead.
- **Do NOT run `go mod tidy` directly.** Use `bazel run @rules_go//go -- mod tidy` instead.
- `RequestIDKey` must remain a public `string` constant — it is used as an HTTP header name in the interceptor, tests, and handler tests.

# Tasks Checklist

## 1. Update `api/internal/tracing/request_id.go`

- [ ] **1a.** Add an unexported context key type and constant:
  ```go
  type contextKey string

  const requestIDContextKey contextKey = "request-id"
  ```
- [ ] **1b.** Add exported helper functions:
  ```go
  // ContextWithRequestID returns a new context with the request ID value stored.
  func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
      return context.WithValue(ctx, requestIDContextKey, requestID)
  }

  // RequestIDFromContext extracts the request ID from the context, returning an empty string if not found.
  func RequestIDFromContext(ctx context.Context) string {
      val, _ := ctx.Value(requestIDContextKey).(string)
      return val
  }
  ```

## 2. Update the interceptor to propagate the request ID

- [ ] **2a.** In [`connect_rpc_request_id_interceptor.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/tracing/v1/pkg/interceptor/connect_rpc_request_id_interceptor.go), after validating the request ID and before calling `next`, add context propagation:
  ```go
  // Before:
  res, err := next(ctx, req)

  // After:
  ctx = tracing.ContextWithRequestID(ctx, reqID)
  res, err := next(ctx, req)
  ```

## 3. Update the interceptor test

- [ ] **3a.** In [`connect_rpc_request_id_interceptor_test.go`](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/tracing/v1/pkg/interceptor/connect_rpc_request_id_interceptor_test.go), update the `next` function to capture and verify the context value. Add a test that verifies `tracing.RequestIDFromContext(ctx)` returns the expected request ID when a valid header is present:
  ```go
  next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
      // Store the request ID from context for verification in tests.
      capturedRequestID = tracing.RequestIDFromContext(ctx)
      return nil, nil
  }
  ```
  Then verify in the `"valid_non_empty_request_id"` test case that `capturedRequestID == "abc123"`.

## 4. Simplify the helloworld handler

- [ ] **4a.** In [`connect_rpc_handler.go` (helloworld)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/helloworld/v1/internal/helloworld/connect_rpc_handler.go), remove the `connect.CallInfoForHandlerContext` boilerplate (lines 34–42) and replace with a `tracing.RequestIDFromContext(ctx)` call:
  ```go
  // Before (lines 34-46):
  callInfo, ok := connect.CallInfoForHandlerContext(ctx)
  if !ok {
      return nil, tracing.ErrRequestIDNotFound
  }

  requestID := callInfo.RequestHeader().Get(tracing.RequestIDKey)
  if requestID == "" {
      return nil, tracing.ErrRequestIDNotFound
  }

  logger := h.logger.With(
      slog.String(logkeys.Rpc, helloworldv1connect.HelloWorldServiceGetHelloWorldProcedure),
      slog.String(tracing.RequestIDKey, requestID),
  )

  // After:
  requestID := tracing.RequestIDFromContext(ctx)

  logger := h.logger.With(
      slog.String(logkeys.Rpc, helloworldv1connect.HelloWorldServiceGetHelloWorldProcedure),
      slog.String(tracing.RequestIDKey, requestID),
  )
  ```
  Note: The empty-request-ID validation is no longer needed in the handler because the interceptor already rejects requests without a valid `request-id` header before the handler runs.

## 5. Update all `ctx.Value(tracing.RequestIDKey)` call sites

- [ ] **5a.** In [`connect_rpc_handler.go` (users)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/connect_rpc_handler.go#L41):
  ```go
  // Before:
  slog.Any(tracing.RequestIDKey, ctx.Value(tracing.RequestIDKey)),

  // After:
  slog.String(tracing.RequestIDKey, tracing.RequestIDFromContext(ctx)),
  ```
- [ ] **5b.** In [`domain.go` (users)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/domain.go#L87):
  ```go
  // Before:
  slog.Any(tracing.RequestIDKey, ctx.Value(tracing.RequestIDKey)),

  // After:
  slog.String(tracing.RequestIDKey, tracing.RequestIDFromContext(ctx)),
  ```
- [ ] **5c.** In [`in_memory_repository.go` (users)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/users/v1/internal/users/in_memory_repository.go#L26):
  ```go
  // Before:
  slog.Any(tracing.RequestIDKey, ctx.Value(tracing.RequestIDKey)),

  // After:
  slog.String(tracing.RequestIDKey, tracing.RequestIDFromContext(ctx)),
  ```
- [ ] **5d.** In [`interactor.go` (helloworld)](file:///Users/jeremymuhia/development/181_create_session_handler_v0/api/internal/helloworld/v1/internal/helloworld/interactor.go#L33):
  ```go
  // Before:
  slog.Any(tracing.RequestIDKey, ctx.Value(tracing.RequestIDKey)),

  // After:
  slog.String(tracing.RequestIDKey, tracing.RequestIDFromContext(ctx)),
  ```

## 6. Sync dependencies and build

- [ ] **6a.** Run `bazel run @rules_go//go -- mod tidy` from the repository root.
- [ ] **6b.** Run `bazel run //:gazelle` from the repository root to regenerate all `BUILD.bazel` files.
- [ ] **6c.** Run `bazel build //api/...` to verify compilation.
- [ ] **6d.** Run `bazel test //api/...` to verify all tests pass.

# Verification

```bash
# Build and test the interceptor
bazel test //api/internal/tracing/...

# Build and test the helloworld service
bazel test //api/internal/helloworld/v1/...

# Build and test the users service
bazel test //api/internal/users/v1/...

# Full API test
bazel test //api/...

# Verify no remaining ctx.Value(tracing.RequestIDKey) calls
grep -rn 'ctx.Value(tracing.RequestIDKey)' api/
# Expected: no output

# Verify no remaining CallInfoForHandlerContext calls in handlers
grep -rn 'CallInfoForHandlerContext' api/ | grep -v '_test.go'
# Expected: no output
```

---

## [Update: 2026-09-04 04:09] - Implementation and Verification Complete
* **Context**: Implemented request ID context propagation in the tracing interceptor, encapsulated context keys with helper functions `ContextWithRequestID` and `RequestIDFromContext`, simplified the helloworld handler to extract the request ID from context, updated helloworld handler tests to mount the interceptor, and updated all `ctx.Value(tracing.RequestIDKey)` callers across users and helloworld services.
* **Objective**: Mark all tasks complete and record successful verification.
* **Changes to Tasks**: All tasks implemented and verified against all Bazel targets.
* **Revised Task Checklist**:
    - [x] Task 1a: Add unexported `contextKey` type and `requestIDContextKey` constant to `api/internal/tracing/request_id.go`
    - [x] Task 1b: Add `ContextWithRequestID` and `RequestIDFromContext` helper functions to `api/internal/tracing/request_id.go`
    - [x] Task 2a: In `connect_rpc_request_id_interceptor.go`, propagate `requestID` into context via `tracing.ContextWithRequestID` before calling `next`
    - [x] Task 3a: Update `connect_rpc_request_id_interceptor_test.go` to capture context and verify `tracing.RequestIDFromContext(ctx)` matches expected request ID
    - [x] Task 4a: Simplify helloworld handler in `connect_rpc_handler.go` by removing `connect.CallInfoForHandlerContext` and using `tracing.RequestIDFromContext(ctx)`
    - [x] Task 5a: Replace `ctx.Value(tracing.RequestIDKey)` with `tracing.RequestIDFromContext(ctx)` in users `connect_rpc_handler.go`
    - [x] Task 5b: Replace `ctx.Value(tracing.RequestIDKey)` with `tracing.RequestIDFromContext(ctx)` in users `domain.go`
    - [x] Task 5c: Replace `ctx.Value(tracing.RequestIDKey)` with `tracing.RequestIDFromContext(ctx)` in users `in_memory_repository.go`
    - [x] Task 5d: Replace `ctx.Value(tracing.RequestIDKey)` with `tracing.RequestIDFromContext(ctx)` in helloworld `interactor.go`
    - [x] Task 6a: Run `bazel run @rules_go//go -- mod tidy`
    - [x] Task 6b: Run `bazel run //:gazelle` to update `BUILD.bazel` files
    - [x] Task 6c: Run `bazel build //api/...`
    - [x] Task 6d: Run `bazel test //api/...`

---

## [Update: 2026-09-04 04:13] - Simplify Interceptor Test to Use tc.headers Direct Comparison
* **Context**: The user requested removing `expectedRequestID` from the test table struct in `connect_rpc_request_id_interceptor_test.go` and instead comparing `capturedRequestID` directly against `tc.headers["request-id"]` when `!tc.err`.
* **Objective**: Remove redundant `expectedRequestID` field from test struct and compare against `tc.headers["request-id"]`.
* **Changes to Tasks**: Refactor `connect_rpc_request_id_interceptor_test.go` to remove `expectedRequestID` from the table struct and check `capturedRequestID` against `tc.headers["request-id"]` when `!tc.err`.
* **Revised Task Checklist**:
    - [ ] Remove `expectedRequestID` from the test table struct in `connect_rpc_request_id_interceptor_test.go`
    - [ ] Update assertion to check `if !tc.err && capturedRequestID != tc.headers["request-id"]`
    - [ ] Run `bazel test //api/internal/tracing/...` to verify

---

## [Update: 2026-09-04 04:14] - Interceptor Test Simplification Complete
* **Context**: Successfully removed `expectedRequestID` from the test table struct in `connect_rpc_request_id_interceptor_test.go` and updated the assertion to compare `capturedRequestID` directly against `tc.headers["request-id"]` when `!tc.err`.
* **Objective**: Mark test simplification tasks complete.
* **Changes to Tasks**: Completed test table refactoring and verification.
* **Revised Task Checklist**:
    - [x] Remove `expectedRequestID` from the test table struct in `connect_rpc_request_id_interceptor_test.go`
    - [x] Update assertion to check `if !tc.err && capturedRequestID != tc.headers["request-id"]`
    - [x] Run `bazel test //api/internal/tracing/...` to verify

---

## [Update: 2026-09-04 05:43] - Refactor Interceptor Test to Use Echo Response Headers (Option 1)
* **Context**: The user approved refactoring `connect_rpc_request_id_interceptor_test.go` to use the Echo Response Header pattern instead of using a closure with the `capturedRequestID` variable.
* **Objective**: Remove `capturedRequestID` closure variable and have `next` echo the request ID from context into the response header (`res.Header().Set(tracing.RequestIDKey, tracing.RequestIDFromContext(ctx))`), then assert directly on `res.Header().Get(...)`.
* **Changes to Tasks**: Update `connect_rpc_request_id_interceptor_test.go` to use response header echo.
* **Revised Task Checklist**:
    - [ ] Update `next` in `connect_rpc_request_id_interceptor_test.go` to return a `connect.Response` with the request ID set on response headers
    - [ ] Remove `capturedRequestID` variable and its reset from the test
    - [ ] Update assertion to check `res.Header().Get(tracing.RequestIDKey)` against `tc.headers["request-id"]`
    - [ ] Run `bazel test //api/internal/tracing/...` to verify

---

## [Update: 2026-09-04 05:45] - Echo Response Header Refactoring Complete
* **Context**: Successfully refactored `connect_rpc_request_id_interceptor_test.go` to use the Echo Response Header pattern. `capturedRequestID` was removed, `next` returns a `connect.Response` with the request ID in `res.Header()`, and the assertion inspects `res.Header().Get(tracing.RequestIDKey)` directly.
* **Objective**: Mark Option 1 test refactoring complete.
* **Changes to Tasks**: Completed all test refactoring steps and verified via Bazel.
* **Revised Task Checklist**:
    - [x] Update `next` in `connect_rpc_request_id_interceptor_test.go` to return a `connect.Response` with the request ID set on response headers
    - [x] Remove `capturedRequestID` variable and its reset from the test
    - [x] Update assertion to check `res.Header().Get(tracing.RequestIDKey)` against `tc.headers["request-id"]`
    - [x] Run `bazel test //api/internal/tracing/...` to verify
