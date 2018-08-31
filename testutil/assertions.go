package testutil

import (
	"fmt"
	"testing"
)

func AssertFalse(t *testing.T, actual bool, what string) {
	if actual {
		t.Errorf("[ %s ]: Expected to be false", what)
	}
}

func AssertTrue(t *testing.T, actual bool, what string) {
	if !actual {
		t.Errorf("[ %s ]: Expected to be true", what)
	}
}

func AssertIntsEqual(t *testing.T, expected int, actual int, what string) {
	if expected != actual {
		t.Errorf("[ %s ]: Expected to be [ %s ], got [ %s ]", what, fmt.Sprintf("%v", expected), fmt.Sprintf("%v", actual))
	}
}

func AssertFloatsEqual(t *testing.T, expected float64, actual float64, what string) {
	if expected != actual {
		t.Errorf("[ %s ]: Expected to be [ %s ], got [ %s ]", what, fmt.Sprintf("%f", expected), fmt.Sprintf("%f", actual))
	}
}

func AssertStringEqual(t *testing.T, expected string, actual string, what string) {
	if expected != actual {
		t.Errorf("[ %s ]: Expected to be [ %s ], got [ %s ]", what, expected, actual)
	}
}
