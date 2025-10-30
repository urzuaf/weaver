package main

import "testing"

func TestProcessFile(t *testing.T) {
	testCases := []struct {
		name       string
		filename   string
		shouldFail bool
	}{
		// Should pass (err == nil)
		{
			name:       "basic",
			filename:   "./tests/basic.wevr",
			shouldFail: false,
		},
		{
			name:       "complex",
			filename:   "./tests/complex.wevr",
			shouldFail: false,
		},
		{
			name:       "lists",
			filename:   "./tests/lists.wevr",
			shouldFail: false,
		},
		{
			name:       "references",
			filename:   "./tests/references.wevr",
			shouldFail: false,
		},
		{
			name:       "magnitudes",
			filename:   "./tests/magnitudes.wevr",
			shouldFail: false,
		}, {
			name:       "integrator1",
			filename:   "./tests/integrator1.wevr",
			shouldFail: false,
		},
		// Should fail (err != nil)
		{
			name:       "syntax error (missing '}')",
			filename:   "./tests/fail1.wevr",
			shouldFail: true,
		},
		{
			name:       "syntax error (unexpected token)",
			filename:   "./tests/fail2.wevr",
			shouldFail: true,
		},
		{
			name:       "invalid magnitude",
			filename:   "./tests/invalid_magnitude.wevr",
			shouldFail: true,
		},
		{
			name:       "invalid reference",
			filename:   "./tests/invalid_reference.wevr",
			shouldFail: true,
		},
		{
			name:       "invalid list",
			filename:   "./tests/invalidList.wevr",
			shouldFail: true,
		},
		{
			name:       "stray token",
			filename:   "./tests/stray_token.wevr",
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := processFile(tc.filename)

			if tc.shouldFail && err == nil {
				t.Errorf("Expected error processing file %s, but got nil", tc.filename)
			}

			if !tc.shouldFail && err != nil {
				t.Errorf("Did not expect error processing file %s, but got: %v", tc.filename, err)
			}
		})
	}
}
