package fs

import (
	"path/filepath"
	"testing"
)

func TestValidateSourcePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantPath string
		wantErr  bool
	}{
		{
			name:     "valid",
			input:    "main.az",
			wantPath: "main.az",
			wantErr:  false,
		},
		{
			name:     "nested path is cleaned",
			input:    filepath.Join("src", ".", "main.az"),
			wantPath: filepath.Join("src", "main.az"),
			wantErr:  false,
		},
		{
			name:    "wrong extension",
			input:   "print.cpp",
			wantErr: true,
		},
		{
			name:    "missing extension",
			input:   "main",
			wantErr: true,
		},
		{
			name:    "only extension",
			input:   ".az",
			wantErr: true,
		},
		{
			name:     "hidden source file",
			input:    ".main.az",
			wantPath: ".main.az",
			wantErr:  false,
		},
		{
			name:     "multiple dots",
			input:    "foo.bar.az",
			wantPath: "foo.bar.az",
			wantErr:  false,
		},
		{
			name:     "uppercase extension",
			input:    "MAIN.AZ",
			wantPath: "MAIN.AZ",
			wantErr:  false,
		},
		{
			name:    "empty path",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, err := validateSourcePath(tt.input)

			if gotErr := err != nil; gotErr != tt.wantErr {
				t.Fatalf("validateSourcePath(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}

			if !tt.wantErr && gotPath != tt.wantPath {
				t.Errorf("validateSourcePath(%q) path = %q, want %q", tt.input, gotPath, tt.wantPath)
			}
		})
	}
}
