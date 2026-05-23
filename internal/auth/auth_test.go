package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	type testCase struct {
		name          string
		inputPassword string
		inputHash     string
		wantMatch     bool
	}

	password := "super-secured-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	if hash == "" {
		t.Fatalf("expected hash to be non empty string")
	}

	tests := []testCase{
		{
			name:          "Correct password",
			inputPassword: password,
			inputHash:     hash,
			wantMatch:     true,
		},
		{
			name:          "Wrong password",
			inputPassword: "wrong-password",
			inputHash:     hash,
			wantMatch:     false,
		},
		{
			name:          "Empty password",
			inputPassword: "",
			inputHash:     hash,
			wantMatch:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tc.inputPassword, tc.inputHash)
			if err != nil {
				t.Fatalf("CheckPasswordHash returned an unexpected error: %v", err)
			}

			if match != tc.wantMatch {
				t.Errorf("CheckPasswordHash() = %v, want %v", match, tc.wantMatch)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "token-secret", time.Minute*5)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			tokenString: validToken,
			tokenSecret: "token-secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotID, err := ValidateJWT(tc.tokenString, tc.tokenSecret)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if gotID != tc.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotID, tc.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Header",
			headers: http.Header{
				"Authorization": []string{"Bearer valid-jwt-token"},
			},
			wantToken: "valid-jwt-token",
			wantErr:   false,
		},
		{
			name:      "Missing Authorization Header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Wrong Auth Scheme (Basic instead of Bearer)",
			headers: http.Header{
				"Authorization": []string{"Basic dXNlcm5hbWU6cGFzc3dvcmQ="},
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed Header (No space after Bearer)",
			headers: http.Header{
				"Authorization": []string{"Bearermytoken"},
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed Header (Empty token after Bearer)",
			headers: http.Header{
				"Authorization": []string{"Bearer "},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tc.headers)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if gotToken != tc.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tc.wantToken)
			}
		})
	}
}
