package tests

import "testing"

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}

func Substract(a, b int) int {
	return a - b
}

func TestSubstract(t *testing.T) {
	result := Substract(5, 3)
	expected := 2

	if result != expected {
		t.Errorf("Substract(5, 3) = %d; want %d", result, expected)
	}
}
