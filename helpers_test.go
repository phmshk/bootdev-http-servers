package main

import (
	"testing"
)

func TestGetCleanedBody(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected string
	}

	tests := []testCase{
		{
			name:     "Clean string without profanities",
			input:    "I had something interesting for breakfast",
			expected: "I had something interesting for breakfast",
		},
		{
			name:     "Single lowercase profanity",
			input:    "This is a kerfuffle opinion",
			expected: "This is a **** opinion",
		},
		{
			name:     "Uppercase profanity (case-insensitivity)",
			input:    "I really need a FORNAX to go to bed",
			expected: "I really need a **** to go to bed",
		},
		{
			name:     "Mixed case profanity",
			input:    "What a Sharbert day",
			expected: "What a **** day",
		},
		{
			name:     "Multiple profanities in one string",
			input:    "sharbert and fornax are bad words",
			expected: "**** and **** are bad words",
		},
		{
			name:     "Profanity with attached punctuation should NOT be replaced",
			input:    "Sharbert! That was a kerfuffle.",
			expected: "Sharbert! That was a kerfuffle.",
		},
		{
			name:     "Profanity with separated punctuation SHOULD be replaced",
			input:    "I go to bed sooner, Fornax !",
			expected: "I go to bed sooner, **** !",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getCleanedBody(tc.input)
			if actual != tc.expected {
				t.Errorf("\nFAIL: %s\nInput:    %q\nExpected: %q\nActual:   %q",
					tc.name, tc.input, tc.expected, actual)
			}
		})
	}
}
