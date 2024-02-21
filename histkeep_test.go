package histkeep

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestAddValue(t *testing.T) {
	hk := NewHistKeep(testDataFile, 2, nil)

	// first new value
	err := hk.AddValue("test-1")
	if err != nil {
		t.Fatal(err)
	}

	defer deleteTestData()

	// second new value
	err = hk.AddValue("test-2")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := readTestData()

	if err != nil {
		t.Fatal(err)
	}

	expected := "test-1\ntest-2"

	if actual != expected {
		t.Fatalf("expected: %s but got %s", expected, actual)
	}

	// adding first new value moves it up in priority
	err = hk.AddValue("test-1")
	if err != nil {
		t.Fatal(err)
	}

	actual, err = readTestData()
	if err != nil {
		t.Fatal(err)
	}

	expected = "test-2\ntest-1"

	if actual != expected {
		t.Fatalf("expected: %s but got %s", expected, actual)
	}

	// new value pushes 2 out
	err = hk.AddValue("test-3")
	if err != nil {
		t.Fatal(err)
	}

	actual, err = readTestData()
	if err != nil {
		t.Fatal(err)
	}

	expected = "test-1\ntest-3"

	if actual != expected {
		t.Fatalf("expected: %s but got %s", expected, actual)
	}

	// same most recent value
	err = hk.AddValue("test-3")
	if err != nil {
		t.Fatal(err)
	}

	actual, err = readTestData()
	if err != nil {
		t.Fatal(err)
	}

	expected = "test-1\ntest-3"

	if actual != expected {
		t.Fatalf("expected: %s but got %s", expected, actual)
	}
}

func deleteTestData() error {
	_, err := os.Stat(testDataFile)
	if os.IsNotExist(err) {
		return nil
	}
	return os.Remove(testDataFile)
}

func readTestData() (string, error) {
	f, err := os.Open(testDataFile)
	if err != nil {
		return "", err
	}
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return "", scanner.Err()
	}

	return strings.Join(lines, "\n"), nil
}

var (
	testDataFile = "histkeep-test-data.txt"
)
