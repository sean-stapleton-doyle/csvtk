package csveditor

import (
	"fmt"
	"strings"

	"sean-stapleton-doyle/csvtk/pkg/csvparser"
)

func MoveColumn(csv *csvparser.CSV, columnName string, targetIndex int) error {
	currentIndex, err := csv.GetColumnIndex(columnName)
	if err != nil {
		return err
	}

	if targetIndex < 0 || targetIndex >= len(csv.Header) {
		return fmt.Errorf("target index %d is out of range (0-%d)", targetIndex, len(csv.Header)-1)
	}

	if currentIndex == targetIndex {
		return nil
	}

	header := make([]string, len(csv.Header))
	copy(header, csv.Header)
	column := header[currentIndex]
	header = append(header[:currentIndex], header[currentIndex+1:]...)
	header = append(header[:targetIndex], append([]string{column}, header[targetIndex:]...)...)
	csv.Header = header

	for i := range csv.Records {
		if len(csv.Records[i]) <= currentIndex {
			continue
		}
		record := make([]string, len(csv.Records[i]))
		copy(record, csv.Records[i])
		value := record[currentIndex]
		record = append(record[:currentIndex], record[currentIndex+1:]...)
		record = append(record[:targetIndex], append([]string{value}, record[targetIndex:]...)...)
		csv.Records[i] = record
	}

	return nil
}

func MoveRow(csv *csvparser.CSV, oldIndex, newIndex int) error {

	if oldIndex < 0 || oldIndex >= len(csv.Records) {
		return fmt.Errorf("old index %d is out of range (0-%d)", oldIndex, len(csv.Records)-1)
	}
	if newIndex < 0 || newIndex >= len(csv.Records) {
		return fmt.Errorf("new index %d is out of range (0-%d)", newIndex, len(csv.Records)-1)
	}

	if oldIndex == newIndex {
		return nil
	}

	row := csv.Records[oldIndex]
	csv.Records = append(csv.Records[:oldIndex], csv.Records[oldIndex+1:]...)
	csv.Records = append(csv.Records[:newIndex], append([][]string{row}, csv.Records[newIndex:]...)...)

	return nil
}

type FilterConfig struct {
	ColumnName string
	Value      string
	Operation  FilterOperation
}

type FilterOperation int

const (
	Equals FilterOperation = iota
	Contains
	StartsWith
	EndsWith
	NotEquals
)

func Filter(csv *csvparser.CSV, config FilterConfig) (*csvparser.CSV, error) {
	columnIndex, err := csv.GetColumnIndex(config.ColumnName)
	if err != nil {
		return nil, err
	}

	filtered := &csvparser.CSV{
		Header:  csv.Header,
		Records: [][]string{},
	}

	for _, record := range csv.Records {
		if len(record) <= columnIndex {
			continue
		}

		value := record[columnIndex]
		match := false

		switch config.Operation {
		case Equals:
			match = value == config.Value
		case Contains:
			match = strings.Contains(value, config.Value)
		case StartsWith:
			match = strings.HasPrefix(value, config.Value)
		case EndsWith:
			match = strings.HasSuffix(value, config.Value)
		case NotEquals:
			match = value != config.Value
		}

		if match {
			filtered.Records = append(filtered.Records, record)
		}
	}

	return filtered, nil
}

func FilterWithStrategy(csv *csvparser.CSV, columnName string, pattern string, strategy FilterStrategy) (*csvparser.CSV, error) {
	columnIndex, err := csv.GetColumnIndex(columnName)
	if err != nil {
		return nil, err
	}

	filtered := &csvparser.CSV{
		Header:  csv.Header,
		Records: [][]string{},
	}

	for _, record := range csv.Records {
		if len(record) <= columnIndex {
			continue
		}

		value := record[columnIndex]
		match, err := strategy.Match(value, pattern)
		if err != nil {
			return nil, fmt.Errorf("filter error on row: %w", err)
		}

		if match {
			filtered.Records = append(filtered.Records, record)
		}
	}

	return filtered, nil
}

func SelectColumns(csv *csvparser.CSV, columnNames []string) (*csvparser.CSV, error) {

	indices := make([]int, len(columnNames))
	for i, name := range columnNames {
		index, err := csv.GetColumnIndex(name)
		if err != nil {
			return nil, err
		}
		indices[i] = index
	}

	selected := &csvparser.CSV{
		Header:  make([]string, len(columnNames)),
		Records: make([][]string, len(csv.Records)),
	}

	for i, index := range indices {
		selected.Header[i] = csv.Header[index]
	}

	for i, record := range csv.Records {
		selectedRecord := make([]string, len(indices))
		for j, index := range indices {
			if index < len(record) {
				selectedRecord[j] = record[index]
			}
		}
		selected.Records[i] = selectedRecord
	}

	return selected, nil
}

type SortConfig struct {
	ColumnName string
	Descending bool
}

func Sort(csv *csvparser.CSV, config SortConfig) error {
	columnIndex, err := csv.GetColumnIndex(config.ColumnName)
	if err != nil {
		return err
	}

	for i := 0; i < len(csv.Records); i++ {
		for j := i + 1; j < len(csv.Records); j++ {
			val1 := ""
			val2 := ""
			if columnIndex < len(csv.Records[i]) {
				val1 = csv.Records[i][columnIndex]
			}
			if columnIndex < len(csv.Records[j]) {
				val2 = csv.Records[j][columnIndex]
			}

			shouldSwap := false
			if config.Descending {
				shouldSwap = val1 < val2
			} else {
				shouldSwap = val1 > val2
			}

			if shouldSwap {
				csv.Records[i], csv.Records[j] = csv.Records[j], csv.Records[i]
			}
		}
	}

	return nil
}
