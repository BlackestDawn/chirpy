package auth

import (
	"net/http"
	"testing"
)

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

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		token    http.Header
		expected string
		wantErr  bool
	}{
		{
			name: "Valid APIKey",
			token: http.Header{
				"Authorization": []string{"ApiKey 1234567890"},
			},
			expected: "1234567890",
			wantErr:  false,
		},
		{
			name:     "Missing APIKey",
			token:    http.Header{},
			expected: "",
			wantErr:  true,
		},
		{
			name: "Invalid APIKey",
			token: http.Header{
				"Authorization": []string{"Invalid 1234567890"},
			},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retVal, err := GetAPIKey(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if retVal != tt.expected {
				t.Errorf("GetAPIKey() retVal = %v, want %v", retVal, tt.expected)
			}
		})
	}
}