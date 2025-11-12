# csvtk

A powerful command-line tool for working with CSV files, built with Go and Cobra.

## Features

- **View** - Interactive terminal viewer with keyboard navigation, horizontal scrolling, row selection, copy, and in-viewer filtering
- **Count** - Count rows and columns
- **Move** - Reorder rows and columns
- **Header** - Display and examine CSV headers
- **Rename** - Rename column headers
- **Convert** - Convert between CSV and TSV formats
- **Lint** - Validate CSV files according to RFC 4180
- **Filter** - Advanced filtering with regex and numeric comparisons (strategy pattern)
- **Select** - Extract specific columns
- **Sort** - Sort data by column values
- **Transform** - Transform data (uppercase, lowercase, replace, trim)
- **Stdin Support** - All commands support stdin for easy command chaining

## Installation

```bash
go build -o csvtk
```

## Usage

All commands support reading from stdin using `-` or by omitting the filename, enabling powerful command chaining.

### View CSV Files

Open an interactive viewer with keyboard navigation, horizontal scrolling, and filtering:

```bash
csvtk view myfile.csv
# Or simply:
csvtk myfile.csv
```

**Keyboard shortcuts:**
- `↑↓/jk`: Move up/down one row
- `←→/hl`: Scroll left/right through columns
- `PgUp/PgDn`: Navigate by page
- `g/home`: Jump to first row
- `G/end`: Jump to last row
- `c`: Copy selected row to clipboard
- `f`: Enter filter mode (search across all columns)
- `r`: Reset/clear filter
- `q`: Quit

**Filter mode:**
- Type to search (case-insensitive, searches all columns)
- `Enter`: Apply filter
- `Esc`: Cancel filter

### Count Operations

Count rows (excluding header):
```bash
csvtk count rows myfile.csv
# Or from stdin:
cat myfile.csv | csvtk count rows
```

Count columns:
```bash
csvtk count columns myfile.csv
```

### Move Operations

Move a column to a different position:
```bash
csvtk move column Email 0 myfile.csv
```

Move a row to a different position:
```bash
csvtk move row 5 0 myfile.csv
```

### Header Operations

Display header:
```bash
csvtk header myfile.csv
# Or from stdin:
cat myfile.csv | csvtk header
```

Display header with column indices:
```bash
csvtk header -n myfile.csv
```

Rename a header:
```bash
csvtk rename Email EmailAddress myfile.csv -o updated.csv
# Or use stdin/stdout:
cat myfile.csv | csvtk rename City Location - > updated.csv
```

### Convert Formats

Convert CSV to TSV:
```bash
csvtk convert myfile.csv --to-tsv
```

Convert TSV to CSV:
```bash
csvtk convert myfile.tsv --to-csv
```

Auto-detect format and convert:
```bash
csvtk convert myfile.csv  # Automatically converts to TSV
```

### Lint CSV Files

Validate CSV file structure:
```bash
csvtk lint myfile.csv
```

Use lazy quotes for less strict parsing:
```bash
csvtk lint --lazy-quotes myfile.csv
```

### Filter Operations

csvtk supports powerful filtering with multiple strategies:

**String Filters:**
```bash
# Exact match (default)
csvtk filter City "New York" data.csv

# Contains
csvtk filter Name "John" data.csv --operator contains

# Starts with
csvtk filter Email "admin" data.csv --operator starts-with

# Ends with
csvtk filter Email ".com" data.csv --operator ends-with

# Not equals
csvtk filter Status "inactive" data.csv --operator not-equals
```

**Regex Filters:**
```bash
# Regex matching
csvtk filter Email "@(gmail|yahoo)\.com" data.csv --regex

# Using operator flag
csvtk filter Phone "^\d{3}-\d{3}-\d{4}$" data.csv --operator regex
```

**Numeric Filters:**
```bash
# Greater than
csvtk filter Age 30 data.csv --operator ">"

# Less than
csvtk filter Price 100 data.csv --operator "<"

# Greater than or equal
csvtk filter Score 85 data.csv --operator ">="

# Less than or equal
csvtk filter Quantity 50 data.csv --operator "<="

# Numeric equality
csvtk filter Count 10 data.csv --operator "=="

# Numeric not equal
csvtk filter Amount 0 data.csv --operator "!="
```

**Output to file:**
```bash
csvtk filter Age 30 data.csv --operator ">" -o adults.csv
```

**Filter from stdin:**
```bash
cat data.csv | csvtk filter Age 30 - --operator ">" > adults.csv
```

### Select Columns

Extract specific columns:
```bash
csvtk select "Name,Email,Age" myfile.csv
```

Save to file:
```bash
csvtk select "Name,Email" myfile.csv -o output.csv
```

From stdin:
```bash
cat data.csv | csvtk select "Name,Age" -
```

### Sort Data

Sort by column (ascending):
```bash
csvtk sort Age myfile.csv -o sorted.csv
```

