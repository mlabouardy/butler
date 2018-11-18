package main

import "testing"

func Test_getSanitizedUrl(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"Has http://", "http://google.de", "http://google.de"},
		{"Has https://", "https://google.de", "https://google.de"},
		{"Has no protocol", "google.de", "http://google.de"},
		{"Is empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSanitizedUrl(tt.url); got != tt.want {
				t.Errorf("getSanitizedUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
