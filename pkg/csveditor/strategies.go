package csveditor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type FilterStrategy interface {
	Match(value string, pattern string) (bool, error)
	Name() string
}

type StringEqualsStrategy struct{}

func (s *StringEqualsStrategy) Match(value string, pattern string) (bool, error) {
	return value == pattern, nil
}

func (s *StringEqualsStrategy) Name() string {
	return "equals"
}

type StringContainsStrategy struct{}

func (s *StringContainsStrategy) Match(value string, pattern string) (bool, error) {
	return strings.Contains(value, pattern), nil
}

func (s *StringContainsStrategy) Name() string {
	return "contains"
}

type StringStartsWithStrategy struct{}

func (s *StringStartsWithStrategy) Match(value string, pattern string) (bool, error) {
	return strings.HasPrefix(value, pattern), nil
}

func (s *StringStartsWithStrategy) Name() string {
	return "starts-with"
}

type StringEndsWithStrategy struct{}

func (s *StringEndsWithStrategy) Match(value string, pattern string) (bool, error) {
	return strings.HasSuffix(value, pattern), nil
}

func (s *StringEndsWithStrategy) Name() string {
	return "ends-with"
}

type StringNotEqualsStrategy struct{}

func (s *StringNotEqualsStrategy) Match(value string, pattern string) (bool, error) {
	return value != pattern, nil
}

func (s *StringNotEqualsStrategy) Name() string {
	return "not-equals"
}

type RegexStrategy struct{}

func (s *RegexStrategy) Match(value string, pattern string) (bool, error) {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}
	return matched, nil
}

func (s *RegexStrategy) Name() string {
	return "regex"
}

type NumericComparisonStrategy struct {
	Operator string
}

func (s *NumericComparisonStrategy) Match(value string, pattern string) (bool, error) {
	val, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {

		return false, nil
	}

	target, err := strconv.ParseFloat(strings.TrimSpace(pattern), 64)
	if err != nil {
		return false, fmt.Errorf("pattern is not a valid number: %s", pattern)
	}

	switch s.Operator {
	case ">":
		return val > target, nil
	case "<":
		return val < target, nil
	case ">=":
		return val >= target, nil
	case "<=":
		return val <= target, nil
	case "==":
		return val == target, nil
	case "!=":
		return val != target, nil
	default:
		return false, fmt.Errorf("unknown operator: %s", s.Operator)
	}
}

func (s *NumericComparisonStrategy) Name() string {
	return fmt.Sprintf("numeric-%s", s.Operator)
}

func NewFilterStrategy(operator string) FilterStrategy {
	switch operator {
	case "equals", "eq":
		return &StringEqualsStrategy{}
	case "contains":
		return &StringContainsStrategy{}
	case "starts-with", "startswith":
		return &StringStartsWithStrategy{}
	case "ends-with", "endswith":
		return &StringEndsWithStrategy{}
	case "not-equals", "ne":
		return &StringNotEqualsStrategy{}
	case "regex", "regexp":
		return &RegexStrategy{}
	case ">", "gt":
		return &NumericComparisonStrategy{Operator: ">"}
	case "<", "lt":
		return &NumericComparisonStrategy{Operator: "<"}
	case ">=", "gte":
		return &NumericComparisonStrategy{Operator: ">="}
	case "<=", "lte":
		return &NumericComparisonStrategy{Operator: "<="}
	case "==":
		return &NumericComparisonStrategy{Operator: "=="}
	case "!=":
		return &NumericComparisonStrategy{Operator: "!="}
	default:
		return &StringEqualsStrategy{}
	}
}
