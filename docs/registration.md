# Storage - database schema

```markdown
TBD
```

# Client/server communication - API schema

The user journey goes:

1. Type in an email address and password on screen 1
2. Backend checks if the email address already exists

This is separated out from the UserService because the RPCs defined there require an authenticated context - making it awkward
to set up CreateUser (which would need to be unauthenticated) next to DeleteUser and UpdateUser that need to be authenticated.

## Service definition

```protobuf
// registration_service.proto

syntax = "proto3";

package fjarm.registration.v1;

service RegistrationService {
    // TODO(2025-03-18): Remove "User" from all these messages. We register credentials so they cannot be reused. We do not register users. Update(2026-09-17): 
    // TODO(2025-03-18): Consider removing the fjarm.users.v1.UserHandle message everywhere. Handles increase the number of indexes needed on the user
    //   table and don't add much benefit when the app doesn't have any social features yet. It also makes the registration process unnecessarily
    //   complicated by increasing the number of fields that need to be checked for duplication - in addition to user ID and user email.
    // After receiving the registration credentials, the service always returns a - maybe invalid - registration token.
    rpc BeginUserRegistration(fjarm.registration.v1.BeginUserRegistrationRequest) returns (fjarm.registration.v1.BeginUserRegistrationResponse);
    // TODO(2025-03-18): This RPC receives the user input from the registration form and the registration token from the first RPC. Then, if the token
    //   exists and is valid, returns an OK and calls the CreateUser RPC.
    rpc CompleteUserRegistration(fjarm.registration.v1.CompleteUserRegistrationRequest) returns (fjarm.registration.v1.CompleteUserRegistrationResponse);

    // TODO(2026-09-17): Consider deleting because a registration token should only be valid for a few minutes. After such a time, we should prompt the user
    // to restart the registration process.
    // This RPC requires the email address that's being registered as well.
    rpc RefreshRegistrationToken(fjarm.registration.v1.RefreshRegistrationTokenRequest) returns (fjarm.registration.v1.RefreshRegistrationTokenResponse);
}
```

## Message definitions

### `begin_user_registration_request.proto`

```protobuf
syntax = "proto3";

package fjarm.registration.v1;

import "fjarm/idempotency/v1/idempotency_key.proto";
import "fjarm/registration/v1/registration_credentials.proto";

message BeginUserRegistrationRequest {
    optional fjarm.idempotency.v1.IdempotencyKey idempotency_key = 1 [(buf.validate.field).required = true];
    optional fjarm.registration.v1.RegistrationCredentials registration_credentials = 2 [(buf.validate.field).required = true];
}

```

### `registration_credentials.proto`

```protobuf
syntax = "proto3";

package fjarm.registration.v1;

import "fjarm/users/v1/user_email_address.proto";
import "fjarm/users/v1/user_password.proto";

message RegistrationCredentials {
    optional fjarm.users.v1.UserEmailAddress user_email_address = 1 [(buf.validate.field).required = true];
    optional fjarm.users.v1.UserPassword user_password = 2 [(buf.validate.field).required = true];
}
```

### `begin_user_registration_response.proto`

```protobuf
syntax = "proto3";

package fjarm.registration.v1;

import "fjarm/registration/v1/registration_token.proto";

message BeginUserRegistrationResponse {
    optional fjarm.registration.v1.RegistrationToken registration_token = 1 [(buf.validate.field).required = true];
}
```

The response is a token that will be sent back with the next request, which is one to complete registration.

### `registration_token.proto`

```protobuf
syntax = "proto3";

package fjarm.registration.v1;

import "buf/validate/validate.proto";
import "google/protobuf/timestamp.proto";

message RegistrationToken {
    // Required but may be gibberish if the registration credentials already exist.
    optional string registration_token = 1 [(buf.validate.field).required = true];

    // Indicates when the registration token above expires. Clients can proactively avoid sending requests with
    // expired tokens.
    // If a client sends a request with an expired token, the request will fail.
    optional google.protobuf.Timestamp expires_at = 2 [(buf.validate.field).required = true];
}
```

### `complete_user_registration_request.proto`

```protobuf
syntax = "proto3";

package fjarm.registration.v1;

```