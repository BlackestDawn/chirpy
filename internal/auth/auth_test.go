package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// base information
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	// base information
	testUUID := uuid.New()
	testToken, _ := MakeJWT(testUUID, "blahaj", 3*time.Second)

	tests := []struct {
		name     string
		token    string
		secret   string
		expected uuid.UUID
		wantErr  bool
	}{
		{
			name:     "Matching token",
			token:    testToken,
			secret:   "blahaj",
			expected: testUUID,
			wantErr:  false,
		},
		{
			name:     "Invalid token",
			token:    "something else",
			secret:   "blahaj",
			expected: uuid.Nil,
			wantErr:  true,
		},
		{
			name:     "Wrong secret",
			token:    testToken,
			secret:   "something else",
			expected: uuid.Nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retVal, err := ValidateJWT(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if retVal != tt.expected {
				t.Errorf("ValidateJWT() retVal = %v, want %v", retVal, tt.expected)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name     string
		token    http.Header
		expected string
		wantErr  bool
	}{
		{
			name: "Valid token",
			token: http.Header{
				"Authorization": []string{"Bearer 1234567890"},
			},
			expected: "1234567890",
			wantErr:  false,
		},
		{
			name:     "Missing token",
			token:    http.Header{},
			expected: "",
			wantErr:  true,
		},
		{
			name: "Invalid token",
			token: http.Header{
				"Authorization": []string{"Invalid 1234567890"},
			},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retVal, err := GetBearerToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if retVal != tt.expected {
				t.Errorf("GetbearerToken() retVal = %v, want %v", retVal, tt.expected)
			}
		})
	}
}
