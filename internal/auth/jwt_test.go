package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

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
