package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CSV struct {
	Header  []string
	Records [][]string
}

type Config struct {
	Delimiter  rune
	LazyQuotes bool
	TrimSpace  bool
	SkipHeader bool
}

func DefaultConfig() *Config {
	return &Config{
		Delimiter:  ',',
		LazyQuotes: false,
		TrimSpace:  false,
		SkipHeader: false,
	}
}

func ParseFile(filename string, config *Config) (*CSV, error) {
	if config == nil {
		config = DefaultConfig()
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return Parse(file, config)
}

func Parse(reader io.Reader, config *Config) (*CSV, error) {
	if config == nil {
		config = DefaultConfig()
	}

	r := csv.NewReader(reader)
	r.Comma = config.Delimiter
	r.LazyQuotes = config.LazyQuotes
	r.TrimLeadingSpace = config.TrimSpace

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return &CSV{
			Header:  []string{},
			Records: [][]string{},
		}, nil
	}

	csvData := &CSV{}
	if config.SkipHeader {
		csvData.Header = []string{}
		csvData.Records = records
	} else {
		csvData.Header = records[0]
		if len(records) > 1 {
			csvData.Records = records[1:]
		} else {
			csvData.Records = [][]string{}
		}
	}

	return csvData, nil
}

func (c *CSV) CountRows() int {
	return len(c.Records)
}

func (c *CSV) CountColumns() int {
	if len(c.Header) > 0 {
		return len(c.Header)
	}
	if len(c.Records) > 0 {
		return len(c.Records[0])
	}
	return 0
}

func (c *CSV) GetColumnIndex(columnName string) (int, error) {
	for i, name := range c.Header {
		if name == columnName {
			return i, nil
		}
	}
	return -1, fmt.Errorf("column %q not found", columnName)
}

func (c *CSV) WriteToFile(filename string, config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return c.Write(file, config)
}

func (c *CSV) Write(writer io.Writer, config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	w := csv.NewWriter(writer)
	w.Comma = config.Delimiter

	if len(c.Header) > 0 {
		if err := w.Write(c.Header); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	for _, record := range c.Records {
		if err := w.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	w.Flush()
	return w.Error()
}
