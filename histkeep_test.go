package histkeep

import (
	"os"
	"strings"
	"testing"
)

func TestWithTwo(t *testing.T) {
	hk := NewHistKeep(testDataFile, 2, nil)

	defer deleteTestData()

	// first new value
	err := hk.AddValue("test-1")
	if err != nil {
		t.Fatal(err)
	}

	verifyData(t, hk, "test-1")

	// second new value
	err = hk.AddValue("test-2")
	if err != nil {
		t.Fatal(err)
	}

	verifyData(t, hk, "test-1\ntest-2")

	// adding first new value moves it up in priority
	err = hk.AddValue("test-1")
	if err != nil {
		t.Fatal(err)
	}

	verifyData(t, hk, "test-2\ntest-1")

	// new value pushes 2 out
	err = hk.AddValue("test-3")
	if err != nil {
		t.Fatal(err)
	}

	verifyData(t, hk, "test-1\ntest-3")

	// same most recent value
	err = hk.AddValue("test-3")
	if err != nil {
		t.Fatal(err)
	}

	verifyData(t, hk, "test-1\ntest-3")

	err = hk.RemoveValue("test-3")
	if err != nil {
		t.Fatal(err)
	}

	verifyData(t, hk, "test-1")

	err = hk.RemoveValue("test-1")
	if err != nil {
		t.Fatal(err)
	}

	verifyData(t, hk, "")

	err = hk.RemoveValue("test-1")
	if err != nil {
		t.Fatal(err)
	}

	verifyData(t, hk, "")

	hk.AddValue("test-1")
	hk.AddValue("test-2")
	hk.AddValue("test-3")

	verifyData(t, hk, "test-2\ntest-3")

	values, err := hk.GetValues()
	if err != nil {
		t.Fatal(err)
	}

	reversed := strings.Join(hk.ReverseValues(values), "\n")
	expectedReversed := "test-3\ntest-2"
	if reversed != expectedReversed {
		t.Fatalf("expected %s but got %s", expectedReversed, reversed)
	}

	err = hk.ClearValues()
	if err != nil {
		t.Fatal(err)
	}

	verifyData(t, hk, "")
}

func verifyData(t *testing.T, hk HistKeep, expected string) {
	actualValues, err := hk.GetValues()
	if err != nil {
		t.Fatal(err)
	}
	actual := strings.Join(actualValues, "\n")

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

var (
	testDataFile = "histkeep-test-data.txt"
)
