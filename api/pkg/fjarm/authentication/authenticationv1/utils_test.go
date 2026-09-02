package authenticationv1

import (
	"context"
	"testing"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"google.golang.org/protobuf/proto"
)

func TestValidateSessionID(t *testing.T) {
	ctx := context.Background()

	t.Run("nil session id", func(t *testing.T) {
		if err := ValidateSessionID(ctx, nil); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty session id field", func(t *testing.T) {
		if err := ValidateSessionID(ctx, &authenticationpb.SessionId{}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("valid session id", func(t *testing.T) {
		id := &authenticationpb.SessionId{SessionId: proto.String("sess-123")}
		if err := ValidateSessionID(ctx, id); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestValidateAccessToken(t *testing.T) {
	ctx := context.Background()

	t.Run("nil access token", func(t *testing.T) {
		if err := ValidateAccessToken(ctx, nil); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty access token field", func(t *testing.T) {
		if err := ValidateAccessToken(ctx, &authenticationpb.AccessToken{}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("valid access token", func(t *testing.T) {
		token := &authenticationpb.AccessToken{AccessToken: proto.String("valid-token")}
		if err := ValidateAccessToken(ctx, token); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestValidateRefreshToken(t *testing.T) {
	ctx := context.Background()

	t.Run("nil refresh token", func(t *testing.T) {
		if err := ValidateRefreshToken(ctx, nil); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty refresh token field", func(t *testing.T) {
		if err := ValidateRefreshToken(ctx, &authenticationpb.RefreshToken{}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("valid refresh token", func(t *testing.T) {
		token := &authenticationpb.RefreshToken{RefreshToken: proto.String("valid-refresh")}
		if err := ValidateRefreshToken(ctx, token); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestValidateSession(t *testing.T) {
	ctx := context.Background()

	t.Run("nil session", func(t *testing.T) {
		if err := ValidateSession(ctx, nil); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		if err := ValidateSession(ctx, &authenticationpb.Session{}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("valid session", func(t *testing.T) {
		session := &authenticationpb.Session{
			SessionId:    &authenticationpb.SessionId{SessionId: proto.String("sess-123")},
			AccessToken:  &authenticationpb.AccessToken{AccessToken: proto.String("access-123")},
			RefreshToken: &authenticationpb.RefreshToken{RefreshToken: proto.String("refresh-123")},
		}
		if err := ValidateSession(ctx, session); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestValidateCreateSessionRequest(t *testing.T) {
	ctx := context.Background()

	t.Run("nil request", func(t *testing.T) {
		if err := ValidateCreateSessionRequest(ctx, nil); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		if err := ValidateCreateSessionRequest(ctx, &authenticationpb.CreateSessionRequest{}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("valid request", func(t *testing.T) {
		req := &authenticationpb.CreateSessionRequest{
			IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174000"),
			EmailAddress:   &userspb.UserEmailAddress{EmailAddress: proto.String("user@example.com")},
			Password:       &userspb.UserPassword{Password: proto.String("password123")},
		}
		if err := ValidateCreateSessionRequest(ctx, req); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestValidateCreateSessionResponse(t *testing.T) {
	ctx := context.Background()

	t.Run("nil response", func(t *testing.T) {
		if err := ValidateCreateSessionResponse(ctx, nil); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("valid empty response", func(t *testing.T) {
		if err := ValidateCreateSessionResponse(ctx, &authenticationpb.CreateSessionResponse{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("valid response with session", func(t *testing.T) {
		res := &authenticationpb.CreateSessionResponse{
			Session: &authenticationpb.Session{
				SessionId:    &authenticationpb.SessionId{SessionId: proto.String("sess-123")},
				AccessToken:  &authenticationpb.AccessToken{AccessToken: proto.String("access-123")},
				RefreshToken: &authenticationpb.RefreshToken{RefreshToken: proto.String("refresh-123")},
			},
		}
		if err := ValidateCreateSessionResponse(ctx, res); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
