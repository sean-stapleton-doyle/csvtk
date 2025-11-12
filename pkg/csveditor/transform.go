package csveditor

import (
	"fmt"
	"strings"

	"sean-stapleton-doyle/csvtk/pkg/csvparser"
)

type TransformFunc func(string) string

func TransformColumn(csv *csvparser.CSV, columnName string, transform TransformFunc) error {
	columnIndex, err := csv.GetColumnIndex(columnName)
	if err != nil {
		return err
	}

	for i := range csv.Records {
		if columnIndex < len(csv.Records[i]) {
			csv.Records[i][columnIndex] = transform(csv.Records[i][columnIndex])
		}
	}

	return nil
}

func TransformAll(csv *csvparser.CSV, transform TransformFunc) error {
	for i := range csv.Records {
		for j := range csv.Records[i] {
			csv.Records[i][j] = transform(csv.Records[i][j])
		}
	}
	return nil
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func Replace(old, new string, n int) TransformFunc {
	return func(s string) string {
		return strings.Replace(s, old, new, n)
	}
}

func ReplaceAll(old, new string) TransformFunc {
	return func(s string) string {
		return strings.ReplaceAll(s, old, new)
	}
}

func Trim(cutset string) TransformFunc {
	return func(s string) string {
		return strings.Trim(s, cutset)
	}
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func RenameHeader(csv *csvparser.CSV, oldName, newName string) error {
	index, err := csv.GetColumnIndex(oldName)
	if err != nil {
		return err
	}

	for i, name := range csv.Header {
		if i != index && name == newName {
			return fmt.Errorf("column %q already exists", newName)
		}
	}

	csv.Header[index] = newName
	return nil
}
