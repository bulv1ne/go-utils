package monaderrors

import (
	"errors"
	"strconv"
	"testing"
)

func TestSome(t *testing.T) {
	safeDiv := func(x, y int) Option[int] {
		if y == 0 {
			return Option[int]{err: errors.New("division by zero")}
		}
		return Some(x / y)
	}

	result := Some(100).
		FlatMap(func(x int) Option[int] {
			return safeDiv(x, 2) // First division
		}).
		FlatMap(func(x int) Option[int] {
			return safeDiv(x, 5) // Second division
		})

	if result.value != 10 {
		t.Fatalf("Expected 10, got %v", result)
	}

	resultWithZeroDiv := Some(100).
		FlatMap(func(x int) Option[int] {
			return safeDiv(x, 0)
		})

	if !resultWithZeroDiv.IsError() {
		t.Fatalf("Expected error, got %v", resultWithZeroDiv)
	}
}

func count(value string) (int, error) {
	nameOption := Some(value)
	numberOption := Map(nameOption, strconv.Atoi)
	calculatedNumber := numberOption.
		Map(func(i int) (int, error) {
			return i + 1, nil
		})
	return calculatedNumber.Unwrap()
}

func TestCount(t *testing.T) {
	calculatedNumber, err := count("Niels")
	if err == nil {
		t.Fatalf("Expected error, got %v", calculatedNumber)
	}

	calculatedNumber, err = count("123")
	if err != nil {
		t.Fatalf("Expected error, got %v", err)
	}
	if calculatedNumber != 124 {
		t.Fatalf("Expected 124, got %v", calculatedNumber)
	}
}
