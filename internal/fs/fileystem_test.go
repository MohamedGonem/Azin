package fs

import (
	"testing"
)
func TestFileSystem(t *testing.T) {
	tests := []struct{
		filepath string
		isValid bool
	}{
		{filepath: "main.az", isValid: true},
		{filepath: "print.cpp", isValid: false},
		{filepath: "", isValid: false},
		{filepath: ".az",isValid: false},	
	} 

	for _, tt := range tests {
		err := validateSourceFile(tt.filepath)
		if (err != nil && tt.isValid == true || err == nil && tt.isValid == false) {
			t.Errorf("validateSourceFile(%q) = isValid(%v), want(%v)",
			tt.filepath, !tt.isValid, tt.isValid)
		}
	}
}