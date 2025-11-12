package csvlint

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		delimiter     rune
		lazyquotes    bool
		wantErrors    int
		wantFatal     bool
		wantNonFatal  bool
	}{
		{
			name:       "valid CSV",
			input:      "Name,Age\nJohn,30\nJane,25",
			delimiter:  ',',
			lazyquotes: false,
			wantErrors: 0,
			wantFatal:  false,
		},
		{
			name:       "inconsistent columns",
			input:      "Name,Age\nJohn,30\nJane,25,Extra",
			delimiter:  ',',
			lazyquotes: false,
			wantErrors: 1,
			wantFatal:  false,
		},
		{
			name:       "multiple errors",
			input:      "Name,Age\nJohn,30,Extra\nJane,25,Extra,TooMany",
			delimiter:  ',',
			lazyquotes: false,
			wantErrors: 2,
			wantFatal:  false,
		},
		{
			name:       "valid TSV",
			input:      "Name\tAge\nJohn\t30\nJane\t25",
			delimiter:  '\t',
			lazyquotes: false,
			wantErrors: 0,
			wantFatal:  false,
		},
		{
			name:       "header only",
			input:      "Name,Age,City",
			delimiter:  ',',
			lazyquotes: false,
			wantErrors: 0,
			wantFatal:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			errors, fatal, err := Validate(reader, tt.delimiter, tt.lazyquotes)

			if err != nil {
				t.Fatalf("Validate() unexpected error = %v", err)
			}

			if fatal != tt.wantFatal {
				t.Errorf("Validate() fatal = %v, want %v", fatal, tt.wantFatal)
			}

			if len(errors) != tt.wantErrors {
				t.Errorf("Validate() got %d errors, want %d", len(errors), tt.wantErrors)
			}
		})
	}
}

func TestCSVError_Error(t *testing.T) {
	err := CSVError{
		Record: []string{"John", "30"},
		Num:    1,
		err:    &testError{"field count mismatch"},
	}

	expected := "Record #1 has error: field count mismatch"
	if err.Error() != expected {
		t.Errorf("Error() = %s, want %s", err.Error(), expected)
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
