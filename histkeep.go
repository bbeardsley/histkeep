package histkeep

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

// HistKeep main struct
type HistKeep interface {
	AddValue(value string) error
	RemoveValue(value string) error
	GetValues() ([]string, error)
	GetFilteredValues(filterFunc func(string) bool) ([]string, error)
	ClearValues() error
	ReverseValues([]string) []string
}

type histKeep struct {
	filename     string
	numberToKeep int
	format       regexp.Regexp
}

// NewHistKeep creates a new HistKeep struct, pass null for format to match any string
func NewHistKeep(filename string, numberToKeep int, format *regexp.Regexp) HistKeep {
	if format == nil {
		format, _ = regexp.Compile(".*")
	}
	return &histKeep{filename, numberToKeep, *format}
}

func (histkeep *histKeep) AddValue(value string) error {
	if !histkeep.format.MatchString(value) {
		return fmt.Errorf("invalid format for value")
	}

	lines, err := readLines(histkeep.filename, value, histkeep.format)
	if err != nil {
		return err
	}

	lines = append(lines, value)

	lines, err = limitSlice(lines, histkeep.numberToKeep)
	if err != nil {
		return err
	}

	err = writeLines(histkeep.filename, lines)
	if err != nil {
		return err
	}
	return nil
}

func (histkeep *histKeep) RemoveValue(value string) error {
	lines, err := readLines(histkeep.filename, value, histkeep.format)
	if err != nil {
		return err
	}

	err = writeLines(histkeep.filename, lines)
	if err != nil {
		return err
	}
	return nil
}

func (histkeep *histKeep) ClearValues() error {
	lines := make([]string, 0)

	err := writeLines(histkeep.filename, lines)
	if err != nil {
		return err
	}

	return nil
}

func (histkeep *histKeep) GetValues() ([]string, error) {
	lines, err := readLines(histkeep.filename, "", histkeep.format)
	if err != nil {
		return nil, err
	}

	lines, err = limitSlice(lines, histkeep.numberToKeep)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func (histkeep *histKeep) GetFilteredValues(filterFunc func(string) bool) ([]string, error) {
	values, err := histkeep.GetValues()
	if err != nil {
		return nil, err
	}
	filtered := []string{}
	for _, line := range values {
		if filterFunc(line) {
			filtered = append(filtered, line)
		}
	}
	return filtered, nil
}

func (histkeep *histKeep) ReverseValues(values []string) []string {
	return reverseValues(values)
}

func readLines(path string, ignoreValue string, format regexp.Regexp) ([]string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return make([]string, 0), nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != ignoreValue && line != "" && format.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func limitSlice(lines []string, lastN int) ([]string, error) {
	linesLen := len(lines)
	if linesLen > lastN {
		return lines[linesLen-lastN : linesLen], nil
	}
	return lines, nil
}

func writeLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for i, line := range lines {
		if i != len(lines)-1 {
			fmt.Fprintln(w, line)
		} else {
			fmt.Fprint(w, line)
		}
	}
	return w.Flush()
}

type stringValue struct {
	value string
	index int
}

type byIndex []stringValue

func (b byIndex) Len() int           { return len(b) }
func (b byIndex) Less(i, j int) bool { return b[i].index < b[j].index }
func (b byIndex) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func reverseValues(values []string) []string {
	stringValues := make([]stringValue, 0)
	for i := 0; i < len(values); i++ {
		stringValues = append(stringValues, stringValue{values[i], i})
	}
	sort.Sort(sort.Reverse(byIndex(stringValues)))
	items := make([]string, 0)
	for i := 0; i < len(stringValues); i++ {
		items = append(items, stringValues[i].value)
	}
	return items
}
