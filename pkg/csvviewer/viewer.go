package csvviewer

import (
	"encoding/csv"
	"fmt"
	"strings"

	"sean-stapleton-doyle/csvtk/pkg/csvparser"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type filterMode int

const (
	normalMode filterMode = iota
	filterInputMode
)

type Model struct {
	csv             *csvparser.CSV
	filteredCSV     *csvparser.CSV
	scrollOffsetRow int
	scrollOffsetCol int
	selectedRow     int
	width           int
	height          int
	filename        string
	mode            filterMode
	filterInput     string
	filterColumn    string
	statusMessage   string
	showingFiltered bool
}

func New(csv *csvparser.CSV, filename string) Model {
	return Model{
		csv:             csv,
		filteredCSV:     nil,
		scrollOffsetRow: 0,
		scrollOffsetCol: 0,
		selectedRow:     0,
		width:           80,
		height:          24,
		filename:        filename,
		mode:            normalMode,
		filterInput:     "",
		filterColumn:    "",
		statusMessage:   "",
		showingFiltered: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode == filterInputMode {
			return m.handleFilterInput(msg)
		}
		return m.handleNormalInput(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) handleNormalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	currentCSV := m.getCurrentCSV()
	maxRows := len(currentCSV.Records)

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "down", "j":
		if m.selectedRow < maxRows-1 {
			m.selectedRow++
			if m.selectedRow >= m.scrollOffsetRow+m.getVisibleRows() {
				m.scrollOffsetRow++
			}
		}
	case "up", "k":
		if m.selectedRow > 0 {
			m.selectedRow--
			if m.selectedRow < m.scrollOffsetRow {
				m.scrollOffsetRow--
			}
		}
	case "right", "l":
		if m.scrollOffsetCol < len(currentCSV.Header)-1 {
			m.scrollOffsetCol++
		}
	case "left", "h":
		if m.scrollOffsetCol > 0 {
			m.scrollOffsetCol--
		}
	case "pgdown":
		visibleRows := m.getVisibleRows()
		m.selectedRow += visibleRows
		if m.selectedRow >= maxRows {
			m.selectedRow = maxRows - 1
		}
		m.scrollOffsetRow = m.selectedRow - visibleRows + 1
		if m.scrollOffsetRow < 0 {
			m.scrollOffsetRow = 0
		}
	case "pgup":
		visibleRows := m.getVisibleRows()
		m.selectedRow -= visibleRows
		if m.selectedRow < 0 {
			m.selectedRow = 0
		}
		m.scrollOffsetRow = m.selectedRow
	case "home", "g":
		m.selectedRow = 0
		m.scrollOffsetRow = 0
	case "end", "G":
		m.selectedRow = maxRows - 1
		m.scrollOffsetRow = maxRows - m.getVisibleRows()
		if m.scrollOffsetRow < 0 {
			m.scrollOffsetRow = 0
		}
	case "c":

		if m.selectedRow < len(currentCSV.Records) {
			var builder strings.Builder
			writer := csv.NewWriter(&builder)
			writer.Write(currentCSV.Records[m.selectedRow])
			writer.Flush()
			clipboard.WriteAll(builder.String())
			m.statusMessage = "Row copied to clipboard"
		}
	case "f":

		m.mode = filterInputMode
		m.filterInput = ""
		m.statusMessage = "Filter mode: enter search text"
	case "r":

		m.showingFiltered = false
		m.filteredCSV = nil
		m.selectedRow = 0
		m.scrollOffsetRow = 0
		m.statusMessage = "Filter cleared"
	}
	return m, nil
}

func (m Model) handleFilterInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":

		if m.filterInput != "" {
			m.applyFilter()
		}
		m.mode = normalMode
	case "esc":

		m.mode = normalMode
		m.filterInput = ""
		m.statusMessage = "Filter cancelled"
	case "backspace":
		if len(m.filterInput) > 0 {
			m.filterInput = m.filterInput[:len(m.filterInput)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.filterInput += msg.String()
		}
	}
	return m, nil
}

func (m *Model) applyFilter() {
	filtered := &csvparser.CSV{
		Header:  m.csv.Header,
		Records: [][]string{},
	}

	searchTerm := strings.ToLower(m.filterInput)
	for _, record := range m.csv.Records {
		for _, cell := range record {
			if strings.Contains(strings.ToLower(cell), searchTerm) {
				filtered.Records = append(filtered.Records, record)
				break
			}
		}
	}

	m.filteredCSV = filtered
	m.showingFiltered = true
	m.selectedRow = 0
	m.scrollOffsetRow = 0
	m.statusMessage = fmt.Sprintf("Filtered: %d of %d rows", len(filtered.Records), len(m.csv.Records))
}

