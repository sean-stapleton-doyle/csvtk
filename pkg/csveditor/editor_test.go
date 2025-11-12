package csveditor

import (
	"testing"

	"sean-stapleton-doyle/csvtk/pkg/csvparser"
)

func TestMoveColumn(t *testing.T) {
	csv := &csvparser.CSV{
		Header: []string{"Name", "Email", "Age", "City"},
		Records: [][]string{
			{"John", "john@example.com", "30", "NYC"},
			{"Jane", "jane@example.com", "25", "LA"},
		},
	}

	err := MoveColumn(csv, "Age", 1)
	if err != nil {
		t.Fatalf("MoveColumn() error = %v", err)
	}

	expectedHeader := []string{"Name", "Age", "Email", "City"}
	for i, name := range csv.Header {
		if name != expectedHeader[i] {
			t.Errorf("Header[%d] = %s, want %s", i, name, expectedHeader[i])
		}
	}

	if csv.Records[0][1] != "30" {
		t.Errorf("First record Age = %s, want 30", csv.Records[0][1])
	}
}

func TestMoveRow(t *testing.T) {
	csv := &csvparser.CSV{
		Header: []string{"Name", "Age"},
		Records: [][]string{
			{"John", "30"},
			{"Jane", "25"},
			{"Bob", "35"},
		},
	}

	err := MoveRow(csv, 0, 2)
	if err != nil {
		t.Fatalf("MoveRow() error = %v", err)
	}

	if csv.Records[2][0] != "John" {
		t.Errorf("Records[2][0] = %s, want John", csv.Records[2][0])
	}

	if csv.Records[0][0] != "Jane" {
		t.Errorf("Records[0][0] = %s, want Jane", csv.Records[0][0])
	}
}

func TestFilter(t *testing.T) {
	csv := &csvparser.CSV{
		Header: []string{"Name", "Age", "City"},
		Records: [][]string{
			{"John", "30", "New York"},
			{"Jane", "25", "New York"},
			{"Bob", "35", "Chicago"},
		},
	}

	tests := []struct {
		name      string
		config    FilterConfig
		wantCount int
	}{
		{
			name: "equals filter",
			config: FilterConfig{
				ColumnName: "City",
				Value:      "New York",
				Operation:  Equals,
			},
			wantCount: 2,
		},
		{
			name: "contains filter",
			config: FilterConfig{
				ColumnName: "Name",
				Value:      "o",
				Operation:  Contains,
			},
			wantCount: 2,
		},
		{
			name: "starts with filter",
			config: FilterConfig{
				ColumnName: "Name",
				Value:      "J",
				Operation:  StartsWith,
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, err := Filter(csv, tt.config)
			if err != nil {
				t.Fatalf("Filter() error = %v", err)
			}

			if len(filtered.Records) != tt.wantCount {
				t.Errorf("Filter() got %d records, want %d", len(filtered.Records), tt.wantCount)
			}
		})
	}
}

func TestSelectColumns(t *testing.T) {
	csv := &csvparser.CSV{
		Header: []string{"Name", "Email", "Age", "City"},
		Records: [][]string{
			{"John", "john@example.com", "30", "NYC"},
			{"Jane", "jane@example.com", "25", "LA"},
		},
	}

	selected, err := SelectColumns(csv, []string{"Name", "Age"})
	if err != nil {
		t.Fatalf("SelectColumns() error = %v", err)
	}

	if len(selected.Header) != 2 {
		t.Errorf("Selected header length = %d, want 2", len(selected.Header))
	}

	if selected.Header[0] != "Name" || selected.Header[1] != "Age" {
		t.Errorf("Selected header = %v, want [Name Age]", selected.Header)
	}

	if len(selected.Records[0]) != 2 {
		t.Errorf("Selected record length = %d, want 2", len(selected.Records[0]))
	}

	if selected.Records[0][0] != "John" || selected.Records[0][1] != "30" {
		t.Errorf("Selected first record = %v, want [John 30]", selected.Records[0])
	}
}

func TestSort(t *testing.T) {
	csv := &csvparser.CSV{
		Header: []string{"Name", "Age"},
		Records: [][]string{
			{"Charlie", "30"},
			{"Alice", "25"},
			{"Bob", "35"},
		},
	}

	config := SortConfig{
		ColumnName: "Name",
		Descending: false,
	}

	err := Sort(csv, config)
	if err != nil {
		t.Fatalf("Sort() error = %v", err)
	}

	expectedOrder := []string{"Alice", "Bob", "Charlie"}
	for i, expected := range expectedOrder {
		if csv.Records[i][0] != expected {
			t.Errorf("Records[%d][0] = %s, want %s", i, csv.Records[i][0], expected)
		}
	}
}
