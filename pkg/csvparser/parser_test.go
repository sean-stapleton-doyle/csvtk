package csvparser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter rune
		wantRows  int
		wantCols  int
		wantErr   bool
	}{
		{
			name:      "basic CSV",
			input:     "Name,Age\nJohn,30\nJane,25",
			delimiter: ',',
			wantRows:  2,
			wantCols:  2,
			wantErr:   false,
		},
		{
			name:      "TSV",
			input:     "Name\tAge\nJohn\t30\nJane\t25",
			delimiter: '\t',
			wantRows:  2,
			wantCols:  2,
			wantErr:   false,
		},
		{
			name:      "empty file",
			input:     "",
			delimiter: ',',
			wantRows:  0,
			wantCols:  0,
			wantErr:   false,
		},
		{
			name:      "header only",
			input:     "Name,Age,City",
			delimiter: ',',
			wantRows:  0,
			wantCols:  3,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			config := &Config{Delimiter: tt.delimiter}

			csv, err := Parse(reader, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if csv.CountRows() != tt.wantRows {
					t.Errorf("CountRows() = %v, want %v", csv.CountRows(), tt.wantRows)
				}
				if csv.CountColumns() != tt.wantCols {
					t.Errorf("CountColumns() = %v, want %v", csv.CountColumns(), tt.wantCols)
				}
			}
		})
	}
}

func TestGetColumnIndex(t *testing.T) {
	csv := &CSV{
		Header: []string{"Name", "Email", "Age"},
	}

	tests := []struct {
		name       string
		columnName string
		wantIndex  int
		wantErr    bool
	}{
		{"found first", "Name", 0, false},
		{"found middle", "Email", 1, false},
		{"found last", "Age", 2, false},
		{"not found", "City", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, err := csv.GetColumnIndex(tt.columnName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetColumnIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if index != tt.wantIndex {
				t.Errorf("GetColumnIndex() = %v, want %v", index, tt.wantIndex)
			}
		})
	}
}

func TestWriteAndParse(t *testing.T) {
	original := &CSV{
		Header: []string{"Name", "Age", "City"},
		Records: [][]string{
			{"John", "30", "NYC"},
			{"Jane", "25", "LA"},
		},
	}

	var buf strings.Builder
	config := DefaultConfig()

	err := original.Write(&buf, config)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	parsed, err := Parse(strings.NewReader(buf.String()), config)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(parsed.Header) != len(original.Header) {
		t.Errorf("Header length mismatch: got %d, want %d", len(parsed.Header), len(original.Header))
	}

	if len(parsed.Records) != len(original.Records) {
		t.Errorf("Records length mismatch: got %d, want %d", len(parsed.Records), len(original.Records))
	}
}
