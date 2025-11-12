package csvlint

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CSVError struct {
	Record []string

	Num int
	err error
}

func (e CSVError) Error() string {
	return fmt.Sprintf("Record #%d has error: %s", e.Num, e.err.Error())
}

func Validate(reader io.Reader, delimiter rune, lazyquotes bool) ([]CSVError, bool, error) {
	r := csv.NewReader(reader)
	r.FieldsPerRecord = -1
	r.LazyQuotes = lazyquotes
	r.Comma = delimiter

	var header []string
	errors := []CSVError{}
	records := 0
	for {
		record, err := r.Read()
		if header != nil {
			records++
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			parsedErr, ok := err.(*csv.ParseError)
			if !ok {
				return errors, true, err
			}
			errors = append(errors, CSVError{
				Record: nil,
				Num:    records,
				err:    parsedErr.Err,
			})
			return errors, true, nil
		}
		if header == nil {
			header = record
			continue
		} else if len(record) != len(header) {
			errors = append(errors, CSVError{
				Record: record,
				Num:    records,
				err:    csv.ErrFieldCount,
			})
		}
	}
	return errors, false, nil
}

func ValidateFile(filename string, delimiter rune, lazyquotes bool) ([]CSVError, bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, false, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return Validate(file, delimiter, lazyquotes)
}
