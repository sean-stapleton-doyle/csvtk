package csveditor

import (
	"testing"

	"sean-stapleton-doyle/csvtk/pkg/csvparser"
)

func TestTransformColumn(t *testing.T) {
	csv := &csvparser.CSV{
		Header: []string{"Name", "Email", "Age"},
		Records: [][]string{
			{"John Doe", "john@example.com", "30"},
			{"Jane Smith", "jane@example.com", "25"},
		},
	}

	err := TransformColumn(csv, "Name", ToUpper)
	if err != nil {
		t.Fatalf("TransformColumn() error = %v", err)
	}

	if csv.Records[0][0] != "JOHN DOE" {
		t.Errorf("TransformColumn() = %q, want %q", csv.Records[0][0], "JOHN DOE")
	}
	if csv.Records[1][0] != "JANE SMITH" {
		t.Errorf("TransformColumn() = %q, want %q", csv.Records[1][0], "JANE SMITH")
	}
}

func TestTransformAll(t *testing.T) {
	csv := &csvparser.CSV{
		Header: []string{"Name", "Email"},
		Records: [][]string{
			{"John Doe", "JOHN@EXAMPLE.COM"},
			{"Jane Smith", "JANE@EXAMPLE.COM"},
		},
	}

	err := TransformAll(csv, ToLower)
	if err != nil {
		t.Fatalf("TransformAll() error = %v", err)
	}

	if csv.Records[0][0] != "john doe" {
		t.Errorf("TransformAll() = %q, want %q", csv.Records[0][0], "john doe")
	}
	if csv.Records[0][1] != "john@example.com" {
		t.Errorf("TransformAll() = %q, want %q", csv.Records[0][1], "john@example.com")
	}
}

func TestReplaceAll(t *testing.T) {
	transform := ReplaceAll("@example.com", "@newdomain.com")

	result := transform("test@example.com")
	expected := "test@newdomain.com"

	if result != expected {
		t.Errorf("ReplaceAll() = %q, want %q", result, expected)
	}
}

func TestTrimSpace(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"  hello  ", "hello"},
		{"test", "test"},
		{"\n\tworld\n", "world"},
	}

	for _, tt := range tests {
		got := TrimSpace(tt.input)
		if got != tt.want {
			t.Errorf("TrimSpace(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRenameHeader(t *testing.T) {
	csv := &csvparser.CSV{
		Header:  []string{"Name", "Email", "Age"},
		Records: [][]string{},
	}

	err := RenameHeader(csv, "Email", "EmailAddress")
	if err != nil {
		t.Fatalf("RenameHeader() error = %v", err)
	}

	if csv.Header[1] != "EmailAddress" {
		t.Errorf("RenameHeader() = %q, want %q", csv.Header[1], "EmailAddress")
	}
}

func TestRenameHeaderErrors(t *testing.T) {
	csv := &csvparser.CSV{
		Header:  []string{"Name", "Email", "Age"},
		Records: [][]string{},
	}

	err := RenameHeader(csv, "Email", "Name")
	if err == nil {
		t.Error("RenameHeader() expected error for duplicate name, got nil")
	}

	err = RenameHeader(csv, "NonExistent", "NewName")
	if err == nil {
		t.Error("RenameHeader() expected error for non-existent column, got nil")
	}
}
