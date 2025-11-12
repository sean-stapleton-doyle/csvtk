package csveditor

import (
	"testing"
)

func TestStringEqualsStrategy(t *testing.T) {
	strategy := &StringEqualsStrategy{}

	tests := []struct {
		value   string
		pattern string
		want    bool
	}{
		{"test", "test", true},
		{"test", "Test", false},
		{"hello", "world", false},
	}

	for _, tt := range tests {
		got, err := strategy.Match(tt.value, tt.pattern)
		if err != nil {
			t.Errorf("Match() error = %v", err)
		}
		if got != tt.want {
			t.Errorf("Match(%q, %q) = %v, want %v", tt.value, tt.pattern, got, tt.want)
		}
	}
}

func TestStringContainsStrategy(t *testing.T) {
	strategy := &StringContainsStrategy{}

	tests := []struct {
		value   string
		pattern string
		want    bool
	}{
		{"hello world", "world", true},
		{"test", "xyz", false},
		{"John Doe", "Doe", true},
	}

	for _, tt := range tests {
		got, err := strategy.Match(tt.value, tt.pattern)
		if err != nil {
			t.Errorf("Match() error = %v", err)
		}
		if got != tt.want {
			t.Errorf("Match(%q, %q) = %v, want %v", tt.value, tt.pattern, got, tt.want)
		}
	}
}

func TestRegexStrategy(t *testing.T) {
	strategy := &RegexStrategy{}

	tests := []struct {
		value   string
		pattern string
		want    bool
		wantErr bool
	}{
		{"john@example.com", `@example\.com`, true, false},
		{"test123", `\d+`, true, false},
		{"hello", `[a-z]+`, true, false},
		{"test", `[`, false, true},
	}

	for _, tt := range tests {
		got, err := strategy.Match(tt.value, tt.pattern)
		if (err != nil) != tt.wantErr {
			t.Errorf("Match(%q, %q) error = %v, wantErr %v", tt.value, tt.pattern, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("Match(%q, %q) = %v, want %v", tt.value, tt.pattern, got, tt.want)
		}
	}
}

func TestNumericComparisonStrategy(t *testing.T) {
	tests := []struct {
		operator string
		value    string
		pattern  string
		want     bool
	}{
		{">", "30", "25", true},
		{">", "20", "25", false},
		{"<", "20", "25", true},
		{"<", "30", "25", false},
		{">=", "25", "25", true},
		{">=", "24", "25", false},
		{"<=", "25", "25", true},
		{"<=", "26", "25", false},
		{"==", "25", "25", true},
		{"==", "24", "25", false},
		{"!=", "24", "25", true},
		{"!=", "25", "25", false},
	}

	for _, tt := range tests {
		strategy := &NumericComparisonStrategy{Operator: tt.operator}
		got, err := strategy.Match(tt.value, tt.pattern)
		if err != nil {
			t.Errorf("Match() error = %v", err)
		}
		if got != tt.want {
			t.Errorf("Match(%s, %q, %q) = %v, want %v", tt.operator, tt.value, tt.pattern, got, tt.want)
		}
	}
}

func TestNumericComparisonStrategyNonNumeric(t *testing.T) {
	strategy := &NumericComparisonStrategy{Operator: ">"}

	got, err := strategy.Match("abc", "25")
	if err != nil {
		t.Errorf("Match() error = %v, expected no error for non-numeric value", err)
	}
	if got != false {
		t.Errorf("Match() = %v, want false for non-numeric value", got)
	}
}

func TestNewFilterStrategy(t *testing.T) {
	tests := []struct {
		operator string
		wantType string
	}{
		{"equals", "equals"},
		{"contains", "contains"},
		{"regex", "regex"},
		{">", "numeric->"},
		{"<", "numeric-<"},
		{">=", "numeric->="},
	}

	for _, tt := range tests {
		strategy := NewFilterStrategy(tt.operator)
		if strategy == nil {
			t.Errorf("NewFilterStrategy(%q) = nil", tt.operator)
		}
	}
}
