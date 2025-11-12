package csvparser

import (
	"fmt"
	"io"
	"os"
)

func ParseFromFileOrStdin(filename string, config *Config) (*CSV, error) {
	if config == nil {
		config = DefaultConfig()
	}

	var reader io.Reader
	if filename == "" || filename == "-" {
		reader = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		return Parse(file, config)
	}

	return Parse(reader, config)
}