func (m Model) getCurrentCSV() *csvparser.CSV {
	if m.showingFiltered && m.filteredCSV != nil {
		return m.filteredCSV
	}
	return m.csv
}

func (m Model) getVisibleRows() int {
	visibleRows := m.height - 10
	if visibleRows < 1 {
		visibleRows = 10
	}
	return visibleRows
}

func (m Model) View() string {

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("63")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	altRowStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Padding(0, 1)

	selectedRowStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("229")).
		Padding(0, 1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(1, 0, 0, 0)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Padding(0, 0, 0, 0)

	filterInputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("63")).
		Padding(0, 1)

	var s strings.Builder

	currentCSV := m.getCurrentCSV()

	title := fmt.Sprintf(" CSV Viewer: %s ", m.filename)
	stats := fmt.Sprintf(" %d rows × %d columns ", len(currentCSV.Records), len(currentCSV.Header))
	s.WriteString(titleStyle.Render(title))
	s.WriteString(" ")
	s.WriteString(titleStyle.Render(stats))
	s.WriteString("\n\n")

	colWidths := make([]int, len(currentCSV.Header))
	for i, header := range currentCSV.Header {
		colWidths[i] = len(header)
	}
	for _, record := range currentCSV.Records {
		for i, cell := range record {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	for i := range colWidths {
		if colWidths[i] > 30 {
			colWidths[i] = 30
		}
		if colWidths[i] < 10 {
			colWidths[i] = 10
		}
	}

	visibleCols := m.getVisibleColumns(colWidths)

	var headerRow strings.Builder
	for i := m.scrollOffsetCol; i < m.scrollOffsetCol+visibleCols && i < len(currentCSV.Header); i++ {
		header := currentCSV.Header[i]
		width := colWidths[i]
		if len(header) > width {
			header = header[:width-3] + "..."
		}
		headerRow.WriteString(headerStyle.Render(fmt.Sprintf("%-*s", width, header)))
	}
	s.WriteString(headerRow.String())
	s.WriteString("\n")

	visibleRows := m.getVisibleRows()
	start := m.scrollOffsetRow
	end := start + visibleRows
	if end > len(currentCSV.Records) {
		end = len(currentCSV.Records)
	}

	for idx := start; idx < end; idx++ {
		record := currentCSV.Records[idx]
		var row strings.Builder
		for i := m.scrollOffsetCol; i < m.scrollOffsetCol+visibleCols && i < len(record); i++ {
			if i >= len(colWidths) {
				break
			}
			cell := record[i]
			width := colWidths[i]
			if len(cell) > width {
				cell = cell[:width-3] + "..."
			}

			var style lipgloss.Style
			if idx == m.selectedRow {
				style = selectedRowStyle
			} else if idx%2 == 1 {
				style = altRowStyle
			} else {
				style = cellStyle
			}
			row.WriteString(style.Render(fmt.Sprintf("%-*s", width, cell)))
		}
		s.WriteString(row.String())
		s.WriteString("\n")
	}

	s.WriteString("\n")
	if m.statusMessage != "" {
		s.WriteString(statusStyle.Render(m.statusMessage))
		s.WriteString("\n")
	}

	if m.mode == filterInputMode {
		s.WriteString(filterInputStyle.Render(fmt.Sprintf("Filter: %s_", m.filterInput)))
		s.WriteString("\n")
	}

	if m.mode == normalMode {
		help := "↑↓/jk: move • ←→/hl: scroll • PgUp/PgDn: page • g/G: top/bottom • c: copy • f: filter • r: reset • q: quit"
		s.WriteString(helpStyle.Render(help))
	} else {
		help := "Type to filter • Enter: apply • Esc: cancel"
		s.WriteString(helpStyle.Render(help))
	}

	if len(currentCSV.Records) > visibleRows {
		scrollInfo := fmt.Sprintf("\nRow %d of %d", m.selectedRow+1, len(currentCSV.Records))
		if m.scrollOffsetCol > 0 {
			scrollInfo += fmt.Sprintf(" • Col %d+", m.scrollOffsetCol+1)
		}
		s.WriteString(helpStyle.Render(scrollInfo))
	}

	return s.String()
}

func (m Model) getVisibleColumns(colWidths []int) int {
	availableWidth := m.width - 4
	totalWidth := 0
	visibleCols := 0

	for i := m.scrollOffsetCol; i < len(colWidths); i++ {
		totalWidth += colWidths[i] + 2
		if totalWidth > availableWidth {
			break
		}
		visibleCols++
	}

	if visibleCols == 0 {
		visibleCols = 1
	}
	return visibleCols
}

func Run(csv *csvparser.CSV, filename string) error {
	p := tea.NewProgram(New(csv, filename), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