Sort in descending order:
```bash
csvtk sort Age myfile.csv -r -o sorted.csv
```

Sort to stdout:
```bash
csvtk sort Name myfile.csv
# Or from stdin:
cat myfile.csv | csvtk sort City -
```

### Transform Data

**Uppercase:**
```bash
# Transform specific column
csvtk transform upper Name data.csv -o output.csv

# Transform all columns
csvtk transform upper --all data.csv
```

**Lowercase:**
```bash
csvtk transform lower Email data.csv
```

**Replace text:**
```bash
# Replace in specific column
csvtk transform replace "@example.com" "@newdomain.com" Email data.csv

# Replace in all columns
csvtk transform replace "old" "new" --all data.csv -o output.csv
```

**Trim whitespace:**
```bash
csvtk transform trim Name data.csv
csvtk transform trim --all data.csv
```

**Transform with stdin:**
```bash
cat data.csv | csvtk transform lower Name - > output.csv
```

## Command Chaining Examples

One of the most powerful features is the ability to chain commands using stdin/stdout:

```bash
# Filter, select columns, and sort in one pipeline
cat data.csv | \
  csvtk filter Age 25 - --operator ">" | \
  csvtk select "Name,Email,Age" - | \
  csvtk sort Age - > result.csv

# Transform, rename, and output
csvtk transform lower Email data.csv | \
  csvtk rename Email EmailAddress - | \
  csvtk transform replace "@example.com" "@newdomain.com" EmailAddress - \
  > transformed.csv

# Count filtered results
csvtk filter City "New York" data.csv | \
  csvtk count rows

# Complex filtering with multiple conditions (using shell)
csvtk filter Age 30 data.csv --operator ">" | \
  csvtk filter City "New York" - | \
  csvtk select "Name,Age" - | \
  csvtk sort Age - -r
```

## Global Flags

Most commands support these flags:

- `-d, --delimiter`: Specify field delimiter (default: `,`)
  - Use `\t` or `\\t` for tab-delimited files
- `-o, --output`: Specify output file (defaults to stdout for most commands)

## Filter Operators

The filter command supports the following operators via the `--operator` flag:

**String operators:**
- `equals` or `eq` - Exact string match (default)
- `contains` - Substring match
- `starts-with` or `startswith` - Prefix match
- `ends-with` or `endswith` - Suffix match
- `not-equals` or `ne` - String inequality
- `regex` or `regexp` - Regular expression match

**Numeric operators:**
- `>` or `gt` - Greater than
- `<` or `lt` - Less than
- `>=` or `gte` - Greater than or equal
- `<=` or `lte` - Less than or equal
- `==` - Numeric equality
- `!=` - Numeric inequality

## Architecture

The project follows clean design principles with clear separation of concerns:

```
csvtk/
├── cmd/              # Cobra CLI commands
│   ├── root.go       # Root command
│   ├── count.go      # Count operations
│   ├── move.go       # Move operations
│   ├── header.go     # Header operations
│   ├── rename.go     # Rename headers
│   ├── convert.go    # Format conversion
│   ├── lint.go       # Validation
│   ├── filter.go     # Filtering with strategy pattern
│   ├── select.go     # Column selection
│   ├── sort.go       # Sorting
│   ├── transform.go  # Data transformation
│   ├── view.go       # Interactive viewer
│   └── utils.go      # Helper functions
├── pkg/
│   ├── csvparser/    # CSV parsing and I/O
│   │   ├── parser.go
│   │   ├── io.go     # Stdin/file reading utilities
│   │   └── parser_test.go
│   ├── csvlint/      # CSV validation
│   │   ├── lint.go
│   │   └── lint_test.go
│   ├── csveditor/    # CSV manipulation
│   │   ├── editor.go
│   │   ├── strategies.go    # Filter strategy pattern
│   │   ├── transform.go     # Transform functions
│   │   ├── editor_test.go
│   │   ├── strategies_test.go
│   │   └── transform_test.go
│   └── csvviewer/    # Interactive TUI viewer
│       └── viewer.go
└── main.go
```

### Design Principles

- **Separation of Concerns**: CLI logic is separated from business logic
- **Strategy Pattern**: Filter operations use the strategy pattern for extensibility
- **Parser Package**: Handles all CSV reading and writing operations
- **Editor Package**: Provides data manipulation capabilities
- **Transform Package**: Implements various data transformation functions
- **Lint Package**: Validates CSV structure and format
- **Viewer Package**: Provides interactive terminal UI with advanced features
- **Stdin Support**: All commands support stdin for Unix-style piping
- **Testability**: Each package has comprehensive unit tests

## Testing

Run all tests:

```bash
go test ./pkg/... -v
```

Run tests for a specific package:

```bash
go test ./pkg/csvparser -v
go test ./pkg/csveditor -v
go test ./pkg/csvlint -v
```

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal output
- [Clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard access

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.
