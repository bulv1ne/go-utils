package channels

import (
	"slices"
	"testing"
)

func TestPipeSeqToChannel(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	seq := slices.Values(slice)

	channel := PipeSeqToChannel(seq, 1)
	value, ok := <-channel
	if !ok {
		t.Fatal("Expected channel to receive value, but it did not")
	}
	if value != 1 {
		t.Fatalf("Expected 1, got %d", value)
	}
	value, ok = <-channel
	if !ok {
		t.Fatal("Expected channel to receive value, but it did not")
	}
	if value != 2 {
		t.Fatalf("Expected 2, got %d", value)
	}
	value, ok = <-channel
	if !ok {
		t.Fatal("Expected channel to receive value, but it did not")
	}
	if value != 3 {
		t.Fatalf("Expected 3, got %d", value)
	}

	remainder := []int{}
	for i := range channel {
		remainder = append(remainder, i)
	}
	if !slices.Equal(remainder, []int{4, 5, 6, 7, 8, 9}) {
		t.Fatalf("Expected remainder, got %v", remainder)
	}
}
